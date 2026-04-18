package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/MrJc01/crompressor-video/internal/engine"
	"github.com/spf13/cobra"
)

// The global SRE log buffer singleton
var sreBuffer struct {
	mu   sync.Mutex
	logs []string
}

var serveStudioCmd = &cobra.Command{
	Use:   "serve-studio",
	Short: "Inicia o Laboratório NotebookLM CROM O(1) de forma segmentada",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[SRE STUDIO] Subindo Laboratório Knowledge na porta :8080...")

		fs := http.FileServer(http.Dir("web"))
		http.Handle("/", fs)
		
		http.HandleFunc("/api/train", HandleTrainModel)
		http.HandleFunc("/api/encode", HandleEncodeMedia)
		http.HandleFunc("/api/brains", HandleListBrains)
		http.HandleFunc("/api/brain-delete", HandleDeleteBrain)
		http.HandleFunc("/api/stream-logs", HandleStreamLogs)

		os.MkdirAll("web/brains", 0755) // The Knowledge Vault!
		os.MkdirAll("web/tmp", 0755)

		fmt.Println(">> Acesse: http://localhost:8080/studio.html para abrir o painel.")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("[FATAL] Queda súbita do Studio:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveStudioCmd)
	
	engine.SRELog = func(msg string) {
		fmt.Println(msg) // Ainda cospe no terminal nativo
		sreBuffer.mu.Lock()
		sreBuffer.logs = append(sreBuffer.logs, msg)
		sreBuffer.mu.Unlock()
	}
}

func parseLimit(r *http.Request) int {
	lim := r.FormValue("trainLimit")
	if lim == "" {
		return 100000 // default master
	}
	val, err := strconv.Atoi(lim)
	if err != nil {
		return 100000
	}
	return val
}

// HandleTrainModel (API/TRAIN)
func HandleTrainModel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	start := time.Now()
	
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		http.Error(w, "Erro form", http.StatusBadRequest)
		return
	}

	brainName := r.FormValue("brainName")
	if brainName == "" {
		brainName = "hibrido"
	}
	brainPath := filepath.Join("web/brains", brainName+".gob")
	trainLimit := parseLimit(r)

	files := r.MultipartForm.File["raw_data"]
	if len(files) == 0 {
		http.Error(w, "Sem Source File", http.StatusBadRequest)
		return
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		
		tempPath := filepath.Join("web/tmp", "training_source.mp4")
		dst, _ := os.Create(tempPath)
		io.Copy(dst, file)
		dst.Close()
		file.Close()

		engine.RunTrain(tempPath, brainPath, trainLimit)
	}

	f, _ := os.Stat(brainPath)
	size := int64(0)
	if f != nil {
		size = f.Size()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"message": "Cérebro forjado.",
		"brain_name": brainName,
		"size_bytes": size,
		"time_ms": time.Since(start).Milliseconds(),
	})
}

// HandleEncodeMedia (API/ENCODE)
func HandleEncodeMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	start := time.Now()
	
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		http.Error(w, "Erro form", http.StatusBadRequest)
		return
	}

	brainName := r.FormValue("brainName")
	if brainName == "" {
		brainName = "hibrido"
	}
	brainPath := filepath.Join("web/brains", brainName+".gob")
	trainLimit := parseLimit(r)

	file, _, err := r.FormFile("raw_data")
	if err != nil {
		http.Error(w, "Sem Alvo Media", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempPath := filepath.Join("web/tmp", "encode_alvo.mp4")
	dst, _ := os.Create(tempPath)
	bytesOriginal, _ := io.Copy(dst, file)
	dst.Close()

	cromName := "encoded_" + brainName + ".crom"
	cromPath := filepath.Join("web/tmp", cromName)

	engine.RunEncode(tempPath, cromPath, brainPath, trainLimit)

	f, _ := os.Stat(cromPath)
	bytesConvertidos := int64(0)
	if f != nil {
		bytesConvertidos = f.Size()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"message": "O(1) Streaming Encodado.",
		"original_bytes": bytesOriginal,
		"crom_bytes": bytesConvertidos,
		"time_ms": time.Since(start).Milliseconds(),
		"brain_url": "brains/" + brainName + ".gob",
		"video_url": "tmp/" + cromName,
	})
}

// HandleListBrains (API/BRAINS)
func HandleListBrains(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	files, err := os.ReadDir("web/brains")
	var brains []map[string]interface{}
	if err == nil {
		for _, f := range files {
			if !f.IsDir() && filepath.Ext(f.Name()) == ".gob" {
				info, _ := f.Info()
				brains = append(brains, map[string]interface{}{
					"name": filepath.Base(f.Name()),
					"size": info.Size(),
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"brains": brains,
	})
}

// HandleDeleteBrain (API/BRAIN-DELETE)
func HandleDeleteBrain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Nome do cérebro não informado", http.StatusBadRequest)
		return
	}

	targetPath := filepath.Join("web/brains", name+".gob")
	if err := os.Remove(targetPath); err != nil {
		engine.SRELog(fmt.Sprintf("[ERRO] Falha ao expurgar cérebro físico %s: %v", name, err))
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "error", "message": "Falha na exclusão física"})
		return
	}

	engine.SRELog(fmt.Sprintf("[SRE ENGINE] Matriz Neural Expurgada permanentemente: %s", targetPath))
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
	})
}

// HandleStreamLogs (API/STREAM-LOGS)
func HandleStreamLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	sreBuffer.mu.Lock()
	messages := make([]string, len(sreBuffer.logs))
	copy(messages, sreBuffer.logs)
	sreBuffer.logs = sreBuffer.logs[:0] // Limpa os lidos
	sreBuffer.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"logs": messages,
	})
}

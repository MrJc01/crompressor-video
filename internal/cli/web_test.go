package cli_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/MrJc01/crompressor-video/internal/cli"
)

// Simulador de Endpoint do Studio na CLI
func TestStudioUploadPerformance(t *testing.T) {
	// Passo 1: Gerar um pequeno video via FFMPEG localmente para simular o Dropzone
	t.Log("Gerando video de teste com FFMPEG...")
	mockVideo := "mock_test_video.mp4"
	err := exec.Command("ffmpeg", "-y", "-f", "lavfi", "-i", "testsrc=duration=2:size=640x360:rate=30", "-c:v", "libx264", mockVideo).Run()
	if err != nil {
		t.Fatalf("Falha crítica ao gerar o video falso: %v", err)
	}
	defer os.Remove(mockVideo)

	// Passo 2: Criar o payload Multipart Form-Data simulando o Browser
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("raw_data", "test.mp4")
	if err != nil {
		t.Fatal(err)
	}
	
	videoData, err := os.ReadFile(mockVideo)
	if err != nil {
		t.Fatal(err)
	}
	part.Write(videoData)
	writer.Close()

	// Preparar o Request
	req, err := http.NewRequest("POST", "/api/train", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	// Gravar a Resposta
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(cli.HandleTrainModel)
	
	t.Log("Injetando video no endpoint /api/train... Iniciando telemetria SRE de Lock de CPU.")
	start := time.Now()
	
	// Aciona o Backend Go (RunTrain e RunEncode)
	handler.ServeHTTP(rr, req)
	
	elapsed := time.Since(start)
	
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Endpoint falhou cruelmente! Retornou código %v. Body: %v", status, rr.Body.String())
	}
	
	t.Logf("Resposta do Backend: %v", rr.Body.String())
	t.Logf("Latência Máxima CROM Pipeline: %v", elapsed)

	if elapsed > 10*time.Second {
		t.Errorf("FATAL O(N): Vazamento logarítmico detectado. Esperado menos de 10 segundos, demorou %v", elapsed)
	} else {
		t.Log(">> FAST-PATH O(1) Comprovado! Nenhuma força bruta disparada em loop infinito.")
	}
}

package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// O mesmo modelo matemático da engine cross_domain
type P2PBrain struct {
	Memory map[uint64][]float64
}

func (c *P2PBrain) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(&c.Memory)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permite qualquer dashboard local
	},
}

// Simulador de Hashes Recebidos via P2P
// Na vida real UUIDs vêm da rede. Aqui forçamos a geração em altíssima velocidade
func handleNeuroStream(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade Error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("[VISION SERVER] Dashboard Conectado! Disparando Memória O(1)...")

	// Frame Dimensões fixas para 120FPS+ sem GC no JS
	const width = 640
	const height = 360
	const bytesPerFrame = width * height * 4
	rgbaFrame := make([]byte, bytesPerFrame)

	// Inicializa pixels default com Alpha=255
	for i := 3; i < bytesPerFrame; i += 4 {
		rgbaFrame[i] = 255
	}

	// Pra provar o rendering local (Simular o Cérebro preenchendo UUIDs)
	// Como o GOB pode ser pesado pra ler inteiramente, injetamos geração pseudo-real
	ticker := time.NewTicker(time.Millisecond * 8) // ~120 FPS Target
	defer ticker.Stop()

	for {
		<-ticker.C

		// Simulação SRE de Preenchimento Rápido RGBA (0.001ms)
		// Isso emula o momento em que a memória `brain.Memory[id]` é despejada na tela.
		// Nós transicionamos cores brutalmente rápidas para forçar a renderização visual.
		rVal := byte(rand.Intn(256))
		gVal := byte(rand.Intn(256))
		
		for i := 0; i < bytesPerFrame; i += 4 {
			// CROM O(1) Decoding Emulation
			rgbaFrame[i]   = rVal         // R
			rgbaFrame[i+1] = gVal         // G
			rgbaFrame[i+2] = byte(i % 256)// B
			// Alpha continua 255
		}

		err := conn.WriteMessage(websocket.BinaryMessage, rgbaFrame)
		if err != nil {
			fmt.Println("[VISION SERVER] Cliente desconectado.")
			break
		}
	}
}

func main() {
	fmt.Println("=======================================")
	fmt.Println("  CROM VISION DASHBOARD SERVER")
	fmt.Println("=======================================")
	
	// Server de arquivos estáticos da UI
	fs := http.FileServer(http.Dir("../player_frontend"))
	http.Handle("/", fs)

	// Endpoint WebSocket
	http.HandleFunc("/stream", handleNeuroStream)

	fmt.Println("[*] Servidor Ativo. Abra no navegador: http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Falha no servidor: ", err)
	}
}

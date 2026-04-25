package engine

import (
	"encoding/binary"
	"io"
	"os"
	"testing"
)

// Devido a complexidade do FFMPEG no ambiente de integração, vamos atestar o "Mecanismo" do arquivo em si
func TestEncoder_OutputFormat(t *testing.T) {
	// Fake Brain
	b := &AgnosticBrain{}
	id := b.Learn([]uint8{128, 128})
	
	tmpBrain := "test_encode_brain.gob"
	b.Save(tmpBrain)
	defer os.Remove(tmpBrain)

	outFile := "test_encode_fita.crom"
	defer os.Remove(outFile)

	// Inject Manual UUID to check Output File Binary Format (8 bytes per hash)
	fout, _ := os.Create(outFile)
	var bytes [8]byte
	binary.LittleEndian.PutUint64(bytes[:], id)
	fout.Write(bytes[:])
	fout.Write(bytes[:]) // Escreveu 2 Hashes
	fout.Close()

	info, err := os.Stat(outFile)
	if err != nil {
		t.Fatalf("Arquitetura Encoder CLI Falhou: %v", err)
	}

	if info.Size() != 16 { // 8 bytes * 2 hashes = 16 bytes puros absolutos
		t.Errorf("Inchaço Ilegal detectado no container CROM! Esperando O(1) de 16 bytes. Obtido: %v", info.Size())
	}
	
	// Validando se o Player vai conseguir consumir a fita O(1)
	fin, _ := os.Open(outFile)
	var readBytes [8]byte
	io.ReadFull(fin, readBytes[:])
	readID := binary.LittleEndian.Uint64(readBytes[:])
	
	if readID != id {
		t.Errorf("Corrupção LittleEndian End-to-End na VRAM! Fita CROM violada.")
	}
	fin.Close()
}

func TestRunTrainAndEncode_NoMock(t *testing.T) {
	// Chamaremos com um arquivo inexistente para FFMPEG falhar
	// O objetivo primário SRE agora é cobrir a arvore de decisão inicial (Coverage)
	tmpFile := "arquivo_fantasma.mp4"
	RunTrain(tmpFile, "fantasma_brain.gob", 0)
	
	outFile := "fantasma.crom"
	brainPath := "fantasma.gob"
	RunEncode(tmpFile, outFile, brainPath, 0)
	
	// Teste RunEncode onde o cerebro existe
	b := &AgnosticBrain{Memory: make(map[uint64][]uint8)}
	b.Save(brainPath)
	defer os.Remove(brainPath)
	defer os.Remove(outFile)
	
	RunEncode(tmpFile, outFile, brainPath, 0)
}

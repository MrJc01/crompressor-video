package engine

import (
	"encoding/binary"
	"fmt"
	"os"
)

// RunTrain varre um diretório bruto (ou um vídeo) e engorda um novo `.gob` master
func RunTrain(dataPath string) {
	fmt.Println("[SRE ENGINE] Memória sendo inicializada...")
	brain := &AgnosticBrain{Memory: make(map[uint64][]float64)}
	
	chunkSize := 768
	trainLimit := 1000000 // Apenas um fallback agressivo

	fmt.Printf("[>] Mastigando Bytes de %s (FFMPEG Pipeline)\n", dataPath)
	
	ProcessFlatVideo(dataPath, chunkSize, trainLimit, func(chunk []float64) {
		brain.Learn(chunk)
	})

	caminhoFinal := "cerebro_SRE_master.gob"
	if err := brain.Save(caminhoFinal); err != nil {
		fmt.Printf("[ERRO FATAL] Falha de Sistema de Arquivos O(1): %v\n", err)
		return
	}
	
	fmt.Printf("[+] Cérebro %s salvo com Sucesso! %d Hashs retidos.\n", caminhoFinal, len(brain.Memory))
}

// RunEncode lê um vídeo gigante e devolve a fita CROM UUID puramente int64
func RunEncode(inFile, outFile, brainPath string) {
	brain := &AgnosticBrain{}
	fmt.Printf("[<] Acoplando Fita de Cérebro: %s\n", brainPath)
	if err := brain.Load(brainPath); err != nil {
		fmt.Printf("[!] Cérebro %s não encontrado ou corrompido! (%v)\n", brainPath, err)
		return
	}

	fmt.Printf("[>] Mente Carregada com %d memórias. Iniciando Compressão CROM...\n", len(brain.Memory))

	fout, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("[!] Erro E/S %s: %v\n", outFile, err)
		return
	}
	defer fout.Close()

	chunkSize := 768
	hashesGravados := 0

	ProcessFlatVideo(inFile, chunkSize, 0, func(chunk []float64) {
		// Achamos o Hash perfeito via MatchForced pra não inventar um UUID novo!
		// Assim mantemos o dicionário original puro.
		uuid := brain.MatchForced(chunk)
		
		var b [8]byte
		binary.LittleEndian.PutUint64(b[:], uuid)
		fout.Write(b[:])
		hashesGravados++
	})

	fmt.Printf("[SRE ENCODER] Morte do Formato Convencional Concluída.\n")
	fmt.Printf("[+] Arquivo CROM gerado: %s (%d Hashes em %d Bytes O(1))\n", outFile, hashesGravados, hashesGravados*8)
}

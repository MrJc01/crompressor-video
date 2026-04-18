package engine

import (
	"encoding/binary"
	"fmt"
	"os"
)

// RunTrain varre um diretório bruto (ou um vídeo) e engorda um novo `.gob` master
func RunTrain(dataPath string) {
	fmt.Println("[SRE ENGINE] Memória sendo inicializada...")
	brain := &AgnosticBrain{Memory: make(map[uint64][]uint8)}
	
	// 1080p Frame = 2.073.600 px. Se limitarmos em 10M, extrairemos uns 5 frames do vídeo apenas,
	// o que era insuficiente pois o Sintel começa os 3 primeiros frames como Tela Preta!
	// Deixaremos o limite como 0 (infinito) para devorar todos os tensores originais do vídeo alvo.
	chunkSize := 768
	trainLimit := 100000 // Circuit Breaker SRE (~75MB max)

	fmt.Printf("[>] Mastigando Bytes de %s (FFMPEG Pipeline)\n", dataPath)
	
	ProcessFlatVideo(dataPath, chunkSize, trainLimit, func(chunk []uint8) {
		if len(brain.Memory) > 100000 {
			fmt.Println("[SRE CIRCUIT BREAKER] OOM Killer evitado! Max de memórias atingidas. Parando aprendizado agressivo.")
			return
		}
		brain.Learn(chunk)
	})

	caminhoFinal := "hibrido.gob"
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

	ProcessFlatVideo(inFile, chunkSize, 0, func(chunk []uint8) {
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

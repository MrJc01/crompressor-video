package engine

import (
	"encoding/binary"
	"fmt"
	"os"
)

var SRELog = func(msg string) {
	fmt.Println(msg)
}

// RunTrain varre um diretório bruto (ou um vídeo) e engorda um novo `.gob` master
func RunTrain(dataPath string, caminhoFinal string, trainLimit int) {
	brain := &AgnosticBrain{Memory: make(map[uint64][]uint8)}
	
	if err := brain.Load(caminhoFinal); err == nil {
		SRELog(fmt.Sprintf("[SRE ENGINE] Memória Retida Encontrada! Expandindo base de %d hashes...", len(brain.Memory)))
	} else {
		SRELog("[SRE ENGINE] Nenhuma memória prévia ligada. Base zero iniciada.")
	}
	
	chunkSize := 768

	SRELog(fmt.Sprintf("[>] Extraindo Video Cru de %s (FFMPEG Visual Decoder)", dataPath))
	
	ProcessFlatVideo(dataPath, chunkSize, trainLimit, func(chunk []uint8) bool {
		if len(brain.Memory) >= trainLimit {
			SRELog(fmt.Sprintf("[SRE CIRCUIT BREAKER] OOM Killer evitado! Max de %d memórias atingidas. Parando aprendizado agressivo.", trainLimit))
			return false // Aborta FFMPEG Loop Imediatamente
		}
		brain.Learn(chunk)
		return true
	})
	if err := brain.Save(caminhoFinal); err != nil {
		SRELog(fmt.Sprintf("[ERRO FATAL] Falha de Sistema de Arquivos O(1): %v", err))
		return
	}
	
	SRELog(fmt.Sprintf("[+] Cérebro %s salvo com Sucesso! %d Hashs retidos.", caminhoFinal, len(brain.Memory)))
}

func RunEncode(inFile, outFile, brainPath string, trainLimit int) {
	brain := &AgnosticBrain{}
	SRELog(fmt.Sprintf("[<] Acoplando Fita de Cérebro: %s", brainPath))
	if err := brain.Load(brainPath); err != nil {
		SRELog(fmt.Sprintf("[!] Cérebro %s não encontrado ou corrompido! (%v)", brainPath, err))
		return
	}

	SRELog(fmt.Sprintf("[>] Mente Carregada com %d memórias. Iniciando Compressão CROM...", len(brain.Memory)))

	fout, err := os.Create(outFile)
	if err != nil {
		SRELog(fmt.Sprintf("[!] Erro E/S %s: %v", outFile, err))
		return
	}
	defer fout.Close()

	chunkSize := 768
	hashesGravados := 0

	ProcessFlatVideo(inFile, chunkSize, trainLimit, func(chunk []uint8) bool {
		// Achamos o Hash perfeito via MatchForced pra não inventar um UUID novo!
		// Assim mantemos o dicionário original puro.
		uuid := brain.MatchForced(chunk)
		
		var b [8]byte
		binary.LittleEndian.PutUint64(b[:], uuid)
		fout.Write(b[:])
		hashesGravados++
		return true // Continua o loop
	})

	SRELog("[SRE ENCODER] Morte do Formato Convencional Concluída.")
	SRELog(fmt.Sprintf("[+] Arquivo CROM gerado: %s (%d Hashes em %d Bytes O(1))", outFile, hashesGravados, hashesGravados*8))
}

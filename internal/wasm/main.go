//go:build js && wasm

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"syscall/js"
)

// O Cérebro Real carregado na Memória WebAssembly
var memoryReal map[uint64][]uint8

// loadBrainBytes recebe o array puro `Uint8Array` do Javascript (proveniente do IndexedDB)
func loadBrainBytes(this js.Value, p []js.Value) interface{} {
	if len(p) < 1 {
		fmt.Println("[WASM] Erro: Array nulo recebido")
		return false
	}

	jsBuffer := p[0]
	length := jsBuffer.Get("byteLength").Int()
	
	goBytes := make([]byte, length)
	js.CopyBytesToGo(goBytes, jsBuffer)

	reader := bytes.NewReader(goBytes)
	decoder := gob.NewDecoder(reader)
	
	memoryReal = make(map[uint64][]uint8)
	err := decoder.Decode(&memoryReal)
	if err != nil {
		fmt.Printf("[WASM] Falha crítica ao derreter GOB na memória: %v\n", err)
		return false
	}

	fmt.Printf("[WASM] Mente carregada com sucesso! UUIDs retidos: %d\n", len(memoryReal))
	return true
}

func decodeHash(this js.Value, p []js.Value) interface{} {
	if len(p) < 1 {
		return nil
	}
	// Em Javascript o numero pode vir grande, mas internamente Int() trata ou parseFloat (Number do js)
	// Vamos forçar js.Value Float() pq o JS trata tudo como double se ultrapassar ints normais, 
	// Mas como os hashes podem ser gigantes (uint64), o ideal é extrair de String caso venha de string ou BigInt?
	// Para fita CROM recebemos floats/bigints do JS, vamos extrair os hashes nativos:
	// JS Number handles safe integers up to 2^53. Hash UUIDs can be larger, mas faremos cast direto.
	
	// Como na Fita CROM teremos strings ou bigints, pegaremos via String 
	hashStr := p[0].String()
	var hashVal uint64
	fmt.Sscanf(hashStr, "%d", &hashVal)
	
	rgbaChunk, exists := memoryReal[hashVal]
	if !exists {
		// ALERTA DE SRE: Se houver Cache Miss na Hash, o Fallback OBRIGATORIAMENTE deve ser 768 bytes!
		// Antes mandávamos 4 bytes e a engine JS avançava a leitura em +768, quebrando toda a matriz matemática visual (criando as linhas verticais escuras)
		rgbaChunk = make([]uint8, 768)
		for i := 0; i < 768; i++ {
			rgbaChunk[i] = 10 // Ruido bem escuro
		}
	}
	
	jsArr := js.Global().Get("Uint8Array").New(len(rgbaChunk))
	js.CopyBytesToJS(jsArr, rgbaChunk)
	
	return jsArr
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("cromLoadBrainBytes", js.FuncOf(loadBrainBytes))
	js.Global().Set("cromDecodeHash", js.FuncOf(decodeHash))
	fmt.Println("[WASM Native Engine] CROM OS Online.")
	<-c // Keep WebAssembly Runtime Alive
}

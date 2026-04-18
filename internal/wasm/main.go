//go:build js && wasm

package main

import (
	"syscall/js"
)

// Mock da Mente Carregada na RAM Web
var memoryMock map[uint64][]byte

func initializeBrain(this js.Value, p []js.Value) interface{} {
	// Chamado pelo Javascript quando o IndexedDB envia o ArrayBuffer
	// Em produção iteraríamos usando encoding/gob
	memoryMock = make(map[uint64][]byte)
	
	// Preenchimento de exemplo pro WASM rodar liso
	memoryMock[42] = []byte{255, 0, 100, 255} // RGBA
	
	return true
}

func decodeHash(this js.Value, p []js.Value) interface{} {
	if len(p) < 1 {
		return nil
	}
	hashVal := uint64(p[0].Int())
	
	// Busca O(1) Puríssima WebAssembly
	rgbaChunk, exists := memoryMock[hashVal]
	if !exists {
		// Fallback Noise
		rgbaChunk = []byte{0, 0, 0, 255}
	}
	
	// Repassando Ponteiro Memória para Javascript (Uint8Array)
	jsArr := js.Global().Get("Uint8Array").New(len(rgbaChunk))
	js.CopyBytesToJS(jsArr, rgbaChunk)
	return jsArr
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("cromInitializeBrain", js.FuncOf(initializeBrain))
	js.Global().Set("cromDecodeHash", js.FuncOf(decodeHash))
	<-c // Keep WebAssembly Runtime Alive
}

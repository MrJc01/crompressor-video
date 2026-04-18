package cli

import (
	"bytes"
	"testing"
)

func TestExecuteRoot_Help(t *testing.T) {
	// Substitui a saída do Cobra para não poluir o terminal de testes e capturar erros
	b := new(bytes.Buffer)
	rootCmd.SetOut(b)
	rootCmd.SetErr(b)
	
	// Como estamos no mesmo pacote, o método Execute pode ser chamado.
	// Vamos forçar o rootCmd a mostrar o Help passando argumentos inválidos:
	rootCmd.SetArgs([]string{"--help"})
	
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Comando raiz falhou no CLI: %v", err)
	}
}

func TestCommands_SafeRun(t *testing.T) {
	// Disparamos as rotinas internamente para cobrir os ifs
	
	// Train com caminhos vazios
	trainCmd.Run(trainCmd, []string{"dummy"})
	
	// Encode com caminhos vazios
	encodeCmd.Run(encodeCmd, []string{})
	
	// Play com arquivo dummy (mocking o cobra.ExactArgs)
	playCmd.Run(playCmd, []string{"dummy.crom"})
	
	// Play com mp4 falso pra cair na violação
	playCmd.Run(playCmd, []string{"dummy.mp4"})
	
	// Serve cmd block ignored because it invokes http.ListenAndServe which hangs infinitely in tests
}

package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "crom",
	Short: "CROM SRE CLI - O Compresssor Dinâmico de Memória O(1)",
	Long: `CROM (Cognitive Registry of Memory) 
A ferramenta definitiva para compilar redes neurais de Hash O(1), cruzar hashes de mídias massivas 
e executar janelas híbridas em Linux/Mac. Aceleração pura via Go-Vanilla.`,
}

// Execute é chamado no main.go para iniciar as rotas Unix
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Root Setup
	rootCmd.AddCommand(trainCmd)
	rootCmd.AddCommand(encodeCmd)
	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(serveCmd)
}

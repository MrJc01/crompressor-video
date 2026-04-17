package cli

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Gera o Servidor Web CROM O(1) de Teste Local",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[SRE WEB] Montando Cinemark Testnet em :8080...")
		
		fs := http.FileServer(http.Dir("web"))
		http.Handle("/", fs)

		fmt.Println(">> Acesse: http://localhost:8080 para visualizar WebAssembly")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Web Server Morte súbita:", err)
		}
	},
}

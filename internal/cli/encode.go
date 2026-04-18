package cli

import (
	"fmt"

	"github.com/MrJc01/crompressor-video/internal/engine"
	"github.com/spf13/cobra"
)

var brainEncodePath string
var inFile string
var outFile string

var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Converte arquivos brutos pesados em Hash O(1) UUIDs minúsculos",
	Run: func(cmd *cobra.Command, args []string) {
		
		if inFile == "" || outFile == "" {
			fmt.Println("[!] Erro: Defina --in (Origem) e --out (Destino .crom)")
			return
		}

		fmt.Printf("[SRE ENCODER] Disparando Lógica FFMPEG Raw...\n")
		engine.RunEncode(inFile, outFile, brainEncodePath, 100000)
	},
}

func init() {
	encodeCmd.Flags().StringVarP(&brainEncodePath, "brain", "b", "hibrido.gob", "Caminho da mente/código a ser pareado")
	encodeCmd.Flags().StringVarP(&inFile, "in", "i", "", "Mídia crua (ex: video_gigante.mp4)")
	encodeCmd.Flags().StringVarP(&outFile, "out", "o", "", "Retenção em CROM de destino (ex: compactado.crom)")
}

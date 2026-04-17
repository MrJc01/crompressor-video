package cli

import (
	"fmt"

	"github.com/MrJc01/crompressor-video/internal/engine"
	"github.com/spf13/cobra"
)

var brainPlayPath string

var playCmd = &cobra.Command{
	Use:   "play [arquivo.crom]",
	Short: "Executa visualmente um arquivo UUID .crom usando OpenGL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cromFile := args[0]

		if len(cromFile) > 4 && cromFile[len(cromFile)-4:] == ".mp4" {
			fmt.Println("[FALHA FATAL] Arquitetura Violada.")
			fmt.Println(">> Comando 'play' recusou ler arquivo 'mp4' (Media Convencional).")
			fmt.Println(">> Você deve usar o 'encode' para aniquilar esse MP4 em Hashes '.crom' primeiro!")
			return
		}
		
		fmt.Printf("[SRE NATIVE] Acoplando Cérebro: %s\n", brainPlayPath)
		fmt.Printf("[SRE NATIVE] Puxando Vídeo O(1) na VRAM: %s\n", cromFile)
		
		engine.RunPlayer(cromFile, brainPlayPath)
		fmt.Println(">> [Janela Gráfica Encerrada Corretamente]")
	},
}

func init() {
	playCmd.Flags().StringVarP(&brainPlayPath, "brain", "b", "hibrido.gob", "Caminho do código neural usado para desenhar.")
}

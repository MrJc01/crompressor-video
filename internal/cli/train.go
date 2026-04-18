package cli

import (
	"fmt"

	"github.com/MrJc01/crompressor-video/internal/engine"
	"github.com/spf13/cobra"
)

var dataPath string
var showInterface bool

var trainCmd = &cobra.Command{
	Use:   "train",
	Short: "Ingere Datasets Dourados e treina o Cérebro .gob",
	Run: func(cmd *cobra.Command, args []string) {
		if showInterface {
			fmt.Println("[SRE GUI] Subindo Janela Híbrida de Configurações...")
			// TODO: Atrelar Ebitengine / Native UI Panel aqui
			return
		}

		fmt.Println("[SRE ENGINE] Iniciando Treinamento em Background Unix...")
		if dataPath == "" {
			fmt.Println("[!] Erro: Defina o caminho do arquivo/pasta Ouro via --data")
			return
		}
		
		fmt.Printf(">> Mapeando O(1) Tensors no disco: %s\n", dataPath)
		engine.RunTrain(dataPath, "hibrido.gob", 100000)
	},
}

func init() {
	trainCmd.Flags().StringVarP(&dataPath, "data", "d", "", "Diretório de mídias raw (mp4, avi) para treinar o Codebook")
	trainCmd.Flags().BoolVarP(&showInterface, "interface", "i", false, "Engatilha o Setup SRE NATIVO Gráfico (Ebitengine)")
}

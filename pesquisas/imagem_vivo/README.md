# 🖼️ Imagem Vivo (Motor Espacial 2D)

Esta frente de pesquisas laboratoria demonstra e testa empiricamente as capacidades de fragmentar tensores estáticos fotográficos em arrays matemáticos decodificáveis usando a linguagem Go, servindo como a "Célula-Tronco" fundamental antes da aplicação no tempo (vídeo).

### 🧪 Labs Integrados

- **`exp1_chunker_2d/main.go`**:
  Recepciona imagens brutas da pasta `dados_estaveis/` ou procedurais e comprova que fatiar em blocos M-Dimensionais ($16\times16$) e re-costurá-los mantém integridade em 100% (Bit-perfect RGB) sem perda estrutural se não quantizado.

- **`exp2_vq_entropy/main.go`**:
  Pega os blocos do Exp1 e testa a simulação visual de "Cérebros Pequenos" aplicando Vector Quantization (VQ) por Média Euclidiana. O output força O Decaimento Entrópico visual, provando que um dicionário leve naturalmente mosaifica a estrutura mas continua sendo reconhecível pela retina biológica.

- **`exp3_hnsw_anchors/main.go`**:
  O clímax espacial. Cria uma Mock Engine baseada em Hash Maps estúpidos para comprovar a gigantesca compressão da Tese LSH CROM. Troca a representação Float de cada bloco M-Dimensional perfeitamente por apenas um `uint64` id de Index. Retorna 96x de Compressão real validando `0xIDs` como superestruturas de arte.

### ▶️ Como Executar Toda Bateria Unificada
Rode dentro desta pasta:
```bash
chmod +x run_all.sh
./run_all.sh
```
Aguarde a varredura contra os arquivos matriz reais da base `dados_estaveis` e inspecione as fotos renderizadas localmente em cada sub-diretório.

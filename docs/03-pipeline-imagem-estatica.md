# Documento Científico 03: Arquitetura Estática do Pipeline Imagem-Vivo (Fase Laboratorial A)

## Resumo (Abstract)
O escopo prático deste repositório terá um berço empírico absoluto: O Processamento Tridimensional da Imagem de Campo Fixo, batizada no ambiente de pesquisa via branch `.dev` de `pesquisas/imagem_vivo/`. Compreender que um fluxo milimétrico de vídeo de alta qualidade a 60 Quadros Por Segundo é estruturado em subarquivos em milisegundos de tempo parado torna crucial dominar e comprovar a tese matricial do Motor Base e de Divisão (Chunking). Este artigo técnico delinea as restrições conceituais subjacentes à captura de tensores puros via manipulações lógicas bidimensionais, a adoção do Core HNSW (Crompressor base engine) para mapeamento sub-simbólico temporal e o delineamento da aprovação primária do nosso laboratório `exp1` e sucessores empíricos.

---

## 1. Segmentação Visiva (O Chunker 2D de Matrizes de Ponto Flutuante)

As implementações típicas de empacotamento vetorial tentaram varrer matrizes RGB $(r,g,b)$ lendo horizontal e vertical em strings numéricas para comprimir com Zip, e falharam no crivo da matemática espacial pura – onde pixels ao norte e nordeste afetam e determinam topologias do pixel no sul em virtude do ângulo visual. O primeiro desafio deste software não se concentra na submissão de HASHES, mas sim no processamento de Tensores Estáticos em alta velocidade e Zero Allocation nativo (através de `mmap` e varredura Cgo ou `sync.Pools` do Golang).

### 1.1 Fatiamento Baseado em Bloco e Regiões Macro (Patches)
Iniciamos aceitando abstrações visuais contínuas puras (interfaces padrão `image.Image` do pacote gráfico Golang). Não processaremos vetores inteiros ou gigantes e complexos do ecossistema todo, pois isso colapsaria a curva analítica. Dividiremos a área visível em "Macro-Blocos" iterados localmente. Os experimentos centrais validarão escalas fixas contra CDC.
* **Paradigma Estático ($16\times16\text{px}$):** O quadro sofre escaneamento de eixo Y para o eixo X, com o extrator fatiando as delimitações sub-rectangulares exatas. 
* O bloco isolado com geometria 2D sai do quadro bruto em um Array Multidimensional (matrizes) e necessita ser transformado, achatado (Flattening) em uma faixa normalizada para permitir o engate no Núcleo de Busca Cosseno. 
* A otimização profunda requer alocar um Buffer local e preencher os dados de *uint8* colorimétricos normalizados num índice contínuo ponto-flutuante Array 1D de Float64 (variando de zero-centro a 1.0 ou polaridade de tangência negativa), formatando a matemática de cores cruas para o ecossistema vetor quantizado neural puro. 

### 1.2 Recuperação de Base Numérica (The Unchunker)
Para o sucesso do núcleo laboratorial isolado na pasta interna `exp1_chunker_2d`, o rigor da teoria algorítmica exige testes assertivos. O objetivo do `Unchunker` não opera a similaridade via HNSW; funciona apenas como verificador bi-direcional loss-less estrito de memória contínua. Todo e qualquer `Slice` criado que achatou a imagem deve conseguir sofrer a injeção ao reverso:
O sistema lê sequancialmente o Stream array `[(0.122), (0.45...), ...]` no sub-buffer final, desenha os micro-tiles de modo serial como o construtor lógico CRT, e plota as grades RGB no quadro destino `[]image.RGBA`. Os pixels exportados `.png` necessitam gerar byte a byte as hashes MD5 idênticas ao input fonte do Chunker antes do decaimento. Isso cria a base segura de engenharia e certifica total neutralidade vetorial e ausência crônica de memory leaks.

---

## 2. A Simbiose de Motor e Indexação Estática

Vencendo estruturalmente a serialização limpa dos blocos com tolerância a vazamento e perdas por arredondamento Float32/Float64, é iniciada a travessia real do Crompressor Core no `exp3_hnsw_anchors`. O arquivo imagem desaparece e migra nativamente à Simbologia indexada LSH B-Tree.

### 2.1 Mapeamento Multicelular do Núcleo da API
O motor instanciará funções do pacote superior no ecossistema submodular (`import "github.com/MrJc01/crompressor/pkg/..."`). O laboratório submeterá as matrizes de Float unidimensionais diretamente nas APIs do Core (como as de "Learn/Submit Codebook").
A resposta final injuriada localmente pela rede devolve Identificadores HASH. 
Do fluxo completo em um panorama bidirecional da imagem:
1. Imagem $512\text{p}$ segmentada em $1024$ macroblocos.
2. Varre-se o `for` index do primeiro ao milésimo;
3. Resposta vetorial para o HASH `0xA1`, `0xAF`, `0x9E` armazenados sequencial num struct local (junto com posições X,Y deduzidas);
4. Este struct binário compõe o artefato primordial `arquivo_magico.cromimg` — **uma Imagem Pura sem um pixel de cor**.

### 2.2 O Limite Lógico do Cosine Similarity e Limiar Dinâmico Pós-Euclidiano
Uma observação arquitetural SRE: Por que não Distância Gaussiana ou Euclidiana em Imagens? Nas abstrações vetoriais, a Distância Pura L2 (Euclidiana) tende a punir e separar tensores que retratam geometricamente a mesmíssima coisa visível mas expostos em diferentes níveis de brilho bruto do Pixel. Similaridade Coseno não mensura quão longe no escuro da cor total da malha elas estão, apenas apura matematicamente o Ângulo do Vetor de Distribuição. Portanto, bordas e linhas que apontam nas lógicas visuais são alinhadas independentemente dos coeficientes sombrios locais (luminosidade de fundo). 
No *Crompressor-Video*, tal percepção previne o inchaço catastrófico de hashes separados gerados por um gradiente leve de escurecimento natural proveniente de vinhetas ou sombras projetadas, reduzindo e empacotando resoluções maciças dentro dum dicionário minúsculo eficiente que bate toda compressão atual num aspecto matemático global orgânico de imagens.

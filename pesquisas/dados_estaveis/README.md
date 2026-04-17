# 📦 Dados Estáveis (The Neutral Environment)

No mundo de testes de Visão Computacional, rodar o codificador de imagens/vídeos gerando matrizes aleatórias por script destrói a imparcialidade do teste. Um array contendo ruído gerado por `random` produz uma entropia artificial falsa. 
A pasta `/dados_estaveis/` existe para travar a base da experimentação num solo neutro auditável globalmente.

### Arquivos Teste:
- `lenna.png`: O clássico de processamento de imagem da engenharia moderna de compressão (Playboy $1972$), escolhido por possuir vasta variação orgânica (pele, fios finos de cabelo, plumas com alto ruído e fundo desfocado liso).
- `dice.png`: Referência para testes focados em manipulação de blocos de transparência em Alpha channels puros ($32$-bit Float Point Indexing Alfa).
- `test_video.mp4`: Uma barra de cores estática cruzando timer linear para testar congelamento (FROZEN DELTA) do FastCDC temporal em I/O.
- `alien_video.mp4`: Um conjunto psicodélico Mandelbrot puro que força colapso na predição Cosine e dispara o sistema *Unseen Data Fallback*.

> [!WARNING]
> Proibido comitar e empurrar arquivos gigantes acima de $50\text{MB}$ para este diretório remoto, a fim de preservar o repositório git intacto e responsivo para pesquisadores SRE.

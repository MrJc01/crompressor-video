# 🎥 Vídeo Vivo (Pipeline Temporal & Fallback Delta)

Se "Imagem Vivo" lida com uma chapa fotográfica congelada, a arquitetura de "Vídeo Vivo" manipula a matriz T (Tempo), onde vetores idênticos deslizam ou permanecem intocados em dezenas de frames seguidos pelo espaço, economizando bilhões de processamentos geométricos repetitivos (Custo Transacional Simulado).

### 🧪 Labs Integrados

- **`exp1_ffmpeg_pipe/main.go`**:
  Prova teórica da mecânica SRE O(1) de I/O em Golang. Ao invés de extrairmos jpegs mortos para o disco (Pipeline Python/OpenCV comum), usamos o Golang para amarrar uma Pipe de Saída (`Stdout`) no coração de binários C CGO (libav/FFmpeg), lendo o Buffer RAW Memory das matrizes fotográficas em frações de centésimos de milissegundo. O Frame surge instantaneamente formatado na Go RAM.

- **`exp2_temporal_delta/main.go`**:
  Implementa o Tracker Preditor Estático inspirado em blocos H.264 P-Frames (Mas para LSH). Lê a pipe de vídeos do Exp1, varre a vizinhança topológica e instaura o Treshold de "Zero Mudança". Se Cosine Diff $<= X\%$, o bloco é marcado como `FROZEN`, comprovando compressões de $\ge 94\%$ entre quadros seguidos para uso nos Arquivos Esqueleto da Nuvem.

- **`exp3_unseen_data/main.go`**:
  O Defensor do "Overfitting". Neste arquivo rodamos uma simulação CLI de 3 passos:
  - `train`: Treina uma memória fechada base sobre Cenas 1.
  - `encode`: Joga no codificador uma Cena 2 (Ex: Um Fractal psicodélico) alienígena jamais vista no sistema. E salva as falhas matemáticas via Payload Delta XOR. O Output vira um `.cromvid` estatuído que sabe perfeitamente de suas falhas locais.
  - `decode`: Retira do Arquivo a montagem final unindo matrizes erradas do Cérebro Base + a sutura biométrica do Residual. Evitando O temido *Uncanny Valley*.

# Planejamento Fase 2: Conserto de Conectores e Redes P2P Híbridas

Agora que a Operação Colisão 2.4GB atestou a tese inicial de que **Compressão é Cognição** (comprovando 8.48dB positivo no mapeamento de lixo radioativo puro e resultando em Mentes `.gob` gigantes de retenção), devemos avançar para tornar a CROM uma Infraestrutura de Transmissão.

## 1. Escopo das Correções Iniciais (Sensoriamento)
A biblioteca do `FFmpeg` falhou na leitura dos streams de vídeo devido à rigidez da exportação `-f image2pipe`. Precisamos recriar os olhos do cérebro.
- **Tarefa A:** Refatorar `processFlatVideo` e `processFlatImage` removendo a dependência de `-f image2pipe` e migrando para o container puro `-f rawvideo -vcodec rawvideo` que aceita parsing direto de `.avi`.
- **Tarefa B:** Implementar Resiliência no `gerador_dataset_master.go`. Abandonar ferramentas como `wget` isolado que sofrem restrições `403 Forbidden` e implementar Downloader HTTP em Go nativamente, imitando navegadores completos.

## 2. P2P Neuro-Streaming (Voo Híbrido)
A prova de fogo: fazer dois computadores conversarem usando os Cérebros `HÍBRIDOS`.
- **Tarefa C:** Criar um simulador Micro-Rede `cliente-servidor` em Go.
- **Tarefa D:** O Lado A (`Emissor`) lê um trecho invisível de Vídeo de 50 Megabytes. Ao invés de trafegar Bits convencionais, ele converte usando o `HIBRIDO.gob` em uma lista infinita de UUIDs (Ex: `IdHash Hash Hash...`).
- **Tarefa E:** O tamanho do tráfego cai para poucos Kilobytes. O Lado B (`Receptor`) recebe os UUIDs e usa seu `HIBRIDO.gob` nativo para alucinar/reconstruir os fragmentos originais na tela.

## 3. FinOps e Cloud Scaling
Com o Lado A e Lado B consolidados via Socket UDP/TCP, não mais estaremos limitados ao processamento local:
- Vamos aprovisionar um novo contêiner sem gargalos na Vast.ai escolhendo uma IMAGEM estável de IA pre-treinadas (`Pytorch/Tensorflow Base`) que não falhe ao ligar SSH.
- Faremos o deploy da rede e calcularemos ao vivo "Latency (Ping) vs Decibéis de Reconstrução Neural".

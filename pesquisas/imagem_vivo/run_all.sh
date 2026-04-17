#!/bin/bash
set -e

echo "=========================================================="
echo "   ORQUESTRADOR DE LABORATÓRIO CROM_ESTÁTICO_V1           "
echo "   Dataset: ../dados_estaveis/lenna.png                   "
echo "=========================================================="

IMG_PATH="../../dados_estaveis/lenna.png"

echo ""
echo "[>>] Iniciando Experimento 1: Tesselador (Chunker 2D Lossless)"
cd exp1_chunker_2d
go run main.go "$IMG_PATH"
cd ..

echo ""
echo "[>>] Iniciando Experimento 2: Entropia de Vetor (MicroBrain)"
cd exp2_vq_entropy
go run main.go "$IMG_PATH"
cd ..

echo ""
echo "[>>] Iniciando Experimento 3: HNSW Codebook (Hash Skeleton)"
cd exp3_hnsw_anchors
go run main.go "$IMG_PATH"
cd ..

echo ""
echo "=========================================================="
echo "   BATERIA DE TESTES CONCLUÍDA COM SUCESSO.               "
echo "   Todos os artefatos visuais foram salvos nas pastas.    "
echo "=========================================================="

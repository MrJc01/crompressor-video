# 🌌 Crompressor-Video (Video-Vivo Engine)

[![Status](https://img.shields.io/badge/Pesquisa-Ativa-blue.svg)]()
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)]()

Bem-vindo ao repositório fundacional do **Crompressor-Video**, o motor de codificação e supercompressão de vídeo baseado na tese de **Quantização Vetorial Sub-simbólica e Árvores Topológicas (LSH/Codebooks)**.

## 📜 O Paradigma CROM
Contrariando os pipelines de estúdios e consórcios analógicos (como MPEG, AV1 e H.265), esta engine *não salva pixels* no sentido convencional de predição espacial/temporal estrita baseada em entropia de transformada de cosseno esparsa (DCT). Nós utilizamos um Dicionário Universal LSH em formato de árvore binária (`Cérebros`) que guarda padrões imutáveis de geometria fractal, faces, gramados, céus, etc. 

Arquivos `.cromvid` são puramente "Esqueletos" preenchidos com apontadores Hash (`uint64`). Este repositório realiza as provas forenses em código Go purista, validando a eficácia em ambiente de borda (Edge/SRE).

## 🗂 Estrutura do Ecossistema

1. **/docs/**:
   Contém 12 tomos de Engenharia de Software e Matemática que governam a arquitetura CROM (Zero Allocation, Estratificação de Codebooks, e Forward Compatibility). Documentação nível acadêmico de Pesquisa.

2. **/pesquisas/**:
   Campo de provas de algoritmos sujos. Nele habitam as premissas testadas e validadas de ponta a ponta sem estragar o design master da pasta `pkg/` de produção. Dividido centralmente em **Imagem Vivo** (Pipeline Espacial 2D) e **Vídeo Vivo** (Pipeline Temporal FastCDC).

## 🚀 Como Explorar a Documentação Técnica
Sugerimos ler sequencialmente a pasta `docs/` para obter a base teórica:
- `docs/01-visao-geral.md` até `docs/05-hash-residuais-e-delta.md` descrevem o pipeline do Zero Absoluto à cura do Uncanny Valley.
- `docs/06-motor-de-memoria-e-mmap.md` ao final cobrem a visão de hardware, go puro, I/O SRE e escala planetária multi-device.

---
*Projetado por Antigravity e MrJc01.*

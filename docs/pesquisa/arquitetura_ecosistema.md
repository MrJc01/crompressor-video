# Arquitetura Mãe do Ecossistema CROM (Vision e Master CLI)

Este documento guarda as diretrizes arquiteturais definidas durante os estudos sobre o ecossistema CROM Web e Nativo para garantir a padronização das implementações futuras.

## 1. O App Nativo Go (Ecossistema SRE de Terminal e GUI O(1))
A principal ferramenta para os pesquisadores e profissionais será um único binário robusto `crompressor-cli`. Ele não usa navegadores web. Ele obedece à filosofia Unix de velocidade absoluta, cruzando comandos de terminal com a Injeção Gráfica via Aceleração de Hardware (OpenGL/Vulkan puro) usando engines minúsculas como `Ebitengine` ou `Fyne`.

### O Fluxo Perfeito do Usuário (CLI + GUI Híbrido)
Você controla 100% via terminal, mas em vez de se limitar a texto, ele transborda para Janelas Nativas Ultra-Rápidas quando o aspecto visual é exigido.

- **Modo Setup/Settings:**
Ao rodar `$ crom train --interface` no terminal, o processo não abre o Chrome. Ele instantaneamente projeta uma Janela Nativa SRE Desktop (Dark/Cyberpunk) contendo o painel de treinamento. Você escolhe as pastas do Dataset Dourado, e acompanha a barra de progresso do Hashing usando toda a fúria dos Cores do seu processador sem gargalos IPC.

- **O Reprodutor Zero-Delay:**
Ao digitar `$ crom play video.crom`, uma janela minimalista preta aparece em `0.01` milissegundos. Sem barra de menus, sem peso operacional (como o VLC nativo). O motor local puxa seu `HIBRIDO.gob` na R.A.M. e empilha os Arrays de Pixel RGB crus puramente no barramento de Vídeo. A janela roda a midia perfeitamente a 144Hz.

*Esta arquitetura elimina a taxa da "SandBox" Web localmente, destila a performance da matriz quântica do Hash O(1) e cria uma ferramenta mágica de "Execução de Terminal que Pula pra Tela Gráfica".*

---

## 2. O Cinemark (WEB WASM END-USER)
Elementos Web projetados para internet de massa, que leem os metadados super leves gerados usando as mesmas lógicas consolidadas pelo App Nativo acima.

---

## O Modelo Web `<crom-video>` Perfeito
Para que a web adote nossos vídeos decodificáveis via Hash, a anatomia de uma Tag HTML no futuro será escrita pelo Desenvolvedor assim:

### Inicialização e Padrões Globais (A Sacada da Escalabilidade Web)
Ao invés do Desenvolvedor Front-end ter que repetir o link para a Mente Mestra (`brain-url`) em 25 tags de vídeo diferentes na mesma página, nós definiremos um **Escopo Global** na raiz da aplicação. O programador anexa o script do nosso Orquestrador apenas no `<head>` usando configs padrão:

```html
<!-- Setup Raiz com Data-Attributes na Tag Script -->
<script src="https://cdn.crom.site/player-v1.js" 
        data-default-brain="/mentes/hibrido_master.bin"
        data-gpu-accel="true"
        data-preload="true"></script>
```

Ou de forma programática para arquiteturas pesadas como NextJS e React (No seu app raiz):
```javascript
import { CromOrchestrator } from 'crom-vision';

CromOrchestrator.init({
    defaultBrainUrl: "https://minha-cdn.com/hibrido_master.bin",
    fallbackToH264: true, // Se falhar o WebGPU, renderiza MP4
    syncCache: "IndexedDB" // Persistência em HD Local
});
```
Dessa forma a tag `<crom-video>` da página passa a ficar 100% limpa, precisando apenas do nome do Hash `src="filme_01.crom"`.

### É possível passar puramente a string de Hash?
Sim. Para pequenos Lotties e Stickers sem necessidade de download, passa-se diretamente no HTML:

```html
<crom-video 
    data-hashes="aGg44Fdd91x..." 
    brain-url="/modelo_local.gob">
</crom-video>
```

---

## Formas Melhores de Fazer (The Ultimate Pipeline WebGPU)
Para garantir que na Web o usuário não sofra gargalos ou engasgos do Chromium/V8 Engine passando dados, a evolução correta englobará três diretrizes cruciais que implementaremos:

### 1. IndexedDB Persistent Caching
Nós **não** baixaremos o Cérebro toda hora. Usaremos a `Service Worker API` do navegador. O arquivo gigantesco `HIBRIDO.gob` (a mente) é baixado apenas no "Acesso 1" do usuário. Ele é guardado no fundo do HD do cache do Browser eternamente. Toda tag `<crom-video>` que usar a mesma mente carregará instantaneamente.

### 2. O Salto para WebGPU (A Trajetória Final)
Lutar contra o renderizador 2D do Canvas em HTML ainda exige que o Javascript mova arrays pela placa-mãe. 
**A solução final definitiva seria:** Pregar o arquivo `HIBRIDO.GOB` inteiro permanentemente na VRAM da Placa de Vídeo do Javascript via uma Ferramenta chamada **WebGPU**. 
Sendo assim, o pacote de dezenas de megabytes só se move da HD para a GPU no Load Page. O Video será apenas uma injeção de ID's Inteiros. Nossa "Placa de vídeo" vira uma dicionário que plota quadros infinitamente mais velozes do que enviar o vídeo já processado porque não há decodificação geométrica (Zero custo de CPU), só textura pré-existente mapeada num array.

## Próximos Passos Imediatos para Fornecimento
- Iniciar construção do **Crom-CLI Master** com flags (`crom train`, `crom encode`, `crom play`).
- Estruturar a renderização Ebitengine para validar a fluidez antes de portá-la para o módulo web CROM Custom Tag (WASM / WebGPU).

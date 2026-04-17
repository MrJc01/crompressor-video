class CromVideoElement extends HTMLElement {
    constructor() {
        super();
        this.attachShadow({ mode: 'open' });
        
        // Cria internamente um Canvas isolado em vez da Tag <video>
        this.canvas = document.createElement('canvas');
        this.canvas.style.width = "100%";
        this.canvas.style.height = "100%";
        this.canvas.width = 1280;
        this.canvas.height = 720;
        this.ctx = this.canvas.getContext('2d', { alpha: false, desynchronized: true });
        
        this.shadowRoot.appendChild(this.canvas);
        this.isPlaying = false;
    }

    connectedCallback() {
        const hashPayload = this.getAttribute('data-hashes') || "";
        const brainUrl = this.getAttribute('brain-url') || window.CROM_DEFAULT_BRAIN;
        
        if (!brainUrl) {
            console.error("<crom-video> requires a brain-url attribute or window.CROM_DEFAULT_BRAIN.");
            return;
        }

        this.initWasmAndPlay(hashPayload, brainUrl);
    }

    async initWasmAndPlay(hashData, brainUrl) {
        // No escopo real, aqui leríamos do IndexedDB
        console.log(`[CROM WASM] Carregando cérebro: ${brainUrl}`);
        
        // Instancia o WASM Compilado do Go
        const go = new Go();
        const wasmResult = await WebAssembly.instantiateStreaming(fetch("crom-player.wasm"), go.importObject);
        go.run(wasmResult.instance);
        
        // Dispara a Memory Inicializer O(1) do WASM
        window.cromInitializeBrain();

        // Extrai a lista de hashes do DataAttribute (Lista Separada por vírgula)
        const hashesToPlay = hashData.split(',').map(n => parseInt(n));
        this.play(hashesToPlay);
    }

    play(hashesList) {
        this.isPlaying = true;
        let frameIndex = 0;

        const renderLoop = () => {
            if (!this.isPlaying || frameIndex >= hashesList.length) return;
            
            const currentHash = hashesList[frameIndex];
            
            // MAGIC: Chamada bloqueante Ultra-Rápida do WebAssembly Módulo
            const rgbaBytes = window.cromDecodeHash(currentHash);
            
            if (rgbaBytes && rgbaBytes.length > 0) {
                // Monta ImageData brutal e preenche o canvas
                // Para 1 UUID retornando um Frame Inteiro (Ex: 1 pixel gigante no teste)
                // Aqui criamos um preenchimento sólido do pixel do Hash!
                this.ctx.fillStyle = `rgba(${rgbaBytes[0]}, ${rgbaBytes[1]}, ${rgbaBytes[2]}, ${rgbaBytes[3] / 255})`;
                this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
            }

            frameIndex++;
            requestAnimationFrame(renderLoop);
        };
        
        renderLoop();
    }
}

// Registrando a tag milagrosa no Browser!
customElements.define('crom-video', CromVideoElement);

console.log("[CROM-SDK] Web Component Tag <crom-video> injetada e pronta para IndexedDB/WASM.");

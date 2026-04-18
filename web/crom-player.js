class CromVideoElement extends HTMLElement {
    constructor() {
        super();
        this.attachShadow({ mode: 'open' });
        
        this.canvas = document.createElement('canvas');
        this.canvas.width = 640;
        this.canvas.height = 360;
        this.canvas.style.width = "100%";
        this.canvas.style.maxWidth = "800px";
        this.canvas.style.borderRadius = "8px";
        
        this.ctx = this.canvas.getContext('2d', { alpha: false, desynchronized: true });
        
        this.statusText = document.createElement('div');
        this.statusText.style.color = "#00f0ff";
        this.statusText.style.fontFamily = "monospace";
        this.statusText.style.padding = "10px";
        this.statusText.innerText = "[CROM SRE] Carregando...";

        this.shadowRoot.appendChild(this.statusText);
        this.shadowRoot.appendChild(this.canvas);
        
        this.isPlaying = false;
        this.chunkSize = 768; // Tensores Gray
        this.pixelsPerFrame = this.canvas.width * this.canvas.height;
        this.hashesPerFrame = Math.ceil(this.pixelsPerFrame / this.chunkSize);
    }

    connectedCallback() {
        const cromFile = this.getAttribute('src');
        const brainUrl = this.getAttribute('brain') || "hibrido.gob";
        
        if (!cromFile) {
            this.statusText.innerText = "[ERRO] Tag <crom-video> requer atributo 'src' apontando para um .crom";
            return;
        }

        this.initPipeline(cromFile, brainUrl);
    }

    async getBrainBytes(brainUrl) {
        // Implementação SRE do IndexedDB - Se o cerebro ja estiver na maquina, carrega em 0ms
        return new Promise((resolve, reject) => {
            const request = indexedDB.open("CromBrainDB", 1);
            
            request.onupgradeneeded = (e) => {
                const db = e.target.result;
                if (!db.objectStoreNames.contains("brains")) {
                    db.createObjectStore("brains");
                }
            };
            
            request.onsuccess = (e) => {
                const db = e.target.result;
                const tx = db.transaction("brains", "readonly");
                const store = tx.objectStore("brains");
                const getReq = store.get(brainUrl);
                
                getReq.onsuccess = async () => {
                    if (getReq.result) {
                        this.statusText.innerText = "[CROM SRE] Cérebro recuperado do Cold Storage (IndexedDB)!";
                        resolve(getReq.result);
                    } else {
                        this.statusText.innerText = `[CROM SRE] Baixando DNA original da rede: ${brainUrl}...`;
                        const res = await fetch(brainUrl);
                        const buffer = await res.arrayBuffer();
                        const uint8 = new Uint8Array(buffer);
                        
                        // Guarda no IDB para a eternidade
                        const txW = db.transaction("brains", "readwrite");
                        txW.objectStore("brains").put(uint8, brainUrl);
                        resolve(uint8);
                    }
                };
            };
            request.onerror = () => reject("IndexedDB falhou");
        });
    }

    async initPipeline(cromFile, brainUrl) {
        // 1. Injeta WebAssembly
        const go = new Go();
        const wasmResult = await WebAssembly.instantiateStreaming(fetch("crom-player.wasm"), go.importObject);
        go.run(wasmResult.instance);
        
        // 2. Extrai ou Baixa Mestre CROM (IndexedDB)
        const brainBytes = await this.getBrainBytes(brainUrl);
        
        // 3. Monta Cérebro no WASM Go
        this.statusText.innerText = "[CROM SRE] Decalibrando .gob massivo em VRAM WASM...";
        const success = window.cromLoadBrainBytes(brainBytes);
        if (!success) {
            this.statusText.innerText = "[ERRO] Cérebro colidiu ou falhou de carregar.";
            return;
        }

        // 4. Recebe a fita Magnética `.crom` p/ Stream
        this.statusText.innerText = `[CROM SRE] Sugando fita de Fótons O(1): ${cromFile}...`;
        const cromRes = await fetch(cromFile);
        const cromBuffer = await cromRes.arrayBuffer();
        
        // Cada UUID é float64/uint64 (8 bytes). O buffer é perfeitamente fragmentado em blocos O(1)
        const hashesList = [];
        const view = new DataView(cromBuffer);
        for (let i = 0; i < cromBuffer.byteLength; i += 8) {
            // LittleEndian (true) - Uint64 BigInt
            const uuidBig = view.getBigUint64(i, true);
            hashesList.push(uuidBig.toString());
        }

        this.statusText.innerText = `[CROM SRE] Sinal Encontrado: ${hashesList.length} UUIDs. Renderizando...`;
        
        // Fadeout do texto de status com delay
        setTimeout(() => this.statusText.style.display = 'none', 1000);

        this.playVirtualVram(hashesList);
    }

    playVirtualVram(hashesList) {
        this.isPlaying = true;
        let frameHashCursor = 0;
        
        // ImageData é pesada de ficar re-criando, criamos 1x e re-pintamos
        const mapBuffer = this.ctx.createImageData(this.canvas.width, this.canvas.height);
        
        const renderLoop = () => {
            if (!this.isPlaying) return;
            
            if (frameHashCursor >= hashesList.length) {
                frameHashCursor = 0; // Loop na fita
            }

            let pixelOffset = 0;
            const dataView = mapBuffer.data;

            // Extrai a cota de IDs para formar Exatos 1 Frame (Costura de Triangulos)
            for (let h = 0; h < this.hashesPerFrame; h++) {
                const cursor = frameHashCursor + h;
                if (cursor >= hashesList.length) break;

                const hashStr = hashesList[cursor];
                
                // MÁGICA: O WASM entrega um Chunk linear instantaneamente do Go
                const tensorUint8 = window.cromDecodeHash(hashStr);

                // Escopo Linear para Bidimensional (Array 1D p/ Uint8ClampedArray RGBA)
                if (tensorUint8 && tensorUint8.length > 0) {
                    for (let t = 0; t < tensorUint8.length; t++) {
                        const bufferPos = (pixelOffset + t) * 4;
                        if (bufferPos + 3 >= dataView.length) break;
                        
                        const val = tensorUint8[t];
                        dataView[bufferPos]     = val; // R
                        dataView[bufferPos + 1] = val; // G
                        dataView[bufferPos + 2] = val; // B
                        dataView[bufferPos + 3] = 255; // Alpha
                    }
                }
                pixelOffset += this.chunkSize;
            }

            // Expele o Quadro p/ a Placa Gráfica do Navegador
            this.ctx.putImageData(mapBuffer, 0, 0);

            frameHashCursor += this.hashesPerFrame;
            
            // Re-aciona o Refresh Rate do Monitor (Usually 60fps/144hz)
            requestAnimationFrame(renderLoop);
        };
        
        renderLoop();
    }
}

customElements.define('crom-video', CromVideoElement);

// Salat WASM Loader v1.6.1
// CDN-ready JavaScript loader for Salat WebAssembly module

(function(global) {
    'use strict';

    // Configuration
    const SALAT_CONFIG = {
        version: '1.6.1',
        wasmURL: null, // Will be auto-detected or set manually
        debug: false
    };

    // Detect CDN URL automatically
    function detectCDNURL() {
        const scripts = document.getElementsByTagName('script');
        for (let script of scripts) {
            if (script.src.includes('salat.js') || script.src.includes('salat-wasm')) {
                const baseURL = script.src.replace(/\/[^\/]*$/, '');
                return baseURL + '/salat.wasm';
            }
        }
        return './salat.wasm'; // fallback
    }

    // WASM loader class
    class SalatWASM {
        constructor() {
            this.isLoaded = false;
            this.loadPromise = null;
            this.go = null;
        }

        async load(wasmURL = null) {
            if (this.loadPromise) return this.loadPromise;

            this.loadPromise = this._loadWASM(wasmURL || SALAT_CONFIG.wasmURL || detectCDNURL());
            return this.loadPromise;
        }

        async _loadWASM(wasmURL) {
            try {
                // Load Go WASM support
                if (!global.Go) {
                    await this._loadGoWASMExec();
                }

                // Initialize Go instance
                this.go = new global.Go();
                
                // Load WASM module
                const result = await WebAssembly.instantiateStreaming(
                    fetch(wasmURL), 
                    this.go.importObject
                );

                // Run the Go program
                this.go.run(result.instance);
                this.isLoaded = true;

                if (SALAT_CONFIG.debug) {
                    console.log('🕌 Salat WASM loaded successfully');
                }

                return this;
            } catch (error) {
                console.error('Failed to load Salat WASM:', error);
                throw error;
            }
        }

        async _loadGoWASMExec() {
            return new Promise((resolve, reject) => {
                const script = document.createElement('script');
                script.src = 'https://cdn.jsdelivr.net/gh/golang/go@latest/misc/wasm/wasm_exec.js';
                script.onload = resolve;
                script.onerror = () => {
                    // Fallback to local wasm_exec.js
                    const fallbackScript = document.createElement('script');
                    fallbackScript.src = './wasm_exec.js';
                    fallbackScript.onload = resolve;
                    fallbackScript.onerror = reject;
                    document.head.appendChild(fallbackScript);
                };
                document.head.appendChild(script);
            });
        }

        // High-level API methods
        async prayerTimes(latitude, longitude, timezone = 'UTC') {
            await this.load();
            if (!global.salat) throw new Error('Salat WASM not loaded');
            return global.salat.prayerTime(latitude, longitude, timezone);
        }

        async version() {
            await this.load();
            if (!global.salat) throw new Error('Salat WASM not loaded');
            return global.salat.version();
        }

        async command(cmd) {
            await this.load();
            if (!global.salatConsole) throw new Error('Salat WASM not loaded');
            const result = global.salatConsole(cmd);
            return typeof result === 'string' ? JSON.parse(result) : result;
        }

        // Terminal integration
        createTerminal(terminalElement) {
            const terminal = {
                element: terminalElement,
                history: [],
                
                async execute(command) {
                    try {
                        const result = await this.command(command);
                        this.addOutput(JSON.stringify(result, null, 2));
                        return result;
                    } catch (error) {
                        this.addOutput(`Error: ${error.message}`);
                        return { error: error.message };
                    }
                }.bind(this),

                addOutput(text) {
                    if (this.element) {
                        const output = document.createElement('div');
                        output.textContent = text;
                        output.className = 'salat-terminal-output';
                        this.element.appendChild(output);
                        this.element.scrollTop = this.element.scrollHeight;
                    }
                    this.history.push(text);
                }
            };

            return terminal;
        }
    }

    // Create global instance
    const salatWASM = new SalatWASM();

    // Public API
    global.Salat = {
        // Core API
        async load(options = {}) {
            if (options.debug) SALAT_CONFIG.debug = true;
            if (options.wasmURL) SALAT_CONFIG.wasmURL = options.wasmURL;
            return salatWASM.load();
        },

        async prayerTimes(lat, lng, timezone) {
            return salatWASM.prayerTimes(lat, lng, timezone);
        },

        async version() {
            return salatWASM.version();
        },

        async command(cmd) {
            return salatWASM.command(cmd);
        },

        // Utilities
        createTerminal(element) {
            return salatWASM.createTerminal(element);
        },

        // Quick methods for common use cases
        async getPrayerTimesForLocation(lat, lng, timezone = 'UTC') {
            const result = await this.prayerTimes(lat, lng, timezone);
            return result.prayers;
        },

        async getCurrentPrayer(lat, lng, timezone = 'UTC') {
            const times = await this.getPrayerTimesForLocation(lat, lng, timezone);
            const now = new Date();
            const currentTime = now.toTimeString().substring(0, 5);
            
            // Simple logic to determine current prayer
            const prayers = Object.entries(times);
            for (let i = 0; i < prayers.length; i++) {
                const [name, time] = prayers[i];
                if (currentTime < time) {
                    return { current: i > 0 ? prayers[i-1][0] : 'isha', next: name, nextTime: time };
                }
            }
            return { current: 'isha', next: 'fajr', nextTime: times.fajr };
        },

        // Configuration
        config: SALAT_CONFIG
    };

    // Auto-load if not in module environment
    if (typeof module === 'undefined' && typeof window !== 'undefined') {
        // Auto-load on DOM ready
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => {
                global.Salat.load().catch(console.error);
            });
        } else {
            global.Salat.load().catch(console.error);
        }
    }

})(typeof window !== 'undefined' ? window : this); 
/**
 * WebSocket utility functions and helpers
 */

export type WebSocketMessage = {
    type: string;
    [key: string]: any;
};

export type WebSocketEventHandler = (data: any) => void;

/**
 * Creates a WebSocket connection with automatic reconnection
 */
export class WebSocketClient {
    private ws: WebSocket | null = null;
    private url: string;
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectDelay = 1000;
    private eventHandlers: Map<string, WebSocketEventHandler[]> = new Map();
    private isManualClose = false;

    constructor(url: string) {
        this.url = url;
    }

    /**
     * Connect to WebSocket server
     */
    connect(): Promise<void> {
        return new Promise((resolve, reject) => {
            try {
                this.ws = new WebSocket(this.url);
                this.isManualClose = false;

                this.ws.onopen = () => {
                    console.log('WebSocket connected');
                    this.reconnectAttempts = 0;
                    this.trigger('open', null);
                    resolve();
                };

                this.ws.onmessage = (event) => {
                    try {
                        const data = JSON.parse(event.data);
                        this.handleMessage(data);
                    } catch (error) {
                        console.error('Failed to parse WebSocket message:', error);
                    }
                };

                this.ws.onerror = (error) => {
                    console.error('WebSocket error:', error);
                    this.trigger('error', error);
                };

                this.ws.onclose = (event) => {
                    console.log('WebSocket closed:', event.code, event.reason);
                    this.trigger('close', event);

                    if (!this.isManualClose) {
                        this.attemptReconnect();
                    }
                };
            } catch (error) {
                reject(error);
            }
        });
    }

    /**
     * Send a message through the WebSocket
     */
    send(data: any): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(data));
        } else {
            console.warn('WebSocket is not connected');
        }
    }

    /**
     * Register an event handler
     */
    on(event: string, handler: WebSocketEventHandler): void {
        if (!this.eventHandlers.has(event)) {
            this.eventHandlers.set(event, []);
        }
        this.eventHandlers.get(event)?.push(handler);
    }

    /**
     * Unregister an event handler
     */
    off(event: string, handler: WebSocketEventHandler): void {
        const handlers = this.eventHandlers.get(event);
        if (handlers) {
            const index = handlers.indexOf(handler);
            if (index > -1) {
                handlers.splice(index, 1);
            }
        }
    }

    /**
     * Trigger event handlers
     */
    private trigger(event: string, data: any): void {
        const handlers = this.eventHandlers.get(event);
        if (handlers) {
            handlers.forEach(handler => handler(data));
        }
    }

    /**
     * Handle incoming messages
     */
    private handleMessage(data: WebSocketMessage): void {
        // Trigger type-specific handler
        if (data.type) {
            this.trigger(data.type, data);
        }

        // Trigger general message handler
        this.trigger('message', data);
    }

    /**
     * Attempt to reconnect to WebSocket
     */
    private attemptReconnect(): void {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.error('Max reconnection attempts reached');
            this.trigger('reconnect_failed', null);
            return;
        }

        this.reconnectAttempts++;
        const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);

        console.log(`Attempting to reconnect in ${delay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);

        setTimeout(() => {
            this.connect().catch((error) => {
                console.error('Reconnection failed:', error);
            });
        }, delay);
    }

    /**
     * Close the WebSocket connection
     */
    close(): void {
        this.isManualClose = true;
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }

    /**
     * Check if WebSocket is connected
     */
    isConnected(): boolean {
        return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
    }

    /**
     * Get connection state
     */
    getState(): number {
        return this.ws?.readyState ?? WebSocket.CLOSED;
    }
}

/**
 * Create a simple WebSocket connection
 */
export function createWebSocket(url: string): WebSocketClient {
    return new WebSocketClient(url);
}

/**
 * Convert relative WebSocket URL to absolute
 */
export function getWebSocketUrl(path: string): string {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    return `${protocol}//${host}${path}`;
}

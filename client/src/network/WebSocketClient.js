import { WS_URL } from '../utils/constants.js';

export class WebSocketClient {
  constructor(onMessage) {
    this.ws = null;
    this.onMessage = onMessage;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 1000;
    this.isIntentionallyClosed = false;
  }

  connect() {
    this.isIntentionallyClosed = false;
    this.ws = new WebSocket(WS_URL);

    this.ws.onopen = () => {
      console.log('Connected to server');
      this.reconnectAttempts = 0;
      if (this.onConnectionChange) {
        this.onConnectionChange(true);
      }
    };

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        if (this.onMessage) {
          this.onMessage(message);
        }
      } catch (error) {
        console.error('Error parsing message:', error);
      }
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    this.ws.onclose = () => {
      console.log('Disconnected from server');
      if (this.onConnectionChange) {
        this.onConnectionChange(false);
      }

      // Attempt to reconnect if not intentionally closed
      if (!this.isIntentionallyClosed && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++;
        console.log(`Reconnecting... (attempt ${this.reconnectAttempts})`);
        setTimeout(() => this.connect(), this.reconnectDelay);
      }
    };
  }

  send(type, payload = {}) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type, payload }));
    } else {
      console.warn('WebSocket not connected, cannot send message');
    }
  }

  close() {
    this.isIntentionallyClosed = true;
    if (this.ws) {
      this.ws.close();
    }
  }

  // Reconnect to pick up new auth state (e.g., after login)
  reconnect() {
    this.isIntentionallyClosed = true;
    if (this.ws) {
      this.ws.close();
    }
    // Small delay to ensure close completes
    setTimeout(() => {
      this.reconnectAttempts = 0;
      this.connect();
    }, 100);
  }

  isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN;
  }

  setConnectionChangeHandler(handler) {
    this.onConnectionChange = handler;
  }
}

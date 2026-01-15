import { WS_URL } from '../utils/constants.js';

export class WebSocketClient {
  constructor(onMessage) {
    this.ws = null;
    this.onMessage = onMessage;
    this.reconnectAttempts = 0;
    this.baseReconnectDelay = 1000;
    this.maxReconnectDelay = 30000; // Cap at 30 seconds
    this.isIntentionallyClosed = false;
    this.reconnectTimeout = null;
  }

  connect() {
    this.isIntentionallyClosed = false;

    // Clear any pending reconnect
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }

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
      // Use exponential backoff with a cap
      if (!this.isIntentionallyClosed) {
        this.reconnectAttempts++;
        const delay = Math.min(
          this.baseReconnectDelay * Math.pow(2, this.reconnectAttempts - 1),
          this.maxReconnectDelay
        );
        console.log(`Reconnecting in ${delay}ms... (attempt ${this.reconnectAttempts})`);
        this.reconnectTimeout = setTimeout(() => this.connect(), delay);
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
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
    if (this.ws) {
      this.ws.close();
    }
  }

  // Reconnect to pick up new auth state (e.g., after login)
  reconnect() {
    this.isIntentionallyClosed = true;
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
    if (this.ws) {
      this.ws.close();
    }
    // Small delay to ensure close completes
    setTimeout(() => {
      this.reconnectAttempts = 0;
      this.connect();
    }, 100);
  }

  // Force an immediate reconnection attempt (used when user takes action)
  forceReconnect() {
    if (this.isConnected()) {
      return; // Already connected
    }
    console.log('Force reconnecting...');
    this.reconnectAttempts = 0;
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
    this.connect();
  }

  isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN;
  }

  setConnectionChangeHandler(handler) {
    this.onConnectionChange = handler;
  }
}

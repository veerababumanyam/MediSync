/**
 * MediSync WebSocket Client
 *
 * WebSocket client for real-time streaming responses with:
 * - Connection management
 * - Automatic reconnection
 * - Message parsing and dispatch
 *
 * @module services/websocket
 */

/**
 * WebSocket message types
 */
export type WSMessageType =
  | 'thinking'
  | 'sql_preview'
  | 'result'
  | 'error'
  | 'clarification'
  | 'done'

/**
 * WebSocket message structure
 */
export interface WSMessage {
  type: WSMessageType
  content?: string
  chartType?: string
  data?: unknown
  confidence?: number
  sql?: string
  suggestions?: string[]
}

/**
 * WebSocket event handlers
 */
export interface WSEventHandlers {
  onMessage?: (message: WSMessage) => void
  onError?: (error: Event) => void
  onClose?: (event: CloseEvent) => void
  onOpen?: () => void
  onReconnecting?: (attempt: number) => void
}

/**
 * WebSocket client configuration
 */
export interface WSClientConfig {
  url?: string
  reconnectAttempts?: number
  reconnectDelay?: number
  heartbeatInterval?: number
}

const DEFAULT_CONFIG: Required<WSClientConfig> = {
  url: `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/v1/chat/stream`,
  reconnectAttempts: 5,
  reconnectDelay: 1000,
  heartbeatInterval: 30000,
}

/**
 * WebSocket Client class
 */
export class WebSocketClient {
  private ws: WebSocket | null = null
  private config: Required<WSClientConfig>
  private handlers: WSEventHandlers
  private reconnectAttempt = 0
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null
  private isIntentionallyClosed = false

  constructor(config: WSClientConfig = {}, handlers: WSEventHandlers = {}) {
    this.config = { ...DEFAULT_CONFIG, ...config }
    this.handlers = handlers
  }

  /**
   * Connect to WebSocket server
   */
  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        resolve()
        return
      }

      this.isIntentionallyClosed = false

      try {
        this.ws = new WebSocket(this.config.url)

        this.ws.onopen = () => {
          this.reconnectAttempt = 0
          this.startHeartbeat()
          this.handlers.onOpen?.()
          resolve()
        }

        this.ws.onmessage = (event) => {
          this.handleMessage(event.data)
        }

        this.ws.onerror = (error) => {
          this.handlers.onError?.(error)
          reject(error)
        }

        this.ws.onclose = (event) => {
          this.stopHeartbeat()
          this.handlers.onClose?.(event)

          if (!this.isIntentionallyClosed && this.reconnectAttempt < this.config.reconnectAttempts) {
            this.scheduleReconnect()
          }
        }
      } catch (error) {
        reject(error)
      }
    })
  }

  /**
   * Send a message to the server
   */
  send(type: string, payload: unknown): void {
    if (this.ws?.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket is not connected')
    }

    this.ws.send(JSON.stringify({ type, payload }))
  }

  /**
   * Send a chat query
   */
  sendQuery(query: string, sessionId?: string, locale?: string): void {
    this.send('query', { query, sessionId, locale })
  }

  /**
   * Close the connection
   */
  close(): void {
    this.isIntentionallyClosed = true
    this.stopHeartbeat()

    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  /**
   * Check if connected
   */
  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }

  /**
   * Handle incoming message
   */
  private handleMessage(data: string): void {
    // Handle done signal
    if (data === '[DONE]') {
      this.handlers.onMessage?.({ type: 'done' })
      return
    }

    try {
      // Parse SSE format: "data: <json>"
      if (data.startsWith('data: ')) {
        const jsonStr = data.slice(6)
        const message = JSON.parse(jsonStr) as WSMessage
        this.handlers.onMessage?.(message)
      } else {
        // Direct JSON
        const message = JSON.parse(data) as WSMessage
        this.handlers.onMessage?.(message)
      }
    } catch (error) {
      console.error('Failed to parse WebSocket message:', error)
    }
  }

  /**
   * Start heartbeat to keep connection alive
   */
  private startHeartbeat(): void {
    this.stopHeartbeat()
    this.heartbeatTimer = setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify({ type: 'ping' }))
      }
    }, this.config.heartbeatInterval)
  }

  /**
   * Stop heartbeat timer
   */
  private stopHeartbeat(): void {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }
  }

  /**
   * Schedule a reconnection attempt
   */
  private scheduleReconnect(): void {
    this.reconnectAttempt++
    this.handlers.onReconnecting?.(this.reconnectAttempt)

    setTimeout(() => {
      this.connect().catch((error) => {
        console.error('Reconnection failed:', error)
      })
    }, this.config.reconnectDelay * this.reconnectAttempt)
  }
}

/**
 * Create a WebSocket client instance
 */
export function createWebSocketClient(
  config?: WSClientConfig,
  handlers?: WSEventHandlers
): WebSocketClient {
  return new WebSocketClient(config, handlers)
}

export default WebSocketClient

import { sendSocketPayload } from './send'
import { useChatModeStore } from '../../chat/store/chatModeStore'

let socket = null

export function connect(token, onMessage) {
  if (!token) {
    console.error('[WS] Connection aborted: no token provided')
    return
  }

  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.close(1000, 'Reconnecting')
  }

  socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`)

  socket.onopen = () => {
    console.log('[WS] Connected')

    const chatMode = useChatModeStore()

    const payload = {
      kind: 'event',
      action: 'init_chat',
      chat_type: chatMode.chatType, // 'global' | 'private'
      target_id: chatMode.receiverID, // can be null
      timestamp: Date.now(),
    }

    sendSocketPayload(payload)
  }

  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      onMessage?.(data)
    } catch (err) {
      console.error('[WS] Parsing error:', err.message);
      console.error('[WS] Invalid message:', event.data)
    }
  }

  socket.onclose = (event) => {
    console.warn(`[WS] Disconnected (${event.code}): ${event.reason || 'no reason'}`)
  }

  socket.onerror = (err) => {
    console.error('[WS] Error:', err)
  }
}

export function disconnect() {
  if (socket) {
    socket.close(1000, 'Manual disconnect')
    socket = null
  }
}

export function getSocket() {
  return socket
}
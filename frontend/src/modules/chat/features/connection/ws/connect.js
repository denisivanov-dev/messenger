let socket = null

export function connect(token, onMessage, mode = 'global', receiverId = null) {
  if (!token) {
    console.error('WebSocket connection aborted: no token provided')
    return
  }

  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.close(1000, 'Reconnecting')
  }

  socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`)

  socket.onopen = () => {
    console.log('WebSocket connected')

    if (mode === 'global') {
      socket.send(JSON.stringify({
        type: 'init_global',
        chat_type: 'global'
      }))
    } else if (mode === 'private' && receiverId) {
      socket.send(JSON.stringify({
        type: 'init_private',
        chat_type: 'private',
        receiver_id: receiverId
      }))
    }
  }

  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      onMessage?.(data)
    } catch (e) {
      console.error('[WS] Invalid message', event.data)
    }
  }

  socket.onclose = () => console.warn('[WS] Disconnected')
  socket.onerror = (err) => console.error('[WS] Error', err)
}

export function disconnect() {
  if (socket) socket.close()
}

export function getSocket() {
  return socket
}
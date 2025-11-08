import { getSocket } from "./connect"

export function sendSocketPayload(payload) {
  const socket = getSocket()
  if (socket?.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(payload))
  } else {
    console.warn('[WS] Cannot send, socket not open')
  }
}

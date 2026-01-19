export function handleMessageEvents(msg, stores) {
  const { messagesStore } = stores
  const payload = msg.payload || {}

  switch (msg.action) {
    case 'send':
      return messagesStore.pushFromWs(msg)

    case 'edit':
      return messagesStore.handleEdited(
        payload.message_id,
        payload.text,
        payload.edited_at
      )

    case 'delete':
      return messagesStore.handleDeleted(payload.message_id)

    case 'pin':
      return messagesStore.handlePinned(payload.message_id, payload.should_pin)

    default:
      console.warn('[WS] Unhandled message action:', msg.action)
  }
}
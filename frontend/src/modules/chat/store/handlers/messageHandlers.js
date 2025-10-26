export function handleMessageDeleted(msg, stores) {
  stores.messagesStore.handleDeleted(msg.message_id)
}

export function handleMessageEdited(msg, stores) {
  stores.messagesStore.handleEdited(msg.message_id, msg.new_text, msg.edited_at)
}

export function handleMessagePinned(msg, stores) {
  stores.messagesStore.handlePinned(msg.message_id, msg.action)
}

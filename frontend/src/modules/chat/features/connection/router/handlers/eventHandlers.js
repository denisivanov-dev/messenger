export function handleTyping(msg, stores) {
  const { typingStore, chatType, receiverID, myId } = stores
  typingStore.handleTypingWs(msg, chatType, myId, receiverID)
}

export function handleUserStatus(msg, stores) {
  const { userStore } = stores
  userStore.applyStatus(msg.user_id, msg.status)
}

export function handleTyping(msg, stores) {
  const { eventsStore, chatType, receiverID, myId } = stores
  eventsStore.handleTypingWs(msg, chatType.value, myId, receiverID.value)
}

export function handleUserStatus(msg, stores) {
  const { userStore } = stores
  userStore.applyStatus(msg.user_id, msg.status)
}

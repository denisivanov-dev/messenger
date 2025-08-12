export function getRoomId(chatType, myId, receiverId) {
  if (chatType === 'global') return '1'
  if (receiverId === undefined || receiverId === null || receiverId === '') return ''

  const a = String(myId)
  const b = String(receiverId)
  return a < b ? `${a}:${b}` : `${b}:${a}`
}

export function isTypingForCurrentRoom(msg, chatType, myId, receiverId) {
  if (msg.type !== 'typing') return false
  if (String(msg.user_id) === String(myId)) return false

  const current = getRoomId(chatType, myId, receiverId)
  if (!msg.chat_id) return false
  
  return msg.chat_id === current
}

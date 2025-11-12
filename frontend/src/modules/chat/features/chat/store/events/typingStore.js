import { defineStore } from 'pinia'
import { ref } from 'vue'
import { sendSocketPayload } from '../../../connection/ws/send'
import { isTypingForCurrentRoom } from '../../utils/chatRoom'

export const useTypingStore = defineStore('events', () => {
  const typingUser = ref(null)

  let typingTimer = null
  let typingCooldown = false

  function setTyping(username) {
    typingUser.value = username
    if (typingTimer) clearTimeout(typingTimer)
    typingTimer = setTimeout(() => {
      typingUser.value = null
    }, 3000)
  }

  function clearTyping() {
    if (typingTimer) clearTimeout(typingTimer)
    typingTimer = null
    typingUser.value = null
}

  function handleTypingWs(msg, chatType, myId, receiverId) {
    if (!isTypingForCurrentRoom(msg, chatType, myId, receiverId)) return
    setTyping(msg.username)
  }

  function sendTyping(chatType, receiverId) {
    if (typingCooldown) return
    sendSocketPayload({
      type: 'typing',
      receiver_id: receiverId,
      text: '',
      chat_type: chatType,
      timestamp: Date.now()
    })
    typingCooldown = true
    setTimeout(() => (typingCooldown = false), 2000)
  }

  return {
    typingUser,
    setTyping,
    clearTyping,
    handleTypingWs,
    sendTyping,
  }
})

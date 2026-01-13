import { defineStore } from 'pinia'
import { ref } from 'vue'
import { sendSocketPayload } from '../../../connection/ws/send'
import { isTypingForCurrentRoom } from '../../utils/chatRoom'
import { useChatModeStore } from '../chatModeStore'
import { useMessagesStore } from '../messagesStore'

export const useTypingStore = defineStore('typing', () => {
  const chatModeStore = useChatModeStore()
  const messagesStore = useMessagesStore()
  
  const typingUsers = ref([]) // { id, username, timer }
  let typingCooldown = false

  function setTyping(userId, username) {
    let index = typingUsers.value.findIndex(u => u.id === userId)

    if (index === -1) {
      typingUsers.value.push({ id: userId, username, timer: null })
      index = typingUsers.value.length - 1
    }

    // reset TTL timer
    clearTimeout(typingUsers.value[index].timer )

    typingUsers.value[index].timer = setTimeout(() => {
      removeTyping(userId)
    }, 3000)

    messagesStore.shouldScroll = true
  }

  function removeTyping(userId) {
    const index = typingUsers.value.findIndex(u => u.id === userId)
    if (index === -1) return

    clearTimeout(typingUsers.value[index].timer)
    typingUsers.value.splice(index, 1)
  }

  // clear typing state (e.g. on chat switch / disconnect)
  function clearTyping() {
    typingUsers.value.forEach(u => clearTimeout(u.timer))
    typingUsers.value = []
  }

  function handleTypingWs(msg, chatType, myId, receiverId) {
    if (!isTypingForCurrentRoom(msg, chatType, myId, receiverId)) return
    setTyping(msg.sender_id, msg.payload.username)
  }

  function sendTyping() {
    const { chatType, receiverID} = chatModeStore

    if (!chatType) return
    if (chatType === 'private' && !receiverID) return
    if (typingCooldown) return

    sendSocketPayload({
      kind: 'event',
      action: 'typing',
      chat_type: chatModeStore.chatType,
      target_id: chatModeStore.receiverID,
      timestamp: Date.now(),
      payload: {},
    })

    typingCooldown = true
    setTimeout(() => {
      typingCooldown = false
    }, 2000)
  }

  return {
    typingUsers,
    setTyping,
    removeTyping,
    clearTyping,
    handleTypingWs,
    sendTyping,
  }
})
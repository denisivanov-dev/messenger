import { defineStore } from 'pinia'
import { ref } from 'vue'
import { sendMessage 
import { useStorage } from '@vueuse/core'
import { useTypingStore } from './events/typingStore'
import { useMessagesStore } from './messagesStore'

export const useChatModeStore = defineStore('chatMode', () => {
  const chatType = ref('global')
  const receiverID = ref(null)
  const imageUrlCache = useStorage('image-url-cache', {})

  const typingStore = useTypingStore()
  const messagesStore = useMessagesStore()

  // --- CHAT MODE MANAGEMENT ---
  function setChatModeGlobal() {
    chatType.value = 'global'
    receiverID.value = null
    typingStore.clearTyping()
    messagesStore.clear()
    localStorage.setItem('chatMode', 'global')

    try {
      sendMessage({
        type: 'init_global',
        chat_type: 'global'
      })
    } catch (err) {
      console.warn('[ChatMode] WS not ready yet:', err.message)
    }
  }

  function setChatModePrivate(targetID) {
    chatType.value = 'private'
    receiverID.value = targetID
    typingStore.clearTyping()
    messagesStore.clear()
    localStorage.setItem('chatMode', 'private')
    localStorage.setItem('receiverId', targetID)

    try {
      sendMessage({
        type: 'init_private',
        chat_type: 'private',
        receiver_id: targetID
      })
    } catch (err) {
      console.warn('[ChatMode] WS not ready yet:', err.message)
    }
  }

  return {
    chatType,
    receiverID,
    imageUrlCache,
    setChatModeGlobal,
    setChatModePrivate,
  }
})
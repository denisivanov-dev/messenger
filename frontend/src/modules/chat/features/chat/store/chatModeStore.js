import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useStorage } from '@vueuse/core'
import { useTypingStore } from './events/typingStore'
import { useMessagesStore } from './messagesStore'

export const useChatModeStore = defineStore('chatMode', () => {
  const chatType = ref(localStorage.getItem('chatMode') || 'global')
  const receiverID = ref(localStorage.getItem('receiverId') || null)
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
    localStorage.removeItem('receiverId')
  }

  function setChatModePrivate(targetID) {
    chatType.value = 'private'
    receiverID.value = targetID
    typingStore.clearTyping()
    messagesStore.clear()
    localStorage.setItem('chatMode', 'private')
    localStorage.setItem('receiverId', targetID)
  }

  return {
    chatType,
    receiverID,
    imageUrlCache,
    setChatModeGlobal,
    setChatModePrivate,
  }
})
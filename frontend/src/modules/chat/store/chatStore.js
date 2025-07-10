import { defineStore } from 'pinia'
import { ref } from 'vue'
import { connect, disconnect, sendMessage, getAllUsers, startPrivateChat } from '../api/chatApi'
import { isNearBottom } from '../utils/chatWindowUtils'
import { useAuthStore } from '../../auth/store/authStore'

export const useChatStore = defineStore('chat', () => {
  const messages = ref([])
  const users = ref([]) 
  const connected = ref(false)
  const shouldScroll = ref(false)
  const chatType = ref("global")
  const receiverID = ref(null)
  const typingUser = ref(null)
  let typingTimer = null
  let typingCooldown = false

  function setChatModeGlobal() {
    chatType.value = "global"
    receiverID.value = null
  }

  function setChatModePrivate(targetID) {
    chatType.value = "private"
    receiverID.value = targetID 
  }

  function startChat(token) {
    messages.value = []
    
    const authStore = useAuthStore() 
    const myId = authStore.getUserId

    connect(token, (msg) => {
      if (msg.event === 'user_status') {
        const user = users.value.find(u => u.id === msg.user_id)
        if (user) user.status = msg.status
        return
      }

      if (msg.type === 'typing') {
        if (String(msg.user_id) === String(myId)) return

        setTyping(msg.username)
        console.info(msg.username)
        return
      }

      typingUser.value = null
      messages.value.push(msg)
      shouldScroll.value = isNearBottom()
    })

    shouldScroll.value = true
    connected.value = true
  }

  function send(text) {
    const msg = {
      text,
      timestamp: Date.now(),
      chat_type: chatType.value,
      receiver_id: receiverID.value
    }

    sendMessage(msg)
    shouldScroll.value = true
  }

  function setTyping(username) {
    typingUser.value = username
    if (typingTimer) clearTimeout(typingTimer)
    typingTimer = setTimeout(() => {
      typingUser.value = null
    }, 3000)
  }

  function sendTyping () {
    if (typingCooldown) return

    sendMessage({
      type:       'typing',
      receiver_id: receiverID.value,
      text:       '',
      chat_type:  chatType.value,
      timestamp:  Date.now()
    })

    typingCooldown = true
    setTimeout(() => (typingCooldown = false), 2000)
  }

  function stopChat() {
    disconnect()
    connected.value = false
  }

  async function fetchUsers() {
    try {
      const data = await getAllUsers()
      users.value = data
      console.info(JSON.stringify(users.value, null, 2))
      } catch (err) {
        console.error("Не удалось загрузить пользователей:", err)
    }
  }

  async function openOrCreatePrivateChat(targetId) {
    try {
      messages.value = []
      const response = await startPrivateChat(targetId)
      setChatModePrivate(targetId)
      receiverID.value = targetId
      console.info(response)
      messages.value = response.messages
      shouldScroll.value = true
      return response.success
    } catch (err) {
        console.error('openOrCreatePrivateChat error:', err)
        throw err
    }
  }

  return {
    messages,
    users,
    connected,
    shouldScroll,
    startChat,
    send,
    stopChat,
    fetchUsers,
    openOrCreatePrivateChat,
    sendTyping,
    typingUser,
    setTyping,
  }
})
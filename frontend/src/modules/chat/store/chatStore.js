import { defineStore } from 'pinia'
import { ref } from 'vue'
import { connect, disconnect, sendMessage, getAllUsers, startPrivateChat } from '../api/chatApi'
import { isNearBottom } from '../utils/chatWindowUtils'
import { useAuthStore } from '../../auth/store/authStore'

export const useChatStore = defineStore('chat', () => {
  // --- STATE ---
  const messages = ref([])
  const users = ref([])
  const connected = ref(false)
  const shouldScroll = ref(false)
  const chatType = ref('global')
  const receiverID = ref(null)
  const typingUser = ref(null)

  let typingTimer = null
  let typingCooldown = false

  // --- MODE MANAGEMENT ---
  function setChatModeGlobal() {
    chatType.value = 'global'
    receiverID.value = null
    messages.value = []

    localStorage.setItem('chatMode', 'global')
     if (connected.value) {
      sendMessage({
        type: 'init_global',
        chat_type: 'global'
      })
    }
  }

  function setChatModePrivate(targetID) {
    chatType.value = 'private'
    receiverID.value = targetID
    messages.value = []

    localStorage.setItem('chatMode', 'private')
    localStorage.setItem('receiverId', targetID)
    if (connected.value) {
      sendMessage({
        type: 'init_private',
        chat_type: 'private',
        receiver_id: targetID
      })
    }
  }

  // --- CONNECTION ---
  function startChat(token, mode = 'global', receiverId = null) {
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
        return
      }

      if (msg.type === 'message_deleted') {
        messages.value = messages.value.filter(m => m.message_id !== msg.message_id)
        return
      }

      if (msg.type === 'message_edited') {
        console.info('message_edited payload:', JSON.stringify(msg))
        const edited = messages.value.find(m => m.message_id === msg.message_id)
        if (edited) {
          edited.text = msg.new_text
          edited.timestamp = msg.edited_at
          edited.edited_at = msg.edited_at
        }
        return
      }

      if (msg.type === 'message_pinned') {
        console.info('message_pinned payload:', JSON.stringify(msg))
        const pinned = messages.value.find(m => m.message_id === msg.message_id)
        if (pinned && msg.action === 'pin') {
          pinned.pinned = true
        } else {
          pinned.pinned = false
        }
        return
      }

      typingUser.value = null
      messages.value.push(msg)
      shouldScroll.value = isNearBottom()
    }, mode, receiverId)

    shouldScroll.value = true
    connected.value = true
  }

  function stopChat() {
    disconnect()
    connected.value = false
  }

  // --- MESSAGING ---
  function send(text, replyToMessage = null) {
    const msg = {
      type: 'send_message',
      text: text.trim(),
      timestamp: Date.now(),
      chat_type: chatType.value,
      receiver_id: receiverID.value,
    }

    if (replyToMessage) {
      msg.reply_to = replyToMessage.message_id
      msg.reply_to_text = replyToMessage.text
      msg.reply_to_user = replyToMessage.username
    }

    console.info(msg.reply_to, msg.reply_to_text, msg.reply_to_user)
    sendMessage(msg)
    shouldScroll.value = true
  }

  function sendTyping() {
    if (typingCooldown) return

    sendMessage({
      type: 'typing',
      receiver_id: receiverID.value,
      text: '',
      chat_type: chatType.value,
      timestamp: Date.now()
    })

    typingCooldown = true
    setTimeout(() => (typingCooldown = false), 2000)
  }

  function setTyping(username) {
    typingUser.value = username
    if (typingTimer) clearTimeout(typingTimer)
    typingTimer = setTimeout(() => {
      typingUser.value = null
    }, 3000)
  }

  // --- USERS ---
  async function fetchUsers() {
    try {
      const data = await getAllUsers()
      users.value = data
      console.info(JSON.stringify(users.value, null, 2))
    } catch (err) {
      console.error('Не удалось загрузить пользователей:', err)
    }
  }

  // --- CHAT CONTROL ---
  async function openOrCreatePrivateChat(targetId) {
    try {
      messages.value = []
      const response = await startPrivateChat(targetId)
      receiverID.value = targetId
      localStorage.setItem('receiverId', targetId)
      console.info(response)
      console.info(receiverID.value)
      messages.value = response.messages
      shouldScroll.value = true
      return response.success
    } catch (err) {
      console.error('openOrCreatePrivateChat error:', err)
      throw err
    }
  }

  function deleteMessage(message) {
    sendMessage({
      type: 'delete_message',
      message_id: message.message_id,
      receiver_id: receiverID.value,
      chat_type: chatType.value
    })
  }

  function editMessage(message, newText) {
    sendMessage({
      type: 'edit_message',
      message_id: message.message_id,
      new_text: newText,
      receiver_id: receiverID.value,
      chat_type: chatType.value
    })
  }

  function pinMessage(message, shouldPin = true) {
    console.info(chatType.value)
    
    sendMessage({
      type: 'pin_message',
      message_id: message.message_id,
      chat_id: message.chat_id,
      action: shouldPin ? 'pin' : 'unpin',
      receiver_id: receiverID.value,
      chat_type: chatType.value
    })
  }

  return {
    // state
    messages,
    users,
    connected,
    shouldScroll,
    typingUser,

    // mode
    setChatModeGlobal,
    setChatModePrivate,

    // connection
    startChat,
    stopChat,

    // actions
    send,
    sendTyping,
    setTyping,
    fetchUsers,
    openOrCreatePrivateChat,
    deleteMessage,
    editMessage,
    pinMessage
  }
})
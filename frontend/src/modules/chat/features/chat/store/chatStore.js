import { defineStore, storeToRefs } from 'pinia'
import { ref } from 'vue'
import { connect, disconnect, sendMessage } from '../api/chatApi'
import { isNearBottom } from '../utils/chatWindowUtils'
import { useAuthStore } from '../../../../auth/store/authStore'
import { useStorage } from '@vueuse/core'

import { useMessagesStore } from './messagesStore'
import { useUserStore } from './userStore'
import { useFriendsStore } from '../../friends/store/friendsStore'
import { useTypingStore } from './events/typingStore'
import { useCallStore } from '../../call/store/callStore'
import { useWebRTCStore } from '../../call/store/webrtcStore'

import { chatHandlersRegistry } from './handlers/chatHandlersRegistry'

export const useChatStore = defineStore('chat', () => {
  const messagesStore = useMessagesStore()
  const userStore = useUserStore()
  const friendsStore = useFriendsStore()
  const typingStore = useTypingStore()
  const callStore = useCallStore()
  const webrtcStore = useWebRTCStore()

  const { messages, shouldScroll } = storeToRefs(messagesStore)
  const { users } = storeToRefs(userStore)
  const { friendStatusCache } = storeToRefs(friendsStore)
  const { typingUser } = storeToRefs(typingStore)

  const connected = ref(false)
  const chatType = ref('global')
  const receiverID = ref(null)
  const imageUrlCache = useStorage('image-url-cache', {})

  // --- MODE MANAGEMENT ---
  function setChatModeGlobal() {
    chatType.value = 'global'
    receiverID.value = null
    typingStore.clearTyping()
    messagesStore.clear()
    localStorage.setItem('chatMode', 'global')
    if (connected.value) {
      sendMessage({ type: 'init_global', chat_type: 'global' })
    }
  }

  function setChatModePrivate(targetID) {
    chatType.value = 'private'
    receiverID.value = targetID
    typingStore.clearTyping()
    messagesStore.clear()
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
    messagesStore.clear()
    const authStore = useAuthStore()
    const myId = authStore.getUserId

    connect(token, (msg) => {
      console.info('[ws message]', JSON.stringify(msg, null, 2))

      const handler = chatHandlersRegistry[msg.type] || chatHandlersRegistry[msg.event]
      if (handler) {
        handler(msg, {
          userStore,
          messagesStore,
          friendsStore,
          typingStore,
          callStore,
          webrtcStore,
          chatType,
          receiverID,
          myId
        })
        return
      }

      // fallback если хендлер не найден
      const shouldAutoScroll = isNearBottom()
      messagesStore.pushFromWs(msg)
    }, mode, receiverId)

    shouldScroll.value = true
    connected.value = true
  }

  function stopChat() {
    disconnect()
    connected.value = false
  }

  // --- MESSAGING ---
  function sendMessageData(payload) {
    messagesStore.sendMessageData(payload, chatType.value, receiverID.value)
  }

  function sendTyping() {
    typingStore.sendTyping(chatType.value, receiverID.value)
  }

  function setTyping(username) {
    typingStore.setTyping(username)
  }

  // --- USERS ---
  async function fetchUsers() {
    try {
      await userStore.fetchUsers()
    } catch (err) {
      console.error('Не удалось загрузить пользователей:', err)
    }
  }

  // --- CHAT CONTROL ---
  async function openOrCreatePrivateChat(targetId) {
    return messagesStore.openOrCreatePrivateChat(targetId, receiverID)
  }

  function deleteMessage(message) {
    messagesStore.deleteMessage(message, chatType.value, receiverID.value)
  }

  function editMessage(message, newText) {
    messagesStore.editMessage(message, newText, chatType.value, receiverID.value)
  }

  function pinMessage(message, shouldPin = true) {
    console.info(chatType.value)
    messagesStore.pinMessage(message, shouldPin, chatType.value, receiverID.value)
  }

  // --- FRIENDS ---
  const getFriends = friendsStore.getFriends
  const sendFriendRequest = friendsStore.sendFriendRequest
  const cancelFriendRequest = friendsStore.cancelFriendRequest
  const acceptFriendRequest = friendsStore.acceptFriendRequest
  const declineFriendRequest = friendsStore.declineFriendRequest
  const deleteFriend = friendsStore.deleteFriend

  return {
    // state
    messages,
    users,
    connected,
    shouldScroll,
    typingUser,
    imageUrlCache,
    friendStatusCache,
    receiverID,

    // mode
    setChatModeGlobal,
    setChatModePrivate,

    // connection
    startChat,
    stopChat,

    // actions
    sendMessageData,
    sendTyping,
    setTyping,
    fetchUsers,
    openOrCreatePrivateChat,
    deleteMessage,
    editMessage,
    pinMessage,
    getFriends,
    sendFriendRequest,
    cancelFriendRequest,
    acceptFriendRequest,
    declineFriendRequest,
    deleteFriend
  }
})

import { defineStore, storeToRefs } from 'pinia'
import { ref } from 'vue'
import { connect, disconnect, sendMessage } from '../api/chatApi'
import { isNearBottom } from '../utils/chatWindowUtils'
import { useAuthStore } from '../../auth/store/authStore'
import { useStorage } from '@vueuse/core'

import { useMessagesStore } from './messagesStore'
import { useUserStore } from './userStore'
import { useFriendsStore } from './friendsStore'
import { useEventsStore } from './eventsStore'
import { useCallStore } from './call/callStore'
import { useWebRTCStore } from './call/webrtcStore'

export const useChatStore = defineStore('chat', () => {
  const messagesStore = useMessagesStore()
  const userStore = useUserStore()
  const friendsStore = useFriendsStore()
  const eventsStore = useEventsStore()
  const callStore = useCallStore()
  const webrtcStore = useWebRTCStore()

  const { messages, shouldScroll } = storeToRefs(messagesStore)
  const { users } = storeToRefs(userStore)
  const { friendStatusCache } = storeToRefs(friendsStore)
  const { typingUser } = storeToRefs(eventsStore)

  const connected = ref(false)
  const chatType = ref('global')
  const receiverID = ref(null)
  const imageUrlCache = useStorage('image-url-cache', {})

  // --- MODE MANAGEMENT ---
  function setChatModeGlobal() {
    chatType.value = 'global'
    receiverID.value = null
    eventsStore.clearTyping()   
    messagesStore.clear()
    localStorage.setItem('chatMode', 'global')
    if (connected.value) {
      sendMessage({ type: 'init_global', chat_type: 'global' })
    }
  }

  function setChatModePrivate(targetID) {
    chatType.value = 'private'
    receiverID.value = targetID
    eventsStore.clearTyping()   
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

      if (msg.event === 'user_status') {
        userStore.applyStatus(msg.user_id, msg.status)
        return
      }

      if (msg.type === 'typing') {
        eventsStore.handleTypingWs(msg, chatType.value, myId, receiverID.value)
        return
      }

      if (msg.type === 'message_deleted') {
        messagesStore.handleDeleted(msg.message_id)
        return
      }

      if (msg.type === 'message_edited') {
        console.info('message_edited payload:', JSON.stringify(msg))
        messagesStore.handleEdited(msg.message_id, msg.new_text, msg.edited_at)
        return
      }

      if (msg.type === 'message_pinned') {
        console.info('message_pinned payload:', JSON.stringify(msg))
        messagesStore.handlePinned(msg.message_id, msg.action)
        return
      }

      if (msg.type === 'friend_request_update') {
        friendsStore.applyFriendRequestUpdate(msg)
        return
      }

      if (msg.type === 'incoming_call') {
        callStore.handleIncomingCall(msg.from_user)
        return
      }

      if (msg.type === 'incoming_cancel_call') {
        callStore.handleCallCanceled(msg.from_user)
        return
      }

      if (msg.type === 'incoming_call_answer') {
        callStore.handleCallAnswer(msg.from_user, msg.accepted)
      }

      if (msg.type === 'incoming_join_call') {
        callStore.handleJoinCall(msg.from_user)
      }

      if (msg.type === 'incoming_leave_call') {
        delete callStore.callMembers[msg.from_user]
      }

      if (msg.type === 'incoming_webrtc_offer') {
        webrtcStore.handleOffer(msg.from_user, msg.offer)
        return
      }

      if (msg.type === 'incoming_webrtc_answer') {
        webrtcStore.handleAnswer(msg.from_user, msg.answer)
        return
      }

      if (msg.type === 'incoming_ice_candidate') {
        webrtcStore.handleIceCandidate(msg.from_user, msg.candidate)
        return
      }

      if (msg.type === 'incoming_camera_status') {
        callStore.updateCameraStatus(msg.from_user, msg.enabled)
      }
            
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
    eventsStore.sendTyping(chatType.value, receiverID.value)
  }
  function setTyping(username) {
    eventsStore.setTyping(username)
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
  const sendFriendRequest  = friendsStore.sendFriendRequest
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
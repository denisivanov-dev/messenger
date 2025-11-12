import { defineStore } from 'pinia'
import { ref } from 'vue'
import { connect, disconnect } from '../ws/connect'
import { useAuthStore } from '../../../../auth/store/authStore'
import { isNearBottom } from '../../chat/utils/chatWindowUtils'
import { wsHandlersMap } from '../router/wsRouter'

import { useMessagesStore } from '../../chat/store/messagesStore'
import { useUserStore } from '../../chat/store/userStore'
import { useFriendsStore } from '../../friends/store/friendsStore'
import { useTypingStore } from '../../chat/store/events/typingStore'
import { useCallStore } from '../../call/store/callStore'
import { useWebRTCStore } from '../../call/store/webrtcStore'
import { useChatModeStore } from '../../chat/store/chatModeStore'

export const useConnectionStore = defineStore('connection', () => {
  const connected = ref(false)
  const shouldScroll = ref(true)

  const messagesStore = useMessagesStore()
  const userStore = useUserStore()
  const friendsStore = useFriendsStore()
  const typingStore = useTypingStore()
  const callStore = useCallStore()
  const webrtcStore = useWebRTCStore()
  const modeStore = useChatModeStore()

  function startChat(token) {
    messagesStore.clear()
    const authStore = useAuthStore()
    const myId = authStore.getUserId

    connect(token, (msg) => {
        console.info('[WS message]', JSON.stringify(msg, null, 2))

        const routeKey = `${msg.kind}_${msg.action}`
        const handler = wsHandlersMap[routeKey] || wsHandlersMap[msg.kind]

        if (handler) {
        handler(msg, {
            userStore,
            messagesStore,
            friendsStore,
            typingStore,
            callStore,
            webrtcStore,
            chatType: modeStore.chatType,
            receiverID: modeStore.receiverID,
            myId
        })
        return
        }

      const shouldAutoScroll = isNearBottom()
      messagesStore.pushFromWs(msg, shouldAutoScroll)
    })

    shouldScroll.value = true
    connected.value = true
  }

  function stopChat() {
    disconnect()
    connected.value = false
  }

  return {
    connected,
    shouldScroll,
    startChat,
    stopChat
  }
})
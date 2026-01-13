import { defineStore } from 'pinia'
import { storeToRefs } from 'pinia'
import { useStorage } from '@vueuse/core'
import { ref } from 'vue'

import { useChatModeStore } from './chatModeStore'
import { useConnectionStore } from '../../connection/store/connectionStore'
// import { useMessageActions } from './messageActions'
import { useTypingStore } from './events/typingStore'
// import { useUserActions } from './userActions'
// import { useFriendsActions } from './friendsActions'
import { useMessagesStore } from './messagesStore'
import { useUserStore } from './userStore'

export const useChatStore = defineStore('chat', () => {
  const mode = useChatModeStore()
  const connection = useConnectionStore()
  // const messageActions = useMessageActions()
  const typing = useTypingStore()
  // const userActions = useUserActions()
  // const friendsActions = useFriendsActions()
  const messagesStore = useMessagesStore()
  const userStore = useUserStore()

  const { messages, shouldScroll } = storeToRefs(messagesStore)
  const { users } = storeToRefs(userStore)
  const { typingUsers } = storeToRefs(typing)
  // const { friendStatusCache } = storeToRefs(friendsActions)
  const { receiverID } = storeToRefs(mode)

  const imageUrlCache = useStorage('image-url-cache', {})

  return {
    // state
    messages,
    users,
    shouldScroll,
    typingUsers,
    // friendStatusCache,
    receiverID,
    imageUrlCache,

    // mode
    ...mode,

    // connection
    ...connection,

    // actions
    // ...messageActions,
    ...typing,
    // ...userActions,
    // ...friendsActions,
  }
})
<template>
  <div class="flex flex-col items-center w-screen h-screen">
    <!-- Панель звонка -->
    <div class="flex flex-col items-center w-full max-w-[1900px] px-4 mt-4">
      <ChatToolbar />
    </div>

    <!-- Чат -->
    <div class="flex justify-center items-start w-full flex-grow">
      <div class="flex flex-col items-center flex-grow max-w-[1900px] px-4 w-full">
        <ChatWindow
          class="h-[730px]"
          @edit-message="handleEditMessage"
          @reply-to-message="handleReplyMessage"
        />
        <MessageUserInput ref="msgInputRef" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { nextTick, onMounted, ref, watch} from 'vue'
import ChatWindow from '../components/chat-window/chatWindow.vue'
import MessageUserInput from '../components/message-user-input/messageUserInput.vue'
import ChatToolbar from '../components/chat-toolbar/chatToolbar.vue'
import { useAuthStore } from '../../../../auth/store/authStore'
import { useChatStore } from '../store/chatStore'
import { useCallStore } from '../../call/store/callStore'
import { useRouter } from 'vue-router'
import { getRoomId } from '../utils/chatRoom'
import { getCallRoomStatus } from '../../call/api/voiceApi'

const authStore = useAuthStore()
const chatStore = useChatStore()
const callStore = useCallStore()
const router = useRouter()
const msgInputRef = ref(null)

onMounted(() => {
  let stop

  stop = watch(
    () => authStore.getUserId,
    async (myId) => {
      if (!myId) return
      if (typeof stop === 'function') stop()

      const receiverId = chatStore.receiverID || localStorage.getItem('receiverId')
      if (!receiverId) {
        chatStore.setChatModeGlobal()
        router.push('/global-chat')
        return
      }

      chatStore.setChatModePrivate(receiverId)
      const roomId = getRoomId('private', myId, receiverId)

      try {
        const response = await getCallRoomStatus(roomId)
        console.info(response)

        if (Object.keys(response).length > 0) {
          console.log('Активный звонок:', response)

          const members = {}
          const cameraMap = {}
          const screenMap = {}
          const micMap = {}

          for (const [key, value] of Object.entries(response)) {
            if (key.startsWith('cam:')) {
              const userId = key.slice(4)
              cameraMap[userId] = value === 'on'
            } else if (key.startsWith('screen:')) {
              const userId = key.slice(7)
              screenMap[userId] = value === 'on'
            } else if (key.startsWith('mic:')) {
              const userId = key.slice(4)
              micMap[userId] = value === 'on'
            } else {
              members[key] = value
            }
          }

          callStore.resetCallState()
          await nextTick()

          callStore.callMembers = members
          callStore.cameraStatusMap = cameraMap
          callStore.screenStatusMap = screenMap
          callStore.micStatusMap = micMap

          if (members[String(myId)] === 'joined') {
            await callStore.leaveCall()
            await nextTick()  
            await callStore.joinCall()
          }  
        }
      } catch (err) {
        console.error('Ошибка при получении статуса звонка:', err)
      }

      chatStore.shouldScroll = true
    },
    { immediate: true }
  )
})

function handleEditMessage(message) {
  msgInputRef.value?.startEdit(message)
}

function handleReplyMessage(message) {
  msgInputRef.value?.startReply(message)
}
</script>
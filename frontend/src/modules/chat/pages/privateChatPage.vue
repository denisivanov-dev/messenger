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
import { onMounted, ref, watch} from 'vue'
import ChatWindow from '../components/ui/chatWindow/chatWindow.vue'
import MessageUserInput from '../components/ui/chatWindow/messageUserInput.vue'
import ChatToolbar from '../components/ui/chatWindow/chatToolbar.vue'
import { useAuthStore } from '../../auth/store/authStore'
import { useChatStore } from '../store/chatStore'
import { useCallStore } from '../store/call/callStore'
import { useRouter } from 'vue-router'
import { getRoomId } from '../utils/chatRoom'
import { getCallRoomStatus } from '../api/chatApi'

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

        if (Object.keys(response).length > 0) {
          console.log('Активный звонок:', response)
          callStore.callMembers = response

          if (response[String(myId)] === 'joined') {
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
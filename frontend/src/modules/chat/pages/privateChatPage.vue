<template>
  <div class="flex justify-center items-center w-screen h-screen">
    <!-- Центр: чат + инпут -->
      <div class="flex flex-col items-center flex-grow">
        <ChatWindow
          @edit-message="handleEditMessage"
          @reply-to-message="handleReplyMessage"
        />
        <MessageUserInput ref="msgInputRef" />
      </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import ChatWindow from '../components/ui/chatWindow/chatWindow.vue'
import MessageUserInput from '../components/ui/chatWindow/messageUserInput.vue'
import { useChatStore } from '../store/chatStore'
import { useRouter } from 'vue-router'

const chatStore = useChatStore()
const router = useRouter()
const msgInputRef = ref(null)

onMounted(() => {
  const receiverId = chatStore.receiverId || localStorage.getItem('receiverId')
  console.info(receiverId)
  if (receiverId) {
    chatStore.setChatModePrivate(receiverId)
  } else {
    chatStore.setChatModeGlobal()
    router.push('/global-chat')
  }

  chatStore.shouldScroll = true
})

function handleEditMessage(message) {
  msgInputRef.value?.startEdit(message)
}

function handleReplyMessage(message) {
  msgInputRef.value?.startReply(message)
}
</script>
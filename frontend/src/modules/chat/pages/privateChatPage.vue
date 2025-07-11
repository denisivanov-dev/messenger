<template>
  <div class="flex justify-center items-center w-screen h-screen">
    <div class="flex flex-col flex-grow max-w-4xl w-full px-4 gap-6">
      <ChatWindow />
      <MessageUserInput />
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import ChatWindow from '../components/ui/chatWindow/chatWindow.vue'
import MessageUserInput from '../components/ui/chatWindow/messageUserInput.vue'
import { useChatStore } from '../store/chatStore'
import { useRouter } from 'vue-router'

const chatStore = useChatStore()
const router = useRouter()

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
</script>
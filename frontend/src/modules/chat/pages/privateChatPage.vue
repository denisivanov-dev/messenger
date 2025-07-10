<template>
  <div class="flex justify-center items-center w-screen h-screen">
    <div class="flex flex-col flex-grow max-w-4xl w-full px-4 gap-6">
      <!-- Центр: чат + инпут -->
      <ChatWindow />
      <MessageUserInput />
    </div>
  </div>
</template>

<script setup>
import { onMounted, watch } from 'vue'
import ChatWindow from '../components/ui/chatWindow/chatWindow.vue'
import MessageUserInput from '../components/ui/chatWindow/messageUserInput.vue'
import { useChatStore } from '../store/chatStore'
import { useAuthStore } from '../../auth/store/authStore'

const authStore = useAuthStore()
const chatStore = useChatStore()

onMounted(() => {
  setTimeout(() => {
    chatStore.shouldScroll = true
  }, 20)
})
</script>  
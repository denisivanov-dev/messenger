<template>
  <div class="mt-6 w-full max-w-[1300px]">
    <input
      v-model="text"
      @input="chatStore.sendTyping"
      @keyup.enter="send"
      type="text"
      placeholder="Type a message..."
      class="block w-full h-[60px] p-4 rounded-2xl border border-gray-300 shadow-md bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useChatStore } from '../../../store/chatStore'

const text = ref('')
const chatStore = useChatStore()

function send() {
  if (text.value.trim() === '') return
  chatStore.send(text.value)
  text.value = ''
}
</script>
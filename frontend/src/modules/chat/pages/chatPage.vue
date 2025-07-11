<template>
  <div class="flex justify-center items-center w-screen h-screen">
    <div class="flex w-full max-w-[1900px] px-4 gap-6 items-start">
      <!-- Левое окно -->
      <SideBarWindow class="shrink-0" />

      <!-- Центр: чат + инпут -->
      <div class="flex flex-col items-center flex-grow">
        <ChatWindow />
        <MessageUserInput />
      </div>

      <!-- Правое окно -->
      <UserListWindow class="shrink-0" />
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import SideBarWindow from '../components/ui/sideBarWindow.vue'
import ChatWindow from '../components/ui/chatWindow/chatWindow.vue'
import MessageUserInput from '../components/ui/chatWindow/messageUserInput.vue'
import UserListWindow from '../components/ui/userList/userListWindow.vue'
import { useChatStore } from '../store/chatStore'

const chatStore = useChatStore()

onMounted(() => {
  chatStore.setChatModeGlobal()
  chatStore.shouldScroll = true
})
</script>
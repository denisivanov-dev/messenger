<template>
  <div class="w-screen h-screen flex justify-center bg-[var(--bg)]">
    <div class="flex w-full max-w-[1920px]">

      <SideBarWindow ref="sidebarRef" class="shrink-0" />

      <div class="flex flex-col flex-grow overflow-hidden w-full max-w-[1400px] relative">

        <ChatZoom />

        <ChatWindow
          class="flex-1"
          @edit-message="handleEditMessage"
          @reply-to-message="handleReplyMessage"
        />

        <ChatContextBar
          :replyingMessage="replyingMessage"
          :editingMessage="editingMessage"
          @cancel-reply="replyingMessage = null"
          @cancel-edit="editingMessage = null"
        />

        <MessageUserInput
          ref="msgInputRef"
          class="shrink-0"
          :editingMessage="editingMessage"
        />
      </div>

      <UserListWindow ref="userListRef" class="shrink-0" />
    </div>
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref, provide } from 'vue'

import SideBarWindow from '../../AppSidebar/SideBarWindow.vue'
import ChatWindow from '../components/chat-window/chatWindow.vue'
import MessageUserInput from '../components/message-user-input/messageUserInput.vue'
import UserListWindow from '../components/user-list/userListWindow.vue'
import ChatZoom from '../components/chat-window/ChatZoom.vue'
import ChatContextBar from '../components/ChatContextBar.vue/ChatContextBar.vue'
import { isNearBottom } from '../utils/chatWindowUtils'

import { useChatStore } from '../store/chatStore'
import { useMessagesStore } from '../store/messagesStore'
import { useAuthStore } from '../../../../auth/store/authStore'

const authStore = useAuthStore()
const chatStore = useChatStore()
const messagesStore = useMessagesStore()

const msgInputRef = ref(null)
const sidebarRef = ref(null)
const userListRef = ref(null)
const replyingMessage = ref(null)
const editingMessage = ref(null)

const zoom = ref(1)
provide('zoom', zoom)

/* message actions */
function handleReplyMessage(message) {
  editingMessage.value = null
  messagesStore.shouldScroll = isNearBottom()
  replyingMessage.value = message.payload
}

function handleEditMessage(message) {
  replyingMessage.value = null
  messagesStore.shouldScroll = isNearBottom()
  editingMessage.value = message.payload
}

/* close side panels */
function handleBodyClick(e) {
  const profileEl = userListRef.value?.profileRef?.profileRootElement
  const friendPanelEl = sidebarRef.value?.friendPanelRef?.$el
  const friendButton = sidebarRef.value?.$el?.querySelector('button')
  const cancelConfirmEl = userListRef.value?.profileRef?.cancelConfirmRef?.value

  if (
    profileEl?.contains(e.target) ||
    friendPanelEl?.contains(e.target) ||
    friendButton?.contains(e.target) ||
    cancelConfirmEl?.contains(e.target)
  ) return

  userListRef.value?.closeProfile?.()
}

/* init */
onMounted(async () => {
  document.body.addEventListener('click', handleBodyClick)

  chatStore.setChatModeGlobal()
  chatStore.shouldScroll = true

  while (!authStore.getUserId)
    await new Promise(r => setTimeout(r, 50))

  chatStore.getFriends(authStore.getUserId)
})

onBeforeUnmount(() => {
  document.body.removeEventListener('click', handleBodyClick)
})
</script>
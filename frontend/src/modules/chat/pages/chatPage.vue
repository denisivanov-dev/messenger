<template>
  <div class="flex justify-center items-center w-screen h-screen">
    <div class="flex w-full max-w-[1900px] px-4 gap-6 items-start">
      <!-- Левое окно -->
      <SideBarWindow ref="sidebarRef" class="shrink-0" />

      <!-- Центр: чат + инпут -->
      <div class="flex flex-col items-center flex-grow h-[865px]">
        <ChatWindow
          @edit-message="handleEditMessage"
          @reply-to-message="handleReplyMessage"
        />
        <MessageUserInput ref="msgInputRef" />
      </div>

      <!-- Правое окно -->
      <UserListWindow ref="userListRef" class="shrink-0" />
    </div>
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref } from 'vue'
import SideBarWindow from '../components/ui/sideBar/sideBarWindow.vue'
import ChatWindow from '../components/ui/chatWindow/chatWindow.vue'
import MessageUserInput from '../components/ui/chatWindow/messageUserInput.vue'
import UserListWindow from '../components/ui/userList/userListWindow.vue'
import { useChatStore } from '../store/chatStore'
import { useAuthStore } from '../../auth/store/authStore'

const authStore = useChatStore()
const chatStore = useChatStore()

const msgInputRef = ref(null)
const sidebarRef = ref(null)
const userListRef = ref(null)

function handleEditMessage(message) {
  msgInputRef.value?.startEdit(message)
}

function handleReplyMessage(message) {
  msgInputRef.value?.startReply(message)
}

function handleBodyClick(e) {
  const profileEl = userListRef.value?.profileRef?.profileRootElement
  const friendPanelEl = sidebarRef.value?.friendPanelRef?.$el
  const friendButton = sidebarRef.value?.$el?.querySelector('button')
  const cancelConfirmEl = userListRef.value?.profileRef?.cancelConfirmRef?.value

  const clickedInsideProfile = profileEl && profileEl.contains(e.target)
  const clickedInsideFriendPanel = friendPanelEl && friendPanelEl.contains(e.target)
  const clickedFriendButton = friendButton && friendButton.contains(e.target)
  const clickedInsideCancelConfirm = cancelConfirmEl && cancelConfirmEl.contains(e.target)

  if (
    clickedInsideProfile ||
    clickedInsideFriendPanel ||
    clickedFriendButton ||
    clickedInsideCancelConfirm
  ) return

  userListRef.value?.closeProfile?.()
}

onMounted(async () => {
  document.body.addEventListener('click', handleBodyClick)
  chatStore.setChatModeGlobal()
  chatStore.shouldScroll = true

  const authStore = useAuthStore()

  while (!authStore.getUserId) {
    await new Promise(resolve => setTimeout(resolve, 50))
  }

  const myId = authStore.getUserId
  console.info(myId)
  
  try {
    chatStore.getFriends(myId)
  } catch (err) {
    console.error('Не удалось загрузить друзей:', err)
  }
})

onBeforeUnmount(() => {
  document.body.removeEventListener('click', handleBodyClick)
})
</script>
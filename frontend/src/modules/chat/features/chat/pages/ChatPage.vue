<template>
  <div class="w-screen h-screen flex justify-center bg-[var(--bg)]">
    <div class="flex w-full max-w-[1900px] px-4 gap-6">

      <!-- Левый сайдбар -->
      <SideBarWindow ref="sidebarRef" class="shrink-0" />

      <!-- Центральная колонка (чаты) -->
      <div class="flex flex-col flex-grow overflow-hidden w-full max-w-[1400px]">

        <!-- Окно чата -->
        <ChatWindow
          class="flex-1"
          @edit-message="handleEditMessage"
          @reply-to-message="handleReplyMessage"
        />

        <!-- Поле ввода -->
        <MessageUserInput 
          ref="msgInputRef" 
          class="shrink-0"
        />

      </div>

      <!-- Список пользователей справа -->
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
import { useChatStore } from '../store/chatStore'
import { useAuthStore } from '../../../../auth/store/authStore'

const authStore = useAuthStore()
const chatStore = useChatStore()

const msgInputRef = ref(null)
const sidebarRef = ref(null)
const userListRef = ref(null)

/* ===== ZOOM ===== */
const zoom = ref(1)
provide("zoom", zoom)

function applyZoom(v) {
  zoom.value = Math.min(2, Math.max(0.5, zoom.value + v))
}

function zoomWheel(e) {
  if (!e.ctrlKey) return
  e.preventDefault()
  applyZoom(e.deltaY < 0 ? +0.05 : -0.05)
}

function zoomKeys(e) {
  if (!e.ctrlKey) return
  if (e.key === '+' || e.key === '=') { applyZoom(+0.05); e.preventDefault() }
  if (e.key === '-') { applyZoom(-0.05); e.preventDefault() }
  if (e.key === '0') { zoom.value = 1; e.preventDefault() }
}

function blockBrowserZoom(e) {
  if (e.ctrlKey) e.preventDefault()
}

function blockBrowserZoomKeys(e) {
  if (e.ctrlKey && ['+', '-', '=', '0'].includes(e.key))
    e.preventDefault()
}

/* ===== MESSAGE ACTIONS ===== */
function handleEditMessage(message) {
  msgInputRef.value?.startEdit(message)
}

function handleReplyMessage(message) {
  msgInputRef.value?.startReply(message)
}

/* ===== CLOSE PANELS ===== */
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

/* ===== LIFECYCLE ===== */
onMounted(async () => {
  document.body.addEventListener('click', handleBodyClick)

  chatStore.setChatModeGlobal()
  chatStore.shouldScroll = true

  while (!authStore.getUserId)
    await new Promise(r => setTimeout(r, 50))

  chatStore.getFriends(authStore.getUserId)

  window.addEventListener("wheel", zoomWheel, { passive: false })
  window.addEventListener("keydown", zoomKeys, { passive: false })
  window.addEventListener("wheel", blockBrowserZoom, { passive: false })
  window.addEventListener("keydown", blockBrowserZoomKeys, { passive: false })
})

onBeforeUnmount(() => {
  document.body.removeEventListener('click', handleBodyClick)

  window.removeEventListener("wheel", zoomWheel)
  window.removeEventListener("keydown", zoomKeys)
  window.removeEventListener("wheel", blockBrowserZoom)
  window.removeEventListener("keydown", blockBrowserZoomKeys)
})
</script>
<template>
  <router-view />
  <div id="global-particle-layer" class="fixed inset-0 pointer-events-none z-[9999]"></div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useAuthStore } from './modules/auth/store/authStore'
import { useChatStore } from './modules/chat/features/chat/store/legacyChatStore'
import { useChatModeStore } from './modules/chat/features/chat/store/chatModeStore'
import { useConnectionStore } from './modules/chat/features/connection/store/connectionStore'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const chatStore = useChatStore()
const chatMode = useChatModeStore()
const connection = useConnectionStore()
const router = useRouter()

let refreshTimer = null

async function initAuthFlow () {
  await authStore.autoLogin()

  const token = authStore.getAccessToken
  if (!token) return

  const mode = localStorage.getItem('chatMode') || 'global'
  const receiverId = localStorage.getItem('receiverId') || null

  if (mode === 'private' && receiverId) {
    chatMode.setChatModePrivate(receiverId)
  } else {
    chatMode.setChatModeGlobal()
  }

  connection.startChat(token)

  // chatStore.fetchUsers()

  const guestPages = ['/', '/login', '/register', '/forgot-password', '/confirm-registration']
  if (guestPages.includes(router.currentRoute.value.path)) {
    router.replace('/global-chat')
  }

  if (!refreshTimer) {
    refreshTimer = setInterval(async () => {
      const ok = await authStore.refreshToken?.()
      if (ok && chatStore.isConnected?.()) {
        chatStore.reconnectIfNeeded?.()
      }
    }, 14 * 60 * 1000)
  }
}

onMounted(initAuthFlow)

onUnmounted(() => {
  if (refreshTimer) clearInterval(refreshTimer)
})

function onGlobalClick(e) {
  const insideMenu = e.target.closest('.message-click_actions')
  if (insideMenu) return

  window.__closeActiveClickMenu?.()
  
}

function onGlobalKeydown(e) {
  if (e.key === 'Escape') {
    window.__closeActiveClickMenu?.()
    window.__cancelChatContext?.()
    return
  }

  if (!window.__closeActiveClickMenu) return

  if (e.key.length === 1 || e.key === 'Backspace' || e.key === 'Enter') {
    window.__closeActiveClickMenu()
  }
}

onMounted(() => {
  document.addEventListener('click', onGlobalClick)
  document.addEventListener('keydown', onGlobalKeydown)
})

onUnmounted(() => {
  document.removeEventListener('click', onGlobalClick)
  document.removeEventListener('keydown', onGlobalKeydown)
})
</script>
<template>
  <div
    ref="userListWindowRef"
    class="user-list-window w-[250px] h-[864px] p-5 bg-gray-100 rounded-2xl shadow-md overflow-y-auto flex flex-col gap-4"
  >
    <!-- Онлайн пользователи -->
    <div>
      <h2 class="text-lg font-semibold mb-2">Online</h2>
      <div id="online-users" class="flex flex-col gap-2">
        <div
          v-for="user in onlineUsers"
          :key="user.id"
          class="bg-green-100 p-2 rounded cursor-pointer hover:bg-green-200"
          @click.stop="openUserProfile(user, $event)"
        >
          {{ user.username }}
        </div>
      </div>
    </div>

    <!-- Офлайн пользователи -->
    <div>
      <h2 class="text-lg font-semibold mt-4 mb-2">Offline</h2>
      <div id="offline-users" class="flex flex-col gap-2">
        <div
          v-for="user in offlineUsers"
          :key="user.id"
          class="bg-gray-200 p-2 rounded cursor-pointer hover:bg-gray-300"
           @click.stop="openUserProfile(user, $event)"
        >
          {{ user.username }}
        </div>
      </div>
    </div>

    <!-- Компонент профиля -->
    <UserProfile
      ref="profileRef"
      v-if="isProfileVisible"
      :user="selectedUser"
      :visible="isProfileVisible"
      :x="profileX"
      :y="profileY"
      :is-self="isSelf" 
      @write="startChatWithUser"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick, computed } from 'vue'
import { useChatStore } from '../../../store/chatStore'
import UserProfile from './userProfile.vue'
import { startPrivateChat } from '../../../api/chatApi'
import { useAuthStore } from '../../../../auth/store/authStore'
import { useRouter } from 'vue-router'

const chatStore = useChatStore()
const authStore = useAuthStore()
const router = useRouter()

const myId = computed(() => String(authStore.user.id))
const onlineUsers  = computed(() => chatStore.users.filter(u => u.status?.trim() === 'online'))
const offlineUsers = computed(() => chatStore.users.filter(u => u.status?.trim() === 'offline'))

const userListWindowRef = ref(null)
const profileRef        = ref(null)
const selectedUser      = ref(null)
const previousSelectedUser = ref(null)
const isProfileVisible  = ref(false)
const profileX = ref(0)
const profileY = ref(0)

const isSelf = computed(() =>
  selectedUser.value && String(selectedUser.value.id) === myId.value
)

function openUserProfile(user, event) {
  const rect = event.currentTarget.getBoundingClientRect()

  if (isProfileVisible.value && previousSelectedUser.value?.id === user.id) {
    closeProfile()
    return
  }

  selectedUser.value = user
  previousSelectedUser.value = user

  isProfileVisible.value = true
  profileX.value = rect.left - 320
  profileY.value = rect.top - 50

  setTimeout(() => {
    document.body.removeEventListener('click', handleBodyClick)
    document.body.addEventListener('click', handleBodyClick)
  }, 0)
}

async function startChatWithUser() {
  try {
    closeProfile()
    const response = await chatStore.openOrCreatePrivateChat(selectedUser.value.id)
    console.info(response)
    if (response === true) router.push('/private-chat')
  } catch (err) {
    console.error('Не удалось открыть чат:', err)
  }
}

function handleBodyClick(e) {
  if (!isProfileVisible.value) return
  
  const profileElement = profileRef.value?.profileRootElement
  if (profileElement && profileElement.contains(e.target)) return

  closeProfile()
}

function closeProfile() {
  isProfileVisible.value = false
  
  document.body.removeEventListener('click', handleBodyClick)
}

// onMounted(() => chatStore.fetchUsers())

onBeforeUnmount(() => document.body.removeEventListener('click', handleBodyClick))
</script>
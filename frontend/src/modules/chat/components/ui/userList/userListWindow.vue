<template>
  <div
    ref="userListWindowRef"
    class="user-list-window w-[250px] h-[864px] p-5 bg-gray-100 rounded-2xl shadow-md overflow-y-auto flex flex-col gap-4 select-none"
  >
    <!-- Онлайн пользователи -->
    <div>
      <h2 class="text-lg font-semibold mb-2">Online</h2>
      <div class="flex flex-col gap-2">
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
      <div class="flex flex-col gap-2">
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

    <!-- Профиль -->
    <transition name="fade-slide-right">
      <UserProfile
        v-if="isProfileVisible"
        :key="selectedUser?.id"
        ref="profileRef"
        :user="selectedUser"
        :visible="isProfileVisible"
        :x="profileX"
        :y="profileY"
        :is-self="isSelf"
        @write="startChatWithUser"
        @add-friend="addFriendHandler"
        @cancel-friend="cancelFriendHandler"
        @accept-friend="acceptFriendHandler"
        @decline-friend="declineFriendHadnler"
      />
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, defineExpose } from 'vue'
import { useRouter } from 'vue-router'
import { useChatStore } from '../../../store/chatStore'
import { useAuthStore } from '../../../../auth/store/authStore'
import UserProfile from './userProfile.vue'

const chatStore = useChatStore()
const authStore = useAuthStore()
const router = useRouter()

const userListWindowRef = ref(null)
const profileRef = ref(null)

const selectedUser = ref(null)
const previousSelectedUser = ref(null)
const isProfileVisible = ref(false)
const profileX = ref(0)
const profileY = ref(0)

const myId = computed(() => String(authStore.user.id))
const onlineUsers = computed(() =>
  Object.values(chatStore.users).filter(u => u.status?.trim() === 'online')
)
const offlineUsers = computed(() =>
  Object.values(chatStore.users).filter(u => u.status?.trim() === 'offline')
)

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
}

function closeProfile() {
  isProfileVisible.value = false
}

async function startChatWithUser() {
  try {
    closeProfile()
    const response = await chatStore.openOrCreatePrivateChat(selectedUser.value.id)
    if (response === true) router.push('/private-chat')
  } catch (err) {
    console.error('Не удалось открыть чат:', err)
  }
}

async function addFriendHandler() {
  if (!selectedUser.value) return
  const myId = authStore.getUserId

  try {
    await chatStore.sendFriendRequest(myId, selectedUser.value.id)
    console.log(`Запрос в друзья отправлен пользователю ${selectedUser.value.username}`)
  } catch (err) {
    console.error('Ошибка при добавлении в друзья:', err)
  }
}

async function cancelFriendHandler() {
  if (!selectedUser.value) return
  const myId = authStore.getUserId

  try {
    await chatStore.cancelFriendRequest(myId, selectedUser.value.id)
    console.log(`Заявка в друзья отменена: ${selectedUser.value.username}`)
  } catch (err) {
    console.error('Ошибка при отмене заявки:', err)
  }
}

async function acceptFriendHandler() {
  if (!selectedUser.value) return
  const myId = authStore.getUserId

  try {
    await chatStore.acceptFriendRequest(myId, selectedUser.value.id)
    console.log(`Заявка в друзья принята: ${selectedUser.value.username}`)
  } catch (err) {
    console.error('Ошибка при принятии друга:', err)
  }
}

async function declineFriendHadnler() {
  if (!selectedUser.value) return
  const myId = authStore.getUserId

  try {
    await chatStore.declineFriendRequest(myId, selectedUser.value.id)
    console.log(`Заявка в друзья принята: ${selectedUser.value.username}`)
  } catch (err) {
    console.error('Ошибка при принятии друга:', err)
  }
}

defineExpose({
  profileRef,
  closeProfile
})
</script>

<style scoped>
.fade-slide-right-enter-active,
.fade-slide-right-leave-active {
  transition: all 0.3s ease;
}
.fade-slide-right-enter-from,
.fade-slide-right-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
.fade-slide-right-enter-to,
.fade-slide-right-leave-from {
  opacity: 1;
  transform: translateX(0);
}

</style>
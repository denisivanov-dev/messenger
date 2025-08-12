<template>
  <div class="space-y-2">
    <!-- Внутренние табы -->
    <div class="flex gap-2 mb-2 no-drag">
      <button
        @click="pendingTab = 'incoming'"
        :class="[
          'text-sm px-3 py-1 rounded-full transition',
          pendingTab === 'incoming'
            ? 'bg-blue-600 text-white'
            : 'text-gray-500 hover:bg-blue-100'
        ]"
      >
        Входящие
      </button>
      <button
        @click="pendingTab = 'outgoing'"
        :class="[
          'text-sm px-3 py-1 rounded-full transition',
          pendingTab === 'outgoing'
            ? 'bg-blue-600 text-white'
            : 'text-gray-500 hover:bg-blue-100'
        ]"
      >
        Отправленные
      </button>
    </div>

    <!-- Входящие заявки -->
    <ul v-if="pendingTab === 'incoming'" class="space-y-2">
      <li
        v-for="user in incomingRequests"
        :key="user.id"
        class="flex justify-between items-center text-sm text-gray-700 border border-gray-200 rounded-lg px-3 py-2"
      >
        <span class="font-medium">{{ user.username }}</span>

        <div class="flex gap-2 no-drag">
          <!-- Принять -->
          <div class="relative group">
            <button
              @click="acceptFriend(user.id)"
              class="w-8 h-8 rounded-md bg-green-500 text-white hover:bg-green-600 transition flex items-center justify-center"
            >
              <Check class="w-4 h-4" />
            </button>
            <div
              class="absolute bottom-full mb-1 left-1/2 -translate-x-1/2 px-2 py-1 text-xs text-white bg-black rounded opacity-0 group-hover:opacity-100 transition pointer-events-none"
            >
              Принять
            </div>
          </div>

          <!-- Отклонить -->
          <div class="relative group">
            <button
              @click="declineFriend(user.id)"
              class="w-8 h-8 rounded-md bg-red-500 text-white hover:bg-red-600 transition flex items-center justify-center"
            >
              <X class="w-4 h-4" />
            </button>
            <div
              class="absolute bottom-full mb-1 left-1/2 -translate-x-1/2 px-2 py-1 text-xs text-white bg-black rounded opacity-0 group-hover:opacity-100 transition pointer-events-none"
            >
              Отклонить
            </div>
          </div>
        </div>
      </li>

      <li
        v-if="incomingRequests.length === 0"
        class="text-sm text-gray-400 italic text-center"
      >
        Нет входящих заявок
      </li>
    </ul>

    <!-- Отправленные заявки -->
    <ul v-else class="space-y-2">
      <li
        v-for="user in outgoingRequests"
        :key="user.id"
        class="flex justify-between items-center text-sm text-gray-700 border border-gray-200 rounded-lg px-3 py-2"
      >
        <span class="font-medium">{{ user.username }}</span>

        <div class="flex gap-2 no-drag">
          <!-- Ожидание -->
          <div class="relative group">
            <div
              class="w-8 h-8 rounded-md bg-gray-300 text-gray-700 flex items-center justify-center"
            >
              <Clock class="w-4 h-4" />
            </div>
            <div
              class="absolute bottom-full mb-1 left-1/2 -translate-x-1/2 px-2 py-1 text-xs text-white bg-black rounded opacity-0 group-hover:opacity-100 transition pointer-events-none"
            >
              Ожидание
            </div>
          </div>

          <!-- Отменить заявку -->
          <div class="relative group">
            <button
              @click="openCancelConfirm(user)"
              class="w-8 h-8 rounded-md bg-red-500 text-white hover:bg-red-600 transition flex items-center justify-center"
            >
              <X class="w-4 h-4" />
            </button>
            <div
              class="absolute bottom-full mb-1 left-1/2 -translate-x-1/2 px-2 py-1 text-xs text-white bg-black rounded opacity-0 group-hover:opacity-100 transition pointer-events-none"
            >
              Отменить
            </div>
          </div>
        </div>
      </li>

      <li
        v-if="outgoingRequests.length === 0"
        class="text-sm text-gray-400 italic text-center"
      >
        Нет отправленных заявок
      </li>
    </ul>

    <!-- Модалка подтверждения отмены -->
    <div
      v-if="showCancelConfirm"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
    >
      <div
        class="bg-white rounded-xl shadow-lg p-6 w-[300px] text-center space-y-4 select-none cursor-default no-drag"
        @click.stop
      >
        <p class="text-lg font-semibold text-gray-800">
          Отменить заявку {{ cancelCandidate?.username }}?
        </p>
        <div class="flex justify-center gap-4">
          <button
            @click="confirmCancel"
            class="px-4 py-2 rounded-lg bg-green-600 text-white hover:bg-green-700 transition"
          >
            Да
          </button>
          <button
            @click="cancelPrompt"
            class="px-4 py-2 rounded-lg bg-red-500 text-white hover:bg-red-600 transition"
          >
            Нет
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '../../../../../auth/store/authStore'
import { useChatStore } from '../../../../store/chatStore'
import { Check, X, Clock } from 'lucide-vue-next'

const emit = defineEmits(['pending-count'])

const pendingTab = ref('incoming')
const showCancelConfirm = ref(false)
const cancelCandidate = ref(null)

const authStore = useAuthStore()
const chatStore = useChatStore()

// Берём реактивные ссылки из Pinia
const { users, friendStatusCache } = storeToRefs(chatStore)

const incomingRequests = computed(() =>
  (users.value || []).filter(u => friendStatusCache.value?.[u.id] === 'incoming')
)
const outgoingRequests = computed(() =>
  (users.value || []).filter(u => friendStatusCache.value?.[u.id] === 'outgoing')
)

async function acceptFriend(userId) {
  const myId = authStore.getUserId
  await chatStore.acceptFriendRequest(myId, userId)
  // Триггерим реактивность useStorage-объекта
  friendStatusCache.value = {
    ...friendStatusCache.value,
    [userId]: 'friends'
  }
}

async function declineFriend(userId) {
  const myId = authStore.getUserId
  await chatStore.declineFriendRequest(myId, userId)
  const copy = { ...friendStatusCache.value }
  delete copy[userId]
  friendStatusCache.value = copy
}

function openCancelConfirm(user) {
  cancelCandidate.value = user
  showCancelConfirm.value = true
}

function cancelPrompt() {
  showCancelConfirm.value = false
  cancelCandidate.value = null
}

async function confirmCancel() {
  const myId = authStore.getUserId
  const user = cancelCandidate.value
  if (!user) return
  try {
    await chatStore.cancelFriendRequest(myId, user.id)
    const copy = { ...friendStatusCache.value }
    delete copy[user.id]
    friendStatusCache.value = copy
  } catch (err) {
    console.error('Ошибка при отмене заявки:', err)
  } finally {
    cancelPrompt()
  }
}

watch(
  () => incomingRequests.value.length,
  count => emit('pending-count', count),
  { immediate: true }
)
</script>

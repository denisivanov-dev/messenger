<template>
  <ul class="space-y-2">
    <li
      v-for="user in friends"
      :key="user.id"
      class="flex justify-between items-center text-sm text-gray-700 border border-gray-200 rounded-lg px-3 py-2"
    >
      <span class="font-medium">{{ user.username }}</span>

      <div class="flex gap-2 items-center no-drag">
        <div class="relative group">
          <button
            @click="startChatWithUser(user.id)"
            class="p-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition flex items-center justify-center"
          >
            <MessageSquareIcon class="w-4 h-4" />
          </button>
          <div
            class="absolute bottom-full mb-1 left-1/2 -translate-x-1/2 px-2 py-1 text-xs text-white bg-black rounded opacity-0 group-hover:opacity-100 transition pointer-events-none"
          >
            Написать
          </div>
        </div>

        <div class="relative group">
          <button
            @click="showPrompt(user)"
            class="p-2 rounded-lg bg-gray-200 text-gray-700 hover:bg-red-100 hover:text-red-600 transition"
          >
            <TrashIcon class="w-4 h-4" />
          </button>
          <div
            class="absolute bottom-full mb-1 left-1/2 -translate-x-1/2 px-2 py-1 text-xs text-white bg-black rounded opacity-0 group-hover:opacity-100 transition pointer-events-none"
          >
            Удалить
          </div>
        </div>
      </div>
    </li>

    <li v-if="friends.length === 0" class="text-sm text-gray-400 italic text-center">
      У тебя нет друзей
    </li>
  </ul>

  <div
    v-if="selectedForDeletion"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
  >
    <div
      class="bg-white rounded-xl shadow-lg p-6 w-[300px] text-center space-y-4 select-none cursor-default no-drag"
      @click.stop
    >
      <p class="text-lg font-semibold text-gray-800">
        Удалить {{ selectedForDeletion.username }} из друзей?
      </p>
      <div class="flex justify-center gap-4">
        <button
          @click="confirmDelete"
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
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { MessageSquareIcon, TrashIcon } from 'lucide-vue-next'
import { useChatStore } from '../../../../store/chatStore'
import { useAuthStore } from '../../../../../auth/store/authStore'
import { useRouter } from 'vue-router'

const router = useRouter()
const chatStore = useChatStore()
const authStore = useAuthStore()

const props = defineProps({
  isModalOpen: Boolean,
  blockDrag: Boolean
})
const emit = defineEmits([
  'close',
  'update:isModalOpen',
  'update:blockDrag'
])

const selectedForDeletion = ref(null)

const friendStatusCache = chatStore.friendStatusCache
const allUsers = computed(() => Object.values(chatStore.users))

const friends = computed(() =>
  allUsers.value.filter(user => friendStatusCache[user.id] === 'friends')
)

async function startChatWithUser(userId) {
  try {
    const response = await chatStore.openOrCreatePrivateChat(userId)
    if (response === true) {
      router.push('/private-chat')
      emit('close')
    }
  } catch (err) {
    console.error('Не удалось открыть чат:', err)
  }
}

function showPrompt(user) {
  selectedForDeletion.value = user
  emit('update:isModalOpen', true)
}

function cancelPrompt() {
  selectedForDeletion.value = null
  emit('update:isModalOpen', false)
}

function confirmDelete() {
  const myId = authStore.getUserId
  const user = selectedForDeletion.value
  chatStore.deleteFriend(myId, user.id)
  selectedForDeletion.value = null
  emit('update:isModalOpen', false)
}

watch(() => props.isModalOpen, val => {
  emit('update:blockDrag', val)
}, { immediate: true })
</script>
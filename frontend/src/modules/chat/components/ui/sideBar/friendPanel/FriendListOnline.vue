<template>
  <ul class="space-y-2">
    <li
      v-for="user in onlineFriends"
      :key="user.id"
      class="flex justify-between items-center text-sm text-gray-700 border border-gray-200 rounded-lg px-3 py-2"
    >
      <span class="font-medium flex items-center gap-2">
        <!-- Онлайн индикатор -->
        <span class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
        {{ user.username }}
      </span>

      <div class="relative group no-drag">
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
    </li>

    <li v-if="onlineFriends.length === 0" class="text-sm text-gray-400 italic text-center">
      Никто не в сети
    </li>
  </ul>
</template>

<script setup>
import { computed } from 'vue'
import { MessageSquareIcon } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { useChatStore } from '../../../../store/chatStore'

const router = useRouter()
const chatStore = useChatStore()

const friendStatusCache = chatStore.friendStatusCache
const allUsers = computed(() => Object.values(chatStore.users))

const friends = computed(() =>
  allUsers.value.filter(user => friendStatusCache[user.id] === 'friends')
)

const onlineFriends = computed(() =>
  friends.value.filter(user => user.status === 'online')
)

async function startChatWithUser(userId) {
  try {
    const response = await chatStore.openOrCreatePrivateChat(userId)
    if (response === true) {
      router.push('/private-chat')
    }
  } catch (err) {
    console.error('Не удалось открыть чат:', err)
  }
}
</script>
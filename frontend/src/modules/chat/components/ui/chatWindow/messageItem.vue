<template>
  <div
    class="relative group px-4 py-2 bg-white rounded-xl shadow hover:bg-gray-50 transition"
  >
    <!-- Индикатор "изменено" и дата -->
    <div v-if="props.message.edited_at" class="absolute top-0 right-0 mt-1 mr-2 text-[10px] text-gray-400">
      изменено • {{ formattedEditDate }}
    </div>

    <!-- Заголовок: имя и время -->
    <div class="mb-1 text-xs text-gray-500">
      {{ props.message.username }} • {{ formattedDate }}
    </div>

    <!-- Текст сообщения -->
    <p class="text-sm text-gray-900">{{ props.message.text }}</p>

    <!-- Ховер-меню -->
    <div
      class="absolute top-0 right-0 mt-1 mr-1 hidden group-hover:flex flex-row bg-white border rounded shadow px-2 py-1 z-10 gap-2"
    >
      <button @click="reply" title="Ответить" class="text-blue-600 hover:text-blue-800">
        <ReplyIcon class="w-4 h-4" />
      </button>

      <button
        v-if="isMyMessage"
        @click="edit"
        title="Редактировать"
        class="text-yellow-600 hover:text-yellow-800"
      >
        <EditIcon class="w-4 h-4" />
      </button>

      <button @click="pin" title="Закрепить" class="text-purple-600 hover:text-purple-800">
        <PinIcon class="w-4 h-4" />
      </button>

      <button
        v-if="isMyMessage"
        @click="remove"
        title="Удалить"
        class="text-red-600 hover:text-red-800"
      >
        <TrashIcon class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { ReplyIcon, EditIcon, PinIcon, TrashIcon } from 'lucide-vue-next'
import { useAuthStore } from '../../../../auth/store/authStore'
import { useChatStore } from '../../../store/chatStore'

const emit = defineEmits(['edit-message'])

const props = defineProps({
  message: {
    type: Object,
    required: true
  }
})

const formattedDate = new Date(props.message.timestamp).toLocaleString('ru-RU', {
  day: '2-digit',
  month: '2-digit',
  year: 'numeric',
  hour: '2-digit',
  minute: '2-digit'
})

const formattedEditDate = computed(() => {
  if (!props.message.edited_at) return ''
  return new Date(props.message.edited_at).toLocaleString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
})

const authStore = useAuthStore()
const chatStore = useChatStore()
const isMyMessage = props.message.user_id == authStore.getUserId

const reply = () => {
  console.log("Ответить:", props.message)
}

const edit = () => {
  console.log("Редактировать:", props.message)
  emit('edit-message', props.message)
}

const pin = () => {
  console.log("Закрепить:", props.message)
}

const remove = () => {
  console.log("Удалить:", props.message)
  chatStore.deleteMessage(props.message)
}
</script>
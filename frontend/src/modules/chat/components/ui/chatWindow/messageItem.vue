<template>
  <div
    :id="`msg-${props.message.message_id}`"
    class="relative group px-4 py-2 bg-white rounded-xl shadow hover:bg-gray-50 transition"
  >
    <!-- Индикатор "изменено" и дата -->
    <div v-if="props.message.edited_at" class="absolute top-0 right-0 mt-1 mr-2 text-[10px] text-gray-400">
      изменено • {{ formattedEditDate }}
    </div>
    
    <!-- Заголовок: имя и время -->
    <div
      class="mb-1 text-xs flex items-center gap-1"
      :class="props.message.pinned ? 'text-purple-800 bg-purple-100 px-1 py-0.5 rounded' : 'text-gray-500'"
    >
      <span>{{ props.message.username }} • {{ formattedDate }}</span>
      <span v-if="props.message.pinned" class="flex items-center gap-1 text-xs">
        📌 <span class="italic m-auto">Закреплено</span>
      </span>
    </div>

    <!-- Ответ на сообщение -->
    <div
      v-if="props.message.reply_to"
      class="mb-1 text-[11px] text-gray-500 border-l-2 border-blue-400 pl-2 cursor-pointer hover:text-blue-600"
      @click="$emit('scroll-to-message', props.message.reply_to)"
    >
      ↩ {{ props.message.reply_to_user }}: 
      <span class="italic text-gray-500">{{ repliedMessageText }}</span>
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

const emit = defineEmits(['reply-to-message', 'edit-message', 'scroll-to-message'])

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

const repliedMessageText = computed(() => {
  const replied = chatStore.messages.find(m => m.message_id === props.message.reply_to)
  if (!replied) return '[сообщение удалено]'
  if (replied.text.trim() !== props.message.reply_to_text.trim()) {
    return replied.text + ' (изменено)'
  }
  return replied.text
})

const authStore = useAuthStore()
const chatStore = useChatStore()
const isMyMessage = props.message.user_id == authStore.getUserId

const reply = () => {
  console.log("Ответить:", props.message)
  emit('reply-to-message', props.message)
}

const edit = () => {
  console.log("Редактировать:", props.message)
  emit('edit-message', props.message)
}

const pin = () => {
  console.log("Закрепить:", props.message)
  const shouldPin = !props.message.pinned
  chatStore.pinMessage(props.message, shouldPin)
}

const remove = () => {
  console.log("Удалить:", props.message)
  chatStore.deleteMessage(props.message)
}
</script>
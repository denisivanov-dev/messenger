<template>
  <div
    :id="`msg-${props.message.message_id}`"
    class="relative group flex items-start gap-3 px-4 py-2 bg-white rounded-xl shadow hover:bg-gray-50 transition"
  >
    <!-- Аватарка -->
    <img
      :src="avatarUrl"
      class="w-8 h-8 rounded-full object-cover mt-0.5"
      alt="аватар"
    />

    <!-- Контейнер контента -->
    <div class="flex-1 relative">
      <!-- Индикатор "изменено" и дата -->
      <div
        v-if="props.message.edited_at"
        class="absolute top-0 right-0 mt-1 mr-2 text-[10px] text-gray-400"
      >
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

      <!-- Атачменты (если есть) -->
      <MessageGallery
        v-if="props.message.attachments && props.message.attachments.some(att => att.type === 'image')"
        :attachments="props.message.attachments.filter(att => att.type === 'image')"
        :imageUrls="attachmentUrls"
        :openImage="(key) => openImage(attachmentUrls[key])"
      />

      <!-- Текст (если есть) -->
      <p v-if="props.message.text" class="text-sm text-gray-900">{{ props.message.text }}</p>

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
  </div>

  <!-- Модалка для увеличенного изображения -->
  <Teleport to="body">
    <div
      v-if="fullscreenImageUrl"
      class="fixed inset-0 bg-black bg-opacity-90 z-50 flex items-center justify-center"
      @click="closeImage"
    >
      <div class="flex flex-col items-center gap-3" @click.stop>
        <img
          :src="fullscreenImageUrl"
          class="max-w-full max-h-[90vh] shadow-xl"
        />
        <a
          :href="fullscreenImageUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="text-sm text-white underline hover:text-blue-300"
        >
          Открыть оригинал
        </a>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { ReplyIcon, EditIcon, PinIcon, TrashIcon, Users } from 'lucide-vue-next'
import { useAuthStore } from '../../../../auth/store/authStore'
import { useChatStore } from '../../../store/chatStore'
import { loadAttachmentUrls } from '../../../utils/attachmentUtils'
import MessageGallery from './messageGallery.vue'

const chatStore = useChatStore()
const authStore = useAuthStore()

const emit = defineEmits(['reply-to-message', 'edit-message', 'scroll-to-message'])
const props = defineProps({
  message: {
    type: Object,
    required: true
  }
})

const user = computed(() => chatStore.users[String(props.message.user_id)] || {})
const avatarUrl = computed(() => user.value.avatar_url)

const attachmentUrls = ref({})
const fullscreenImageUrl = ref(null)

function openImage(url) {
  fullscreenImageUrl.value = url
}

function closeImage() {
  fullscreenImageUrl.value = null
}

function onEsc(event) {
  if (event.key === 'Escape') {
    closeImage()
  }
}

onMounted(() => {
  window.addEventListener('keydown', onEsc)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onEsc)
})

const isMyMessage = props.message.user_id == authStore.getUserId
const onlyImage = computed(() => {
  return (
    props.message.text.trim() === '' &&
    Array.isArray(props.message.attachments) &&
    props.message.attachments.some(att => typeof att.key === 'string' && att.key.endsWith('.png'))
  )
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
  if (!replied || typeof replied.text !== 'string') return '[сообщение удалено]'

  if (replied.text != '') return replied.text
  if (onlyImage) return 'Изображение'

  const isEdited = replied.edited_at != null
  if (isEdited && replied.text.trim() !== props.message.reply_to_text?.trim()) {
    return replied.text + ' (изменено)'
  }

  return replied.text
})

function reply() {
  emit('reply-to-message', props.message)
}

function edit() {
  emit('edit-message', props.message)
}

function pin() {
  const shouldPin = !props.message.pinned
  chatStore.pinMessage(props.message, shouldPin)
}

function remove() {
  chatStore.deleteMessage(props.message)
}

watch(() => props.message.attachments, () => {
  loadAttachmentUrls([props.message], chatStore.imageUrlCache, attachmentUrls)
}, { immediate: true })
</script>
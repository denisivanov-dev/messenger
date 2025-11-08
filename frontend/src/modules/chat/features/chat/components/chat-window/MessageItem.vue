<template>
  <!-- Системное сообщение -->
  <div
    v-if="isSystemMessage"
    :id="`msg-${props.message.message_id}`"
    class="relative flex items-start gap-3 px-4 py-2 bg-white rounded-xl shadow hover:bg-gray-200 transition my-2"
  >
    <!-- Иконка звонка -->
    <div class="mt-0.5 w-8 h-8 flex items-center justify-center rounded-full bg-white border border-gray-300">
      <PhoneIcon
        v-if="props.message.call_info.status === 'ongoing'"
        class="w-4 h-4 text-green-600"
      />
      <PhoneOffIcon
        v-else-if="props.message.call_info.status === 'ended'"
        class="w-4 h-4 text-gray-400"
      />
      <PhoneMissedIcon
        v-else-if="props.message.call_info.status === 'missed'"
        class="w-4 h-4 text-red-500"
      />
      <XIcon
        v-else-if="props.message.call_info.status === 'cancelled'"
        class="w-4 h-4 text-gray-400"
      />
    </div>

    <!-- Контент -->
    <div class="flex-1 relative">
      <!-- Дата -->
      <div class="absolute top-0 right-0 mt-1 mr-2 text-[11px] text-gray-400">
        {{ formattedDate }}
      </div>

      <!-- Заголовок + стрелка -->
      <div class="mb-1 text-sm text-gray-800 font-medium flex items-center gap-1">
        <span>
          {{
            props.message.call_info.status === 'ongoing'
              ? 'Звонок начался'
              : props.message.call_info.status === 'ended'
                ? 'Звонок завершён'
                : props.message.call_info.status === 'missed'
                  ? 'Пропущенный звонок'
                  : 'Звонок отменён'
          }}
        </span>

        <!-- Стрелка + попап -->
        <div class="relative inline-block">
          <button
            @click="toggleParticipants"
            class="text-gray-500 hover:text-gray-800"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 transition-transform"
                :class="{ 'rotate-180': showParticipants }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>

          <!-- Попап участников (над стрелкой) -->
          <transition name="fade">
            <div
              v-if="showParticipants"
              class="absolute bottom-full mb-2 left-0 bg-white border border-gray-200 shadow-lg rounded-lg w-64 max-h-40 overflow-y-auto z-50"
            >
              <!-- Заголовок -->
              <div class="px-2 py-1 border-b text-xs font-semibold text-gray-600 bg-gray-50">
                Список участников
              </div>

              <!-- Участники -->
              <div class="p-2">
                <div
                  v-for="uid in props.message.call_info.participants"
                  :key="uid"
                  class="flex items-center gap-2 p-1 hover:bg-gray-50 rounded"
                >
                  <img
                    :src="chatStore.users[uid]?.avatar_url || '/default-avatar.png'"
                    class="w-6 h-6 rounded-full object-cover"
                  />
                  <span class="text-xs text-gray-700">
                    {{ chatStore.users[uid]?.username || 'Неизвестно' }}
                  </span>
                </div>
              </div>
            </div>
          </transition>
        </div>
      </div>

      <!-- Длительность звонка -->
      <p class="text-xs text-gray-600">
        {{
          props.message.call_info.status === 'ongoing'
            ? 'В звонке: ' + liveDuration
            : props.message.call_info.duration
              ? 'Длительность: ' + formatDuration(props.message.call_info.duration)
              : ''
        }}
      </p>
    </div>
  </div>

  <!-- Обычное сообщение -->
  <div
    v-else
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
import {
  ref,
  computed,
  onMounted,
  onBeforeUnmount,
  watch
} from 'vue'
import {
  ReplyIcon,
  EditIcon,
  PinIcon,
  TrashIcon,
  Users,
  PhoneIcon,
  PhoneOffIcon,
  PhoneMissedIcon,
  XIcon,
  PhoneCallIcon
} from 'lucide-vue-next'

import { useAuthStore } from '../../../../../auth/store/authStore'
import { useChatStore } from '../../../store/chatStore'
import { useCallStore } from '../../../store/call/callStore'
import { loadAttachmentUrls } from '../../utils/attachmentUtils'
import MessageGallery from './messageGallery.vue'

const chatStore = useChatStore()
const authStore = useAuthStore()
const callStore = useCallStore()

const emit = defineEmits(['reply-to-message', 'edit-message', 'scroll-to-message'])

const props = defineProps({
  message: {
    type: Object,
    required: true
  }
})

const showParticipants = ref(false)

function toggleParticipants() {
  showParticipants.value = !showParticipants.value
}

// ========== SYSTEM MESSAGE ==========
const isSystemMessage = computed(() =>
  props.message.user_id === '0' &&
  props.message.username === 'system' &&
  props.message.type.startsWith('call_')
)

const isParticipant = computed(() =>
  props.message.call_info?.participants.includes(String(authStore.getUserId))
)

// ========== LIVE DURATION ==========
const liveDuration = ref('')
let interval = null

function updateLiveDuration() {
  if (!props.message.call_info?.started_at) return
  const now = Date.now()
  const startedAt = props.message.call_info.started_at * 1000
  const diff = now - startedAt
  const totalSeconds = Math.floor(diff / 1000)
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60
  liveDuration.value = `${minutes}м ${seconds < 10 ? '0' : ''}${seconds}с`
}

// ========== JOIN CALL ==========
function joinCall() {
  callStore.joinCall()
}

// ========== USER INFO ==========
const user = computed(() => chatStore.users[String(props.message.user_id)] || {})
const avatarUrl = computed(() => user.value.avatar_url)
const isMyMessage = props.message.user_id == authStore.getUserId

// ========== ATTACHMENTS ==========
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

watch(() => props.message.attachments, () => {
  loadAttachmentUrls([props.message], chatStore.imageUrlCache, attachmentUrls)
}, { immediate: true })

// ========== FORMATTING ==========
function formatDuration(seconds) {
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  return `${m}м ${s < 10 ? '0' : ''}${s}с`
}

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

// ========== ACTIONS ==========
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

// ========== MOUNT/UNMOUNT ==========
onMounted(() => {
  window.addEventListener('keydown', onEsc)

  if (props.message.call_info?.status === 'ongoing') {
    updateLiveDuration()
    interval = setInterval(updateLiveDuration, 1000)
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onEsc)
  if (interval) clearInterval(interval)
})
</script>
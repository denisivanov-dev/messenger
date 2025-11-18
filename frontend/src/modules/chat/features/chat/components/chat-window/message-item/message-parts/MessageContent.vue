<template>
  <div class="flex-1 relative text-[#E4E6EB]">

    <!-- Имя + дата + pinned -->
    <div class="flex items-center gap-2 mb-1">
      <span class="font-semibold text-[#F2F3F5]">{{ message.payload.username }}</span>

      <span class="text-xs text-[#9CA3AF]">{{ formattedDate }}</span>

      <span
        v-if="message.payload.pinned"
        class="flex items-center gap-1 ml-2 text-xs text-[#C084FC]"
      >
        📌 <span class="italic">Закреплено</span>
      </span>
    </div>

    <!-- Reply -->
    <div
      v-if="message.payload.reply_to"
      class="text-[12px] text-[#9CA3AF] border-l-2 border-[#5865F2]
             pl-2 mb-1 cursor-pointer hover:text-[#A5B4FC]"
      @click="$emit('scroll-to-message', message.payload.reply_to)"
    >
      ↩ {{ repliedMessageUser }}:
      <span class="italic opacity-75">{{ repliedMessageText }}</span>
    </div>

    <!-- Text -->
    <p
      v-if="message.payload.text"
      class="text-sm leading-relaxed whitespace-pre-wrap"
    >
      {{ message.payload.text }}
    </p>

    <!-- Attachments -->
    <MessageImageGallery
      v-if="message.payload.attachments?.length"
      :attachments="message.payload.attachments"
      :imageUrls="attachmentUrls"
      :openImage="openImage"
      class="mt-2"
    />

    <!-- Edited -->
    <div
      v-if="message.payload.is_edited"
      class="text-[10px] text-[#9CA3AF] italic mt-1"
    >
      изменено • {{ formattedEditDate }}
    </div>

  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useChatStore } from '../../../../store/chatStore'
import MessageImageGallery from './MessageImageGallery.vue'

const props = defineProps({
  message: Object,
  attachmentUrls: Object
})

const emit = defineEmits(['scroll-to-message', 'open-image'])

const chatStore = useChatStore()

// ===== FORMATTING =====

const formattedDate = computed(() =>
  new Date(props.message.timestamp).toLocaleString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
)

const formattedEditDate = computed(() => {
  if (!props.message.payload?.edited_at) return ''
  return new Date(props.message.payload.payload.edited_at).toLocaleString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
})

// ===== REPLY TEXT =====

const repliedMessage = computed(() =>
  chatStore.messages.find(
    m => m.message_id === props.message.payload.reply_to
  )
)

const repliedMessageUser = computed(() =>
  props.message.payload.reply_to_user || 'неизвестный'
)

const repliedMessageText = computed(() => {
  const m = repliedMessage.value
  if (!m) return '[сообщение удалено]'
  if (m.text) return m.text

  const hasImage = Array.isArray(m.attachments) &&
                    m.attachments.some(a => a.key.endsWith('.png'))

  return hasImage ? 'Изображение' : '[пусто]'
})

// ===== OPEN IMAGE =====

function openImage(url) {
  emit('open-image', url)
}

</script>
<template>
  <div
    class="bubble"
    :class="bubbleClasses"
    @mouseenter="hover = true"
    @mouseleave="hover = false"
    style="position: relative;"
  >

    <!-- USERNAME + DATE (only for first in group) -->
    <div
      v-if="isFirst"
      class="flex items-center gap-2 mb-1"
      :style="headerStyle"
    >
      <span class="font-semibold text-[#F2F3F5]" :style="usernameStyle">
        {{ message.payload.username }}
      </span>

      <span class="text-xs text-[#9CA3AF]" :style="dateStyle">
        {{ formattedDate }}
      </span>

      <span
        v-if="message.payload.pinned"
        class="flex items-center gap-1 ml-2 text-xs text-[#C084FC]"
        :style="pinStyle"
      >
        📌 <span class="italic">Pinned</span>
      </span>
    </div>

    <!-- REPLY PREVIEW -->
    <div
      v-if="isFirst && message.payload.reply_to"
      class="text-[12px] text-[#9CA3AF] border-l-2 border-[#5865F2] pl-2 mb-1 cursor-pointer hover:text-[#A5B4FC]"
      :style="replyStyle"
      @click="$emit('scroll-to-message', message.payload.reply_to)"
    >
      ↩ {{ repliedMessageUser }}:
      <span class="italic opacity-75">{{ repliedMessageText }}</span>
    </div>

    <!-- TEXT -->
    <p
      v-if="message.payload.text"
      class="text-sm whitespace-pre-wrap"
      :style="textStyle"
      style="margin-right: 40px;"
    >
      {{ message.payload.text }}
    </p>

    <!-- INLINE HOVER TIME -->
    <span
      v-if="bubbleClasses !== 'bubble--first' && bubbleClasses !== 'bubble--solo'"
      :style="{
        position: 'absolute',
        right: '10px',
        bottom: '6px',
        opacity: hover ? 1 : 0,
        transition: 'opacity .15s',
        fontSize: '10px',
        color: '#9CA3AF',
        userSelect: 'none',
        pointerEvents: 'none'
      }"
    >
      {{ formattedShortTime }}
    </span>

    <!-- ATTACHMENTS -->
    <MessageImageGallery
      v-if="message.payload.attachments?.length"
      :attachments="message.payload.attachments"
      :imageUrls="attachmentUrls"
      :openImage="openImage"
      class="mt-2"
      :zoom="zoom"
    />

    <!-- EDITED LABEL -->
    <div
      v-if="isFirst && message.payload.is_edited"
      class="text-[10px] text-[#9CA3AF] italic mt-1"
      :style="editedStyle"
    >
      edited • {{ formattedEditDate }}
    </div>

  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useChatStore } from '../../../../store/chatStore'
import MessageImageGallery from './MessageImageGallery.vue'

const props = defineProps({
  message: Object,
  attachmentUrls: Object,
  isFirst: Boolean
})

const emit = defineEmits(['scroll-to-message', 'open-image'])

const chatStore = useChatStore()
const hover = ref(false)

/* ===== ZOOM ===== */
const zoom = inject("zoom")

/* ===== TEXT SCALE ===== */
const textStyle = computed(() => ({
  fontSize: `${14 * zoom.value}px`,
  lineHeight: `${16 * zoom.value}px`,
  marginTop: `${2 * zoom.value}px`,
  marginBottom: `${2 * zoom.value}px`
}))

/* ===== USERNAME ===== */
const usernameStyle = computed(() => ({
  fontSize: `${15 * zoom.value}px`,
}))

const dateStyle = computed(() => ({
  fontSize: `${12 * zoom.value}px`,
}))

/* ===== DYNAMIC TIME POSITION ===== */
const timeStyle = computed(() => ({
  fontSize: `${10 * zoom.value}px`,
  left: `${-40 * zoom.value}px`,
  top: `${(16 * zoom.value) / 4}px`,
}))

/* ===== PIN ===== */
const pinStyle = computed(() => ({
  fontSize: `${12 * zoom.value}px`,
}))

/* ===== HEADER ===== */
const headerStyle = computed(() => ({
  gap: `${8 * zoom.value}px`
}))

/* ===== REPLY ===== */
const replyStyle = computed(() => ({
  fontSize: `${12 * zoom.value}px`,
  borderLeftWidth: `${2 * zoom.value}px`,
  paddingLeft: `${6 * zoom.value}px`
}))

/* ===== EDITED ===== */
const editedStyle = computed(() => ({
  fontSize: `${10 * zoom.value}px`
}))

/* ===== FORMATTED DATES ===== */

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

const formattedShortTime = computed(() =>
  new Date(props.message.timestamp).toLocaleString('ru-RU', {
    hour: '2-digit',
    minute: '2-digit'
  })
)

/* ===== REPLY DATA ===== */
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

  const hasImage = Array.isArray(m.attachments)
    && m.attachments.some(a => a.key.endsWith('.png'))

  return hasImage ? 'Изображение' : '[пусто]'
})

/* ===== OPEN IMAGE ===== */
function openImage(url) {
  emit('open-image', url)
}

const bubbleClasses = computed(() => {
  const m = props.message

  if (m.groupIndex === 0 && m.isLastFromSender) {
    return "bubble--solo"
  }

  if (m.groupIndex === 0) {
    return "bubble--first"
  }

  if (m.isLastFromSender) {
    return "bubble--last"
  }

  return "bubble--middle"
})
</script>
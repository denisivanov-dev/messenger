<template>
  <div
    :id="`msg-${message.payload.message_id}`"
    class="relative group flex items-start gap-3 px-4 py-3
           bg-[#2B2D31] hover:bg-[#323439] transition-colors
           rounded-xl shadow-sm border border-[#383A3F]"
  >
    <!-- AVATAR -->
    <img
      :src="avatarUrl"
      class="w-8 h-8 rounded-full object-cover mt-0.5 ring-1 ring-[#3A3B3F]"
      alt="avatar"
    />

    <!-- CONTENT -->
    <MessageContent
      :message="message"
      :attachmentUrls="attachmentUrls"
      @open-image="openImage"
      @scroll-to-message="$emit('scroll-to-message', $event)"
    />

    <!-- ACTIONS -->
    <MessageActions
      class="absolute top-0 right-0 mt-1 mr-1"
      :isMy="isMyMessage"
      @reply="reply"
      @edit="edit"
      @pin="pin"
      @delete="remove"
    />
  </div>

  <!-- FULLSCREEN IMAGE -->
  <Teleport to="body">
    <div
      v-if="fullscreenImageUrl"
      class="fixed inset-0 bg-black bg-opacity-90 z-50 flex items-center justify-center"
      @click="closeImage"
    >
      <div @click.stop class="flex flex-col items-center gap-3">
        <img :src="fullscreenImageUrl" class="max-w-full max-h-[90vh] shadow-xl" />
        <a
          :href="fullscreenImageUrl"
          target="_blank"
          class="text-sm text-white underline hover:text-blue-300"
        >
          Открыть оригинал
        </a>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useAuthStore } from '../../../../../../../auth/store/authStore'
import { useChatStore } from '../../../../store/chatStore'

import { loadAttachmentUrls } from '../../../../utils/attachmentUtils'

import MessageContent from '../message-parts/MessageContent.vue'
import MessageActions from '../message-parts/MessageActions.vue'

const props = defineProps({
  message: { type: Object, required: true }
})

const emit = defineEmits([
  'reply-to-message',
  'edit-message',
  'scroll-to-message'
])

const authStore = useAuthStore()
const chatStore = useChatStore()

/* USER / AVATAR */
const avatarUrl = computed(() => {
  const u = chatStore.users[String(props.message.user_id)]
  return u?.avatar_url || '/default-avatar.png'
})

const isMyMessage = computed(() =>
  props.message.sender_id == authStore.getUserId
)

/* ATTACHMENTS */
const attachmentUrls = ref({})
watch(() => props.message.payload?.attachments, () => {
  loadAttachmentUrls([props.message.payload], chatStore.imageUrlCache, attachmentUrls)
}, { immediate: true })

/* IMAGE MODAL */
const fullscreenImageUrl = ref(null)
function openImage(url) { fullscreenImageUrl.value = url }
function closeImage() { fullscreenImageUrl.value = null }

/* ACTIONS */
function reply() { emit('reply-to-message', props.message) }
function edit() { emit('edit-message', props.message) }
function pin() { chatStore.pinMessage(props.message, !props.message.pinned) }
function remove() { chatStore.deleteMessage(props.message) }
</script>
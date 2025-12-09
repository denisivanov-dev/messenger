<template>
  <div
    :id="`msg-${message.payload.message_id}`"
    class="block w-full message-user relative"
    :class="{
      my: isMyMessage,
      'group-end': message.isLastFromSender
    }"
  >

    <div class="flex items-end gap-3 px-1 w-full group">

      <!-- Аватар -->
      <template v-if="message.isLastFromSender">
        <img
          :src="avatarUrl"
          :style="avatarStyle"
          class="rounded-full object-cover mt-0.5 ring-1 ring-[#3A3B3F]"
        />
      </template>

      <template v-else>
        <div :style="avatarSpacerStyle"></div>
      </template>

      <!-- САМ ПУЗЫРЬ — БЕЗ ВРАПЕРОВ -->
      <MessageContent
        ref="bubble"
        class="msg-bubble"
        :message="message"
        :attachmentUrls="attachmentUrls"
        :isFirst="isFirst"
        :zoom="zoom"
        @open-image="openImage"
        @scroll-to-message="$emit('scroll-to-message', $event)"
      />

      <MessageActions
        class="msg-actions absolute opacity-0 transition -right-2 mt-[-4px]"
        :isMy="isMyMessage"
        @reply="reply"
        @edit="edit"
        @pin="pin"
        @delete="remove"
      />
    </div>

    <!-- fullscreen image -->
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
  </div>
</template>

<script setup>
import { ref, computed, watch, inject } from 'vue'
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

/* ZOOM */
const zoom = inject("zoom")

/* FIRST IN GROUP */
const isFirst = computed(() => props.message.groupIndex === 0)

/* USER / AVATAR */
const avatarUrl = computed(() => {
  const u = chatStore.users[String(props.message.user_id)]
  return u?.avatar_url || '/default-avatar.png'
})

const isMyMessage = computed(() =>
  props.message.sender_id == authStore.getUserId
)

/* ZOOM STYLES */
const avatarStyle = computed(() => ({
  width: `${32 * zoom.value}px`,
  height: `${32 * zoom.value}px`,
}))

const avatarSpacerStyle = computed(() => ({
  width: `${32 * zoom.value}px`,
  height: `${32 * zoom.value}px`,
}))

/* ATTACHMENTS */
const attachmentUrls = ref({})
watch(
  () => props.message.payload?.attachments,
  () => {
    loadAttachmentUrls(
      [props.message.payload],
      chatStore.imageUrlCache,
      attachmentUrls
    )
  },
  { immediate: true }
)

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
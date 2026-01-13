<template>
  <div
    :id="`msg-${message.payload.message_id}`"
    class="block w-full message-user relative"
    :class="{
      my: isMyMessage,
      'group-end': message.isLastFromSender
    }"
  > 
    <div class="flex items-end gap-3 px-1 w-full">

      <!-- user profile pic -->
      <template v-if="message.isLastFromSender">
        <img
          :src="avatarUrl"
          :style="avatarStyle"
          class="rounded-full object-cover mt-0.5 ring-1 ring-[#3A3B3F] select-none"
        />
      </template>

      <template v-else>
        <div :style="avatarSpacerStyle"></div>
      </template>

      <!-- message wrapper for content and actions -->
      <div
        class="relative flex items-start w-fit"
        @mouseenter="onHoverStart"
        @mouseleave="onHoverEnd"
        @contextmenu.prevent="openClickMenu"
      >
        <MessageContent
          ref="bubble"
          class="msg-bubble"
          :message="message"
          :attachmentUrls="attachmentUrls"
          :isFirstMessage="isFirstMessage"
          :zoom="zoom"
          @open-image="openImage"
          @scroll-to-message="$emit('scroll-to-message', $event)"
        />

        <MessageHoverActions
          class="message-user-actions absolute top-[-14px] right-4 translate-x-full transition-all"
          :class="hoverMenu
            ? 'pointer-events-auto'
            : 'pointer-events-none'"
          :isMy="isMyMessage"
          @reply="reply"
          @edit="edit"
          @pin="pin"
          @delete="remove"
        />
      </div>
    </div>

    <!-- actions when user clicks on message -->
    <Teleport to="body">
      <Transition name="click-menu">
        <MessageClickActions
          v-if="clickMenu.open"
          :x="clickMenu.x"
          :y="clickMenu.y"
          :openUp="clickMenu.openUp"
          :isMy="isMyMessage"
          @reply="() => { reply(); closeClickMenu() }"
          @edit="() => { edit(); closeClickMenu() }"
          @pin="() => { pin(); closeClickMenu() }"
          @copy="() => { copyMessage(message.payload.text); closeClickMenu() }"
          @forward="() => { forward(); closeClickMenu() }"
          @delete="() => { remove(); closeClickMenu() }"
          @select="() => { selectMessage(); closeClickMenu() }"
        />
      </Transition>
    </Teleport>

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
import { useClipboard } from '@vueuse/core'
import { useAuthStore } from '../../../../../../../auth/store/authStore'
import { useChatStore } from '../../../../store/chatStore'
import { loadAttachmentUrls } from '../../../../utils/attachmentUtils'

import MessageContent from '../message-parts/MessageContent.vue'
import MessageHoverActions from '../message-parts/MessageHoverActions.vue'
import MessageClickActions from '../message-parts/MessageClickActions.vue'

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

const zoom = inject('zoom')

const isFirstMessage = computed(() => props.message.groupIndex === 0)

const avatarUrl = computed(() => {
  const u = chatStore.users[String(props.message.user_id)]
  return u?.avatar_url || '/default-avatar.png'
})

const isMyMessage = computed(() =>
  props.message.sender_id == authStore.getUserId
)

/* avatar scaling */
const avatarStyle = computed(() => ({
  width: `${32 * zoom.value}px`,
  height: `${32 * zoom.value}px`,
}))

const avatarSpacerStyle = computed(() => ({
  width: `${32 * zoom.value}px`,
  height: `${32 * zoom.value}px`,
}))

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

/* fullscreen image */
const fullscreenImageUrl = ref(null)

function openImage(url) {
  fullscreenImageUrl.value = url
}

function closeImage() {
  fullscreenImageUrl.value = null
}

/* message actions */
function reply() {
  emit('reply-to-message', props.message)
}

function edit() {
  emit('edit-message', props.message)
}

function pin() {
  chatStore.pinMessage(props.message, !props.message.pinned)
}

function remove() {
  chatStore.deleteMessage(props.message)
}

/* hover actions */
const hoverMenu = ref(false)
let hoverTimer = null
let hideTimer = null

const isAnyClickMenuOpen = inject('isAnyClickMenuOpen', ref(false))

function onHoverStart() {
  if (isAnyClickMenuOpen.value) return
  clearTimeout(hideTimer)

  hoverTimer = setTimeout(() => {
    hoverMenu.value = true
  }, 300)
}

function onHoverEnd() {
  if (isAnyClickMenuOpen.value) return
  clearTimeout(hoverTimer)

  hideTimer = setTimeout(() => {
    hoverMenu.value = false
  }, 300)
}

/* click menu */
const clickMenu = ref({ open: false, x: 0, y: 0 })
const registerClickMenu = inject('registerClickMenu')

function openClickMenu(e) {
  registerClickMenu?.(closeClickMenu)

  clearTimeout(hoverTimer)
  clearTimeout(hideTimer)
  hoverMenu.value = false

  const MENU_HEIGHT = 290
  const OFFSET = 10

  const spaceBelow = window.innerHeight - e.clientY
  const openUp = spaceBelow < MENU_HEIGHT + OFFSET

  clickMenu.value = {
    open: true,
    x: e.clientX + 10,
    y: e.clientY + OFFSET,
    openUp
  }
}

function closeClickMenu() {
  clickMenu.value.open = false
}

/* clipboard */
const { copy, isSupported } = useClipboard()

const copyMessage = (text) => {
  if (!isSupported.value) return
  copy(text)
  closeClickMenu()
}
</script>
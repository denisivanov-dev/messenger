<template>
  <div
    ref="chatWindowRef"
    class="flex-1 min-h-0 w-full max-w-[1300px]
          p-6 overflow-y-auto flex flex-col gap-4 relative
          chat-window"
  >
    <MessageItem
      v-for="msg in messages"
      :key="msg.message_id"
      :message="msg"
      @edit-message="emit('edit-message', $event)"
      @reply-to-message="emit('reply-to-message', $event)"
      @scroll-to-message="scrollToMessage"
    />

  <div
    v-if="!messages.length"
    class="flex-1 flex flex-col items-center justify-center text-[#9CA3AF] gap-6 py-8 select-none"
  >
    <div class="text-xl font-medium">Тут пока пусто...</div>

    <img
      src="/images/pythonus/lonely_pythonus.png"
      alt="empty"
      class="w-80 h-auto opacity-80 pointer-events-none select-none"
      draggable="false"
    />
  </div>

    <div
      v-if="typingUser"
      class="text-sm italic text-[#A1A1AA] px-2 py-1"
    >
      {{ typingUser }} печатает…
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted  } from 'vue'
import { useChatStore } from '../../store/chatStore'
import MessageItem from './message-item/MessageItem.vue'
import { waitForImagesAndThen } from '../../utils/chatWindowUtils'

const chatWindowRef = ref(null)
const chatStore = useChatStore()

const messages = computed(() => chatStore.messages)
const typingUser = computed(() => chatStore.typingUser)
const emit = defineEmits(['reply-to-message', 'edit-message', 'scroll-to-message'])

watch(() => chatStore.shouldScroll, async (val) => {
  if (!val) return
  await nextTick()
  waitForImagesAndThen(() => {
    if (chatWindowRef.value) {
      chatWindowRef.value.scrollTop = chatWindowRef.value.scrollHeight
    }
    chatStore.shouldScroll = false
  }, chatWindowRef.value)
})

watch(
  typingUser,
  async (val) => {
    if (val) {
      await nextTick()
      chatWindowRef.value.scrollTop = chatWindowRef.value.scrollHeight
    }
  }
)

function scrollToMessage(messageID) {
  const el = document.getElementById(`msg-${messageID}`)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'center' })
    el.classList.add('ring', 'ring-blue-400', 'transition')
    setTimeout(() => {
      el.classList.remove('ring', 'ring-blue-400')
    }, 1000)
  } else {
    console.warn(`Сообщение msg-${messageID} не найдено`)
  }
}

onMounted(async () => {
  await nextTick()
  if (chatWindowRef.value) {
    chatWindowRef.value.scrollTop = chatWindowRef.value.scrollHeight
  }
})
</script>
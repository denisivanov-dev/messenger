<template>
  <div
    ref="chatWindowRef"
    class="chat-window flex-1 min-h-0 w-full max-w-[1300px]
           p-5 bg-[#232428] text-[#E4E6EB]
           rounded-2xl shadow-[0_0_10px_rgba(0,0,0,0.3)]
           overflow-y-auto flex flex-col gap-3
           border border-[#2E2F33]"
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
      class="flex-1 flex items-center justify-center text-[#9CA3AF] italic"
    >
      Пока тут пусто 👀
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
import MessageItem from './MessageItem.vue'
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
<template>
  <div
    ref="chatWindowRef"
    class="chat-window w-full max-w-[1300px] h-[780px] p-5 bg-gray-100 rounded-2xl shadow-md overflow-y-scroll flex flex-col gap-2"
  >
    <!-- сообщения -->
    <MessageItem
      v-for="msg in messages"
      :key="msg.timestamp"
      :message="msg"
      @edit-message="emit('edit-message', $event)"
    />

    <!-- индикатор печати -->
    <div
      v-if="typingUser"
      class="text-sm italic text-gray-500 px-2 py-1"
    >
      {{ typingUser }} печатает…
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted  } from 'vue'
import { useChatStore } from '../../../store/chatStore'
import MessageItem from './messageItem.vue'

const chatWindowRef = ref(null)
const chatStore = useChatStore()

const messages = computed(() => chatStore.messages)
const typingUser = computed(() => chatStore.typingUser)

const emit = defineEmits(['edit-message'])

watch(() => chatStore.shouldScroll, async (val) => {
  if (val) {
    await nextTick()
    setTimeout(() => {
      if (chatWindowRef.value) {
        chatWindowRef.value.scrollTop = chatWindowRef.value.scrollHeight
      }
      chatStore.shouldScroll = false
    }, 20)  
  }
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


onMounted(async () => {
  await nextTick()
  if (chatWindowRef.value) {
    chatWindowRef.value.scrollTop = chatWindowRef.value.scrollHeight
  }
})
</script>
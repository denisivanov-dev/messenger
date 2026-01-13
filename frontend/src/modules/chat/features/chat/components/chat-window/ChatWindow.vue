<template>
  <div
    id="chat-window"
    ref="chatWindowRef"
    class="flex-1 overflow-y-auto flex flex-col relative py-2 min-h-0 px-2 sm:px-3"
    @scroll="onUserScroll"
  >
    <MessageItem
      v-for="msg in messagesWithGroup"
      :key="msg.message_id"
      :message="msg"
      @edit-message="emit('edit-message', $event)"
      @reply-to-message="emit('reply-to-message', $event)"
      @scroll-to-message="scrollToMessage"
    />

    <!-- scroll to bottom button -->
    <button
      v-if="showScrollToBottom"
      @click="scrollToBottom"
      class="sticky bottom-3 mr-2 ml-auto z-[50]
            bg-[#2A2B2F]
            text-white
            p-2 rounded-full
            shadow-lg
            hover:bg-[#3A3B3E]
            transition"
    >
      <ChevronDown class="w-5 h-5" />
    </button>
  </div>
</template>


<script setup>
import { ref, computed, watch, nextTick, onMounted, provide } from 'vue'
import { useChatStore } from '../../store/chatStore'
import { useMessagesStore } from '../../store/messagesStore'
import MessageItem from './message-item/MessageItem.vue'
import { waitForImagesAndThen } from '../../utils/chatWindowUtils'
import { isNearBottom } from '../../utils/chatWindowUtils'
import { ChevronDown } from 'lucide-vue-next'

const chatWindowRef = ref(null)
const chatStore = useChatStore()
const messagesStore = useMessagesStore()


const emit = defineEmits([
  'reply-to-message',
  'edit-message',
  'scroll-to-message'
])

const messages = computed(() => chatStore.messages)

const showScrollToBottom = ref(false)

/* click-menu coordination */
const isAnyClickMenuOpen = ref(false)
const activeClickMenuCloser = ref(null)

provide('registerClickMenu', (closeFn) => {
  if (activeClickMenuCloser.value) {
    activeClickMenuCloser.value()
  }
  
  activeClickMenuCloser.value = () => {
    closeFn()
    isAnyClickMenuOpen.value = false
  }

  window.__closeActiveClickMenu = activeClickMenuCloser.value
  isAnyClickMenuOpen.value = true
})

provide('isAnyClickMenuOpen', isAnyClickMenuOpen)

provide('closeActiveClickMenu', () => {
  if (activeClickMenuCloser.value) {
    activeClickMenuCloser.value()
    activeClickMenuCloser.value = null
    isAnyClickMenuOpen.value = false
  }
})

/* user interactions */
function onUserScroll() {
  // close click menu on scroll
  if (activeClickMenuCloser.value) {
    activeClickMenuCloser.value()
    activeClickMenuCloser.value = null
    isAnyClickMenuOpen.value = false
  }

  // toggle scroll-to-bottom button
  showScrollToBottom.value = !isNearBottom(500)
}

/* initial scroll */
onMounted(async () => {
  await nextTick()
  console.info("scrolling")
  if (chatWindowRef.value) {
    chatWindowRef.value.scrollTop = chatWindowRef.value.scrollHeight
  }
})

/* auto-scroll logic */
watch(
  () => messagesStore.shouldScroll,
  async (val) => {
    if (!val) return
    await nextTick()
    console.info("auto-scrolling")

    waitForImagesAndThen(() => {
      if (chatWindowRef.value) {
        chatWindowRef.value.scrollTop = chatWindowRef.value.scrollHeight
        showScrollToBottom.value = false
      }
    }, chatWindowRef.value)

    chatWindowRef.value.scrollTop = chatWindowRef.value.scrollHeight
    messagesStore.shouldScroll = false
  }
)

/* scrolling functions */
function scrollToMessage(messageID) {
  const el = document.getElementById(`msg-${messageID}`)
  if (!el) return

  el.scrollIntoView({ behavior: 'smooth', block: 'center' })
  el.classList.add('ring', 'ring-blue-400', 'transition')

  setTimeout(() => {
    el.classList.remove('ring', 'ring-blue-400')
  }, 1000)
}

function scrollToBottom() {
  const el = chatWindowRef.value
  if (!el) return

  el.scrollTo({
    top: el.scrollHeight,
    behavior: 'smooth'
  })

  showScrollToBottom.value = false
}

/* message grouping */
const rawMessages = computed(() => chatStore.messages)

const messagesWithGroup = computed(() => {
  let lastSender = null
  let groupIndex = 0

  return rawMessages.value.map((m, index) => {
    if (m.sender_id === lastSender) {
      groupIndex++
    } else {
      lastSender = m.sender_id
      groupIndex = 0
    }

    const nextMsg = rawMessages.value[index + 1]
    const isLast = !nextMsg || nextMsg.sender_id !== m.sender_id

    return {
      ...m,
      groupIndex,
      isLastFromSender: isLast
    }
  })
})
</script>
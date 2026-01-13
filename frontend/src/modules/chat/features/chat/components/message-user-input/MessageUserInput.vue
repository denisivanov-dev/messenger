<template>
  <div class="relative w-full">
    <!-- input container -->
    <div
      class="relative flex items-center w-full h-[60px]
             chat-input chat-input-focus chat-input-drag"
      @dragover.prevent="onDragOver"
      @dragleave.prevent="onDragLeave"
      @drop.prevent="onDrop"
      :class="{ 'dark-panel-drag': isDragging }"
    >
      <!-- image previews -->
      <div
        v-if="previewUrls.length > 0"
        class="absolute bottom-full left-0 mb-2
               bg-white border border-gray-300
               shadow-xl rounded-xl p-2 z-50
               flex gap-2 flex-wrap max-w-[90%]"
      >
        <div
          v-for="(url, index) in previewUrls"
          :key="index"
          class="relative"
        >
          <img
            :src="url"
            class="max-h-32 max-w-[100px] rounded-md object-cover"
          />
          <button
            @click="removePreview(index)"
            class="absolute top-0 right-0
                   bg-white bg-opacity-80
                   hover:bg-opacity-100
                   rounded-full
                   text-gray-700 hover:text-red-500"
            title="Remove"
          >
            <XIcon class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- hidden file input -->
      <input
        ref="fileInput"
        type="file"
        accept="image/*"
        multiple
        @change="handleFileUpload"
        hidden
      />

      <!-- attach button -->
      <button
        @click="triggerFileSelect"
        class="mr-3 p-1 text-gray-600 hover:text-gray-500"
      >
        <PaperclipIcon class="w-5 h-5" />
      </button>

      <!-- text input -->
      <input
        ref="inputRef"
        v-model="text"
        type="text"
        spellcheck="false"
        placeholder="Enter a message…"
        class="flex-1 h-full ml-3
               border-none outline-none bg-transparent"
        @input="typingStore.sendTyping"
        @keyup.enter="send"
        @keyup.esc="onEsc"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { PaperclipIcon, XIcon } from 'lucide-vue-next'
import { useMessagesStore } from '../../store/messagesStore'
import { useChatModeStore } from '../../store/chatModeStore'
import { useTypingStore } from '../../store/events/typingStore'
import { uploadAllImages } from '../../utils/messageUtils'

const text = ref('')
const inputRef = ref(null)
const fileInput = ref(null)
const selectedFiles = ref([])
const previewUrls = ref([])
const isDragging = ref(false)

const messagesStore = useMessagesStore()
const chatModeStore = useChatModeStore()
const typingStore = useTypingStore()

async function send() {
  const trimmedText = text.value.trim()
  const hasText = trimmedText !== ''
  const hasImages = selectedFiles.value.length > 0

  if (!hasText && !hasImages) return

  let attachments = []
  if (hasImages) {
    attachments = await uploadAllImages(selectedFiles.value)
  }

  messagesStore.sendMessageData({
    text: trimmedText,
    attachments,
    replyToMessage: messagesStore.replyToMessage,
  }, chatModeStore.chatType, chatModeStore.receiverID)

  // chatStore.cancelReply()
  clearPreview()
  text.value = ''
}

/* file upload */
function handleFileUpload(event) {
  const files = Array.from(event.target.files || [])

  for (const file of files) {
    if (!file.type.startsWith('image/')) continue
    if (selectedFiles.value.length >= 5) break

    selectedFiles.value.push(file)
    previewUrls.value.push(URL.createObjectURL(file))
  }

  fileInput.value.value = null
}

function triggerFileSelect() {
  fileInput.value.value = null
  fileInput.value?.click()
}

function removePreview(index) {
  selectedFiles.value.splice(index, 1)
  previewUrls.value.splice(index, 1)
}

function clearPreview() {
  selectedFiles.value = []
  previewUrls.value = []
  fileInput.value.value = null
}

/* drag & drop */
function onDragOver() {
  isDragging.value = true
}

function onDragLeave() {
  isDragging.value = false
}

function onDrop(event) {
  isDragging.value = false
  const files = Array.from(event.dataTransfer.files || [])

  for (const file of files) {
    if (!file.type.startsWith('image/')) continue
    if (selectedFiles.value.length >= 5) break

    selectedFiles.value.push(file)
    previewUrls.value.push(URL.createObjectURL(file))
  }
}

/* keyboard handling */
function onEsc() {
  if (inputRef.value === document.activeElement) {
    inputRef.value.blur()
  }
}

function focusInputOnKeyPress(event) {
  if (event.ctrlKey || event.key === 'Escape') return

  const el = document.activeElement
  const tag = el?.tagName?.toLowerCase()
  const isTypingElement = ['input', 'textarea'].includes(tag)

  if (!isTypingElement && inputRef.value) {
    inputRef.value.focus()
  }
}

/* init */
onMounted(() => {
  window.addEventListener('keydown', focusInputOnKeyPress)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', focusInputOnKeyPress)
})
</script>
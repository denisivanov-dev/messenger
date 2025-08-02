<template>
  <div class="relative mt-6 w-full max-w-[1300px]">
    <!-- Всплывающее окно isEditing/isReplying -->
    <div
      v-if="isEditing || isReplying"
      class="absolute bottom-full left-0 mb-1 px-4 py-2 bg-white border border-gray-300 rounded-xl shadow z-40 text-sm text-gray-600 max-w-[80%]"
    >
      <template v-if="isEditing">
        Вы редактируете сообщение "{{ editingMessage.text }}"
        <button @click="cancelEdit" class="ml-2 text-blue-600 hover:underline">Отменить • Esc</button>
      </template>
      <template v-else-if="isReplying">
        Вы отвечаете на "{{ replyingMessage.text }}" от {{ replyingMessage.username }}
        <button @click="cancelReply" class="ml-2 text-blue-600 hover:underline">Отменить • Esc</button>
      </template>
    </div>

    <!-- Контейнер инпута -->
    <div
      class="relative flex items-center w-full h-[60px] rounded-2xl border border-gray-300 shadow-md bg-white px-4 focus-within:ring-2 focus-within:ring-blue-500 transition"
      @dragover.prevent="onDragOver"
      @dragleave.prevent="onDragLeave"
      @drop.prevent="onDrop"
      :class="{ 'ring-2 ring-blue-400 bg-blue-50': isDragging }"
    >
      <!-- Превью выбранных файлов -->
      <div
        v-if="previewUrls.length > 0"
        class="absolute bottom-full left-0 mb-2 bg-white border border-gray-300 shadow-xl rounded-xl p-2 z-50 flex gap-2 flex-wrap max-w-[90%]"
      >
        <div
          v-for="(url, index) in previewUrls"
          :key="index"
          class="relative"
        >
          <img :src="url" class="max-h-32 max-w-[100px] rounded-md object-cover" />
          <button
            @click="removePreview(index)"
            class="absolute top-0 right-0 bg-white bg-opacity-80 hover:bg-opacity-100 rounded-full text-gray-700 hover:text-red-500"
            title="Удалить"
          >
            <XIcon class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Текстовое поле -->
      <input
        ref="inputRef"
        v-model="text"
        @input="chatStore.sendTyping"
        @keyup.enter="send"
        @keyup.esc="onEsc"
        type="text"
        placeholder="Type a message..."
        class="flex-1 h-full border-none outline-none bg-transparent"
      />

      <!-- Выбор файлов -->
      <input
        type="file"
        ref="fileInput"
        accept="image/*"
        multiple
        @change="handleFileUpload"
        style="display: none"
      />
      <button @click="triggerFileSelect" class="mr-3 p-1 text-gray-600 hover:text-black">
        <PaperclipIcon class="w-5 h-5" />
      </button>
    </div>
  </div>
</template>

<script setup>
import { PaperclipIcon, XIcon } from 'lucide-vue-next'
import { ref, nextTick, onMounted, onBeforeUnmount } from 'vue'
import { useChatStore } from '../../../store/chatStore'
import { uploadAllImages } from '../../../utils/messageUtils'

const text = ref('')
const isEditing = ref(false)
const editingMessage = ref(null)
const isReplying = ref(false)
const replyingMessage = ref(null)

const inputRef = ref(null)
const fileInput = ref(null)
const selectedFiles = ref([])
const previewUrls = ref([])
const isDragging = ref(false)

const chatStore = useChatStore()

async function send() {
  const trimmedText = text.value.trim()
  const hasText = trimmedText !== ''
  const hasImages = selectedFiles.value.length > 0

  if (!hasText && !hasImages) return

  if (isEditing.value && editingMessage.value && hasText && !hasImages) {
    chatStore.editMessage(editingMessage.value, trimmedText)
    cancelEdit()
    text.value = ''
    return
  }

  let attachments = []
  if (hasImages) {
    attachments = await uploadAllImages(selectedFiles.value)
  }

  chatStore.sendMessageData({
    text: trimmedText,
    attachments,
    replyToMessage: replyingMessage.value,
  })

  cancelReply()
  text.value = ''
  clearPreview()
}

function handleFileUpload(event) {
  const newFiles = Array.from(event.target.files || [])

  for (const file of newFiles) {
    if (!file.type.startsWith('image/')) continue
    if (selectedFiles.value.length >= 5) break

    selectedFiles.value.push(file)
    previewUrls.value.push(URL.createObjectURL(file))
  }

  fileInput.value.value = null
}

function removePreview(index) {
  selectedFiles.value.splice(index, 1)
  previewUrls.value.splice(index, 1)
}

function triggerFileSelect() {
  fileInput.value.value = null
  fileInput.value?.click()
}

function onDrop(e) {
  isDragging.value = false
  const dropped = Array.from(e.dataTransfer.files)

  for (const file of dropped) {
    if (!file.type.startsWith('image/')) continue
    if (selectedFiles.value.length >= 5) break

    selectedFiles.value.push(file)
    previewUrls.value.push(URL.createObjectURL(file))
  }
}

function onDragOver() {
  isDragging.value = true
}

function onDragLeave() {
  isDragging.value = false
}

function clearPreview() {
  selectedFiles.value = []
  previewUrls.value = []
  fileInput.value.value = null
}

function startEdit(msg) {
  if (isReplying.value) cancelReply()
  text.value = msg.text
  editingMessage.value = msg
  isEditing.value = true
  nextTick(() => inputRef.value?.focus())
}

function cancelEdit() {
  text.value = ''
  editingMessage.value = null
  isEditing.value = false
}

function startReply(msg) {
  if (isEditing.value) cancelEdit()
  replyingMessage.value = msg
  isReplying.value = true
  nextTick(() => inputRef.value?.focus())
}

function cancelReply() {
  replyingMessage.value = null
  isReplying.value = false
}

function onEsc() {
  if (isEditing.value) cancelEdit()
  if (isReplying.value) cancelReply()
  if (inputRef.value === document.activeElement) {
    inputRef.value.blur()
  }
}

function focusInputOnKeyPress(event) {
  const tag = document.activeElement.tagName.toLowerCase()
  const isTypingElement = ['input', 'textarea'].includes(tag)
  if (!isTypingElement && inputRef.value) {
    inputRef.value.focus()
  }
}

onMounted(() => {
  window.addEventListener('keydown', focusInputOnKeyPress)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', focusInputOnKeyPress)
})

defineExpose({ startEdit, startReply })
</script>
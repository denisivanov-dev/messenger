<template>
  <div class="mt-6 w-full max-w-[1300px]">
    <div v-if="isEditing" class="text-sm text-gray-600 mb-1">
      Вы редактируете сообщение "{{ editingMessage.text}}"
      <button @click="cancelEdit" class="ml-2 text-blue-600 hover:underline">Отменить • Esc</button>
    </div>

    <div v-if="isReplying" class="text-sm text-gray-600 mb-1">
      Вы отвечаете на "{{ replyingMessage.text }}" от {{ replyingMessage.username }}
      <button @click="cancelReply" class="ml-2 text-blue-600 hover:underline">Отменить • Esc</button>
    </div>

    <input
      ref="inputRef"
      v-model="text"
      @input="chatStore.sendTyping"
      @keyup.enter="send"
      @keyup.esc="onEsc"
      type="text"
      placeholder="Type a message..."
      class="block w-full h-[60px] p-4 rounded-2xl border border-gray-300 shadow-md bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref, nextTick } from 'vue'
import { useChatStore } from '../../../store/chatStore'

const text = ref('')
const isEditing = ref(false)
const editingMessage = ref(null)

const isReplying = ref(false)
const replyingMessage = ref(null)

const inputRef = ref(null)
const chatStore = useChatStore()

function send() {
  if (text.value.trim() === '') return

  if (isEditing.value && editingMessage.value) {
    chatStore.editMessage(editingMessage.value, text.value)
    cancelEdit()
  } else {
    chatStore.send(text.value, replyingMessage.value)
    cancelReply()
  }

  text.value = ''
}

function startEdit(msg) {
  if (isReplying.value) cancelReply()
  
  text.value = msg.text
  editingMessage.value = msg
  isEditing.value = true
  nextTick(() => inputRef.value?.focus())
}

function startReply(msg) {
  if (isEditing.value) cancelEdit()
  
  isReplying.value = true
  replyingMessage.value = msg
  nextTick(() => inputRef.value?.focus())
}

function cancelEdit() {
  text.value = ''
  editingMessage.value = null
  isEditing.value = false
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
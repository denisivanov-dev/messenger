<template>
  <div
    v-if="hasContextPanel"
    class="chat-context-menu h-[48px] px-4 flex items-center gap-3"
  >
    <!-- context -->
    <div class="flex items-center gap-2 min-w-0 select-none">
      <template v-if="editingMessage">
        <Pencil class="w-4 h-4 text-blue-400 shrink-0" />
        <span class="truncate">
          Вы редактируете: {{ editingMessage.text }}
        </span>
      </template>

      <template v-else-if="replyingMessage">
        <Reply class="w-4 h-4 text-blue-400 shrink-0" />
        <span class="truncate">
          Вы отвечаете {{ replyingMessage.username }}:
          {{ replyingMessage.text }}
        </span>
      </template>

      <!-- context cancel action -->
      <button
        class="group ml-auto relative flex items-center justify-center
              w-6 h-6  text-blue-500 hover:bg-blue-500/10 rounded
              shrink-0"
        @click="cancel"
      >
        <X class="w-4 h-4" />

        <!-- hint on shortcut key -->
        <span
          class="absolute bottom-7 right-0 px-1.5 py-0.5 text-[10px]
                text-gray-100 bg-gray-800/90 rounded opacity-0 scale-95
                transition-all duration-150 pointer-events-none
                group-hover:opacity-100 group-hover:scale-100"
        >
          Esc
        </span>
      </button>
    </div>

    <!-- typing pushed to the right when answering or editing msg -->
   <div
      v-if="hasTyping"
      class="ml-auto flex items-center gap-2 text-xs text-gray-500 whitespace-nowrap"
      >
      <span class="typing-dots">
        <i></i><i></i><i></i>
      </span>
      <span>{{ typingText }}</span>
    </div>
  </div>

  <div
    v-else-if="hasTyping"
    class="px-4 py-1 flex items-center gap-2 text-xs text-gray-500"
  >
    <span class="typing-dots">
      <i></i><i></i><i></i>
    </span>
    <span>{{ typingText }}</span>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted , computed } from 'vue'
import { Pencil, Reply, X } from 'lucide-vue-next'
import { useTypingStore } from '../../store/events/typingStore'

/* props */
const props = defineProps({
  editingMessage: { type: Object, default: null },
  replyingMessage: { type: Object, default: null },
})

const emit = defineEmits(['cancel-edit', 'cancel-reply'])

const typingStore = useTypingStore()
const typingUsers = computed(() => typingStore.typingUsers)

const hasTyping = computed(() => typingUsers.value.length > 0)
const hasContextPanel = computed(() =>
  !!props.editingMessage || !!props.replyingMessage
)

const typingText = computed(() => {
  const users = typingUsers.value.map(u => u.username)

  if (users.length === 1) return `${users[0]} печатает`
  if (users.length === 2) return `${users[0]} и ${users[1]} печатают`
  return `${users[0]} и ещё ${users.length - 1} печатают`
})

/* actions */
function cancel() {
  if (props.editingMessage) emit('cancel-edit')
  if (props.replyingMessage) emit('cancel-reply')
}

onMounted(() => {
  window.__cancelChatContext = cancel
})

onUnmounted(() => {
  if (window.__cancelChatContext === cancel) {
    window.__cancelChatContext = null
  }
})
</script>

<style>
  .typing-dots {
  display: inline-flex;
  align-items: center;
  gap: 3px;
}

.typing-dots i {
  width: 4px;
  height: 4px;
  background-color: #9CA3AF;
  border-radius: 50%;
  opacity: 0.4;
  animation: typing-bounce 1.4s infinite ease-in-out;
}

.typing-dots i:nth-child(1) {
  animation-delay: 0s;
}
.typing-dots i:nth-child(2) {
  animation-delay: 0.2s;
}
.typing-dots i:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing-bounce {
  0%, 80%, 100% {
    transform: translateY(0);
    opacity: 0.4;
  }
  40% {
    transform: translateY(-4px);
    opacity: 1;
  }
}

</style>
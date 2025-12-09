<template>
  <!-- SYSTEM MESSAGE (temporarily disabled) -->
  <!-- <MessageSystem
    v-if="false && isSystemMessage"
    :message="message"
  /> -->

  <!-- USER MESSAGE -->
  <MessageUser
    :message="message"
    @reply-to-message="$emit('reply-to-message', $event)"
    @edit-message="$emit('edit-message', $event)"
    @scroll-to-message="$emit('scroll-to-message', $event)"
  />
</template>

<script setup>
import { computed } from 'vue'

import MessageUser from './message-types/MessageUser.vue'
// import MessageSystem from './MessageSystem.vue'

const props = defineProps({
  message: { type: Object, required: true }
})

const emit = defineEmits([
  'reply-to-message',
  'edit-message',
  'scroll-to-message'
])

/* Determine message type (kept, but unused for now) */
const isSystemMessage = computed(() =>
  props.message.sender_id === '0'
  // props.message.username === 'system' &&
  // props.message.type?.startsWith('call_')
)
</script>

<template>
  <!-- SYSTEM MESSAGE -->
  <MessageSystem
    v-if="isSystemMessage"
    :message="message"
  />

  <!-- USER MESSAGE -->
  <MessageUser
    v-else
    :message="message"
    @reply-to-message="$emit('reply-to-message', $event)"
    @edit-message="$emit('edit-message', $event)"
    @scroll-to-message="$emit('scroll-to-message', $event)"
  />
</template>

<script setup>
import { computed } from 'vue'

import MessageUser from './message-types/MessageUser.vue'
// import MessageSystem from './MessageSystem.vue'  // тебе я сделаю позже

const props = defineProps({
  message: { type: Object, required: true }
})

const emit = defineEmits([
  'reply-to-message',
  'edit-message',
  'scroll-to-message'
])

/* Determine message type */
const isSystemMessage = computed(() =>
  props.message.user_id === '0' &&
  props.message.username === 'system' &&
  props.message.type?.startsWith('call_')
)
</script>

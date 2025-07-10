<template>
  <div class="message p-2 bg-gray-200 w-[100%] rounded-md break-words">
    <strong>{{ message.username || 'неизвестный' }}:</strong> {{ message.text }}
  </div>
</template>

<script setup>
defineProps({
  message: {
    type: Object,
    required: true
  }
})
</script>
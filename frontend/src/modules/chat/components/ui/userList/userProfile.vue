<template>
  <div
    ref="rootElement"
    v-if="visible"
    class="fixed z-50 bg-white shadow-lg rounded-2xl p-6 flex flex-col items-center gap-4 border border-gray-300"
    :style="{ top: y + 'px', left: x + 'px', width: '288px' }"
  >
    <h2 class="text-xl font-semibold break-words text-center">{{ user.username }}</h2>
    <button
      v-if="!isSelf"
      class="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700"
      @click="write"
    >
      Написать
    </button>
  </div>
</template>

<script setup>
import { ref, defineExpose, defineProps, defineEmits } from 'vue'

const rootElement = ref(null)

defineExpose({
  profileRootElement: rootElement
})

const props = defineProps({
  user: { type: Object, required: true },
  visible: { type: Boolean, required: true },
  x: { type: Number, required: true },
  y: { type: Number, required: true },
  isSelf:   Boolean 
})
const emit = defineEmits(['write'])

function write() {
  emit('write')
}
</script>
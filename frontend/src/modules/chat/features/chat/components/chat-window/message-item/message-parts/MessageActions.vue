<template>
  <div
    class="absolute top-0 right-0 mt-1 mr-1 hidden group-hover:flex flex-row
          bg-[#202225]/90 backdrop-blur-sm border border-[#3A3B3E]
          rounded-lg shadow-lg px-2 py-1 z-10 gap-2"
  >
    <!-- Ответить -->
    <button
      @click="$emit('reply')"
      title="Ответить"
      class="text-[#60A5FA] hover:text-[#93C5FD]"
    >
      <ReplyIcon class="w-4 h-4" />
    </button>

    <!-- Редактировать -->
    <button
      v-if="isMy"
      @click="$emit('edit')"
      title="Редактировать"
      class="text-[#FACC15] hover:text-[#FDE68A]"
    >
      <EditIcon class="w-4 h-4" />
    </button>

    <!-- Закрепить -->
    <button
      @click="$emit('pin')"
      title="Закрепить"
      class="text-[#C084FC] hover:text-[#E9D5FF]"
    >
      <PinIcon class="w-4 h-4" />
    </button>

    <!-- Удалить -->
    <button
      v-if="isMy"
      @click="$emit('delete')"
      title="Удалить"
      class="text-[#F87171] hover:text-[#FCA5A5]"
    >
      <TrashIcon class="w-4 h-4" />
    </button>
  </div>
</template>

<script setup>
import { ReplyIcon, EditIcon, PinIcon, TrashIcon } from 'lucide-vue-next'

defineProps({
  isMy: Boolean
})

defineEmits(['reply', 'edit', 'pin', 'delete'])
</script>
<template>
  <div
    class="absolute min-w-[220px]
          bg-[#202225]/95 backdrop-blur-md
          border border-[#3A3B3E]
          rounded-xl shadow-xl
          py-1 z-[9999]
          select-none
          message-click_actions"
    :class="{ 'open-up': openUp }"
    :style="{
      top: `${y}px`,
      left: `${x}px`,
      transform: openUp
        ? 'translateY(-100%) scale(1)'
        : 'scale(1)'
    }"
    @click.stop
  >
    <button
      @click="$emit('reply')"
      class="message-click_actions_item"
    >
      <ReplyIcon class="message-click_actions_icon" />
      <span>Ответить</span>
    </button>

    <button
      v-if="isMy"
      @click="$emit('edit')"
      class="message-click_actions_item"
    >
      <EditIcon class="message-click_actions_icon" />
      <span>Редактировать</span>
    </button>

    <button
      @click="$emit('pin')"
      class="message-click_actions_item"
    >
      <PinIcon class="message-click_actions_icon" />
      <span>Закрепить</span>
    </button>

    <div class="message-click_actions_divider" />

    <button
      @click="$emit('copy')"
      class="message-click_actions_item"
    >
      <CopyIcon class="message-click_actions_icon" />
      <span>Копировать текст</span>
    </button>

    <button
      @click="$emit('forward')"
      class="message-click_actions_item"
    >
      <ForwardIcon class="message-click_actions_icon" />
      <span>Переслать</span>
    </button>

    <div class="message-click_actions_divider" />

    <button
      v-if="isMy"
      @click="$emit('delete')"
      class="message-click_actions_item message-click_actions_item--danger"
    >
      <TrashIcon class="message-click_actions_icon" />
      <span>Удалить</span>
    </button>

    <button
      @click="$emit('select')"
      class="message-click_actions_item"
    >
      <CheckSquareIcon class="message-click_actions_icon" />
      <span>Выделить</span>
    </button>
  </div>
</template>

<script setup>
import {
  ReplyIcon,
  EditIcon,
  PinIcon,
  CopyIcon,
  ForwardIcon,
  TrashIcon,
  CheckSquareIcon
} from 'lucide-vue-next'

defineProps({
  x: Number,
  y: Number,
  openUp: Boolean,
  isMy: Boolean
})

defineEmits([
  'reply',
  'edit',
  'pin',
  'copy',
  'forward',
  'delete',
  'select'
])
</script>
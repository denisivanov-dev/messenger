<template>
  <div
    class="absolute top-0 right-0 mt-1 mr-1 flex flex-row
           bg-[#202225]/90 backdrop-blur-sm
           border border-[#3A3B3E]
           rounded-lg shadow-lg
           px-2 py-1 gap-2 z-10"
  >
    <button
      class="hover-menu-btn"
      @pointerenter.stop="show('Ответить', $event)"
      @pointerleave.stop="hide"
      @click="$emit('reply')"
    >
      <ReplyIcon class="hover-menu-icon" />
    </button>

    <button
      v-if="isMy"
      class="hover-menu-btn"
      @pointerenter.stop="show('Редактировать', $event)"
      @pointerleave.stop="hide"
      @click="$emit('edit')"
    >
      <EditIcon class="hover-menu-icon" />
    </button>

    <button
      class="hover-menu-btn"
      @pointerenter.stop="show('Закрепить', $event)"
      @pointerleave.stop="hide"
      @click="$emit('pin')"
    >
      <PinIcon class="hover-menu-icon" />
    </button>

    <button
      v-if="isMy"
      class="hover-menu-btn hover-menu-btn--danger"
      @pointerenter.stop="show('Удалить', $event)"
      @pointerleave.stop="hide"
      @click="$emit('delete')"
    >
      <TrashIcon class="hover-menu-icon" />
    </button>

    <Teleport to="body">
      <div
        v-if="tooltip.visible"
        class="tooltip"
        :style="tooltip.style"
      >
        {{ tooltip.text }}
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ReplyIcon, EditIcon, PinIcon, TrashIcon } from 'lucide-vue-next'

defineProps({ isMy: Boolean })
defineEmits(['reply', 'edit', 'pin', 'delete'])

let timer = null
const DELAY = 300

const tooltip = ref({
  visible: false,
  text: '',
  style: {}
})

function show(text, e) {
  clearTimeout(timer)

  const rect = e.currentTarget.getBoundingClientRect()

  timer = setTimeout(() => {
    tooltip.value = {
      visible: true,
      text,
      style: {
        top: `${rect.top - 30}px`,
        left: `${rect.left + rect.width / 2}px`
      }
    }
  }, DELAY)
}

function hide() {
  clearTimeout(timer)
  tooltip.value.visible = false
}
</script>
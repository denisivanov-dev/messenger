<template>
  <div :class="galleryClass" class="w-full max-w-[600px] mx-0">
    <div
      v-for="(att, index) in attachments"
      :key="att.key"
      :class="getImageClass(attachments.length, index)"
      class="overflow-hidden cursor-pointer rounded-lg relative"
      @click="showImage(index)"
    >
      <img
        :src="imageUrls[att.key]"
        class="w-full h-full object-cover"
        alt="Изображение"
      />
    </div>
  </div>

  <!-- Модалка -->
  <Teleport to="body">
    <div
      v-if="fullscreenIndex !== null"
      class="fixed inset-0 bg-black bg-opacity-90 z-50 flex flex-col items-center justify-center"
      @click="closeImage"
    >
      <div class="relative max-w-[90vw] max-h-[90vh] flex flex-col items-center" @click.stop>
        <img
          :src="imageUrls[attachments[fullscreenIndex].key]"
          class="max-h-[80vh] rounded-xl shadow-xl"
        />

        <!-- Навигация -->
        <div class="mt-2 flex items-center justify-center gap-6 text-white text-xl">
          <button
            @click="prevImage"
            :class="fullscreenIndex === 0 ? 'opacity-50 cursor-default' : 'hover:text-blue-400'"
            :disabled="fullscreenIndex === 0"
          >
            ‹
          </button>

          <span class="text-sm text-white">
            {{ fullscreenIndex + 1 }} / {{ attachments.length }}
          </span>

          <button
            @click="nextImage"
            :class="fullscreenIndex === attachments.length - 1 ? 'opacity-50 cursor-default' : 'hover:text-blue-400'"
            :disabled="fullscreenIndex === attachments.length - 1"
          >
            ›
          </button>
        </div>

        <!-- Ссылка на оригинал -->
        <a
          :href="imageUrls[attachments[fullscreenIndex].key]"
          target="_blank"
          rel="noopener noreferrer"
          class="mt-2 text-sm text-white underline hover:text-blue-300"
        >
          Открыть оригинал
        </a>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  attachments: Array,
  imageUrls: Object,
  openImage: Function,
})

const fullscreenIndex = ref(null)

function showImage(index) {
  fullscreenIndex.value = index
}

function closeImage() {
  fullscreenIndex.value = null
}

function prevImage() {
  if (fullscreenIndex.value > 0) fullscreenIndex.value--
}

function nextImage() {
  if (fullscreenIndex.value < props.attachments.length - 1) fullscreenIndex.value++
}

function onKeydown(e) {
  if (fullscreenIndex.value === null) return
  if (e.key === 'Escape') closeImage()
  if (e.key === 'ArrowLeft') prevImage()
  if (e.key === 'ArrowRight') nextImage()
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown))

const galleryClass = computed(() => {
  const len = props.attachments.length
  if (len === 5) return 'grid grid-cols-6 gap-1'
  if (len === 1) return 'grid grid-cols-1'
  if (len === 2) return 'grid grid-cols-2 gap-1'
  if (len === 3) return 'grid grid-cols-2 gap-1'
  if (len === 4) return 'grid grid-cols-2 gap-1'
  return 'grid grid-cols-3 gap-1'
})

function getImageClass(length, index) {
  if (length === 1) return 'w-full h-[300px]'
  if (length === 2) return 'h-[200px]'
  if (length === 3) return index < 2 ? 'h-[150px]' : 'col-span-2 h-[200px]'
  if (length === 4) return 'h-[150px]'
  if (length === 5) {
    if (index === 0 || index === 1) return 'col-span-3 h-[140px]'  // 2 сверху = 3+3 = 6
    return 'col-span-2 h-[110px]'                                  // 3 снизу = 2+2+2 = 6
  }
  return ''
}
</script>
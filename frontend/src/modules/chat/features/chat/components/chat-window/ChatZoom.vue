<template>
  <!-- ZOOM HUD that shows on ctrl, expands on hover -->
  <div
    v-if="showZoomHud"
    @mouseenter="onZoomHudEnter"
    @mouseleave="onZoomHudLeave"
    class="group absolute top-3 right-4 z-[60]
           flex items-center
           bg-[#2A2B2F]
           rounded-md
           text-sm text-white
           shadow-md
           overflow-hidden"
  >
    <!-- % zoom percent -->
    <div class="px-2 py-1 whitespace-nowrap">
      {{ Math.round(zoom * 100) }}%  
    </div>

    <!-- actions -->
    <div
      class="flex items-center gap-1 px-1
             max-w-0 overflow-hidden
             transition-all duration-200 ease-out
             group-hover:max-w-[120px]"
    >
      <button class="hover:text-blue-400" @click="applyZoom(-0.05)">
        <Minus class="w-4 h-4" />
      </button>

      <button class="hover:text-blue-400" @click="applyZoom(+0.05)">
        <Plus class="w-4 h-4" />
      </button>

      <button class="hover:text-blue-400" @click="resetZoom">
        <RotateCcw class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, inject, watch } from 'vue'
import { Plus, Minus, RotateCcw } from 'lucide-vue-next'

const zoom = inject('zoom')

const STORAGE_KEY = 'app:zoom'

const showZoomHud = ref(false)
const isHoveringZoomHud = ref(false)
let zoomHudTimer = null

/* show hud */
function showZoomTemporarily() {
  showZoomHud.value = true
  clearTimeout(zoomHudTimer)

  zoomHudTimer = setTimeout(() => {
    if (!isHoveringZoomHud.value) {
      showZoomHud.value = false
    }
  }, 1200)
}

function onZoomHudEnter() {
  isHoveringZoomHud.value = true
  clearTimeout(zoomHudTimer)
}

function onZoomHudLeave() {
  isHoveringZoomHud.value = false
  showZoomTemporarily()
}

/* zoom actions */
function applyZoom(v) {
  zoom.value = Math.min(2, Math.max(0.6, zoom.value + v))
  showZoomTemporarily()
}

function resetZoom() {
  zoom.value = 1
  showZoomTemporarily()
}

/* mouse wheel zoom */
function zoomWheel(e) {
  if (!e.ctrlKey) return
  e.preventDefault()
  applyZoom(e.deltaY < 0 ? +0.05 : -0.05)
}

/* keyboard zoom */
function zoomKeys(e) {
  if (!e.ctrlKey) return

  if (e.key === '+' || e.key === '=') {
    applyZoom(+0.05)
    e.preventDefault()
  }

  if (e.key === '-') {
    applyZoom(-0.05)
    e.preventDefault()
  }

  if (e.key === '0') {
    resetZoom()
    e.preventDefault()
  }
}

/* init */
onMounted(() => {
  // load zoom
  const saved = Number(localStorage.getItem(STORAGE_KEY))
  if (!Number.isNaN(saved) && saved >= 0.6 && saved <= 2) {
    zoom.value = saved
  }

  window.addEventListener('wheel', zoomWheel, { passive: false })
  window.addEventListener('keydown', zoomKeys)
})

onBeforeUnmount(() => {
  window.removeEventListener('wheel', zoomWheel)
  window.removeEventListener('keydown', zoomKeys)
})

/* save zoom */
watch(zoom, (v) => {
  localStorage.setItem(STORAGE_KEY, v)
})
</script>
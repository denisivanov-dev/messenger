  <template>
    <div
      class="fixed bg-white rounded-xl shadow-2xl p-4 z-50 cursor-move border border-gray-300"
      ref="panel"
      :style="{ top: `${props.position.top}px`, left: `${props.position.left}px` }"
      @mousedown="startDrag"
    >
      <!-- Заголовок + кнопки -->
      <div class="flex justify-between items-center mb-3">
        <h2 class="text-lg font-semibold">Друзья</h2>
        <div class="flex gap-2 items-center">
          <button
            @click="resetPosition"
            draggable="false"
            @mousedown.stop
            class="text-gray-500 hover:text-blue-600 transition"
            title="Сбросить позицию"
          >
            <RotateCcw class="w-5 h-5" />
          </button>
          <button
            @click.stop="$emit('close')"
            draggable="false"
            @mousedown.stop
            class="text-gray-500 hover:text-red-600 transition"
            title="Закрыть"
          >
            <X class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- Фильтры -->
      <div class="flex gap-3 mb-4">
        <button
          v-for="filter in filters"
          :key="filter"
          @click="selectedFilter = filter"
          draggable="false"
          @mousedown.stop
          :class="[
            'px-3 py-1 rounded-full text-sm transition',
            selectedFilter === filter
              ? 'bg-gray-800 text-white'
              : 'text-gray-500 hover:text-white hover:bg-gray-700'
          ]"
        >
          {{ getLabel(filter) }}
        </button>
      </div>

      <!-- Подключение компонентов -->
      <FriendListPending v-if="selectedFilter === 'pending'" />
      <FriendListOnline v-else-if="selectedFilter === 'online'" />
      <FriendListAll
        v-if="selectedFilter === 'all'"
        v-model:isModalOpen="isModalOpen"
        v-model:blockDrag="blockDrag"
      />
    </div>
  </template>

  <script setup>
  import { ref } from 'vue'
  import { X, RotateCcw } from 'lucide-vue-next'
  import FriendListPending from './FriendListPending.vue'
  import FriendListOnline from './FriendListOnline.vue'
  import FriendListAll from './FriendListAll.vue'

  const props = defineProps({
    position: Object
  })
  const emit = defineEmits(['update:position', 'close'])

  const isModalOpen = ref(false)
  const blockDrag = ref(false)

  const panel = ref(null)
  let offsetX = 0
  let offsetY = 0
  let dragging = false

  function resetPosition() {
    emit('update:position', { top: 30, left: 300 })
  }

  function startDrag(e) {
    if (blockDrag.value) return

    const notDraggable = ['button', 'svg', 'path']
    const tag = e.target.tagName.toLowerCase()
    const insideNoDragZone = e.target.closest('.no-drag')

    if (notDraggable.includes(tag) || insideNoDragZone) return

    dragging = true
    offsetX = e.clientX - panel.value.offsetLeft
    offsetY = e.clientY - panel.value.offsetTop
    document.addEventListener('mousemove', drag)
    document.addEventListener('mouseup', stopDrag)
  }

  function drag(e) {
    if (blockDrag.value || !dragging) return
    emit('update:position', { top: e.clientY - offsetY, left: e.clientX - offsetX })
  }

  function stopDrag(e) {
    dragging = false
    emit('update:position', { top: e.clientY - offsetY, left: e.clientX - offsetX })
    document.removeEventListener('mousemove', drag)
    document.removeEventListener('mouseup', stopDrag)
  }

  const filters = ['online', 'all', 'pending']
  const selectedFilter = ref('online')
  function getLabel(key) {
    return { online: 'В сети', all: 'Все', pending: 'Ожидание' }[key]
  }
  </script>
<template>
  <div
    ref="sidebarWindowRef"
    class="sidebar-window w-[250px] h-[864px] p-5 bg-gray-100 rounded-2xl shadow-md overflow-y-auto flex flex-col gap-2"
  >
    <!-- Кнопка "Друзья" с бейджем -->
    <div class="relative">
      <button
        ref="friendButtonRef"
        class="bg-white text-black px-3 py-2 rounded-xl shadow hover:bg-gray-200 transition w-full text-left"
        @click="toggleFriendPanel"
      >
        Друзья
      </button>

      <!-- Бейдж -->
      <span
        v-if="pendingCount > 0"
        class="absolute -top-1 -right-1 bg-red-500 text-white text-xs min-w-[20px] h-[20px] px-1 flex items-center justify-center rounded-full"
      >
        {{ pendingCount }}
      </span>
    </div>

    <!-- Остальной контент сайдбара -->
    <div>Нигер</div>

    <!-- Плавающее окно друзей -->
    <Transition name="fade-slide-left" appear>
      <DraggableFriendPanel
        ref="friendPanelRef"
        v-if="showFriendPanel"
        :position="friendPanelPosition"
        @update:position="friendPanelPosition = $event"
        @close="showFriendPanel = false"
      />
    </Transition>

    <!-- Скрытый компонент для учёта входящих -->
    <FriendListPending
      @pending-count="pendingCount = $event"
      style="display: none"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import DraggableFriendPanel from './friendPanel/draggableFriendPanel.vue.vue'
import FriendListPending from './friendPanel/FriendListPending.vue'

const showFriendPanel = ref(false)
const friendPanelPosition = ref({ top: 30, left: 300 })
const pendingCount = ref(0)

const sidebarWindowRef = ref(null)
const friendPanelRef = ref(null)
const friendButtonRef = ref(null)

function toggleFriendPanel() {
  showFriendPanel.value = !showFriendPanel.value
}

function closeFriendPanel() {
  showFriendPanel.value = false
}

defineExpose({
  sidebarWindowRef,
  friendPanelRef,
  friendButtonRef,
  closeFriendPanel
})
</script>

<style scoped>
.fade-slide-left-enter-active,
.fade-slide-left-leave-active {
  transition: all 0.3s ease;
}
.fade-slide-left-enter-from,
.fade-slide-left-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}
.fade-slide-left-enter-to,
.fade-slide-left-leave-from {
  opacity: 1;
  transform: translateX(0);
}
</style>
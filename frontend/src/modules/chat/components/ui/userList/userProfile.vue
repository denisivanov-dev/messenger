<template>
  <div
    ref="rootElement"
    v-if="visible"
    class="fixed z-50 bg-white shadow-lg rounded-2xl p-6 flex flex-col gap-4 border border-gray-300"
    :style="{ top: y + 'px', left: x + 'px', width: '288px' }"
  >
    <!-- Ник -->
    <h2 class="text-xl font-semibold break-words select-text">{{ user.username }}</h2>

    <!-- Панель кнопок -->
    <div v-if="!isSelf" class="flex items-center gap-2">
      <!-- === Вариант: Входящий запрос === -->
      <template v-if="friendStatus === 'incoming'">
        <!-- Написать -->
        <button
          class="flex-grow flex items-center justify-between gap-1 px-3 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 text-sm"
          @click="write"
        >
          <span class="truncate">Написать</span>
          <MessageSquareIcon class="w-4 h-4" />
        </button>

        <!-- Принять -->
        <div class="relative group w-10 h-10 shrink-0">
          <button
            class="w-10 h-10 flex items-center justify-center bg-emerald-500 hover:bg-emerald-600 text-white rounded-lg"
            @click.stop="emit('accept-friend')"
          >
            <UserCheckIcon class="w-5 h-5" />
          </button>
          <div
            class="absolute bottom-full mb-2 left-1/2 -translate-x-1/2 px-3 py-1 rounded-md bg-black text-white text-sm whitespace-nowrap opacity-0 group-hover:opacity-100 transition"
          >
            Принять
          </div>
        </div>

        <!-- Отклонить -->
        <div class="relative group w-10 h-10 shrink-0">
          <button
            class="w-10 h-10 flex items-center justify-center bg-red-500 hover:bg-red-600 text-white rounded-lg"
            @click.stop="emit('decline-friend')"
          >
            <UserXIcon class="w-5 h-5" />
          </button>
          <div
            class="absolute bottom-full mb-2 left-1/2 -translate-x-1/2 px-3 py-1 rounded-md bg-black text-white text-sm whitespace-nowrap opacity-0 group-hover:opacity-100 transition"
          >
            Отклонить
          </div>
        </div>
      </template>

      <!-- === Остальные статусы: outgoing, none, friends === -->
      <template v-else>
        <!-- Написать -->
        <button
          class="flex-grow flex items-center justify-between gap-1 px-3 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 text-sm"
          @click="write"
        >
          <span class="truncate">Написать</span>
          <MessageSquareIcon class="w-4 h-4" />
        </button>

        <!-- Кнопка действия -->
        <div class="relative group w-10 h-10 shrink-0">
          <button
            class="w-10 h-10 rounded-lg flex items-center justify-center transition-colors"
            :class="{
              'bg-red-500 hover:bg-red-600 text-white': friendStatus === 'outgoing',
              'bg-emerald-500 hover:bg-emerald-600 text-white': friendStatus === 'none',
              'bg-blue-500 text-white cursor-default': friendStatus === 'friends'
            }"
          >
            <component
              :is="friendStatus === 'outgoing'
                ? UserXIcon
                : friendStatus === 'friends'
                ? UserCheckIcon
                : UserPlusIcon"
              class="w-5 h-5"
            />
          </button>

          <!-- Подсказка -->
          <div
            class="absolute bottom-full mb-2 left-1/2 -translate-x-1/2 px-3 py-1 rounded-md bg-black text-white text-sm whitespace-nowrap opacity-0 group-hover:opacity-100 transition"
          >
            {{
              friendStatus === 'outgoing'
                ? 'Отменить заявку'
                : friendStatus === 'friends'
                ? 'Друзья'
                : 'Добавить в друзья'
            }}
          </div>

          <!-- Обработка клика -->
          <div
            class="absolute inset-0"
            @click.stop="friendStatus !== 'friends' && handleFriendAction()"
          ></div>
        </div>
      </template>
    </div>

    <!-- Подтверждение отмены -->
    <div
      v-if="showCancelConfirm"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
    >
      <div
        ref="cancelConfirmRef"
        @click.stop
        class="bg-white rounded-xl shadow-lg p-6 w-[300px] text-center space-y-4 select-none cursor-default no-drag"
      >
        <p class="text-lg font-semibold text-gray-800">
          Отменить заявку {{ user.username }}?
        </p>
        <div class="flex justify-center gap-4">
          <button
            @click="confirmCancelFriend"
            class="px-4 py-2 rounded-lg bg-green-600 text-white hover:bg-green-700 transition"
          >
            Да
          </button>
          <button
            @click="cancelPrompt"
            class="px-4 py-2 rounded-lg bg-red-500 text-white hover:bg-red-600 transition"
          >
            Нет
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, defineExpose, defineProps, defineEmits, onMounted, onUnmounted } from 'vue'
import {
  MessageSquareIcon,
  UserPlusIcon,
  UserXIcon,
  UserCheckIcon
} from 'lucide-vue-next'
import { eventBus } from '../../../utils/eventBus'
import { useChatStore } from '../../../store/chatStore'

const rootElement = ref(null)
const cancelConfirmRef = ref(null)
const showCancelConfirm = ref(false)

defineExpose({ profileRootElement: rootElement, cancelConfirmRef  })

const props = defineProps({
  user: { type: Object, required: true },
  visible: { type: Boolean, required: true },
  x: { type: Number, required: true },
  y: { type: Number, required: true },
  isSelf: Boolean,
})

const emit = defineEmits(['write', 'add-friend', 'cancel-friend', 'accept-friend', 'decline-friend'])

function write() {
  emit('write')
}

function handleFriendAction() {
  if (friendStatus.value === 'outgoing') {
    showCancelConfirm.value = true
  } else if (friendStatus.value === 'none') {
    emit('add-friend')
  }
}

function confirmCancelFriend() {
  emit('cancel-friend')
  showCancelConfirm.value = false
}

function cancelPrompt() {
  showCancelConfirm.value = false
}

const chatStore = useChatStore()
const friendStatus = ref(null)

function updateFriendStatus(status) {
  friendStatus.value = status
}

function handleFriendUpdate(msg) {
  if (msg.user_id === props.user.id) {
    updateFriendStatus(msg.status)
  }
}

onMounted(() => {
  if (props.user?.id) {
    friendStatus.value = chatStore.friendStatusCache[props.user.id] ?? 'none'
  }
  eventBus.on('friend-update', handleFriendUpdate)
})

onUnmounted(() => {
  eventBus.off('friend-update', handleFriendUpdate)
})
</script>
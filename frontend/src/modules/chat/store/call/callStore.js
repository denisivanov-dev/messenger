import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useChatStore } from '../chatStore'
import { useAuthStore } from '../../../auth/store/authStore'
import { useWebRTCStore } from './webrtcStore'
import { useMediaStore } from './mediaStore'
import { sendMessage } from '../../api/chatApi'

export const useCallStore = defineStore('call', () => {
  const authStore = useAuthStore()
  const chatStore = useChatStore()
  const webrtcStore = useWebRTCStore()
  const mediaStore = useMediaStore()

  const inCall = ref(false)
  const isCalling = ref(false)
  const hasJoined = ref(false)
  const incomingCallFrom = ref(null)
  const callMembers = ref({})

  const cameraStatusMap = ref({})
  const screenStatusMap = ref({})

  const getMyId = () => String(authStore.getUserId)

  const setMemberStatus = (userId, status) => {
    callMembers.value[String(userId)] = status
  }

  const removeMember = (userId) => {
    delete callMembers.value[String(userId)]
    delete cameraStatusMap.value[String(userId)]
    delete screenStatusMap.value[String(userId)]
  }

  function updateCameraStatus(userId, isEnabled) {
    cameraStatusMap.value = {
      ...cameraStatusMap.value,
      [String(userId)]: isEnabled
    }
  }

  function updateScreenStatus(userId, isEnabled) {
    screenStatusMap.value = {
      ...screenStatusMap.value,
      [String(userId)]: isEnabled
    }
  }

  function sendCameraStatusUpdate(isEnabled) {
    sendMessage({
      type: 'camera_status',
      chat_type: 'private',
      receiver_id: chatStore.receiverID,
      enabled: isEnabled,
      user_id: String(authStore.getUserId)
    })

    updateCameraStatus(String(authStore.getUserId), isEnabled)
  }

  function sendScreenStatusUpdate(isEnabled) {
    sendMessage({
      type: 'screen_status',
      chat_type: 'private',
      receiver_id: chatStore.receiverID,
      enabled: isEnabled,
      user_id: String(authStore.getUserId)
    })

    updateScreenStatus(String(authStore.getUserId), isEnabled)
  }

  async function startRequestCall() {
    const myId = getMyId()

    inCall.value = true
    isCalling.value = true
    hasJoined.value = true

    callMembers.value = {
      [myId]: 'joined',
      [String(chatStore.receiverID)]: 'calling'
    }

    sendMessage({
      type: 'start_call',
      chat_type: 'private',
      receiver_id: chatStore.receiverID
    })
  }

  function cancelRequestCall() {
    if (!isCalling.value) return

    sendMessage({
      type: 'cancel_call',
      chat_type: 'private',
      receiver_id: chatStore.receiverID
    })

    resetCallState()
  }

  function handleIncomingCall(fromUserID) {
    incomingCallFrom.value = fromUserID
    setMemberStatus(fromUserID, 'joined')
  }

  function handleCallCanceled(fromUserID) {
    if (incomingCallFrom.value === fromUserID) {
      incomingCallFrom.value = null
      removeMember(fromUserID)
    }
  }

  async function acceptCall(fromUserID) {
    incomingCallFrom.value = null
    hasJoined.value = true
    setMemberStatus(getMyId(), 'joined')

    sendMessage({
      type: 'call_answer',
      chat_type: 'private',
      receiver_id: fromUserID,
      accepted: true
    })

    await webrtcStore.startCall(fromUserID, false)
    inCall.value = true
  }

  function declineCall(fromUserID) {
    sendMessage({
      type: 'call_answer',
      chat_type: 'private',
      receiver_id: fromUserID,
      accepted: false
    })

    incomingCallFrom.value = null
  }

  async function handleCallAnswer(fromUserID, accepted) {
    isCalling.value = false
    if (!accepted) {
      removeMember(fromUserID)
      return
    }

    setMemberStatus(fromUserID, 'joined')
    await webrtcStore.startCall(fromUserID, true)
    inCall.value = true
  }

  async function joinCall() {
    const myId = getMyId()

    mediaStore.loadMicSettings()
    mediaStore.loadCamSettings()

    if (!callMembers.value[myId]) {
      setMemberStatus(myId, 'joined')
    }

    hasJoined.value = true
    inCall.value = true

    sendCameraStatusUpdate(mediaStore.camSettings.enabled)

    sendMessage({
      type: 'join_call',
      chat_type: 'private',
      receiver_id: chatStore.receiverID
    })

    const others = Object.keys(callMembers.value).filter(id => id !== myId)
    await Promise.all(others.map(id => webrtcStore.startCall(id, false)))
  }

  function handleJoinCall(userId) {
    setMemberStatus(userId, 'joined')

    if (hasJoined.value) {
      webrtcStore.startCall(userId, true)
    }
  }

  function leaveCall() {
    const myId = getMyId()

    if (isCalling.value) {
      cancelRequestCall()
      return
    }

    sendMessage({
      type: 'leave_call',
      chat_type: 'private',
      receiver_id: chatStore.receiverID
    })

    Object.keys(callMembers.value).forEach(id => {
      if (id !== myId) {
        webrtcStore.endCallWith(id)
      }
    })

    removeMember(myId)
    hasJoined.value = false

    if (Object.keys(callMembers.value).length === 0) {
      resetCallState()
    }
  }

  function resetCallState() {
    inCall.value = false
    isCalling.value = false
    hasJoined.value = false
    incomingCallFrom.value = null
    callMembers.value = {}
    webrtcStore.endAllCalls()
  }

  return {
    inCall,
    isCalling,
    hasJoined,
    incomingCallFrom,
    callMembers,
    cameraStatusMap,
    screenStatusMap,

    startRequestCall,
    cancelRequestCall,
    acceptCall,
    declineCall,
    handleIncomingCall,
    handleCallCanceled,
    handleCallAnswer,
    joinCall,
    handleJoinCall,
    leaveCall,
    resetCallState,

    setMemberStatus,
    removeMember,
    updateCameraStatus,
    updateScreenStatus,
    sendCameraStatusUpdate,
    sendScreenStatusUpdate
  }
})
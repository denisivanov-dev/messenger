import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useChatStore } from '../chatStore'
import { useAuthStore } from '../../../auth/store/authStore'
import { useWebRTCStore } from './webrtcStore'
import { useMediaStore } from './mediaStore'
import { sendSocketPayload } from '../../api/chatApi'

export const useCallStore = defineStore('call', () => {
  const authStore = useAuthStore()
  const chatStore = useChatStore()
  const webrtcStore = useWebRTCStore()
  const mediaStore = useMediaStore()

  // for UI
  const inCall = ref(false)
  const isCalling = ref(false)
  const incomingCallFrom = ref(null)
  // for WebRTC
  const hasJoined = ref(false)

  const callMembers = ref({})

  const micStatusMap = ref({})
  const cameraStatusMap = ref({})
  const screenStatusMap = ref({})

  const showCamPopup = ref(false)

  const getMyId = () => String(authStore.getUserId)

  const setMemberStatus = (userId, status) => {
    callMembers.value[String(userId)] = status
  }

  const removeMember = (userId) => {
    delete callMembers.value[String(userId)]
    delete cameraStatusMap.value[String(userId)]
    delete screenStatusMap.value[String(userId)]
  }

  function updateMicStatus(userId, isEnabled) {
    micStatusMap.value = {
      ...micStatusMap.value,
      [String(userId)]: isEnabled
    }
  }

  function sendMicStatusUpdate(isEnabled) {
    sendSocketPayload({
      type: 'mic_status',
      chat_type: 'private',
      receiver_id: chatStore.receiverID,
      enabled: isEnabled,
      user_id: String(authStore.getUserId)
    })

    updateMicStatus(String(authStore.getUserId), isEnabled)
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
    sendSocketPayload({
      type: 'camera_status',
      chat_type: 'private',
      receiver_id: chatStore.receiverID,
      enabled: isEnabled,
      user_id: String(authStore.getUserId)
    })

    updateCameraStatus(String(authStore.getUserId), isEnabled)
  }

  function sendScreenStatusUpdate(isEnabled) {
    sendSocketPayload({
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

    sendSocketPayload({
      type: 'start_call',
      chat_type: 'private',
      receiver_id: chatStore.receiverID
    })
  }

  function cancelRequestCall() {
    if (!isCalling.value) return

    sendSocketPayload({
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

    sendSocketPayload({
      type: 'call_answer',
      chat_type: 'private',
      receiver_id: fromUserID,
      accepted: true
    })

    await webrtcStore.startCall(fromUserID, false)
    inCall.value = true
  }

  function declineCall(fromUserID) {
    sendSocketPayload({
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
    const prevCamSettings = { ...mediaStore.loadCamSettings() }
    console.info('prevCamSettings:', prevCamSettings)

    if (!callMembers.value[myId]) {
      setMemberStatus(myId, 'joined')
    }

    hasJoined.value = true
    inCall.value = true

    mediaStore.camSettings.enabled = false
    sendCameraStatusUpdate(false)
    mediaStore.saveCamSettings()

    sendSocketPayload({
      type: 'join_call',
      chat_type: 'private',
      receiver_id: chatStore.receiverID
    })

    const others = Object.keys(callMembers.value).filter(id => id !== myId)
    await Promise.all(others.map(id => webrtcStore.startCall(id, false)))

    setTimeout(() => {
      webrtcStore.renegotiateWithAll()
    }, 200)

    if (prevCamSettings?.enabled) {
      showCamPopup.value = true
    }

    console.info('----------------------')
    console.info(mediaStore.remoteCameraStreams)
    console.info(cameraStatusMap.value)
    console.info('showCamPopup:', showCamPopup.value)
    console.info('----------------------')
  }


  function handleJoinCall(userId) {
    setMemberStatus(userId, 'joined')

    if (hasJoined.value) {
      webrtcStore.startCall(userId, true)
      webrtcStore.renegotiateWithAll()
    }
  }

  function handleLeaveCall(userId) {
    removeMember(userId)            
    webrtcStore.endCallWith(userId)
  }

  async function leaveCall() {
    const myId = getMyId()

    if (isCalling.value) {
      cancelRequestCall()
      return
    }

    sendSocketPayload({
      type: 'leave_call',
      chat_type: 'private',
      receiver_id: chatStore.receiverID
    })

    // дожидаемся завершения всех endCallWith
    await Promise.all(
      Object.keys(callMembers.value)
        .filter(id => id !== myId)
        .map(id => webrtcStore.endCallWith(id))
    )

    removeMember(myId)
    hasJoined.value = false
    inCall.value = false
    
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
    micStatusMap,
    cameraStatusMap,
    screenStatusMap,
    showCamPopup,

    startRequestCall,
    cancelRequestCall,
    acceptCall,
    declineCall,
    handleIncomingCall,
    handleCallCanceled,
    handleCallAnswer,
    joinCall,
    handleJoinCall,
    handleLeaveCall,
    leaveCall,
    resetCallState,

    setMemberStatus,
    removeMember,
    updateMicStatus,
    updateCameraStatus,
    updateScreenStatus,
    sendMicStatusUpdate,
    sendCameraStatusUpdate,
    sendScreenStatusUpdate
  }
})
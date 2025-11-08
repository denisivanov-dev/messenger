export function handleCallEvents(msg, stores) {
  const { callStore, messagesStore } = stores

  switch (msg.type) {
    case 'incoming_call':
      return callStore.handleIncomingCall(msg.from_user)
    case 'incoming_cancel_call':
      return callStore.handleCallCanceled(msg.from_user)
    case 'incoming_call_answer':
      return callStore.handleCallAnswer(msg.from_user, msg.accepted)
    case 'incoming_join_call':
      return callStore.handleJoinCall(msg.from_user)
    case 'incoming_leave_call':
      return callStore.handleLeaveCall(msg.from_user)
    case 'incoming_camera_status':
      return callStore.updateCameraStatus(msg.from_user, msg.enabled)
    case 'incoming_screen_status':
      return callStore.updateScreenStatus(msg.from_user, msg.enabled)
    case 'incoming_mic_status':
      return callStore.updateMicStatus(msg.from_user, msg.enabled)
    case 'call_started':
      if (msg.call_info.status === 'ongoing') {
        messagesStore.pushFromWs(msg)
      } else {
        const existing = messagesStore.findById(msg.message_id)
        if (existing) messagesStore.updateSystemMessage(msg)
        else messagesStore.pushFromWs(msg)
      }
      break
  }
}
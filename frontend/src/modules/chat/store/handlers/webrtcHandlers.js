export function handleWebRTCEvents(msg, stores) {
  const { webrtcStore } = stores

  switch (msg.type) {
    case 'incoming_webrtc_offer':
      return webrtcStore.handleOffer(msg.from_user, msg.offer)
    case 'incoming_webrtc_answer':
      return webrtcStore.handleAnswer(msg.from_user, msg.answer)
    case 'incoming_ice_candidate':
      return webrtcStore.handleIceCandidate(msg.from_user, msg.candidate)
  }
}
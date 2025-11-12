import { handleUserStatus, handleTyping } from './handlers/eventHandlers'
import { handleMessageEvents } from './handlers/messageHandlers'
import { handleFriendRequestUpdate } from './handlers/friendHandlers'
import { handleCallEvents } from './handlers/callHandlers'
import { handleWebRTCEvents } from './handlers/webrtcHandlers'

export const wsHandlersMap = {
  // --- MESSAGE EVENTS ---
  'message_send': handleMessageEvents,
  'message_edit': handleMessageEvents,
  'message_delete': handleMessageEvents,
  'message_pin': handleMessageEvents,

  // --- USER / TYPING EVENTS ---
  'event_user_status': handleUserStatus,
  'event_typing': handleTyping,

  // --- FRIEND EVENTS ---
  'event_friend_request_update': handleFriendRequestUpdate,

  // --- CALL EVENTS ---
  'call_incoming': handleCallEvents,
  'call_cancel': handleCallEvents,
  'call_answer': handleCallEvents,
  'call_join': handleCallEvents,
  'call_leave': handleCallEvents,
  'call_camera_status': handleCallEvents,
  'call_screen_status': handleCallEvents,
  'call_mic_status': handleCallEvents,
  'call_started': handleCallEvents,

  // --- WEBRTC EVENTS ---
  'webrtc_offer': handleWebRTCEvents,
  'webrtc_answer': handleWebRTCEvents,
  'webrtc_ice': handleWebRTCEvents,
}
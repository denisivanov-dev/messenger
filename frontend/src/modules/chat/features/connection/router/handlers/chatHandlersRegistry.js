import { handleUserStatus } from './eventHandlers'
import {
  handleMessageDeleted,
  handleMessageEdited,
  handleMessagePinned
} from './messageHandlers'
import { handleFriendRequestUpdate } from './friendHandlers'
import { handleCallEvents } from './callHandlers'
import { handleWebRTCEvents } from './webrtcHandlers'
import { handleTyping } from './eventHandlers'

export const chatHandlersRegistry = {
  user_status: handleUserStatus,
  typing: handleTyping,
  message_deleted: handleMessageDeleted,
  message_edited: handleMessageEdited,
  message_pinned: handleMessagePinned,
  friend_request_update: handleFriendRequestUpdate,
  incoming_call: handleCallEvents,
  incoming_cancel_call: handleCallEvents,
  incoming_call_answer: handleCallEvents,
  incoming_join_call: handleCallEvents,
  incoming_leave_call: handleCallEvents,
  incoming_webrtc_offer: handleWebRTCEvents,
  incoming_webrtc_answer: handleWebRTCEvents,
  incoming_ice_candidate: handleWebRTCEvents,
  incoming_camera_status: handleCallEvents,
  incoming_screen_status: handleCallEvents,
  incoming_mic_status: handleCallEvents,
  call_started: handleCallEvents
}

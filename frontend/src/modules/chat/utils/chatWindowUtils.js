function getChatContainer () {
  return document.querySelector('.chat-window')
}

/**
 * Проверяет, находится ли пользователь рядом с нижним краем чата.
 * @param {number} threshold — запас в пикселях
 */
export function isNearBottom (threshold = 120) {
  const el = getChatContainer()
  if (!el) return false
  return el.scrollTop + el.clientHeight >= el.scrollHeight - threshold
}

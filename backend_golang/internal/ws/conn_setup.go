package ws

import "time"

func (c *Client) prepareConn() {
	c.Conn.SetReadLimit(10 * 1024 * 1024)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}
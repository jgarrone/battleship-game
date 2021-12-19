package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type ClientConnection struct {
	mu       sync.Mutex
	conn     net.Conn
	loggedIn bool
	username string
}

func NewClientConnection(conn net.Conn) *ClientConnection {
	return &ClientConnection{
		conn: conn,
	}
}

func (c *ClientConnection) Username() string {
	return c.username
}

func (c *ClientConnection) IsLoggedIn() bool {
	return c.loggedIn
}

func (c *ClientConnection) LoggedInAs(username string) {
	c.mu.Lock()
	c.loggedIn = true
	c.username = username
	c.mu.Unlock()
}

func (c *ClientConnection) LoggedOut() {
	c.mu.Lock()
	c.loggedIn = false
	c.mu.Unlock()
}

func (c *ClientConnection) SendMessage(msg string) {
	_, _ = c.conn.Write([]byte(fmt.Sprintf("%s\n", msg)))
}

func (c *ClientConnection) ReceiveMessage() (string, error) {
	raw, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(raw), nil
}

func (c *ClientConnection) Close() error {
	return c.conn.Close()
}

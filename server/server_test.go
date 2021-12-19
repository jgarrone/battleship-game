package server

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"testing"
	"time"
)

func TestServer_Run(t *testing.T) {
	// Seed random generator with current time.
	now := time.Now().UTC().UnixNano()
	rand.Seed(now)

	servAddr := "localhost:8888"
	sv, err := NewServer(servAddr, "random")
	if err != nil {
		t.Fatalf("error initializing the server: %v", err)
	}

	go func() { _ = sv.Run() }()

	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		t.Fatalf(fmt.Sprintf("ResolveTCPAddr failed: %v", err.Error()))
	}

	quit := make(chan error)

	clients := map[string]*net.TCPConn{}
	for _, username := range []string{"user-test-1", "user-test-2", "user-test-3"} {
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			t.Fatalf(fmt.Sprintf("Dial failed: %v", err.Error()))
		}

		clients[username] = conn

		go handleConnection(conn, quit)

		_, err = conn.Write([]byte(fmt.Sprintf("login test-user-%d\n", rand.Intn(9999))))
		if err != nil {
			t.Fatalf(fmt.Sprintf("Write to server failed: %v", err.Error()))
		}
	}

	attackCount := 0
	for {
		select {
		case err := <-quit:
			if err != nil {
				t.Fatalf("error attacking: %v", err)
			}
			fmt.Printf("Game won after %d random attacks\n", attackCount)
			return
		default:
			for _, conn := range clients {
				if err := doRandomAttack(conn); err != nil {
					t.Fatalf(fmt.Sprintf("Write to server failed: %v", err.Error()))
				}
				attackCount += 1
			}
		}
	}

}

func doRandomAttack(conn *net.TCPConn) error {
	attackStr := fmt.Sprintf("attack %d %d\n", rand.Intn(BoardLengthX), rand.Intn(BoardLengthY))
	_, err := conn.Write([]byte(attackStr))
	return err
}

func handleConnection(conn *net.TCPConn, ch chan error) {
	defer func() { _ = conn.Close() }()
	for {
		raw, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			ch <- err
			return
		}
		msg := strings.TrimSpace(raw)
		if msg == "Game won" {
			ch <- nil
			return
		}
	}
}

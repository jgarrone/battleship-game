package server

import (
	"fmt"
	"github.com/jgarrone/battleship-game/server/client"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/jgarrone/battleship-game/server/board"
	"github.com/jgarrone/battleship-game/server/enum"
)

// Board size is fixed but easily parametrised if needed.
const (
	BoardLengthX = 10
	BoardLengthY = 10
)

type Attack struct {
	At       *board.Cell
	Username string
}

type CommandHandler func(conn *client.ClientConnection, args ...string) error

type BattleshipServer struct {
	address     string
	attackCh    chan Attack
	broadcastCh chan string
	board       board.BattleshipBoard
	clients     map[string]*client.ClientConnection
	handlers    map[enum.ClientCommand]CommandHandler
	muClients   sync.Mutex
}

func NewServer(address, boardGenStrategy string) (*BattleshipServer, error) {
	var cellSelector board.CellSelector
	switch boardGenStrategy {
	case "fixed":
		cellSelector = board.NewDummySelector()
	case "random":
		cellSelector = board.NewRandomSelector()
	default:
		return nil, fmt.Errorf("unknown strategy for board generation %q", boardGenStrategy)
	}

	fmt.Printf("Using %s strategy\n", boardGenStrategy)

	battleBoard, err := board.NewBattleshipBoard(BoardLengthX, BoardLengthY, cellSelector)
	if err != nil {
		return nil, fmt.Errorf("error initializing board: %v", err)
	}

	return &BattleshipServer{
		address: address,
		board:   battleBoard,
		clients: make(map[string]*client.ClientConnection),
	}, nil
}

func (s *BattleshipServer) Run() error {
	// Add new handlers here.
	s.handlers = map[enum.ClientCommand]CommandHandler{
		enum.ClientCommandAttack: s.handleAttack,
		enum.ClientCommandLogin:  s.handleLogin,
		enum.ClientCommandLogout: s.handleLogout,
	}

	// Listen for incoming connections.
	l, err := net.Listen("tcp", s.address)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer func() { _ = l.Close() }()

	fmt.Printf("Listening at %s\n", s.address)

	s.attackCh = make(chan Attack)
	s.broadcastCh = make(chan string)
	go s.doAttacks()
	go s.doBroadcasts()

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connection in a new goroutine.
		clientConn := client.NewClientConnection(conn)
		go s.handleConnection(clientConn)
	}

	//TODO: better handle of interruptions.

	return nil
}

func (s *BattleshipServer) doAttacks() {
	for {
		attack := <-s.attackCh
		s.attack(attack.Username, attack.At)
	}
}

func (s *BattleshipServer) doBroadcasts() {
	for {
		msg := <-s.broadcastCh
		s.broadcast(msg)
	}
}

func (s *BattleshipServer) enqueueAttack(attack Attack) {
	s.attackCh <- attack
}

func (s *BattleshipServer) enqueueBroadcast(msg string) {
	s.broadcastCh <- msg
}

func (s *BattleshipServer) attack(username string, cell *board.Cell) {
	x, y := cell.X, cell.Y
	broadcastMsg := ""
	outcome := s.board.Attack(cell)
	switch outcome {
	case enum.AttackOutcomeHit:
		broadcastMsg = fmt.Sprintf("%s: Hit %d %d", username, x, y)
	case enum.AttackOutcomeMiss:
		broadcastMsg = fmt.Sprintf("%s: Miss %d %d", username, x, y)
	case enum.AttackOutcomeAlreadyHit:
		broadcastMsg = fmt.Sprintf("%s: Already hit %d %d", username, x, y)
	case enum.AttackOutcomeHitAndWin:
		broadcastMsg = "Game won"
		s.board.Restart()
	}

	if broadcastMsg != "" {
		s.enqueueBroadcast(broadcastMsg)
	}

	fmt.Println(broadcastMsg)
}

func (s *BattleshipServer) broadcast(msg string) {
	for _, conn := range s.clients {
		if !conn.IsLoggedIn() {
			continue
		}
		conn.SendMessage(msg)
	}
}

func (s *BattleshipServer) handleConnection(conn *client.ClientConnection) {
	for {
		msg, err := conn.ReceiveMessage()
		if err != nil {
			break
		}
		if msg == "" {
			continue
		}
		msgParts := strings.Split(msg, " ")
		msgCmd, msgArgs := msgParts[0], msgParts[1:]

		if handler, exists := s.handlers[enum.ClientCommand(msgCmd)]; exists {
			if err := handler(conn, msgArgs...); err != nil {
				conn.SendMessage(fmt.Sprintf("Error executing %s command: %v", msgCmd, err))
			}
			continue
		}

		conn.SendMessage(fmt.Sprintf("Unknown command %q.", msgCmd))
	}

	_ = conn.Close()

}

func (s *BattleshipServer) handleAttack(conn *client.ClientConnection, args ...string) error {
	if !conn.IsLoggedIn() {
		return fmt.Errorf("you must login to be able to attack")
	}

	if len(args) != 2 {
		return fmt.Errorf("expected two arguments, got %d", len(args))
	}

	x, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("error parsing X coordinate: %v", err)
	}
	y, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("error parsing Y coordinate: %v", err)
	}

	pos := &board.Cell{
		X: x,
		Y: y,
	}

	if !s.board.Exists(pos) {
		conn.SendMessage(fmt.Sprintf("Invalid: %d %d", x, y))
		return nil
	}

	s.enqueueAttack(Attack{
		At:       pos,
		Username: conn.Username(),
	})

	return nil
}

func (s *BattleshipServer) handleLogin(conn *client.ClientConnection, args ...string) error {
	if conn.IsLoggedIn() {
		return fmt.Errorf("already logged in")
	}

	if len(args) != 1 {
		return fmt.Errorf("expected one argument, got %d", len(args))

	}
	username := args[0]

	s.muClients.Lock()
	defer s.muClients.Unlock()

	if _, exists := s.clients[username]; exists {
		return fmt.Errorf("player named %q already exists", username)
	}

	s.clients[username] = conn
	conn.LoggedInAs(username)
	s.enqueueBroadcast(fmt.Sprintf("%s: Logged in", username))

	fmt.Printf("Player %q logged in\n", username)

	return nil
}

func (s *BattleshipServer) handleLogout(conn *client.ClientConnection, args ...string) error {
	if !conn.IsLoggedIn() {
		return fmt.Errorf("not logged in")
	}

	if len(args) != 0 {
		return fmt.Errorf("expected zero arguments, got %d", len(args))
	}

	s.muClients.Lock()
	defer s.muClients.Unlock()

	username := conn.Username()

	if _, exists := s.clients[username]; !exists {
		return fmt.Errorf("player %q does not exist", username)
	}

	delete(s.clients, username)
	conn.LoggedOut()
	s.enqueueBroadcast(fmt.Sprintf("%s: Logged out", username))

	fmt.Printf("Player %q logged out\n", username)

	return nil
}

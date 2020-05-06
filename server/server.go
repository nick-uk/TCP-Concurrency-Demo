package server

import (
	"log"
	"net"
	"sync"
	"time"
)

var port = ":8081"
var debug = false

// timeout in secs
var timeout = time.Second * 15

// PongMsg is a data struct for pong messages with error field
type PongMsg struct {
	Msg       string `json:"msg"`
	ErrorCode int    `json:"error"`
}

// NotifyMsg is the data struct to send when a user status is change
type NotifyMsg struct {
	Online bool `json:"online"`
}

// CliRequest is the expected client request data struct
type CliRequest struct {
	UserID  int   `json:"user_id"`
	Friends []int `json:"friends"`
}

type userData struct {
	Friends []int
	Online  bool
}
type usersTableType struct {
	mutex sync.RWMutex
	users map[int]userData
}

// UsersTable is an in memory users table with mutex locking to avoid race conditions for parallel requests.
// For demo/tests purposes only. You can't use this for distributed production environments.
var UsersTable usersTableType

func init() {
	log.Println("Initiate [server]")
	UsersTable.users = make(map[int]userData)
}

// StartServer ...
func StartServer() {
	log.Println("Welcome to TCP server. Accept connections on:", port)
	// Listen TCP
	listenerTCP, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}
	// Listen UDP
	// listenUDP, err := net.ListenPacket("udp", port)
	// go handleUDPconn(listenUDP)

	for {
		conn, err := listenerTCP.Accept()
		if err != nil {
			log.Panicf("ERROR: failed to accept listener: %v", err)
		}
		if debug {
			log.Printf("Accepted connection %s -> %s\n", conn.RemoteAddr().String(), conn.LocalAddr().String())
		}
		go handleTCPconn(conn)
	}
}

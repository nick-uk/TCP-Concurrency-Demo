package server

import (
	"errors"
	"log"
	"net"
	"time"
)

func handleTCPconn(tcpconn net.Conn) {

	recvbuf := make([]byte, 256)
	// this connections is owned by a userid
	userid := 0
	var friends []int
	onlineNoticed := false
	offlineNoticed := false
	var err *PongMsg

	// I'm expecting clients to ping server every 1 sec with their IDs and their friends IDs.
	// When there is no friend status change, server will send a generic "pong" message back to the clients, otherwise
	// server will send one {"online": bool} messages if a friend joins/exits the app.
	// The issue here is the {"online": bool} it's too generic so this implementation will send only one online status message
	// back to the user even if more than one of his friends are online at the same monent.
	for {
		onlinestatussend := false

		// SetReadDeadline could return an error in some rare cases
		tcpconn.SetReadDeadline(time.Now().Add(timeout))
		// Read client packet. If there is no error returns the userid, the friends from the packet
		// and update the map table.
		userid, friends, _, err = readClientMsg(tcpconn, nil, userid, recvbuf)
		if err != nil {
			if debug {
				if err.Msg != "[server] EOF" {
					log.Println(err.Msg)
				}
			}
			break
		}

		// Check the online status of friends from the map //
		for friendindex := range friends {
			if friends[friendindex] == userid {
				// You can't be friend with yourself Error. Send TCP feedback.
				errMsg := AppErr(errors.New("You can't be friend with yourself"), 6)
				sendAppErr(errMsg, tcpconn, nil, nil)
				return
			}
			UsersTable.mutex.Lock()
			userdata, friendexist := UsersTable.users[friends[friendindex]]
			UsersTable.mutex.Unlock()
			if friendexist {
				if userdata.Online == true && !onlineNoticed {
					// A friend of this user is online
					if err := notifyOnlineStatus(tcpconn, nil, nil, true); err != nil {
						return
					}
					onlineNoticed = true
					offlineNoticed = false
					onlinestatussend = true
					// Here we're breaking the loop and prevent a "spamming" with online messages in case
					// user has a lot of friends.
					break
				} else if userdata.Online == false && !offlineNoticed {
					// A friend of this user is offline
					if err := notifyOnlineStatus(tcpconn, nil, nil, false); err != nil {
						return
					}
					onlineNoticed = false
					offlineNoticed = true
					onlinestatussend = true
					// Here we're breaking the loop and prevent a "spamming" with online messages in case
					// user has a lot of friends.
					break
				}
			}
		} // for friends

		if !onlinestatussend {
			// SetWriteDeadline could return an error in some rare cases
			tcpconn.SetWriteDeadline(time.Now().Add(timeout))
			// When there is no friend online status change, server will send a generic "pong" message back to the clients
			if err := writePong(tcpconn, nil, nil); err != nil {
				if debug {
					if err.Msg != "[server] EOF" {
						log.Println(err.Msg)
					}
				}
				break
			}
		}
	} // main TCP loop

	defer func() {
		// Connection closed, timeout or other connection error
		log.Println("=> [server] Userid:", userid, "exit")
		if userid != 0 {
			UsersTable.mutex.Lock()
			currentUserState := UsersTable.users[userid]
			currentUserState.Online = false
			UsersTable.users[userid] = currentUserState
			UsersTable.mutex.Unlock()
		} else {
			log.Println("=> [server] Client is dead before known userid")
		}
		tcpconn.Close()
	}()
}

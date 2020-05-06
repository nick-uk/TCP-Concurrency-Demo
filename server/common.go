package server

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"
)

func readClientMsg(tcpconn net.Conn, udpconn net.PacketConn, userid int, recvbuf []byte) (int, []int, net.Addr, *PongMsg) {

	thisuserid := userid
	var friends []int
	var readLen int
	var err error
	var logtag string
	var udpaddr net.Addr

	if tcpconn != nil {
		logtag = "TCP"
		readLen, err = tcpconn.Read(recvbuf)
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				// Client timeout
				return thisuserid, friends, nil, AppErr(err, 1)
			}
			// Unknown error
			return thisuserid, friends, nil, AppErr(err, 7)
		}
	} else {
		logtag = "UDP"
		readLen, udpaddr, err = udpconn.ReadFrom(recvbuf)
		if err != nil {
			// Unknown error
			return thisuserid, friends, udpaddr, AppErr(errors.New("UDP Read: "+err.Error()), 7)
		}
	}

	//// Proc request ////
	if debug {
		log.Println("["+logtag+"server] Recv:", string(recvbuf[:readLen]), ", Len:", readLen)
	}

	clientReq := CliRequest{}
	if err := json.Unmarshal(recvbuf[:readLen], &clientReq); err != nil {
		// JSON decoding Error, TCP feedback.
		return thisuserid, friends, udpaddr, sendAppErr(AppErr(err, 3), tcpconn, udpconn, udpaddr)
	}
	if clientReq.UserID == 0 {
		// invalid userID Error, TCP feedback.
		return thisuserid, friends, udpaddr, sendAppErr(AppErr(err, 5), tcpconn, udpconn, udpaddr)
	}
	thisuserid = clientReq.UserID
	friends = clientReq.Friends

	// Update user online status
	UsersTable.mutex.Lock()
	defer UsersTable.mutex.Unlock()

	// print joins messages
	if _, exists := UsersTable.users[thisuserid]; exists {
		if !UsersTable.users[thisuserid].Online {
			log.Println("=> ["+logtag+"server] Userid:", thisuserid, "joins")
		}
	} else {
		log.Println("=> ["+logtag+"server] New Userid:", thisuserid, "joins")
	}

	UsersTable.users[thisuserid] = userData{
		Online:  true,
		Friends: clientReq.Friends,
	}
	return thisuserid, friends, udpaddr, nil
}

// Create & send a pong response
func writePong(tcpconn net.Conn, udpconn net.PacketConn, udpaddr net.Addr) *PongMsg {

	pongResp := PongMsg{
		Msg:       "are you still there?",
		ErrorCode: 0,
	}
	respBytes, _ := json.Marshal(pongResp)

	if tcpconn != nil {
		if _, err := tcpconn.Write(respBytes); err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				// Client timeout
				return AppErr(err, 1)
			}
			// Unknown error
			return AppErr(err, 7)
		}
	} else {
		_, err := udpconn.WriteTo(respBytes, udpaddr)
		if err != nil {
			// log.Println("UDP Write:", err, udpaddr.String(), "client is gone?")
			// Unknown error
			return AppErr(err, 7)
		}
	}
	return nil
}

// Create & send a online/offline response
func notifyOnlineStatus(tcpconn net.Conn, udpconn net.PacketConn, udpaddr net.Addr, online bool) *PongMsg {

	notifResp := NotifyMsg{
		Online: online,
	}
	respBytes, _ := json.Marshal(notifResp)

	if tcpconn != nil {
		if _, err := tcpconn.Write(respBytes); err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				// Client timeout
				return AppErr(err, 1)
			}
			// Unknown error
			return AppErr(err, 7)
		}
	} else {
		_, err := udpconn.WriteTo(respBytes, udpaddr)
		if err != nil {
			// log.Println("UDP Write:", err, addr.String(), "client is gone?")
			// Unknown error
			return AppErr(err, 7)
		}
	}
	return nil
}

// AppErr returns a formated error in a PongMsg struct
func AppErr(err error, code int) *PongMsg {
	resp := PongMsg{
		Msg:       "[server] " + err.Error(),
		ErrorCode: code,
	}
	return &resp
}

// TCP error feedback
func sendAppErr(errMsg *PongMsg, tcpconn net.Conn, udpconn net.PacketConn, udpaddr net.Addr) *PongMsg {
	errMsgBytes, err := json.Marshal(errMsg)
	if err != nil {
		// JSON encoding Error
		return AppErr(err, 4)
	}
	if tcpconn != nil {
		tcpconn.SetWriteDeadline(time.Now().Add(timeout))
		if _, err := tcpconn.Write(errMsgBytes); err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				// Client timeout
				return AppErr(err, 1)
			}
			// Unknown error
			return AppErr(err, 7)
		}
	} else {
		_, err := udpconn.WriteTo(errMsgBytes, udpaddr)
		if err != nil {
			// log.Println("UDP Write:", err, addr.String(), "client is gone?")
			// Unknown error
			return AppErr(err, 7)
		}
	}
	return errMsg
}

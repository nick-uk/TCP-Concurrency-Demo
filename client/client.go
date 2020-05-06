package myclient

import (
	"encoding/json"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// Definition of possible errors
var mErrors = map[int]string{
	0: "No Error",
	1: "Client timeout",
	2: "Unknown Client timeout",
	3: "JSON decoding Error",
	4: "JSON encoding Error",
	5: "Invalid UserID",
	6: "You can't be friend with yourself",
	7: "Unknown Error",
	8: "Server Connection Error",
}

var timeout = time.Second * 15
var debug = false

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
	UserID  int   `json:"user_id,omitempty"`
	Friends []int `json:"friends,omitempty"`
}

// TCPClient accepts max alive client time, hoin delay time and returns total number of offlines / onlines events.
func TCPClient(targetHost string, userid int, friends []int, joindelay, stayalive time.Duration) (int, int, *PongMsg) {

	recvBuf := make([]byte, 256)
	logtag := "[user" + strconv.Itoa(userid) + "] "
	totalofflines := 0
	totalonlines := 0
	pong := &PongMsg{}

	time.Sleep(joindelay)

	client, err := net.DialTimeout("tcp", targetHost, time.Second*15)
	if err != nil {
		// Server Connection Error
		return totalofflines, totalonlines, &PongMsg{Msg: logtag + err.Error(), ErrorCode: 8}
	}

	// Prepare request
	request := CliRequest{
		UserID:  userid,
		Friends: friends,
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		// JSON encoding error
		return totalofflines, totalonlines, &PongMsg{Msg: logtag + err.Error(), ErrorCode: 4}
	}

	go func() {
		for {
			pong = &PongMsg{Msg: "", ErrorCode: 0}

			// Send request
			client.SetWriteDeadline(time.Now().Add(timeout))
			n, err := client.Write(requestBytes)
			if err != nil {
				log.Println(logtag, err)
				return
			}

			// Read request
			client.SetReadDeadline(time.Now().Add(timeout))
			n, err = client.Read(recvBuf)
			if err != nil {
				log.Println(logtag, err)
				return
			}
			if debug {
				log.Println(logtag, "Recv:", string(recvBuf[:n]))
			}

			if strings.Contains(string(recvBuf[:n]), "msg") {
				// is a pong
				if err := json.Unmarshal(recvBuf[:n], pong); err != nil {
					log.Println(logtag + "Unmarshal: " + err.Error() + " (" + string(recvBuf[:n]) + ")")
					return
				}
				if pong.ErrorCode != 0 {
					return
				}
			} else {
				friendStatus := NotifyMsg{}
				if err := json.Unmarshal(recvBuf[:n], &friendStatus); err != nil {
					log.Println(logtag + "Unmarshal: " + err.Error())
					return
				}
				if friendStatus.Online {
					totalonlines++
				} else {
					totalofflines++
				}
			}
			// jres := make(map[string]interface{})
			// merr := json.Unmarshal(recvBuf[:n], &jres)

			time.Sleep(time.Second * 1)
		}
	}()

	select {
	case <-time.After(stayalive):
		client.Close()
	}

	return totalofflines, totalonlines, pong
}

package myclient

import (
	"lockwood/server"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestTCPClient(t *testing.T) {
	type args struct {
		targetHost string
		userid     int
		friends    []int
		joindelay  time.Duration
		stayalive  time.Duration
	}
	tests := []struct {
		name              string
		args              args
		wanttotalofflines int
		wanttotalonlines  int
		wantthiserr       *PongMsg
	}{
		{"Test server connection error handling", args{
			targetHost: "127.0.0.1:8082",
			userid:     1,
			friends:    []int{2, 3, 4}, joindelay: 0, stayalive: time.Second * 1,
		}, 0, 0, &PongMsg{ErrorCode: 8}},
		{"user1 joins for 5 secs", args{
			targetHost: "127.0.0.1:8081",
			userid:     1,
			friends:    []int{2, 3, 4}, joindelay: 0, stayalive: time.Second * 5,
		}, 2, 2, &PongMsg{ErrorCode: 0}},
		{"user5 joins with wrong payload", args{
			targetHost: "127.0.0.1:8081",
			userid:     5,
			friends:    []int{2, 3, 4, 5}, joindelay: time.Millisecond * 0, stayalive: time.Second * 1,
		}, 0, 0, &PongMsg{ErrorCode: 6}},
		{"user2 joins for 1 sec", args{
			targetHost: "127.0.0.1:8081",
			userid:     2,
			friends:    []int{1, 3, 4}, joindelay: time.Millisecond * 300, stayalive: time.Second * 1,
		}, 0, 1, &PongMsg{ErrorCode: 0}},
		{"user3 joins for 2 secs", args{
			targetHost: "127.0.0.1:8081",
			userid:     3,
			friends:    []int{1, 4}, joindelay: time.Millisecond * 1300, stayalive: time.Second * 2,
		}, 0, 1, &PongMsg{ErrorCode: 0}},
	}

	go server.StartServer()
	time.Sleep(time.Second * 1)
	var wg sync.WaitGroup

	for _, tt := range tests {
		wg.Add(1)

		go func(t *testing.T, tt struct {
			name              string
			args              args
			wanttotalofflines int
			wanttotalonlines  int
			wantthiserr       *PongMsg
		}, wg *sync.WaitGroup) {

			defer wg.Done()

			t.Run(tt.name, func(t *testing.T) {
				totalofflines, totalonlines, err := TCPClient(tt.args.targetHost, tt.args.userid, tt.args.friends,
					tt.args.joindelay, tt.args.stayalive)

				if err.ErrorCode != tt.wantthiserr.ErrorCode {
					t.Errorf("TCPClient() ERROR = %+v, but I want ErrorCode %+v", err, tt.wantthiserr.ErrorCode)
				}
				if totalofflines != tt.wanttotalofflines {
					t.Errorf("TCPClient() totalofflines = %v, but I want %v", totalofflines, tt.wanttotalofflines)
				}
				if totalonlines != tt.wanttotalonlines {
					t.Errorf("TCPClient() totalonlines = %v, but I want %v", totalonlines, tt.wanttotalonlines)
				}
				log.Println("=> User"+strconv.Itoa(tt.args.userid)+" client quit with the expeted Total Online messages:", totalonlines,
					"and Total Offline messages:", totalofflines, "("+err.Msg+")")
			})
		}(t, tt, &wg)
	} // tests

	wg.Wait()
}

# TCP-Concurrency-Demo
TCP server and a client for tests that demonstrates a very basic online/offline clients management

### What is this repository for?
This project demonstrates a low-level TCP server/client communication and unit testing with in-memory Key-Value storage. Each client sends a JSON payload packet with ({"user_id": 1, "friends": [2, 3, 4]}) and while is online is being notified if a friend is changing online/offline status.

### How to test it?
```go test -v ./client/```

### Expecting output
```bash
2020/05/06 13:38:59 Initiate [server]
=== RUN   TestTCPClient
2020/05/06 13:38:59 Welcome to TCP server. Accept connections on: :8081
=== RUN   TestTCPClient/user3_joins_for_2_secs
=== RUN   TestTCPClient/user1_joins_for_5_secs
=== RUN   TestTCPClient/user5_joins_with_wrong_payload
=== RUN   TestTCPClient/user2_joins_for_1_sec
=== RUN   TestTCPClient/Test_server_connection_error_handling
2020/05/06 13:39:00 => User1 client quit with the expeted Total Online messages: 0 and Total Offline messages: 0 ([user1] dial tcp 127.0.0.1:8082: connect: connection refused)
2020/05/06 13:39:00 => [TCPserver] New Userid: 1 joins
2020/05/06 13:39:00 => [TCPserver] New Userid: 5 joins
2020/05/06 13:39:01 => [TCPserver] New Userid: 2 joins
2020/05/06 13:39:01 => User5 client quit with the expeted Total Online messages: 0 and Total Offline messages: 0 ([server] You can't be friend with yourself)
2020/05/06 13:39:02 => User2 client quit with the expeted Total Online messages: 1 and Total Offline messages: 0 ()
2020/05/06 13:39:02 => [TCPserver] New Userid: 3 joins
2020/05/06 13:39:02 => [server] Userid: 2 exit
2020/05/06 13:39:02 [user2]  write tcp 127.0.0.1:55214->127.0.0.1:8081: use of closed network connection
2020/05/06 13:39:04 => User3 client quit with the expeted Total Online messages: 1 and Total Offline messages: 0 (are you still there?)
2020/05/06 13:39:04 => [server] Userid: 3 exit
2020/05/06 13:39:04 [user3]  write tcp 127.0.0.1:55216->127.0.0.1:8081: use of closed network connection
2020/05/06 13:39:05 => User1 client quit with the expeted Total Online messages: 2 and Total Offline messages: 2 ()
2020/05/06 13:39:05 => [server] Userid: 1 exit
--- PASS: TestTCPClient (6.01s)
    --- PASS: TestTCPClient/Test_server_connection_error_handling (0.00s)
    --- PASS: TestTCPClient/user5_joins_with_wrong_payload (1.00s)
    --- PASS: TestTCPClient/user2_joins_for_1_sec (1.30s)
    --- PASS: TestTCPClient/user3_joins_for_2_secs (3.30s)
    --- PASS: TestTCPClient/user1_joins_for_5_secs (5.00s)
PASS
ok  	lockwood/client	6.495s
```


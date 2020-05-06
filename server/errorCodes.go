package server

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

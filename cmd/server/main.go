package main

import "github.com/pmouraguedes/battleship/internal/server"

func main() {
	s := server.NewServer(":8000")
	s.Start()
}

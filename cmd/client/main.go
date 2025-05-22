package main

import (
	"github.com/pmouraguedes/battleship/internal/client"
)

func main() {
	// box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	// if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
	// 	panic(err)
	// }
	// if true {
	// 	return
	// }

	c, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	println("Client started", c)
	if err := c.Run(); err != nil {
		panic(err)
	}
}

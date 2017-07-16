package main

import (
	"fmt"

	maxim "github.com/TeaMeow/Maxim"
)

type HelloWorld struct {
	message string
}

func main() {
	e := maxim.Default()
	chunkHandler := maxim.ChunkHandler()
	e.OnFile("Photo", chunkHandler, func(c *maxim.Context) {
		fmt.Printf("%+v", c.File)
	})
	e.Run(":5000")
}

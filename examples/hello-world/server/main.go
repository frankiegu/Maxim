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
	e.On("Hello", func(c *maxim.Context) {
		var h HelloWorld
		err := c.Bind(&h)
		if err != nil {
			panic(err)
		}
		fmt.Printf("來自客戶端的招呼：%s", h.message)
		c.Respond(maxim.StatusOK, maxim.H{
			"message": "hello, world!",
		})
	})
	e.Run(":5000")
}

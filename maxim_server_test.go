package maxim

import "testing"

func TestMain(t *testing.T) {
	e := Default()
	e.On("Ping", func(c *Context) {
		c.Respond(StatusOK, H{
			"pong": "Hello, world!",
		})
	})
	e.Run(":5000")
}

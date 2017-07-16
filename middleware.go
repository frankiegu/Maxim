package maxim

func Logger() HandlerFunc {
	return func(c *Context) {
		c.Next()
	}
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		c.Next()
	}
}

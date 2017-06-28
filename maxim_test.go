package maxim

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	m := New()
	m.Function("SayHello", func(c *Context) {
		var data map[string]string
		if c.Bind(&data) == nil {
			// 輸出：world!
			fmt.Println(data["hello"])
			// 將訊息傳遞回去給使用者。
			c.Send(StatusOK, H{
				"foo": "bar!",
			})
		}
	})
	m.Run(":8080")
}

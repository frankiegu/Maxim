# Maxim

Maxim 是一個基於 Golang

```go
package main

import "github.com/TeaMeow/Maxim"

func main() {
    e := maxim.Default()
    e.On("CreateUser", func(c *maxim.Context) {
        c.Respond(maxim.StatusOK, maxim.H{
            "hello": "world",
        })
    })
    e.Run(":5000")
}
```

## 基本內容

```go
e.On("GetUser", func(c *maxim.Context) {
    c.Respond(maxim.StatusOK, maxim.H{
        "username": "YamiOdymel",
        "birthday": "1998-07-13",
    })
})
```

### 綁定模型

```go
e.On("GetSession", func(c *maxim.Context) {
    var form LoginForm
    if err := c.Bind(&form); err == nil {
        if form.Username == "YamiOdymel" && form.Password == "test" {
            c.Respond(maxim.StatusOK, maxim.H{
                "message": "Logged in successfully!",
            })
        } else {
            c.Respond(maxim.StatusNotAuthorized, maxim.H{
                 "message": "The password is incorrect!",
            })
        }
    }
})
```

### 回應模型

```go
e.On("GetBook", func(c *maxim.Context) {
    var book struct {
        Title       string `json:"t"`
        Description string `json:"d"`
    }
    book.Title = "世界上最好的語言：PHP"
    book.Description = "這本書將帶領你理解為什麼 PHP 能夠領先任何程式語言十多年。"
    // 輸出：{"t": "世界上最...", "d": "這本書將帶領你理..."}
    c.Respond(maxim.StatusOK, book)
})
```

### 回應其他人

```go
e.On("CreateMessage", func(c *maxim.Context) {
    c.RespondOthers(maxim.StatusOK, maxim.H{
        "message": "Pong!",
    })
})
```

## 檔案接收

```go
e.OnFile("Avatar", maxim.DefaultChunker, func(c *maxim.File) {

})
```

## 中間件

```go
e.On("CreateUser", myMiddleware, anotherMiddleware, func(c *maxim.Context) {

})
```

## 主動式回應

```go
func main() {
    e := maxim.Default()
    go func() {
        for {
            <- time.After(1 * time.Second)
            e.Respond(maxim.StatusOK, maxim.H{
                "time": time.Now(),
            })
        }
    }()
    e.Run(":5000")
}
```
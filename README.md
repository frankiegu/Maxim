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

### 主動式回應

```go
func main() {
    e := maxim.Default()
    go func() {
        for {
            <-time.After(1 * time.Second)
            e.Respond(maxim.StatusOK, maxim.H{
                "time": time.Now(),
            })
        }
    }()
    e.Run(":5000")
}
```

## 檔案接收

在 Maxim 中，上傳檔案和送出資料是分開的，這意味著當你想要上傳帶有圖片的表單時，你需要先上傳圖片，接著取得已上傳圖片的以檔案編號的方式夾帶到另一個表單方可傳遞相關資訊。這在上傳大型檔案如影片時非常有用。

Maxim 會自動在上傳時將檔案切分成塊（基於客戶端區塊大小而定），上傳區塊之前不會和伺服器請求建立任何資訊，而是直接在客戶端建立不重複金鑰（非加密用），然後與二進制資料一同直接傳遞給伺服器並呼叫指定函式。伺服端能夠在區塊組合完畢之後進行縮圖、轉檔等工作，並且傳遞相關資訊回去給客戶端。

```
[Client]      [Server]
    |             |      客戶端：建立不重複金鑰，將檔案切分成塊。
    |------------>|      客戶端：傳送區塊 1／2 與金鑰。
    |<------------|      伺服器：完成 #1，請傳送下一塊。
    |------------>|      客戶端：傳送區塊 2／2 與金鑰。
    |             |      伺服器：組合所有區塊。
    |<------------|      伺服器：上傳程序完成，呼叫完成函式進行檔案處理（如：縮圖、轉檔）並回傳檔案資料。
    |             |      客戶端：取得檔案編號，存至新表單資料。
   ~~~~~~新請求~~~~~~~
    |------------>|      客戶端：將帶有檔案編號的表單資料傳遞至伺服器。
    v             v
```


```go
e.OnFile("Avatar", maxim.DefaultChunker, func(c *maxim.Context) {

})
```

## 中間件

```go
e.On("CreateUser", myMiddleware, anotherMiddleware, func(c *maxim.Context) {

})
```


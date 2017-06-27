# Maxim

Maxim 是一套前端與後端以 JSON 作為基礎並透過 MessagePack 壓縮且經由 WebSocket 的通訊方式，用以取代傳統的 REST。

## 使用範例

JavaScript

```js
import maxim from "maxim"

conn = maxim.open("wss://www.example.com/")
resp = await conn.execute("SayHello", {
    hello: "world!"
})
data = await resp.data()
// 輸出：bar!
console.log(data.foo)
```

Golang

```go
package main

import (
	"github.com/TeaMeow/Maxim"
)

func main() {
    m := maxim.New()
    m.Function("SayHello", func(c *maxim.Context) {
        var data map[string]string
        if c.Bind(&data) == nil {
            // 輸出：world!
            fmt.Println(data["hello"])
            // 將訊息傳遞回去給使用者。
            c.Send(maxim.H{
                "foo": "bar!"
            })
        }
    })
    m.Run()
}

```

## 內容結構

透過 WebSocket，傳遞一個最基本的 Maxim 請求格式如下。

* `fun`：欲呼叫的函式名稱。
* `met`：中繼資料。
* `dat`：傳遞的資料內容。
* `col`：請求欄位。

### 上傳

欲要上傳檔案，須先透過 Maxim 傳遞檔案資訊至遠端伺服器，並由遠端建立相對應的金鑰（辨識非加密用），接著透過此金鑰上傳片段二進制內容，再由遠端伺服器進行組合以達到可暫停式的上傳方式。

第一次的金鑰交換完畢之後，請透過下列結構上傳片段二進制內容。

* `fil`：檔案資料。
    * `key`：辨識金鑰。
    * `bin`：片段二進制內容。
    * `inf`：資訊。
        * `cur`：目前片段
        * `tol`：總計片段數

### 回傳

Maxim 的回傳內容帶有資料和錯誤資訊以方便進行除錯，當錯誤發生時 `dat` 應保持空白。

* `met`：中繼資料（如：版本、網路速率）。
* `dat`：成功資料內容。
* `cod`：本次成功、錯誤回傳碼。
* `err`：錯誤資料。

## 使用方式

Maxim 目前僅支援 Golang 與前端 JavaScript 進行連動，這意味著後端必須要是 Golang。前端的 JavaScript 連線方式十分簡易，透過 `open()` 之後會直接連線至遠端伺服器。

```js
import maxim from "maxim"

maxim.open("ws://www.example.com/")
maxim.disconnect() // 斷開連線
maxim.connect()    // 重新連線
```

### 建立請求

以 Maxim 在 JavaScript 向遠端 Golang 伺服器建立一個請求時需透過一連串的函式鍊。此函式會回傳 Promise，這令你能夠與 ASync/Await 一同搭配。

```js
maxim.meta({
    session: "AJjMC39xO1cpELfbGC8H4Os21G"
}).data({
    title: "十億人都驚呆的真相！",
    description: "到玉山看岩漿，你絕對不會相信 43 秒時我所看見的東西。"
}).send()
```

#### 檔案上傳

如果要上傳的是檔案，傳入一個 `FileReader` 至 `file()` 函式中。

```js
reader = new FileReader()
maxim.file(reader).send()
```

### 全域配置

Maxim 提供了 `setup()` 函式可進行全域設置，以省去每次發送時都還需要呼叫特定函式、配置的困擾。

```js
maxim.setup({
    // 每次 `send()` 執行時所會自動夾帶的中繼資料。
    meta: {
    },
    // 每次 `send()` 執行之前所會呼叫的函式，當此函式回傳的是 `false` 則停止繼續。
    beforeSend: () => {
        return true
    },
    // 每次 `send()` 執行完畢之後所會呼叫的函式。
    afterSend: () => {
    }
})
```

## 後端結構

Maxim 已有後端套件供 Golang 使用。

```go
import "github.com/TeaMeow/Maxim"

// 定義相對應的函式。
maxim.Function("GetUser", func (c *Maxim.Context) {
    c.Data(OK, Maxim.H{
        "Users": "Wow"
    }).Send()
})
maxim.Run(":8080")
```

```go
maxim.Function("GetUser", func (c *Maxim.Context) {
    c.Data(OK, Maxim.H{
        "Users": "Wow"
    }).Send()
})
```

###

```go
maxim.File(func (c *Maxim.Context) {

})
```

# Maxim [![GoDoc](https://godoc.org/github.com/TeaMeow/Maxim?status.svg)](https://godoc.org/github.com/TeaMeow/Maxim)

用以取代傳統 [RESTful API](https://zh.wikipedia.org/zh-tw/REST) 的前後端雙向溝通協定，基於 [WebSocket](https://developer.mozilla.org/zh-TW/docs/WebSockets/WebSockets_reference/WebSocket) 且由 [MessagePack](http://msgpack.org/) 對資料進行編碼、解碼令傳遞的資料更加輕量。同時支援檔案可暫停式上傳。

# 這是什麼？

馬克沁是一個由 [Golang](https://golang.org/) 所撰寫的前後端雙向溝通方式，這很適合用於聊天室、即時連線、單頁應用程式。目前僅支援前端 JavaScript 與後端 Golang 的溝通（若有時間將會加入其他語言的版本）。

* 可取代傳統 RESTful API。
* 以 MessagePack 編碼資料與達到更輕量的傳遞。
* 可雙向溝通與廣播以解決傳統僅單向溝通的障礙。
* 可供暫停式的全自動檔案分塊上傳方式。
* 支援中介軟體（Middleware）。
* 可使用 ES7 Async/Await 處理前端的呼叫。
* 友善的錯誤處理環境。

# 為什麼？

在部份時候傳統 RESTful API 確實能派上用場，但久而久之就會為了如何命名、遵循 REST 風格而產生困擾，例如網址必須是名詞而非動詞，登入必須是 `GET /token` 而非 `POST /login`，但使用 `GET` 傳遞機密資料極為不安全，因此只好更改為 `POST /session` 等，陷入如此地窘境。

且多個時候 RESTful API 都是單向，而非雙向。這造成了在即時連線、雙向互動上的溝通困擾，不得不使用 Comet、Long Polling 達到相關的要求，令 API 更加分散難以管理。

Maxim 試圖以 WebSocket 解決雙向溝通和多數重複無用的 HTTP 標頭導致浪費頻寬的問題，Maxim 同時亦支援透過 WebSocket 上傳檔案的功能。

# 單頻道與多頻道廣播？

WebSocket 簡單說就是開一個共享廣播令大家接收到相同的訊息，但這可不能發生在一般網頁溝通上（除非你想把機密資料傳去別人電腦裡），為此 Maxim 的廣播方式如同一般 RESTful API，僅會回傳給指定客戶端（Session），當然你也能夠透過 Maxim 將某個資料廣播給所有客戶端，或者除了某客戶端以外的所有人，這能令你在製作如聊天室廣播時更加得心應手。

# 檔案上傳？

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

# 效能如何？

# 索引

* [安裝方式](#安裝方式)
* [命名方式](#命名方式)
* [使用方式](#使用方式)
* [後端](#後端)
	* [監聽事件](#監聽事件)
	* [回應](#回應)
		* [回應模型](#回應模型)
		* [回應其他人](#回應其他人)
		* [主動式回應](#主動式回應)
	* [綁定資料](#綁定資料)
	* [呼叫](#呼叫)
		* [呼叫其他人](#呼叫其他人)
		* [主動式呼叫](#主動式呼叫)
	* [中介軟體](#中介軟體)
		* [單函式中介](#單函式中介)
		* [自造中介軟體](#自造中介軟體)
		* [中介軟體裡的 Goroutine](#中介軟體裡的-goroutine)
	* [檔案處理](#檔案處理)
		* [預設區塊處理函式](#預設區塊處理函式)
	* [中繼資料](#中繼資料)
* [前端](#前端)
	* [開啟連線](#開啟連線)
* [狀態碼](#狀態碼)

# 安裝方式

打開終端機並且透過 `go get` 安裝此套件即可。

```bash
$ go get github.com/TeaMeow/Maxim
```

# 命名方式

在 Maxim 中你可以依照你的喜好替函式命名，但如果你覺得不曉得從何起頭，這裡提供了一個比起傳統 RESTful API 還要更有彈性的命名規範。傳統的 API 中我們有 `GET`、`POST`、`PATCH`⋯等，但多數情況下那並不是很適用，例如你想要替某項東西按讚，變成了 `POST Like` 並不是很直覺。

比起使用 CRUD（Create、Read、Update、Delete），Maxim 更推薦使用 BGUCDS，他們的含義如下。

| 名詞    | 簡稱  | 說明                                     | 範例                       | 傳統範例              |
|--------|------|------------------------------------------|---------------------------|----------------------|
| Browse | 瀏覽  | 瀏覽資源的清單，例如影片列表、好友列表。        | BrowseUsers               | GET /users           |
| Get    | 取得  | 取得單個資源。                             | GetComment                | GET /comment/1       |
| Update | 更新  | 更新、編輯單個或多個資源。                   | UpdatePost                | PUT /post/1          |
| Create | 建立  | 新增、建立單個或多個資源。                   | CreatePhotos              | POST /photos         |
| Delete | 刪除  | 移除單個或多個資源。                        | DeleteFriend              | DELETE /friend/1     |
| Search | 搜尋  | 搜尋資源。                                | SearchPosts               | GET /posts/?count=30 |

如果你覺得這好像蠻難記的，那麼你就可以創造出屬於你自己的命名方式，因為這不是固定的。

# 使用方式

Maxim 的使用方式參考令人容易上手的 [Gin 網站框架](https://github.com/gin-gonic/gin)，令你撰寫 WebSocket 就像在撰寫普通的網站框架ㄧ樣容易。

###### Golang

```go
import "github.com/TeaMeow/Maxim"

func main() {
	engine := maxim.Default()
	engine.On("Ping", func(c *maxim.Context) {
		c.Respond(maxim.StatusOK, maxim.H{
			"message": "pong!",
		})
	})
	engine.Run(":5000")
}
```

###### JavaScript

```javascript
import "maxim"

var conn   = new Maxim("ws://localhost:5000/"),
    result = await conn.execute("Ping")

console.log(result.data().pong) // 輸出：pong!
```

## 後端

### 監聽事件

透過 `On` 監聽事件，這是最基本的用法，就像是在網站框架中定義一個 API 進入點一樣，這個範例定義了一個 `Ping` 函式可供客戶端呼叫。

```go
engine.On("Ping", func(c *maxim.Context) {
	// ...
})
```

### 回應

在 Maxim 裏，回應的方式和傳統的 RESTful API 一樣都帶有一個狀態碼，但這個狀態碼可以是自訂的，或者你也能使用 Maxim 中現有的狀態碼。最基本的回應方式就是透過 `maxim.H`（實際上就是 `map[string]interface{}`）。

```go
engine.On("Ping", func(c *maxim.Context) {
	c.Respond(maxim.StatusOK, maxim.H{
		"message": "pong!",
	})
})
// 回傳：{"message": "pong!"}
```

#### 回應模型

回應不一定要是一個 `map`，例如建構體也是可行的。你亦能在建構體中透過標籤指定回傳的鍵名。

```go
engine.On("GetBook", func(c *maxim.Context) {
	var book struct {
		Title       string `json:"t"`
		Description string `json:"d"`
	}
	book.Title = "世界上最好的語言：PHP"
	book.Description = "這本書將帶領你理解為什麼 PHP 能夠領先任何程式語言十多年。"

	c.Respond(maxim.StatusOK, book)
	// 輸出：{"t": "世界上最...", "d": "這本書將帶領你理..."}
})
```

#### 回應其他人

當在建立一個類似聊天室的服務時，我們不希望接收到自己剛才發送的訊息，但卻希望其他所有人接收到這則訊息，此時可以透過 `RespondOthers` 來將訊息傳遞給自己以外的所有人（這可能不適用於有負載平衡的情況）。

```go
engine.On("SendMessage", func(c *maxim.Context) {
	var m Message
	if err := c.Bind(&m); err != nil {
		c.RespondOthers(maxim.StatusOK, maxim.H{
			"message": m.Content,
		})
	}
})
```

#### 主動式回應

你可以主動式的直接向所有客戶端進行回應廣播，而不需要等到客戶端發送請求。這很適合用於即時聊天或者線上遊戲。

```go
engine.Respond(maxim.H{
	"message": "Hello, world!",
})
```

### 綁定資料

當接收到來自客戶端的資料時，可以透過 `Bind` 將其資料直接映射在本地端的特定建構體或 `map`。

```go
engine.On("Login", func(c *maxim.Context) {
	var u User
	if err := c.Bind(&u); err == nil {
		c.Respond(maxim.StatusOK, maxim.H{
			"message":  "已接收到使用者資料！",
			"username": u.Username,
			"password": u.Password,
		})
	}
})
```

### 呼叫

由於 Maxim 是雙向的，這意味著伺服端也能夠透過 `Execute` 函式呼叫客戶端的方式，但這可能會令程式結構變的複雜。

```go
// 建立一個檢查更新的函式。
engine.On("CheckUpdate", func(c *maxim.Context) {
	// 如果有更新的話⋯⋯。
	if hasUpdate {
		data := maxim.H{"version": "1.0.0-beta1"}
		// 就呼叫客戶端的 Update 函式來進行更新。
		c.Execute("Update", data, func(c *maxim.Context) {
			// ...處理執行客戶端 Update 所回傳的資訊...
		})
	}
})
```

#### 呼叫其他人

和「回應其他人」相同的功能，你可以呼叫除了這個客戶端以外的其他人函式。

```go
engine.On("SetName", func(c *maxim.Context) {
	c.ExecuteOthers("Update", data, func(c *maxim.Context) {
		// ...
	})
})
```

#### 主動式呼叫

你也能主動直接向所有客戶端呼叫指定函式，但這是向所有客戶端進行廣播，因此無法取得回傳的資訊。

```go
engine.Execute("SetColor", maxim.H{
	"color": "#00ADEA",
})
```

### 中介軟體

中介軟體讓你在執行 API 進入點函式之前安插一些中介函式，例如紀錄、監聽、身份驗證⋯等，如此一來你就能夠在使用者身份不對時，直接結束請求。透過 `Use` 可以安插全域中介軟體。

```go
engine.Use(maxim.Logger())
engine.Use(myMiddleware())
```

#### 單函式中介

有時候你只需要讓指定函式有中介軟體，而非所有函式，這麼做就行了。

```go
engine.On("Login", myMiddleware(), func(c *maxim.Context) {
	// ...
})
```

#### 自造中介軟體

除了使用 Maxim 內建的中介軟體外，你當然也能自己寫一個。

```go
// myMiddleware 是一個自己撰寫的中介軟體，會在接收的資料中安插一個新的 `Foo` 資料。
func myMiddleware() maxim.HandlerFunc {
	return func(c *maxim.Context) {
		// 在資料中安插自訂資料。
		c.Set("Foo", "Bar")
		// 呼叫下一個中介軟體，或者繼續。
		c.Next()
		// 之後還可以繼續執行程式。
		fmt.Println("myMiddleware 已執行！")
	}
}

engine.On("Login", myMiddleware(), func(c *maxim.Context) {
	fmt.Println(c.Get("Foo").(string))
})
// 輸出：Bar
//      myMiddleware 已執行！
```

#### 中介軟體裡的 Goroutine

如果你的中介軟體裡會用到 Goroutine，那麼你就不應該把原本的 `maxim.Context` 傳入 Goroutine 中。你必須透過 `Copy` 函式複製一份唯讀的 `Context` 避免資料起衝突。

```go
engine.On("Login", func(c *maxim.Context) {
	go myGoroutineFunc(c.Copy())
}, func(c *maxim.Context) {
	// ...
})
```

### 檔案處理

透過 `OnFile` 建立一個基於檔案處理的事件監聽器，和基本的 `On` 函式類似，但第二個參數是區塊處理函式。

```go
chunker := maxim.NewChunker()
engine.OnFile("Video", chunker, func(c *maxim.Context) {
	// 當上傳完畢後輸出檔案名稱。
	fmt.Println(c.File.Name)
})
```

#### 預設區塊處理函式

你不需要手動處理區塊分割的問題，因此你的檔案上傳處理也變得異常簡單。透過 Maxim 中的 `NewChunker` 新建一個內建的區塊處理器 `maxim.Chunker`，這個處理器會將區塊暫時擺放於 `/tmp` 中，並且最終組合檔案到指定位置。以此就能直接實作一個檔案處理函式。你也能自行建立一個區塊處理器來將區塊上傳至 Amazon S3 或者以其他方式處理。

```go
chunker := maxim.NewChunker(maxim.ChunkerOption{
	Path:    "/uploads",
	MaxSize: 5000000,
})
engine.OnFile("Photo", chunker, func(c *maxim.Context) {
	// ...
})
```

### 中繼資料

透過 `Metadata` 函式從接收到的資訊中取得客戶端所傳遞的中繼資料。

```go
engine.On("CreatePost", chunker, func(c *maxim.Context) {
	metadata := c.Metadata()
	fmt.Println(metadata["token"].(token))
})
```

## 前端

### 開啟連線

```javascript

```

## 狀態碼

在 Maxim 中有內建這些狀態碼，優於傳統 RESTful API 之處在於你可以自訂你自己的狀態碼。

| 狀態碼              | 說明                                                                      |
|---------------------|--------------------------------------------------------------------------|
| StatusOK            | 任何事情都很 Okay，如果刪除了早就不存在的事物，也可以算是 OK。                    |
| StatusError         | 內部錯誤發生，可能是非預期的錯誤。                                             |
| StatusProcessing    | 已成功傳送請求，但現在還不會完成，將會在背景中進行。                              |
| StatusFull          | 請求因為已滿而被拒絕，例如：好友清單已滿無法新增、聊天室人數已滿無法加入。            |
| StatusExists        | 已經有相同的事物存在而被拒絕。                                                |
| StatusInvalid       | 格式無效而被拒絕，通常是不符合表單驗證要求。                                    |
| StatusNotFound      | 找不到指定資源。                                                           |
| StatusNotAuthorized | 請求被拒。沒有驗證的身份，需要登入以進行驗證。                                  |
| StatusNoPermission  | 已驗證身份，但沒有權限提出此請求因而被拒。                                     |
| StatusNoChanges     | 這項請求沒有改變什麼事情，通常來說可以用 StatusOK 即可。                        |
| StatusUnimplemented | 此功能尚未實作完成，如果呼叫了一個不存在的函式即會回傳此狀態。                     |
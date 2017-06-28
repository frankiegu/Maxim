import msgpack from "msgpack-lite"

class Maxim
    constructor: (url) ->
        @connection           = new WebSocket url
        @connection.onmessage = @_messageHandler
        @meta     = {}
        @tid      = -1
        @col      = []
        @resolves = {}
        @rejects  = {}
        @options  =
            beforeExecute: -> true
            afterExecute:  ->
            onClose:       ->
            onError:       ->
            onMessage:     ->
            meta: {}

    # _messageHandler 負責處理接收到的任何訊息，並且 resolve 相對應的 Promise。
    _messageHandler: (event) ->
        # 如果接收到有工作編號的回應，則呼叫並解決相對應的 Promise。
        if event.data.tid isnt undefined
            @resolves[event.data.tid](event.data)

        # 呼叫使用者自訂的 onMessage 回呼函式。
        @options.onMessage(event)

    # _clean 會清除上次的 Maxim 傳遞內容。
    _clean: ->
        @meta = {}
        @col  = []

    # _dataResolver 會回傳一個新的 Data 類別用來輔助使用者取得 Maxim 所回傳的資料。
    _dataResolver: (data, func) ->
        class Data
            constructor: (@data, @func) ->
            meta:  -> if @data.met isnt undefined then @data.met else null
            data:  -> if @data.dat isnt undefined then @data.dat else null
            code:  -> @data.cod || @data.err?.cod
            func:  -> @func
            err:   -> if @data.err isnt undefined then @data.err else null
            isErr: -> @data.err is undefined
        new Data(data, func)

    # columns 會設置傳入的欄位，並且夾帶到資料中傳遞至遠端伺服器請求相關欄位。
    columns: ->
        @col = arguments

    # meta 會設置中繼資料並且夾帶於資料中傳遞至遠端伺服器。
    meta: (data) ->
        @meta = data

    # setup 用來設置這個 Maxim 連線的全域設定。
    setup: (options) ->
        @options = Object.assign(@options, options)

    # execute 呼叫遠端函式，並且將本地資料上傳至遠端執行。
    # 參數中的 data 通常是 Object 型態。
    execute: (func, data) ->
        that = this
        body =
            tid: @tid
            dat: data
            fun: func

        # 如果有中繼資料的話就設置。
        if Object.getOwnPropertyNames(@meta).length > 0 or Object.getOwnPropertyNames(@options.meta).length > 0
            body.met = Object.assign @meta, @options.meta
        # 如果有指定欄位請求就將欄位放入資料中。
        if @col.length > 0
            body.col = @col
        # 建立一個 Promise 並在稍後會傳以利於使用 Async/Await。
        promise = new Promise (resolve, reject) ->
            that.resolves[that.tid] = (value) -> resolve(that._dataResolver(value, func))

        # 透過 MessagePack 壓縮資料內容並且透過 WebSocket 傳遞資料。
        @connection.send msgpack.encode(data)
        # 清除本次暫存資料。
        @_clean()
        # 增加工作編號。
        @tid++
        
        return promise
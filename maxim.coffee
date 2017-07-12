# MaximData 是接收到的資料函式庫，用以轉換資料供開發者使用。
class MaximData
    constructor: (@event) ->

    data: ->
    error: ->
    metadata: ->

# Maxim 是主要的連線函式庫。
class Maxim
    constructor: (url, options) ->
        @connection       = new WebSocket(url)
        @columns          = ''
        @options          =
            metadata: {}
        @metadata         = {}
        @taskID           = 0
        @tasks            = {}
        @messageListeners = {}
        @errorListeners   = {}
        @openListeners    = {}
        @reopenListeners  = {}
        @closeListeners   = {}

        @addListener 'message', @_taskCaller

    # metadata 會設置單次性的指定中繼資料。
    metadata: (data) ->
        @metadata = data

    # _taskCaller 會監聽所有接收到的訊息，如果帶有工作編號則呼叫指定的工作 Resolve 函式。
    _taskCaller: (event) ->
        # 將接收到的資料以 MessagePack 解碼。
        data = msgpack.decode(event.data)
        # 如果該資料含有工作編號，就呼叫指定的工作 Resolve 函式。
        @task[data.tid](event) if data.tid isnt undefined

    # execute 會透過 WebSocket 呼叫遠端函式。
    execute: (func, data) ->
        # 集合所有的內容，稍後傳遞到遠端伺服器。
        payload =
            tid: @taskID
            fun: func
            col: @columns
            met: {
                @metadata...
                @options.metadata...
            }
            dat: data
        # 遞增工作編號供下次呼叫使用。
        @taskID++
        # 回傳一個 Promise 以方便使用 ES7 的 Async/Await。
        return new Promise (resolve, reject) ->
            # 保存這個 Promise 的 Resolve 到本地。
            @task[data.tid] = resolve
            # 透過 MessagePack 編碼資料，並透過 WebSocket 傳送。
            @connection.send msgpack.encode payload

    # addListener 會新增指定的事件監聽函式。
    addListener: (event, callback, name = null) ->
        switch event
            when 'message' then @messageListeners = array.push(@messageListeners, callback)
            when 'close'   then @closeListeners   = array.push(@closeListeners, callback)
            when 'open'    then @openListeners    = array.push(@openListeners, callback)
            when 'reopen'  then @reopenListeners  = array.push(@reopenListeners, callback)
            when 'error'   then @errorListeners   = array.push(@errorListeners, callback)
        @

    removeListener: (event, name) ->
        switch event
            when 'message' then @messageListeners = array.push(@messageListeners, callback)
            when 'close'   then @closeListeners   = array.push(@closeListeners, callback)
            when 'open'    then @openListeners    = array.push(@openListeners, callback)
            when 'reopen'  then @reopenListeners  = array.push(@reopenListeners, callback)
            when 'error'   then @errorListeners   = array.push(@errorListeners, callback)



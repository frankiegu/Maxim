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
        @listenerID       = 0
        @tasks            = {}
        @messageListeners = {}
        @errorListeners   = {}
        @openListeners    = {}
        @reopenListeners  = {}
        @closeListeners   = {}

        _registerListeners()
        @addListener 'message', @_taskCaller


    # metadata 會設置單次性的指定中繼資料。
    metadata: (data) ->
        @metadata = data


    # _registerListeners 會註冊所有的事件監聽器，用以呼叫使用者自訂的所有監聽器。
    _registerListeners: ->
        that = @

        @connection.onmessage = (event) ->
            for listener in that.messageListeners
                listener(event) if listener isnt undefined

        @connection.onclose = (event) ->
            for listener in that.closeListeners
                listener(event) if listener isnt undefined

        @connection.onopen = (event) ->
            for listener in that.openListeners
                listener(event) if listener isnt undefined

        #@connection.onmessage = (event) ->
        #    for listener in that.messageListeners
        #        listener(event) if listener isnt undefined

        @connection.onerror = (event) ->
            for listener in that.errorListeners
                listener(event) if listener isnt undefined

    # _taskCaller 會監聽所有接收到的訊息，如果帶有工作編號則呼叫指定的工作 Resolve 函式。
    _taskCaller: (event) ->
        # 將接收到的資料以 MessagePack 解碼。
        data = msgpack.decode(event.data)
        #
        if data.tid is undefined
            return
        #
        maximData = new MaximData(data)
        # 如果該資料含有工作編號，就呼叫指定的工作 Resolve 函式。
        @task[data.tid](maximData) if data.tid isnt undefined

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
            when 'message' then if name is null then @messageListeners["maxim_#{@listenerID}"] = callback else @messageListeners[name] = callback
            when 'close'   then if name is null then @closeListeners["maxim_#{@listenerID}"]   = callback else @closeListeners[name]   = callback
            when 'open'    then if name is null then @openListeners["maxim_#{@listenerID}"]    = callback else @openListeners[name]    = callback
            when 'reopen'  then if name is null then @reopenListeners["maxim_#{@listenerID}"]  = callback else @reopenListeners[name]  = callback
            when 'error'   then if name is null then @errorListeners["maxim_#{@listenerID}"]   = callback else @errorListeners[name]   = callback
        @

    # removeListener 會移除指定事件的所有監聽器，或指定名稱的事件監聽器。
    removeListener: (event, name = null) ->
        switch event
            when 'message' then if name is null then @messageListeners = {} else @messageListeners[name] = undefined
            when 'close'   then if name is null then @closeListeners   = {} else @closeListeners[name]   = undefined
            when 'open'    then if name is null then @openListeners    = {} else @openListeners[name]    = undefined
            when 'reopen'  then if name is null then @reopenListeners  = {} else @reopenListeners[name]  = undefined
            when 'error'   then if name is null then @errorListeners   = {} else @errorListeners[name]   = undefined
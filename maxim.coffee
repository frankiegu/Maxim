class MaximData
    constructor: ->

    data: ->
    error: ->
    metadata: ->


class Maxim
    constructor: (url, options) ->
        @connection       = new WebSocket(url)
        @columns          = ''
        @options          =
            metadata: {}
        @metadata         = {}
        @taskID           = 0
        @tasks            = {}
        @messageListeners = []
        @errorListeners   = []
        @openListeners    = []
        @reopenListeners  = []
        @closeListeners   = []

    # metadata 會設置單次性的指定中繼資料。
    metadata: (data) ->
        @metadata = data

    _taskCaller: (event) ->
        # 將接收到的資料以 MessagePack 解碼。
        data = msgpack.decode(event.data)

        if data.tid isnt undefined


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
        # 透過 MessagePack 編碼資料，並透過 WebSocket 傳送。
        @connection.send msgpack.encode payload
        # 遞增工作編號供下次呼叫使用。
        @taskID++

        return new Promise (resolve, reject) ->





    # addListener 會新增指定的事件監聽函式。
    addListener: (event, callback) ->
        switch event
            when 'message' then @messageListeners = array.push(@messageListeners, callback)
            when 'close'   then @closeListeners   = array.push(@closeListeners, callback)
            when 'open'    then @openListeners    = array.push(@openListeners, callback)
            when 'reopen'  then @reopenListeners  = array.push(@reopenListeners, callback)
            when 'error'   then @errorListeners   = array.push(@errorListeners, callback)
        @



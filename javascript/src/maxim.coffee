import msgpack from "msgpack-lite"

class Maxim {
    constructor(url) {
        this.connection = new WebSocket(url)
        this.meta = {}
        this.tid = -1
        this.col = []
        this.resolves = {}
        this.rejects = {}
        this.options = {
            beforeExecute: () => true,
            afterExecute: () => {},
            onClose: () => {},
            onError: () => {},
            onMessage: () => {},
            meta: {}
        }
    }
    _clean() {
        this.meta = {}
        this.col = []
    }
    _dataResolver(data) {
        class Data {
            constructor(data) {
                this.data = data
            }
            meta() {
                return this.data.met
            }
            data() {
                return this.data.dat
            }
            code() {
                if (this.data.cod !== undefined) {
                    return this.data.cod
                } else if (this.data.err.cod !== undefined) {
                    return this.data.err.cod
                }
            }
            func() {

            }
            err() {

            }
            isErr() {

            }
        }
        return new Data(data)
    }
    columns() {
        this.col = arguments
    }
    execute(func, data) {
        var body, that
        that = this
        this.tid++
            body = {
                tid: this.tid,
                dat: data,
                met: Object.assign(this.meta, this.options.meta),
                col: this.col,
                fun: func
            }
        this.connection.send(msgpack.encode(data))
        this._clean()
        return new Promise((resolve, reject) => {
            that.resolves[that.tid] = (value) => resolve(that.)
        })
    }
    meta(data) {
        this.meta = data
    }
    setup(options) {
        this.options = Object.assign(this.options, options)
    }
}

console.log(msgpack.encode({ "hello": "World" }))
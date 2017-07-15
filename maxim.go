package maxim

import (
	"net/http"

	"github.com/olahol/melody"
	"github.com/vmihailenco/msgpack"
)

const (
	StatusOK            = "MaximOK"
	StatusError         = "MaximError"
	StatusProcessing    = "MaximProcessing"
	StatusFull          = "MaximFull"
	StatusExists        = "MaximExists"
	StatusInvalid       = "MaximInvalid"
	StatusNotFound      = "MaximNotFound"
	StatusNotAuthorized = "MaximNotAuthorized"
	StatusNoPermission  = "MaximNoPermission"
	StatusNoChanges     = "MaximNoChanges"
	StatusUnimplemented = "MaximUnimplemented"

	StatusFileDone         = "MaximFileDone"
	StatusFileNext         = "MaximFileNext"
	StatusFileRetry        = "MaximFileRetry"
	StatusFileAbort        = "MaximFileAbort"
	StatusFileEmpty        = "MaximFileEmpty"
	StatusFileIncomplete   = "MaximFileIncomplete"
	StatusFileTimeout      = "MaximFileTimeout"
	StatusFileNoPermission = "MaximFileNoPermission"

	//ErrChunkRetry = errors.New("Please resend the chunk.")
	//ErrChunkAbort = errors.New("Abort the entire upload process.")
)

type H map[string]interface{}
type HandlerFunc func(*Context)

func New() *Engine {
	var e Engine
	return &e
}

func Default() *Engine {
	var e Engine
	return &e
}

type Engine struct {
	functions     map[string][]HandlerFunc
	fileFunctions map[string][]HandlerFunc
}

type mainReceiver struct {
	melody *melody.Melody
}

func (h *mainReceiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.melody.HandleRequest(w, r)
}

func (e *Engine) Run(port string) {
	// Create the main WebSocket connection by the melody framework.
	m := melody.New()
	// Handle the main path.
	http.Handle("/", &mainReceiver{melody: m})

	// Process the received data from the client via WebSocket.
	m.HandleMessage(func(s *melody.Session, raw []byte) {
		// Unmarshal the data from message pack to the struct.
		var d rawData
		err := msgpack.Unmarshal(raw, &d)
		if err != nil {
			panic(err)
		}

		// Create the context so we can pass it to the first handler.
		ctx := Context{
			session:  s,
			melody:   m,
			metadata: d.metadata,
		}

		var handlers []HandlerFunc
		var ok bool
		// Get the handlers.
		if d.file.key == "" {
			handlers, ok = e.functions[d.function]
		} else {
			handlers, ok = e.fileFunctions[d.function]
		}
		// Return the unimplemented error message if the function is not found.
		if !ok {

		}

		// Save the handlers to the context so
		// we could use `Next()` in the context to call the next handler.
		ctx.handlers = handlers
		// Call the first handler with the context.
		handlers[0](&ctx)
	})
	http.ListenAndServe(port, nil)
}

func (e *Engine) Use(handler HandlerFunc) {

}

func (e *Engine) On(function string, handlers ...HandlerFunc) {

}

func (e *Engine) OnFile(function string, handlers ...HandlerFunc) {

}

func (e *Engine) Boardcast(data interface{}) {

}

func (e *Engine) Execute(function string, data interface{}) {

}

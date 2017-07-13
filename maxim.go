package maxim

import (
	"net/http"

	"github.com/olahol/melody"
	"github.com/vmihailenco/msgpack"
)

const (
	// StatusOK means everything's okay. It's a ok if the user is trying to delete something that has already been deleted.
	StatusOK = "OK"
	// StatusError means a common or an internal error has occurred.
	StatusError = "Error"
	// StatusProcessing means the request is now processing and won't be done in just few seconds.
	StatusProcessing = "Processing"
	// StatusFull means the request is not acceptable because something is full (for example: The friend list, group).
	StatusFull = "Full"
	// StatusExists means something has already been existed, like the username or the email address.
	StatusExists = "Exists"
	// StatusInvalid means the format of the request is invalid.
	StatusInvalid = "Invalid"
	// StatusNotFound means the resource which the user was requested is not found.
	StatusNotFound = "NotFound"
	// StatusNotAuthorized means the user should be logged in to make the request.
	StatusNotAuthorized = "NotAuthorized"
	// StatusNoPermission means the user has logged in but has no permission to do something.
	StatusNoPermission = "NoPermission"
	// StatusNoChanges means the request has changed nothing, it's the same as what the request trying to change for.
	StatusNoChanges = "NoChanges"
	//
	StatusUnimplemented = "Unimplemented"

	// FileStatusDone means the file chunks were all uploaded, and it's combined, this is the final result.
	FileStatusDone = "Done"
	// FileStatusNext means the current file chunk has been processed, please upload the next chunk.
	FileStatusNext = "Next"
	// FileStatusRetry means the error occurred while processing the current file chunk, please resend the chunk.
	FileStatusRetry = "Retry"
	// FileStatusAbort means to abort the entire upload process.
	FileStatusAbort = "Abort"

	//ErrChunkRetry = errors.New("Please resend the chunk.")
	//ErrChunkAbort = errors.New("Abort the entire upload process.")
)

// Engine represents the main Maxim engine.
type Engine struct {
	// Handlers stores the normal event handlers.
	Handlers map[string]Handler
	// FileHandlers stores the file upload handlers.
	FileHandlers map[string]Handler
}

type Handler func(*Context)

type H map[string]interface{}

func (e *Engine) On(event string, handler Handler) {
	e.Handlers[event] = handler
}

func (e *Engine) OnFile(event string, chunkHandler ChunkHandler, handler Handler) {

}

type helloHandler struct {
	melody *melody.Melody
}

func (h *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.melody.HandleRequest(w, r)
}

type Request struct {
	TaskID   int                    `msgpack:"tid"`
	Function string                 `msgpack:"fun"`
	Columns  string                 `msgpack:"col"`
	Data     []byte                 `msgpack:"dat"`
	Metadata map[string]interface{} `msgpack:"met"`
}

type ChunkRequest struct {
	TaskID   int                    `msgpack:"tid"`
	Metadata map[string]interface{} `msgpack:"met"`
	Chunk    Chunk                  `msgpack:"chk"`
}

type Chunk struct {
	Part   uint8  `msgpack:"par"`
	Total  uint8  `msgpack:"tol"`
	Key    string `msgpack:"key"`
	Binary []byte `msgpack:"bin"`
}

type Respond struct {
	TaskID   int                    `msgpack:"tid"`
	Metadata map[string]interface{} `msgpack:"met"`
	Code     string                 `msgpack:"cod"`
	Data     []byte                 `msgpack:"dat"`
	Error    []byte                 `msgpack:"err"`
}

type ChunkRespond struct {
}

func (e *Engine) Run(port string) {
	m := melody.New()
	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	http.Handle("/", &helloHandler{
		melody: m,
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var req Request
		err := msgpack.Unmarshal(msg, &req)
		if err != nil {

		}

		handler, ok := e.Handlers[req.Function]
		if !ok {
			resp, err := msgpack.Marshal(Respond{
				TaskID: req.TaskID,
				Code:   StatusUnimplemented,
			})
			if err != nil {

			}
			s.Write(resp)
		}

		context := &Context{
			TaskID:     req.TaskID,
			Function:   req.Function,
			Columns:    req.Columns,
			Data:       req.Data,
			Metadata:   req.Metadata,
			session:    s,
			connection: m,
		}
		handler(context)
	})
	http.ListenAndServe(port, nil)
}

func Default() *Engine {
	return &Engine{}
}

func New() *Engine {
	return &Engine{}
}

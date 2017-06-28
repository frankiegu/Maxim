package maxim

import "errors"

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

	// FileStatusDone means the file chunks were all uploaded, and it's combined, this is the final result.
	FileStatusDone = "Done"
	// FileStatusNext means the current file chunk has been processed, please upload the next chunk.
	FileStatusNext = "Next"
	// FileStatusRetry means the error occurred while processing the current file chunk, please resend the chunk.
	FileStatusRetry = "Retry"
	// FileStatusAbort means to abort the entire upload process.
	FileStatusAbort = "Abort"

	ErrChunkRetry = errors.New("Please resend the chunk.")
	ErrChunkAbort = errors.New("Abort the entire upload process.")
)

// H is the alias for map[string]interface{}
type H map[string]interface{}

type Context struct {
}

// Bind binds the payload to the local struct or the map.
func (c *Context) Bind(dest interface{}) (err error) {
	return
}

// Send sends the data back to the client.
func (c *Context) Send(code string, data interface{}) (err error) {
	return
}

type Data struct {
	Tid int
	Met interface{}
	Dat interface{}
	Cod string
	Err ErrorData
}

type ErrorData struct {
	Dat interface{}
	Cod string
}

// Engine represents the main websocket service engine.
type Engine struct {
	Options Options
}

// NewHandler creates a new handler to handle the specified function.
func (e *Engine) NewHandler(name string, handler func(*Context)) {

}

// Run runs the websocket service.
func (e *Engine) Run(port string) {

}

// Options stores the options to create the new engine.
type Options struct {
	Chunker    Chunker
	NameFormat string
	ExpireTime int
}

// New creates a main engine with the specified options.
func New(opts Options) *Engine {
	return &Engine{}
}

// Chunk represents a uploaded chunk from the client.
type Chunk struct {
	Key         string
	Bin         string
	CurrentPart int
	TotalParts  int
}

// Chunker defines the interface to process the uploaded chunk.
type Chunker interface {
	Save(Chunk) (error, H)
	Clean()
	Done(Chunk) (error, H)
}

// Default creates a new engine with the default options.
func Default() *Engine {
	return &Engine{Options: Options{
		Chunker:    &DefaultChunker{},
		ExpireTime: 86400,
	}}
}

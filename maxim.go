package maxim

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
)

type H map[string]interface{}

type Context struct {
}

func (c *Context) Bind(dest interface{}) (err error) {
	return
}

func (c *Context) Send(code string, data interface{}) (err error) {
	return
}

type Engine struct {
}

func (e *Engine) Function(name string, handler func(*Context)) {

}

func (e *Engine) Run(port string) {

}

func New() *Engine {
	return &Engine{}
}

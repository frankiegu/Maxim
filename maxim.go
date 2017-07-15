package maxim

type H map[string]interface{}
type HandlerFunc func(*Context)

func New() {

}

func Default() {

}

type Engine struct {
}

func (e *Engine) Run(port string) {

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

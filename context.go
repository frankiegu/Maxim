package maxim

import (
	"encoding/json"
	"time"

	"github.com/olahol/melody"
	"github.com/vmihailenco/msgpack"
)

type Context struct {
	File     File
	Chunk    Chunk
	data     string
	taskID   uint
	function string
	metadata map[string]interface{}
	store    map[string]interface{}
	handlers []HandlerFunc
	index    int8
	session  *melody.Session
	melody   *melody.Melody
}

type File struct {
	Name      string
	Size      int64
	Extension string
	Path      string
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.store[key]
	return
}

func (c *Context) GetBool(key string) (b bool) {
	if v, ok := c.Get(key); ok && v != nil {
		b, _ = v.(bool)
	}
	return
}

func (c *Context) GetDuration(key string) (d time.Duration) {
	if v, ok := c.Get(key); ok && v != nil {
		d, _ = v.(time.Duration)
	}
	return
}

func (c *Context) GetFloat64(key string) (f64 float64) {
	if v, ok := c.Get(key); ok && v != nil {
		f64, _ = v.(float64)
	}
	return
}

func (c *Context) GetInt(key string) (i int) {
	if v, ok := c.Get(key); ok && v != nil {
		i, _ = v.(int)
	}
	return
}

func (c *Context) GetInt64(key string) (i64 int64) {
	if v, ok := c.Get(key); ok && v != nil {
		i64, _ = v.(int64)
	}
	return
}

func (c *Context) GetString(key string) (s string) {
	if v, ok := c.Get(key); ok && v != nil {
		s, _ = v.(string)
	}
	return
}

func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
	if v, ok := c.Get(key); ok && v != nil {
		sm, _ = v.(map[string]interface{})
	}
	return
}

func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if v, ok := c.Get(key); ok && v != nil {
		sms, _ = v.(map[string]string)
	}
	return
}

func (c *Context) GetStringSlice(key string) (ss []string) {
	if v, ok := c.Get(key); ok && v != nil {
		ss, _ = v.([]string)
	}
	return
}

func (c *Context) GetTime(key string) (t time.Time) {
	if v, ok := c.Get(key); ok && v != nil {
		t, _ = v.(time.Time)
	}
	return
}

func (c *Context) Set(key string, value interface{}) {
	if c.store == nil {
		c.store = make(map[string]interface{})
	}
	c.store[key] = value
}

func (c *Context) Metadata() map[string]interface{} {
	if c.metadata == nil {
		c.metadata = make(map[string]interface{})
	}
	return c.metadata
}

func (c *Context) SetMetadata(key string, value interface{}) {
	if c.metadata == nil {
		c.metadata = make(map[string]interface{})
	}
	c.metadata[key] = value
}

func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Bind(destination interface{}) error {
	return json.Unmarshal([]byte(c.data), destination)
}

func (c *Context) Abort() {
}

func (c *Context) Copy() *Context {
	ctx := *c
	return &ctx
}

func (c *Context) RespondStatus(status string) error {
	// Build the respond.
	resp := respond{
		taskID:   c.taskID,
		function: c.function,
		metadata: c.metadata,
		code:     status,
	}
	// Convert the respond struct to the message pack binary.
	msg, err := msgpack.Marshal(resp)
	if err != nil {
		return err
	}
	// Wrtie the message pack to the specified websocket session.
	c.session.Write(msg)

	return nil
}

func (c *Context) RespondError(status string, errData interface{}) error {
	// Conver the error data to json.
	errJSON, err := json.Marshal(errData)
	if err != nil {
		return err
	}

	// Build the respond.
	resp := respond{
		taskID:   c.taskID,
		function: c.function,
		metadata: c.metadata,
		code:     status,
		err:      string(errJSON),
	}
	// Convert the respond struct to the message pack binary.
	msg, err := msgpack.Marshal(resp)
	if err != nil {
		return err
	}
	// Wrtie the message pack to the specified websocket session.
	c.session.Write(msg)

	return nil
}

func (c *Context) Respond(status string, data interface{}) error {
	// Conver the data to json.
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Build the respond.
	resp := respond{
		taskID:   c.taskID,
		function: c.function,
		metadata: c.metadata,
		code:     status,
		data:     string(dataJSON),
	}
	// Convert the respond struct to the message pack binary.
	msg, err := msgpack.Marshal(resp)
	if err != nil {
		return err
	}
	// Wrtie the message pack to the specified websocket session.
	c.session.Write(msg)

	return nil
}

func (c *Context) RespondOthers(status string, data interface{}) error {
	// Conver the data to json.
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Build the respond.
	resp := respond{
		metadata: c.metadata,
		code:     status,
		data:     string(dataJSON),
	}
	// Convert the respond struct to the message pack binary.
	msg, err := msgpack.Marshal(resp)
	if err != nil {
		return err
	}
	// Wrtie the message pack to the specified websocket session.
	c.melody.BroadcastOthers(msg, c.session)

	return nil
}

func (c *Context) Execute() error {
	return nil
}

func (c *Context) ExecuteOthers() error {
	return nil
}

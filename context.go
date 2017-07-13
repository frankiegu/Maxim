package maxim

import (
	"encoding/json"

	"github.com/olahol/melody"
	"github.com/vmihailenco/msgpack"
)

type Context struct {
	// TaskID is the ID of the request, keep the same task id in the respond so the client can call the correct callback.
	TaskID int
	// Function is the name of the function which the request is trying to call.
	Function string
	// Columns is the fields that the request is trying to get.
	Columns string
	// Data is the main data from the request.
	Data []byte
	// Metadata is the metadata from the request, it could be the API version, network speed, token.. etc.
	Metadata map[string]interface{}
	//
	connection *melody.Melody
	//
	session *melody.Session
	//
	respondMetadata map[string]interface{}
}

type File struct {
	Key       string
	Name      string
	Size      string
	Extension string
	TmpPath   string
}

func (c *Context) Metadata(metadata map[string]interface{}) {
	c.respondMetadata = metadata
}

func (c *Context) Bind(destination interface{}) {
	err := json.Unmarshal(c.Data, &destination)
	if err != nil {

	}
}

func (c *Context) Respond(status string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {

	}
	resp, err := msgpack.Marshal(Respond{
		TaskID:   c.TaskID,
		Metadata: c.respondMetadata,
		Code:     status,
		Data:     jsonData,
	})
	if err != nil {

	}
	c.session.Write(resp)
}

func (c *Context) RespondOthers(status string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {

	}
	resp, err := msgpack.Marshal(Respond{
		TaskID:   -1,
		Metadata: c.respondMetadata,
		Code:     status,
		Data:     jsonData,
	})
	if err != nil {

	}
	c.connection.BroadcastOthers(resp, c.session)
}

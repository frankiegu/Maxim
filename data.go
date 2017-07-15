package maxim

type rawData struct {
	taskID   uint                   `msgpack:"tid"`
	metadata map[string]interface{} `msgpack:"met"`
	data     string                 `msgpack:"dat"`
	function string                 `msgpack:"fun"`
	file     rawFileInfo            `msgpack:"fil"`
}

type rawFileInfo struct {
	key         string `msgpack:"key"`
	currentPart uint   `msgpack:"par"`
	totalPart   uint   `msgpack:"tol"`
	bin         []byte `msgpack:"bin"`
}

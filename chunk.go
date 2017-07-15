package maxim

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

type Chunk struct {
	Binary    []byte
	Key       string
	Part      int8
	TotalPart int8
}

type ChunkHandlerOption struct {
	TmpPath string
	Path    string
	Timeout uint
	MaxSize uint
	OnError func(error, *Context)
}

func ChunkHandler(option ...ChunkHandlerOption) {
	// Create a default option if the user didn't pass a chunk handler option.
	var o ChunkHandlerOption
	if len(option) == 0 {
		o = ChunkHandlerOption{
			TmpPath: "/tmp",
			Path:    "/tmp",
			Timeout: 3000,     // 5 minutes
			MaxSize: 40000000, // 40 MB
		}
	} else {
		o = option[0]
	}

	//
	return func(c *Context) {
		// Make sure the key is a valid UUID format.
		isUUID := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString
		// Abort the upload process when there's no chunk or the key name is invalid.
		if c.Chunk.Key == "" || !isUUID(c.Chunk.Key) {
			c.RespondStatus(StatusFileEmpty)
			return
		}

		// Make sure the previous chunk does exist if the chunk is not the first one.
		if c.Chunk.Part != 1 {
			// Create the file name for the previous chunk.
			// Format: maxim_chunk_[KEY]_[TIMESTAMP]_[CHUNK-NUMBER]
			previousChunk := fmt.Sprintf("maxim_chunk_%s_%d_%d", c.Chunk.Key, time.Now().Unix(), c.Chunk.Part-1)
			// Make sure the previous chunk does exist.
			_, err := os.Stat(fmt.Sprintf("%s/%s", o.TmpPath, previousChunk))
			if err != nil && os.IsNotExist(err) {
				c.RespondStatus(StatusFileEmpty)
				return
			}
		}

		// Is it the last chunk?
		isLastChunk := c.Chunk.Part == c.Chunk.TotalPart

		// Create the file name for the current chunk with the specified format.
		// Format: maxim_chunk_[KEY]_[TIMESTAMP]_[CHUNK-NUMBER]
		chunk := fmt.Sprintf("maxim_chunk_%s_%d_%d", c.Chunk.Key, time.Now().Unix(), c.Chunk.Part)

		// Create an empty file if the file doesn't exist.
		_, err := os.Stat()
		if err != nil && os.IsNotExist(err) {
			// ..
		}

		// Write the chunk to the temporary file.
		c.Chunk.Binary

		// If we've done with the last chunk,
		// we continue to the next handler with the file information.
		if isLastChunk {
			// Move the file to the specified directory.

			// Don't move the file if the file directory is the same as the temporary directory.

			// Put the file information to the context.

			c.Next()
		}
	}
}

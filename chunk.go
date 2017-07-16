package maxim

import (
	"fmt"
	"os"
	"regexp"

	"golang.org/x/sys/unix"
)

type Chunk struct {
	Binary    []byte
	Key       string
	Part      int8
	Name      string
	Extension string
	TotalPart int8
}

type ChunkHandlerOption struct {
	TmpPath string
	MaxSize int64
	OnError func(error, *Context)
}

func ChunkHandler(option ...ChunkHandlerOption) HandlerFunc {
	// Create a default option if the user didn't pass a chunk handler option.
	var o ChunkHandlerOption
	if len(option) == 0 {
		o = ChunkHandlerOption{
			TmpPath: "/tmp",
			MaxSize: 40000000, // 40 MB
		}
	} else {
		o = option[0]
	}

	// The real chunk handler.
	return func(c *Context) {
		// Make sure the key is a valid UUID format.
		isUUID := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString
		// Abort the upload process when there's no chunk or the key name is invalid.
		if c.Chunk.Key == "" || !isUUID(c.Chunk.Key) {
			c.RespondStatus(StatusFileEmpty)
			return
		}
		// Make sure the directories is writable.
		if err := unix.Access(o.TmpPath, unix.W_OK); err != nil {
			o.OnError(err, c)
			c.RespondStatus(StatusFileNoPermission)
			return
		}
		// Create the filename and the path for the file.
		filename := fmt.Sprintf("%s/%s.%s", o.TmpPath, c.Chunk.Key, c.Chunk.Extension)
		// Create the file if it's the first chunk.
		if c.Chunk.Part == 1 {
			f, err := os.Create(filename)
			if err != nil {
				o.OnError(err, c)
				c.RespondStatus(StatusFileAbort)
				return
			}
			defer f.Close()
		}

		// Open the file.
		f, err := os.Open(filename)
		if err != nil {
			o.OnError(err, c)
			c.RespondStatus(StatusFileAbort)
			return
		}
		defer f.Close()
		// Write the chunk binary to the file.
		_, err = f.Write(c.Chunk.Binary)
		if err != nil {
			o.OnError(err, c)
			c.RespondStatus(StatusFileAbort)
			return
		}
		f.Sync()

		// Get the information of the file.
		info, err := os.Stat(filename)
		if err != nil {
			o.OnError(err, c)
			c.RespondStatus(StatusFileAbort)
			return
		}
		size := info.Size()

		// Calculate the size of the file, abort and delete the file if it's too large.
		if size > o.MaxSize {
			err := os.Remove(filename)
			if err != nil {
				o.OnError(err, c)
				c.RespondStatus(StatusFileAbort)
				return
			}
			c.RespondStatus(StatusFileSize)
			return
		}

		// Get the next chunk if this is not the last chunk.
		if c.Chunk.Part != c.Chunk.TotalPart {
			c.RespondStatus(StatusFileNext)
			return
		}

		// We continue to the next handler with the file information
		// if this is the last chunk and we just nailed it.
		c.File = File{
			Name:      c.Chunk.Name,
			Size:      size,
			Extension: c.Chunk.Extension,
			Path:      filename,
		}
		// Here we go!
		c.Next()
	}
}

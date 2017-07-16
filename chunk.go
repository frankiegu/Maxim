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
	Timeout uint
	MaxSize int64
	OnError func(error, *Context)
}

func generateChunkName(key string, part int8) string {
	return fmt.Sprintf("maxim_chunk_%s_%d", key, part)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err != nil
}

func ChunkHandler(option ...ChunkHandlerOption) {
	// Create a default option if the user didn't pass a chunk handler option.
	var o ChunkHandlerOption
	if len(option) == 0 {
		o = ChunkHandlerOption{
			TmpPath: "/tmp",
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

		// Make sure the directories is writable.
		if unix.Access(o.TmpPath, unix.W_OK) != nil {
			c.RespondStatus(StatusFileNoPermission)
			return
		}

		filename := fmt.Sprintf("%s/%s.%s", o.TmpPath, c.Chunk.Key, c.Chunk.Extension)

		// Create the file if it's the first chunk.
		if c.Chunk.Part == 1 {
			f, err := os.Create(filename)
			if err != nil {
				// ...
			}
			defer f.Close()
		}

		// Get the information of the file.
		info, err := os.Stat(filename)
		if err != nil {

		}

		// Abort and delete the file if it's not done and timeouted.

		// Open the file.
		f, err := os.Open(filename)
		if err != nil {
			// ...
		}
		defer f.Close()
		// Write the chunk binary to the file.
		_, err = f.Write(c.Chunk.Binary)
		if err != nil {
			// maxim.StatusFileNoPERMISSION
			// ...
		}
		f.Sync()

		// Calculate the size of the file, abort and delete the file if it's too large.
		if info.Size() > o.MaxSize {
			err := os.Remove(filename)
			if err != nil {

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
		//Name: c.Chunk.Name,
		//Size:
		//Extension: c.Chunk.Extension,
		//Path: filename,
		//Duration:
		}

		// Move the file to the specified directory.

		// Don't move the file if the file directory is the same as the temporary directory.

		// Put the file information to the context.

		c.Next()

	}
}

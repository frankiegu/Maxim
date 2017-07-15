package maxim

import (
	"fmt"
	"os"
	"regexp"

	uuid "github.com/satori/go.uuid"

	"golang.org/x/sys/unix"
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

		// Make sure the directories is writable.
		if unix.Access(o.TmpPath, unix.W_OK) != nil || unix.Access(o.Path, unix.W_OK) != nil {
			c.RespondStatus(StatusFileNoPermission)
			return
		}

		// Make sure the previous chunk does exist if the chunk is not the first one.
		//if c.Chunk.Part != 1 {
		//	// Create the file name for the previous chunk.
		//	// Format: maxim_chunk_[KEY]_[CHUNK-NUMBER]
		//	previousChunk := fmt.Sprintf("maxim_chunk_%s_%d", c.Chunk.Key, c.Chunk.Part-1)
		//	// Make sure the previous chunk does exist.
		//	_, err := os.Stat(fmt.Sprintf("%s/%s", o.TmpPath, previousChunk))
		//	if err != nil && os.IsNotExist(err) {
		//		c.RespondStatus(StatusFileEmpty)
		//		return
		//	}
		//}

		// Create the file name for the current chunk with the specified format.
		// Format: maxim_chunk_[KEY]_[CHUNK-NUMBER]
		chunkName := generateChunkName(c.Chunk.Key, c.Chunk.Part)

		// Create an empty file for the new chunk.
		f, err := os.Create(fmt.Sprintf("%s/%s", o.TmpPath, chunkName))
		if err != nil {
			// ...
		}
		defer f.Close()

		// Write the binary to the chunk file.
		_, err = f.Write(c.Chunk.Binary)
		if err != nil {
			// maxim.StatusFileNoPERMISSION
			// ...
		}
		f.Sync()

		// If we've done with the last chunk...
		if c.Chunk.Part == c.Chunk.TotalPart {
			var filename string

			// Combines the chunks to a single file. (1 ~ TotalPart)
			for i := int8(1); i <= c.Chunk.TotalPart; i++ {
				// Create the name for the chunk.
				chunkName = generateChunkName(c.Chunk.Key, i)

				// Return an error if the chunk doesn't exist.
				if !fileExists(fmt.Sprintf("%s/%s", o.TmpPath, chunkName)) {
					c.RespondStatus(StatusFileIncomplete)
					return
				}

				// Create the real file so we have the destination for the chunks.
				if i == 1 {
					// Generate a UUID for the file name.
					filename = uuid.NewV4().String()
					// Just create another UUID if the file name does exist in the destination.
					for fileExists(fmt.Sprintf("%s/%s", o.Path, filename)) {
						filename = uuid.NewV4().String()
					}

					// And we create the file.
					f, err := os.Create(fmt.Sprintf("%s/%s", o.Path, filename))
					if err != nil {
						// ...
					}
					defer f.Close()
				}

				// Now we open the file.
				f, err := os.Open(fmt.Sprintf("%s/%s", o.Path, filename))
				if err != nil {
					// ...
				}
				defer f.Close()

				// And we open the chunk.
				chk, err := os.Open(fmt.Sprintf("%s/%s", o.TmpPath, generateChunkName(c.Chunk.Key, i)))
				if err != nil {
					// ...
				}
				defer chk.Close()

				// Write the chunk
				_, err = f.Write(c.Chunk.Binary)
				if err != nil {
					// maxim.StatusFileNoPERMISSION
					// ...
				}
				f.Sync()
			}

			// we continue to the next handler with the file information.

			// Move the file to the specified directory.

			// Don't move the file if the file directory is the same as the temporary directory.

			// Put the file information to the context.

			c.Next()
		}
	}
}

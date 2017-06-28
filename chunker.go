package maxim

type DefaultChunker struct {
	TmpPath  string
	FilePath string
}

func (c *DefaultChunker) Save(chunk Chunk) (error, H) {

}

func (c *DefaultChunker) Clean() {

}

func (c *DefaultChunker) Done(chunk Chunk) (error, H) {

}

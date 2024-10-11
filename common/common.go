package common

type BlockType struct {
	Name       string
	LinePrefix string
	Handler    func(Block) (interface{}, string, error)
}

type Block struct {
	Type     BlockType
	Command  interface{}
	Queries  []string
	NextPage string
}

func (b Block) Handle() (interface{}, string, error) {
	return b.Type.Handler(b)
}

type PipetApp struct {
	Blocks    []Block
	Data      []interface{}
	MaxPages  int
	Separator []string
}

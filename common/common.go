package common

type Block struct {
	Type     string
	Command  interface{}
	Queries  []string
	NextPage string
}

type PipetApp struct {
	Blocks    []Block
	Data      []interface{}
	MaxPages  int
	Separator []string
}

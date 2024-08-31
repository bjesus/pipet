package common

type Block struct {
	Type     string
	Command  string
	Queries  []string
	NextPage string
}

type ErolApp struct {
	Blocks    []Block
	Data      []interface{}
	MaxPages  int
	Separator []string
}

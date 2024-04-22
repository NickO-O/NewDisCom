package parser

type Node struct {
	Left     *Node
	Right    *Node
	Operator string
	Value    float64
}

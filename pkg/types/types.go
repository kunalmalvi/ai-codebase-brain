package types

type NodeType string

const (
	NodeTypeFile     NodeType = "file"
	NodeTypeFunction NodeType = "function"
	NodeTypeClass    NodeType = "class"
	NodeTypeVariable NodeType = "variable"
)

type ProjectInfo struct {
	Name        string
	Language    string
	Framework   string
	BuildSystem string
	Root        string
}
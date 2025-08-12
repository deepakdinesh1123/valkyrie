package jsonschema

type File struct {
	Name    string
	Content string
}

type ExecReqFiles struct {
	Files []File
}

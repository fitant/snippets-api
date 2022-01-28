package contract

type CreateSnippet struct {
	Metadata Metadata
	Snippet  string
}

type CreateSnippetResponse struct {
	URL string
}

type Metadata struct {
	Ephemeral bool
	Language  string
}

var CS CreateSnippet

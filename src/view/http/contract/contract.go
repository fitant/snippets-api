package contract

type CreateSnippet struct {
	Metadata Metadata `json:"metadata"`
	Data     Data     `json:"data"`
}

type CreateSnippetResponse struct {
	URL string
}

type Data struct {
	Snippet string `json:"snippet"`
}

type Metadata struct {
	Ephemeral bool   `json:"ephemeral"`
	Language  string `json:"language"`
}

var CS CreateSnippet

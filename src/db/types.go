package db

type FindSnippetQuery struct {
	ID string `bson:"_id"`
}

type InsertSnippetQuery struct {
	ID       string `bson:"_id"`
	Snippet  string `bson:"snippet"`
	Language string `bson:"language"`
}

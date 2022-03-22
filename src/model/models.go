package model

import (
	"fmt"
	"time"

	"github.com/fitant/xbin-api/src/db"
	"github.com/fitant/xbin-api/src/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type SnippetController interface {
	NewSnippet(name, snippet, language string, ephemeral bool) (*Snippet, error)
	FindSnippet(name string) (*Snippet, error)
}

type Snippet struct {
	ID        string
	Ephemeral bool
	Snippet   string
	Language  string
}

type mongoSnippetController struct {
	db *db.MongoFindInsert
}

func NewMongoSnippetController(db *db.MongoFindInsert) SnippetController {
	return &mongoSnippetController{
		db: db,
	}
}

func (msc *mongoSnippetController) NewSnippet(id, snip, language string, ephemeral bool) (*Snippet, error) {
	data := db.InsertSnippetQuery{
		ID:        id,
		Snippet:   snip,
		Language:  language,
		CreatedAt: time.Now().Unix(),
	}

	doc, err := db.StructToBSON(data)
	if err != nil {
		utils.Logger.Debug(fmt.Sprintf("%s : %v", "[Models] [MongoSnippetController] [NewSnippet] [toBSON]", err))
		return nil, err
	}

	_, err = msc.db.InsertOne(doc, ephemeral)
	if err != nil {
		if err == db.ErrDuplicateKey {
			return nil, err
		}
		utils.Logger.Debug(fmt.Sprintf("%s : %v", "[Models] [MongoSnippetController] [NewSnippet] [toBSON]", err))
		return nil, err
	}

	// Create new Snippet and return
	return &Snippet{
		ID:        data.ID,
		Snippet:   data.Snippet,
		Language:  data.Language,
		Ephemeral: ephemeral,
	}, nil
}

func (msc *mongoSnippetController) FindSnippet(id string) (*Snippet, error) {
	rawQuery := db.FindSnippetQuery{
		ID: id,
	}

	query, err := db.StructToBSON(rawQuery)
	if err != nil {
		utils.Logger.Debug(fmt.Sprintf("%s : %v", "[Models] [MongoSnippetController] [FindSnippet] [toBSON]", err))
		return nil, err
	}

	ephemeral := false
	res, err := msc.db.FindOne(query, ephemeral)
	if err != nil && err != db.ErrNoDocuments {
		utils.Logger.Debug(fmt.Sprintf("%s : %v", "[Models] [MongoSnippetController] [FindSnippet] [static] [toBSON]", err))
		return nil, err
	}

	if err == db.ErrNoDocuments {
		ephemeral = true
		res, err = msc.db.FindOne(query, ephemeral)
		if err != nil {
			if err == db.ErrNoDocuments {
				return nil, err
			}
			utils.Logger.Debug(fmt.Sprintf("%s : %v", "[Models] [MongoSnippetController] [FindSnippet] [ephemeral] [toBSON]", err))
			return nil, err
		}
	}

	if res == nil {
		return nil, nil
	}

	raw, err := res.DecodeBytes()
	if err != nil {
		utils.Logger.Debug(fmt.Sprintf("%s : %v", "[Models] [MongoSnippetController] [FindSnippet] [DecodeBytes]", err))
		return nil, err
	}
	var snippet Snippet
	bson.Unmarshal(raw, &snippet)

	snippet.Ephemeral = ephemeral
	// Create new Snippet and return
	return &snippet, nil
}

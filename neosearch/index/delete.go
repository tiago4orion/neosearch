package index

import (
	"fmt"
	"net/http"

	"github.com/NeowayLabs/neosearch"
	"github.com/NeowayLabs/neosearch/neosearch/handler"
)

type DeleteIndexHandler struct {
	handler.DefaultHandler
	search *neosearch.NeoSearch
}

func NewDeleteHandler(search *neosearch.NeoSearch) *DeleteIndexHandler {
	return &DeleteIndexHandler{
		search: search,
	}
}

func (handler *DeleteIndexHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	handler.ProcessVars(req)
	indexName := handler.GetIndexName()

	if exists, err := handler.search.IndexExists(indexName); exists == false && err == nil {
		response := map[string]string{
			"error": "Index '" + indexName + "' doesn't exists.",
		}

		handler.WriteJSONObject(res, response)
		return
	} else if exists == false && err != nil {
		handler.Error(res, err.Error())
		return
	}

	err := handler.search.DeleteIndex(indexName)

	if err != nil {
		handler.Error(res, err.Error())
		return
	}

	handler.WriteJSON(res, []byte(fmt.Sprintf("{\"status\": \"Index '%s' deleted.\"}", indexName)))
}

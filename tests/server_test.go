package tests

import (
	"context"
	"go_rpc_kit_gateway_example/proto"
	"go_rpc_kit_gateway_example/server/omdb"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestOMDBGetMovies(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// our database of articles
	movies := map[string]interface{}{"Response": "mock result"}

	// mock to list out the articles
	httpmock.RegisterResponder("GET", "http://www.omdbapi.com/",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, movies)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	server := omdb.Server{}
	request := proto.MovieSearchRequest{
		Apikey:     "abcde12345",
		Searchword: "search this",
		Pagination: "1",
	}
	res, err := server.GetMovies(context.Background(), &request)
	if err != nil {
		t.Fail()
	}

	if res.Response != "mock result" {
		t.Fail()
	}
}

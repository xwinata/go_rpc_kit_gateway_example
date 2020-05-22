package omdb

import (
	"context"
	"encoding/json"
	"fmt"
	"go_rpc_kit_gateway_example/proto"
	"io/ioutil"
	"net/http"
)

// Server struct
type Server struct{}

// GetMovies func
func (s *Server) GetMovies(ctx context.Context, request *proto.MovieSearchRequest) (*proto.MovieSearchResult, error) {
	apikey, search, page := request.GetApikey(), request.GetSearchword(), request.GetPagination()

	urlstring := fmt.Sprintf("http://www.omdbapi.com/?apikey=%s", apikey)
	if len(search) > 0 {
		urlstring = urlstring + fmt.Sprintf("&s=%s", search)
	}
	if len(page) > 0 {
		urlstring = urlstring + fmt.Sprintf("&page=%s", page)
	}

	response, err := http.Get(urlstring)
	if err != nil {
		return &proto.MovieSearchResult{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &proto.MovieSearchResult{}, err
	}

	result := &proto.MovieSearchResult{}

	json.Unmarshal(body, result)

	return result, nil
}

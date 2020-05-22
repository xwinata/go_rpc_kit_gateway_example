package movie

import (
	"context"
	"encoding/json"
	"go_rpc_kit_gateway_example/client/responsehandler"
	"go_rpc_kit_gateway_example/proto"
	"io/ioutil"
	"log"
	"net/http"

	logkit "github.com/go-kit/kit/log"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

// Movie type
type Movie struct {
	Title  string `json:"Title,omitempty"`
	Year   string `json:"Year,omitempty"`
	ImdbID string `json:"imdbID,omitempty"`
	Type   string `json:"Type,omitempty"`
	Poster string `json:"Poster,omitempty"`
}

type SearchMovieRequest struct {
	SearchWord string `json:"searchword"`
	Pagination string `json:"pagination"`
}

// SearchMovieResult type
type SearchMovieResult struct {
	Search       []*Movie `json:"Search,omitempty"`
	TotalResults string   `json:"totalResults,omitempty"`
	Response     string   `json:"Response,omitempty"`
}

// Service type
type Service interface {
	SearchMovieService(ctx context.Context, s string, page string) interface{}
}

// MovieService type
type MovieService struct{}

type Endpoints struct {
	SearchMovieEndpoint endpoint.Endpoint
}

func (MovieService) SearchMovieService(ctx context.Context, s string, page string) interface{} {
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := proto.NewOMDBapiClient(conn)

	request := &proto.MovieSearchRequest{
		Apikey:     "faf7e5bb",
		Searchword: s,
		Pagination: page,
	}

	response, err := client.GetMovies(ctx, request)
	if err != nil {
		log.Println(err.Error())
	}
	// result := SearchMovieResult{}
	// json.Unmarshal([]byte(fmt.Sprint(response)), &result)
	return response
}

func MakeSearchMovieEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SearchMovieRequest)

		res := svc.SearchMovieService(ctx, req.SearchWord, req.Pagination)

		return res, nil
	}

}

func decodeSearchMovieRequest(_ context.Context, r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)

	sr := SearchMovieRequest{}
	err = json.Unmarshal(body, &sr)

	return sr, err
}

// MakeHttpHandler func
func MakeHttpHandler(endpoint Endpoints, logger logkit.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(responsehandler.EncodeError),
	}

	r.Methods("POST").Path("/search").Handler(httptransport.NewServer(
		endpoint.SearchMovieEndpoint,
		decodeSearchMovieRequest,
		responsehandler.EncodeResponse,
		options...,
	))

	return r
}

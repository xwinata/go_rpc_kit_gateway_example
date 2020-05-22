package main

import (
	"fmt"
	"go_rpc_kit_gateway_example/client/movie"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	logkit "github.com/go-kit/kit/log"
)

// ===These works. not using gokit===
// var client proto.OMDBapiClient

// func main() {
// 	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
// 	if err != nil {
// 		panic(err)
// 	}

// 	client = proto.NewOMDBapiClient(conn)

// 	http.HandleFunc("/search", GetMoviesHandler)

// 	log.Fatal(http.ListenAndServe(":8001", nil))
// }

// // GetMoviesHandler endpoint handler
// func GetMoviesHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		json.NewEncoder(w).Encode("only POST method")
// 		return
// 	}
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		json.NewEncoder(w).Encode("error read body")
// 		return
// 	}

// 	request := &proto.MovieSearchRequest{}

// 	json.Unmarshal(body, request)
// 	request.Apikey = "faf7e5bb"

// 	if response, err := client.GetMovies(r.Context(), request); err == nil {
// 		json.NewEncoder(w).Encode(response)
// 	} else {
// 		json.NewEncoder(w).Encode(err.Error())
// 	}
// }

// ================

func main() {
	errChan := make(chan error)

	var svc movie.Service
	svc = movie.MovieService{}
	endpoint := movie.Endpoints{
		SearchMovieEndpoint: movie.MakeSearchMovieEndpoint(svc),
	}

	// Logging domain.
	var logger logkit.Logger

	r := movie.MakeHttpHandler(endpoint, logger)

	// HTTP transport
	go func() {
		fmt.Println("Starting server at port 8001")
		handler := r
		errChan <- http.ListenAndServe(":8001", handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<-errChan)
}

package main

import (
	"go_rpc_kit_gateway_example/proto"
	"go_rpc_kit_gateway_example/server/omdb"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	server := omdb.Server{}

	srv := grpc.NewServer()
	proto.RegisterOMDBapiServer(srv, &server)
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}

}

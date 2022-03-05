package main

import (
	"fmt"
	"google.golang.org/grpc"
	"micro-client/communicate"
	"micro-client/controller"
	"net"
)

func main() {

	// port := os.Getenv("PORT")
	port := 5000
	host := fmt.Sprintf("0.0.0.0:%v", port)

	listener, err := net.Listen("tcp", host)

	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	communicate.RegisterClientCommunicateServer(grpcServer, &controller.ClientServer{})
	communicate.RegisterDestinationCommunicateServer(grpcServer, &controller.DestinationServer{})
	communicate.RegisterProductCommunicateServer(grpcServer, &controller.ProductServer{})

	fmt.Printf("[x] - Server logistic listen http://localhost:%v\n", port)

	if err := grpcServer.Serve(listener); err != nil {
		panic(err.Error())
	}
}

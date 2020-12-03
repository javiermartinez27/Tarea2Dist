package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/tutorialedge/go-grpc-tutorial/chat"
	"google.golang.org/grpc"
)

var propuesta string

func escuchar() {
	puerto := 9000
	fmt.Println("DataNode escuchando en el puerto " + strconv.Itoa(puerto))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", puerto))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := chat.Server{}

	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func main() {

	escuchar()

}

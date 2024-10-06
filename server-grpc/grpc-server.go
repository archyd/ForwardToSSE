package main

import (
	"log"
	"net"
	"time"

	pb "stream.go/protos"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct{}

func (s *server) InitTimer(a *emptypb.Empty, srv pb.TimeService_InitTimerServer) error {
	log.Println("init stream")

	for x := range time.Tick(1 * time.Second) {
		resp := pb.ResponseMessage{Sometext: "timetick", Timestamp: x.String()}

		if err := srv.Send(&resp); err != nil {
			log.Println("error generating response")
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":5001") //the grpc server start on port 5001
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	//Create new gRPC server handler
	serv := grpc.NewServer()
	pb.RegisterTimeServiceServer(serv, &server{})

	log.Println("start server")

	if err := serv.Serve(lis); err != nil {
		panic("error building server: " + err.Error())
	}
}

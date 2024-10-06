package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	pb "stream.go/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	//init stream
	conn, err := grpc.NewClient(":5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}
	client := pb.NewTimeServiceClient(conn)
	streamer, err := client.InitTimer(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatal(err) // dont use panic in your real project
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
		w.Header().Set("Transfer-Encoding", "identity")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		dataCh := make(chan string)
		_, cancel := context.WithCancel(r.Context())
		defer cancel()
		go func() {
			for data := range dataCh {
				w.Write([]byte(data))
				w.(http.Flusher).Flush()
			}
		}()
		for {
			resp, err := streamer.Recv()
			if err != nil {
				log.Fatalf("cannot receive %v", err)
				//<-closeNotify.Done()
				return
			}
			a, b := resp.GetSometext(), resp.GetTimestamp()
			dataCh <- fmt.Sprintf("event: %s\ndata: %v\n\n", a, b)
		}
	})
	http.ListenAndServe(":8000", nil)
	//log.Printf("finished")
}

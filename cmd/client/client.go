package main

import (
	"context"
	"fmt"
	"go-aluno/pb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	// AddUser(client)

	// AddUserVerbose(client);

	// AddUsers(client)

	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Rodolfo",
		Email: "rodolfo@email.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Rodolfo",
		Email: "rodolfo@email.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}

		fmt.Println("Status:", stream.Status)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "1",
			Name:  "rodolfo",
			Email: "email",
		},
		{
			Id:    "2",
			Name:  "rodolfo2",
			Email: "email2",
		},
		{
			Id:    "3",
			Name:  "rodolfo3",
			Email: "email3",
		},
		{
			Id:    "4",
			Name:  "rodolfo4",
			Email: "email4",
		},
		{
			Id:    "5",
			Name:  "rodolfo5",
			Email: "email5",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving request: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	reqs := []*pb.User{
		{
			Id:    "1",
			Name:  "rodolfo",
			Email: "email",
		},
		{
			Id:    "2",
			Name:  "rodolfo2",
			Email: "email2",
		},
		{
			Id:    "3",
			Name:  "rodolfo3",
			Email: "email3",
		},
		{
			Id:    "4",
			Name:  "rodolfo4",
			Email: "email4",
		},
		{
			Id:    "5",
			Name:  "rodolfo5",
			Email: "email5",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error when receiving: %v", err)
				break
			}
			fmt.Printf("Recebendo user %v, com status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}

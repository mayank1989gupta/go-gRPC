package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"../greetpb"

	"google.golang.org/grpc"
)

type server struct{}

//from the generated file from proto
func (*server) Greet(context context.Context, r *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet func was invoked with request: %v/n", r)
	firstName := r.GetGreeting().GetFirstName()
	//creating the response
	response := greetpb.GreetResponse{
		Response: "Hello" + firstName,
	}

	return &response, nil
}

// GreetManyTimes : Server streaming API
func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes() function invoked: %v", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	// Creating response from request &, adding to stream for server streaming!
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " " + lastName + ": " + strconv.Itoa(i)
		res := greetpb.GreetManyTimesResponse{
			Response: result,
		}

		stream.Send(&res)           //sending data to stream.
		time.Sleep(1 * time.Second) //just to show streaming
	}

	return nil //
}

// LongGreet -> Server implementation for client streaming
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet() function invoked")
	result := "Hello "
	for {
		//Req will of type : LongGreetRequest
		req, err := stream.Recv()
		if err == io.EOF {
			//Finished reading client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{Response: result})
		}

		if err != nil {
			log.Fatalf("Error while reading the stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += firstName + ": "
	}
}

func main() {
	fmt.Println("Hello!!")
	//50051 is the default port for grpc - binding port: 50051
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	//Creating new gRPC server
	s := grpc.NewServer()

	// RegisterGreetServiceServer from the generated go file.
	// Registering the ser
	greetpb.RegisterGreetServiceServer(s, &server{})

	// Binding the port to grpc server
	// Below code if short hand code for if in go
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

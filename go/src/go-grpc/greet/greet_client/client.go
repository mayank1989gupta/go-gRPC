package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"../greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello! Client")

	// we will use inscure as by default grpc has ssl
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}

	// As we want the connection to be closed, but we want it to be used in the end hence using defer
	// when whole main will be done the defer statement would be called
	defer conn.Close()

	// NewGreetServiceClient - takes the client connection
	// here c is the service client
	c := greetpb.NewGreetServiceClient(conn)
	doUnary(c) //unary function call to server

	// Server Streaming func call to server
	doServerStreaming(c)
}

// Method to perform Unary operation
func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Invkoing the Unary RPC")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mayank",
			LastName:  "Gupta",
		},
	}
	//Auto generated method
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Greet rpc: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Response)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Invking the Server Streaming RPC func")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mayank",
			LastName:  "Gupta",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while invkoing server streaming: %v", err)
	}
	//This returns the stream of data &, once it reaches the end -> it gives EOF error
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we have reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading the stream: %v", err)
		}

		//All is well
		log.Printf(msg.GetResponse())
	}

}

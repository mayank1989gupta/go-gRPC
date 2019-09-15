package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"../greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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
	//doUnary(c) //unary function call to server

	// Server Streaming func call to server
	//doServerStreaming(c)

	//Client Streaming
	//doClientStreaming(c)

	//Client: Bi Directional Streaming
	//doBiDiStreaming(c)

	// Deeadlines - To test the server is implemented to sleep for 3 seconds!
	doUnaryWithDeadline(c, 5*time.Second) // should complete
	doUnaryWithDeadline(c, 1*time.Second) // should timeout
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

// func to implement Server Streaming
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

// func to implement Client Streaming
func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Client streaming RPC call")
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mayank",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ABC",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "XYZ",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ABCXYZ",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ABCAXYTEST",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}

	//Sending the requests
	for _, req := range requests {
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while recieving response for the requested stream: %v", err)
	}

	fmt.Printf("Response from server: %v", response)
	fmt.Println("")
}

// doBiDiStreaming streming
func doBiDiStreaming(c greetpb.GreetServiceClient) {

	//To call: GreetEveryone(ctx context.Context) (GreetService_GreetEveryoneClient, error)

	fmt.Println("Client streaming RPC call")
	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mayank",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ABC",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "XYZ",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ABCXYZ",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ABCAXYTEST",
			},
		},
	}

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while reading stream from server: %v", err)
		return
	}

	//channel
	waitc := make(chan struct{})

	// Func to send bunch of messages: using go routines
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending req: %v", req)
			stream.Send(req)
			time.Sleep(time.Second)
		}

		stream.CloseSend() // once all requests are done
	}()

	// to recieve bunch of messages
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf(" Error while reading response: %v", err)
			}

			fmt.Printf("Response from server: %v\n", res.GetResponse())
		}

		close(waitc) // once the waitc is closed it will unblock
	}()

	// block until everything is done - using channels
	<-waitc // we will for channel to be closed
}

// doUnaryWithDeadline - Unary RPC with Deadline
func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Invkoing the Unary RPC")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mayank",
			LastName:  "Gupta",
		},
	}
	//Auto generated method
	// building on top of parent context background
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, req)

	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeline was exceeded!")
			} else {
				fmt.Printf("Unexpected Error: %v", statusErr)
			}
		} else {
			log.Fatalf("Error while calling Greet Deadline RPC: %v", err)
		}
		return
	}

	log.Printf("Response from Greet: %v", res.Result)
}

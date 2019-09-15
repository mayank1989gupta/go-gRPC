package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"../calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello! Calculator Client!!")

	// we will use inscure as by default grpc has ssl
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
	}

	defer conn.Close()

	// NewGreetServiceClient - takes the client connection
	// here c is the service client
	c := calculatorpb.NewCalculatorServiceClient(conn)
	//Unary RPC Calls
	//doSumUnary(c)
	//doSubsUnary(c)

	//Server Streaming RPC Calls
	//doServerStreaming(c)

	// Client Stream RPC
	//doClientStreaming(c)

	// Bi Directional Streaming
	//doBiDiStreaming(c)

	// Error Unary
	doErrorUnary(c)
}

//Util func to call the method implemented on server side
func doSumUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Invoking the Calculator Sum Unary Call!")

	req := calculatorpb.SumRequest{
		First_Number:  10,
		Second_Number: 20,
	}

	res, err := c.Sum(context.Background(), &req)

	if err != nil {
		log.Fatalf("Error while calling server: %v", err)
	}

	log.Printf("Response from server for Sum call: %v", res)
}

//method to implement call to substract service
func doSubsUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Invoking the Calculator Substract Unary Call!")

	req := calculatorpb.SubstractRequest{
		First_Number:  30,
		Second_Number: 10,
	}

	res, err := c.Substract(context.Background(), &req)

	if err != nil {
		log.Fatalf("Error while executing the substract req: %v", err)
	}

	log.Printf("Response from server for substarct call: %v", res)
}

// doServerStreaming - Server Streaming
func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	//PrimeNumberDecomposition

	fmt.Printf("Invkoing the client call for server streaming")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12,
	}
	//Invkoking the client call
	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while invoking server streaming request: %v", err)
	}
	for {
		//MSG is of type: PrimeNumberDecompositionResponse
		msg, err := resStream.Recv()

		if err == io.EOF {
			//Reached end of stream
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		// No errors
		log.Printf("Factor: %v", msg.GetPrimeFactor())
	}
}

// doClientStreaming - Client Streaming RPC
func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Client Streaming RPC for Calculator Service")

	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 4,
		},
	}
	// stream
	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Error while recieving result from server: %v", err)
	}

	for _, req := range requests {
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}
	// Recieving the response
	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error from server: %v", err)
	}

	//All good
	fmt.Println("Response from server: ", res)
}

// Bi Directioanl Streaming - client
func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Bi Directaional Streaming RPC for Calculator Service")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while reading stream from server: %v", err)
		return
	}

	// channel
	waitc := make(chan struct{})

	// Sending stream to server - using go routines
	go func() {
		numbers := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			fmt.Printf("Sending Req to server: %v \n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{Number: number})
			time.Sleep(time.Second)
		}
	}()

	// Recieving stream from server - using go routines
	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error while reading stream from server: %v", err)
				return
			}

			fmt.Printf("Maximum Number: %v\n", res)
		}
		close(waitc) //unblocking the process
	}()

	// block until everything is done - using channels
	<-waitc // we will for channel to be closed
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Invkoing doErrorUnary()")
	// Correct call
	checkError(c, 16)

	//Error Call
	checkError(c, -1)
}

// checkError Util func
func checkError(c calculatorpb.CalculatorServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: number})
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			// atcual error from GRPC
			if resErr.Code() == codes.InvalidArgument {
				fmt.Printf("%v\n", resErr.Message())
			}
		} else {
			// BIG Error calling sqrt
			log.Fatalf("Error while calling server RPC: %v\n", err)
		}
	}

	fmt.Printf("Square Root: %v\n", res)
}

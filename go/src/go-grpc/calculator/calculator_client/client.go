package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"../calculatorpb"

	"google.golang.org/grpc"
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
	doSumUnary(c)
	doSubsUnary(c)

	//Server Streaming RPC Calls
	doServerStreaming(c)
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

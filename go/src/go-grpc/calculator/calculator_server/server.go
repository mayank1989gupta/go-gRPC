package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"../calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

//Implementing the sum method for the interface
func (*server) Sum(cont context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	sum := req.GetFirst_Number() + req.GetSecond_Number()
	//Converting to response type
	response := calculatorpb.SumResponse{
		Sum_Result: sum,
	}

	return &response, nil
}

// Implementing the Substract method for the interface
func (*server) Substract(ctx context.Context, req *calculatorpb.SubstractRequest) (*calculatorpb.SubstractResponse, error) {
	subs := req.GetFirst_Number() - req.GetSecond_Number()

	res := calculatorpb.SubstractResponse{
		Substract_Result: subs,
	}

	return &res, nil
}

// Implementing ComputeAverage() from the interface
func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("Inside ComputeAverage() - server")
	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: float64(sum) / float64(count),
			})
		}

		if err != nil {
			log.Fatalf("Error while reading the stream from client: %v", err)
		}
		sum += req.GetNumber()
		count++
	}
}

// Server streaming func implementation
func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest,
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Invkoing the Prime Number Decomposition RPC: %v", req)
	number := req.GetNumber()
	//Prime Number decomposition
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			//Adding to stream
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{PrimeFactor: divisor})
			number = number / divisor
		} else {
			divisor++
			fmt.Println("Icrementing divisor")
		}
	}

	return nil
}

// FindMaximum - Bi Di Streaming
func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("Invoking FindMaximum() - Bi-Directional Streaming!")
	maximum := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading stream from client: %v", err)
			return err
		}

		number := req.GetNumber()
		if number > maximum {
			maximum = number
			// Sending to stream once max found
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{Maximum: maximum})
			if sendErr != nil {
				log.Fatalf("Error while sending response to client: %v", sendErr)
				return sendErr
			}
		}
	}
}

func main() {
	fmt.Println("Hello!! This is calculator gRPC!")

	//50051 is the default port for grpc - binding port: 50051
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//creating grpc server
	s := grpc.NewServer()

	//Register the service client
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// Binding the port to grpc server
	// Below code if short hand code for if in go
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

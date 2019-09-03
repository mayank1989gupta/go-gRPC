package main

import (
	"context"
	"fmt"
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

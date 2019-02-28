package main

import (
	"context"
	"fmt"

	pb "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := pb.NewValidatorServiceClient(conn)

	ctx := context.Background()

	for {
		req := &pb.AllValidatorEpochAssignmentsRequest{}

		res, err := client.AllValidatorEpochAssignments(ctx, req)

		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", res)
		return
	}

}

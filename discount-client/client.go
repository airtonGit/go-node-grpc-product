package main

import (
	"context"
	"log"
	"time"

	pb "github.com/airtonGit/go-node-grpc-product/product"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func Discount(userID, productID string) (float32, int32, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
		return 0, 0, err
	}
	defer conn.Close()
	c := pb.NewDiscountServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Discount(ctx, &pb.DiscountRequest{ProductId: "1", UserId: "1"})
	if err != nil {
		log.Printf("could not ask discount: %v", err)
		return 0, 0, err
	}
	log.Printf("Discount: %v %v", r.Pct, r.ValueInCents)
	return r.Pct, r.ValueInCents, nil
}

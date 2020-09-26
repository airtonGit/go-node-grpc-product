package main

import (
	"context"
	"log"
	"time"

	pb "github.com/airtonGit/go-node-grpc-product/product"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

type meuDiscountServer struct {
}

func Discount(userID, productID int) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDiscountClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Discount(ctx, &pb.DiscountRequest{ProductId: 1, UserId: 1})
	if err != nil {
		log.Fatalf("could not ask discount: %v", err)
	}
	log.Printf("Discount: %v %v", r.Pct, r.ValueInCents)
}

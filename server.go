package main

import (
	"context"

	pb "github.com/airtonGit/go-node-grpc-product/product"
)

type discountServer struct {
}

func (s *discountServer) Discount(ctx context.Context, req *pb.DicountRequest) *pb.DiscountResponse {

}

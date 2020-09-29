package dc

import (
	"context"
	"errors"
	"log"
	"time"

	pb "github.com/airtonGit/go-node-grpc-product/product"
	"google.golang.org/grpc"
)

type DiscountClient struct {
	address           string
	conn              *grpc.ClientConn
	timeout           time.Duration
	pbDiscountService pb.DiscountServiceClient
}

func NewDiscountClient(address string, timeout time.Duration) *DiscountClient {
	return &DiscountClient{address: address, timeout: timeout}
}

func (c *DiscountClient) Close() {
	c.conn.Close()
}

func (c *DiscountClient) Dial() error {
	var err error
	c.conn, err = grpc.Dial(c.address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println("DiscountClient fail grpc.Dial err:", err.Error())
	}
	c.pbDiscountService = pb.NewDiscountServiceClient(c.conn)
	return err
}

//DiscountAsk makes a single request for discount,
func (c *DiscountClient) DiscountAsk(ctx context.Context, userID, productID string) (float32, int32, error) {
	if c.conn == nil {
		return 0, 0, errors.New("no connection, dial first")
	}
	requestCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	r, err := c.pbDiscountService.Discount(requestCtx, &pb.DiscountRequest{ProductId: productID, UserId: userID})
	if err != nil {
		log.Printf("could not ask discount: %v", err)
		return 0, 0, err
	}
	return r.Pct, r.ValueInCents, nil
}

// func Discount(ctx contex.Context, userID, productID string) (pb.DiscountReply, error) {
// 	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
// 	if err != nil {
// 		log.Printf("did not connect: %v", err)
// 		return 0, 0, err
// 	}
// 	defer conn.Close()
// 	c := pb.NewDiscountServiceClient(conn)

// 	ctx, cancel := context.WithTimeout(ctx, time.Second)
// 	defer cancel()
// 	r, err := c.Discount(ctx, &pb.DiscountRequest{ProductId: "1", UserId: "1"})
// 	if err != nil {
// 		log.Printf("could not ask discount: %v", err)
// 		return 0, 0, err
// 	}
// 	log.Printf("Discount: %v %v", r.Pct, r.ValueInCents)
// 	return r.Pct, r.ValueInCents, nil
// }

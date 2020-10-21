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

func (c *DiscountClient) Dial() (func(), error) {
	var err error
	c.conn, err = grpc.Dial(c.address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println("DiscountClient fail grpc.Dial err:", err.Error())
		return func() {}, err
	}
	c.pbDiscountService = pb.NewDiscountServiceClient(c.conn)

	closeFunc := func() {
		c.Close()
	}

	return closeFunc, err
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
		log.Printf("DicountClient could not ask discount: %v", err)
		return 0, 0, err
	}

	return r.GetPct(), r.GetValueInCents(), nil
}

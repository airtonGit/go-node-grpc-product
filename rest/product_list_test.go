package rest

import (
	"context"
	"log"
	"testing"
)

type discountAskFunc func(ctx context.Context, userID string, productID string) (float32, int32, error)

func (t discountAskFunc) DiscountAsk(ctx context.Context, userID string, productID string) (float32, int32, error) {
	return t(ctx, userID, productID)
}

func discountAskMock(ctx context.Context, userID string, productID string) (float32, int32, error) {
	//time.Sleep(75 * time.Millisecond)
	log.Println("discountAskMock response")
	return 10, 150, nil
}

func TestAskDiscount(t *testing.T) {
	var discountMock discountAskFunc
	discountMock = discountAskMock
	ctx := context.TODO() //context.WithTimeout(context.TODO(), 15*time.Second)
	got := askDiscount(ctx, discountMock, productsFixtures, "1")
	t.Fatal("Forced fail: ", got)
}

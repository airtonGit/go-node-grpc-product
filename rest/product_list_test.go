package rest

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

type discountClientMock func(ctx context.Context, userID string, productID string) (float32, int32, error)

func (t discountClientMock) DiscountAsk(ctx context.Context, userID string, productID string) (float32, int32, error) {
	return t(ctx, userID, productID)
}

var discountMockFunc = func(ctx context.Context, userID string, productID string) (float32, int32, error) {
	time.Sleep(70 * time.Millisecond)
	log.Println("discountAskMock response")
	return 10, 150, nil
}

func generateJobs() <-chan jobEnvelope {
	rc := make(chan jobEnvelope, 10)
	go func() {
		defer close(rc)

		for _, prod := range productsFixtures {
			p := jobEnvelope{"1", prod}
			rc <- p
		}
	}()
	return rc
}

func generateResults() <-chan product {
	rc := make(chan product, 10)
	go func() {
		defer close(rc)

		for _, prod := range productsFixtures {
			p := prod
			p.Discount.Percent = 0.9
			p.Discount.ValueInCents = 120
			rc <- p
		}
	}()
	return rc
}

func TestDiscountReceiver(t *testing.T) {
	type input struct {
		ctx     context.Context
		prodMap map[string]product
		results <-chan product
		done    chan struct{}
	}
	for _, testcase := range []struct {
		name      string
		input     input
		ctxCancel context.CancelFunc
		want      string
	}{
		{
			name: "Test results closed finish",
			input: input{
				prodMap: make(map[string]product),
				results: generateResults(),
				done:    make(chan struct{}),
			},
			want: "results ch closed",
		},
		{
			name: "Ctx done",
			input: input{
				prodMap: make(map[string]product),
				results: generateResults(),
				done:    make(chan struct{}),
			},
			want: "ctx.done",
		},
		{
			name: "products map len1",
			input: input{
				prodMap: make(map[string]product),
				results: generateResults(),
				done:    make(chan struct{}),
			},
			want: "len3",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.input.ctx, testcase.ctxCancel = context.WithCancel(context.Background())
			switch testcase.name {
			case "Test results closed finish":
				err := discountsReceiver(
					testcase.input.ctx,
					testcase.input.prodMap,
					testcase.input.results, testcase.input.done)
				got := err.Error()
				want := testcase.want
				//wait 3s for done
				tmout := time.After(300 * time.Second)
				shouldStop := false
				for !shouldStop {
					select {
					case _, ok := <-testcase.input.done:
						if !ok && err.Error() != "results ch closed" {
							t.Fatalf("Err return want %s got %s", want, got)
						}
						//test ok
						shouldStop = true
					case <-tmout:
						t.Fatal("Discount want done closed in 3s, take too long, timeout!")
					}
				}
			case "Ctx done":
				testcase.ctxCancel()
				err := discountsReceiver(testcase.input.ctx, testcase.input.prodMap, testcase.input.results, testcase.input.done)
				got := err.Error()
				want := testcase.want

				if err.Error() != "ctx.done" {
					t.Fatalf("Err return want %s got %s", want, got)
				}

				//test ok
			case "products map len1":
				err := discountsReceiver(
					testcase.input.ctx,
					testcase.input.prodMap,
					testcase.input.results, testcase.input.done)
				got := err.Error()
				want := testcase.want
				//pre-condition
				if got != "results ch closed" {
					t.Fatalf("want %s got %s", "results ch closed", got)
				}
				//test
				gotMap := testcase.input.prodMap
				gotStr := fmt.Sprintf("len%d", len(gotMap))
				if gotStr != want {
					t.Fatalf("got %s want %s", gotStr, want)
				}
			}

		})

	}
}

func TestWorker(t *testing.T) {
	type input struct {
		cMock   discountClientMock
		prodMap map[string]product
		jobs    <-chan jobEnvelope
		results chan<- product
		done    chan struct{}
	}
	type want struct {
		desc    string
		prodMap map[string]product
	}
	for _, testcase := range []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "Ctx done",
			input: input{
				cMock: discountMockFunc,
				prodMap: map[string]product{
					"1": {
						ID: "1",
					},
				},
				jobs:    generateJobs(),
				results: make(chan product, 10),
				done:    make(chan struct{}),
			},
			want: want{
				desc: "ctx.done",
				prodMap: map[string]product{
					"1": {
						ID: "1",
					},
				},
			},
		},
		{
			name: "normal finish",
			input: input{
				cMock: discountMockFunc,
				prodMap: map[string]product{
					"1": {
						ID: "1",
					},
				},
				jobs:    generateJobs(),
				results: make(chan product, 10),
				done:    make(chan struct{}),
			},
			want: want{
				desc: "normal finish len3",
				prodMap: map[string]product{
					"1": {
						ID: "1",
					},
				},
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			switch testcase.name {

			case "Ctx done":
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				cancel()

				err := worker(ctx, testcase.input.cMock, testcase.input.jobs, testcase.input.results)
				got := err.Error()
				want := testcase.want.desc
				if want != got {
					t.Fatalf("want %s got %s", want, got)
				}
			case "normal finish":
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				err := worker(ctx, testcase.input.cMock, testcase.input.jobs, testcase.input.results)
				want := fmt.Sprintf("normal finish len%d", len(testcase.input.prodMap))
				if err != nil && want != err.Error() {
					t.Fatalf("want %s got %s", want, err.Error())
				}
			}

		})
	}
}

func TestAskDiscount(t *testing.T) {
	type input struct {
		cMock    discountClientMock
		products []product
		userID   string
	}
	type want struct {
		dsc   float32
		value int32
	}
	for _, testcase := range []struct {
		name  string
		input input
		want  []want
	}{
		{
			name: "Test1",
			input: input{
				cMock:    discountMockFunc,
				products: productsFixtures,
				userID:   "1",
			},
			want: []want{
				{
					dsc:   0.85,
					value: 150,
				},
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			got := askDiscount(context.TODO(), testcase.input.cMock, testcase.input.products, testcase.input.userID)

			if want, have := len(testcase.want), len(got); want != have {
				t.Errorf("Results want %d have %d", want, have)
			}
		})
	}
}

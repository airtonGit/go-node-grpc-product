package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/airtonGit/go-node-grpc-product/dc"
)

const (
	discountAddress    = "localhost:50051"
	askDiscountTimeout = time.Millisecond * 100
)

func jsonEncode(products []product) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(products); err != nil {
		log.Println("jsonEncode fail err:", err.Error())
		return nil, err
	}
	return buf.Bytes(), nil
}

type jobParams struct {
	UserID  string
	Product product
}

func worker(ctx context.Context, clt discountClient, jobs <-chan jobParams, results chan<- product) {
	for j := range jobs {
		pResp := j.Product
		var err error
		pResp.Discount.Percent, pResp.Discount.ValueInCents, err = clt.DiscountAsk(ctx, j.UserID, j.Product.ID)
		if err != nil {
			log.Println("worker askDiscount err:", err.Error())
		}
		results <- pResp
	}
}

type discountClient interface {
	DiscountAsk(context.Context, string, string) (float32, int32, error)
}

func askDiscount(ctx context.Context, dscClient discountClient, products []product, userID string) []product {
	var (
		resp           []product
		respMap        map[string]product
		wgWorkers      sync.WaitGroup
		allResultsDone chan struct{}
	)
	allResultsDone = make(chan struct{})
	respMap = make(map[string]product)

	for _, prod := range products {
		respMap[prod.ID] = prod
	}

	results := make(chan product, 50)
	discountAskCh := make(chan jobParams, 50)
	//workers results collection
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("results agregation gr ctx.done")
				return
			case prod, ok := <-results:
				if !ok {
					log.Println("results agregation ch closed")
					allResultsDone <- struct{}{}
					return
				}
				//updating map with discount
				respMap[prod.ID] = prod
			default:
				log.Println("workser result aggregation sleep waiting results")
				time.Sleep(time.Second)
			}
		}
	}()

	//4 workers launch
	for i := 0; i < 3; i++ {
		go func() {
			wgWorkers.Add(1)
			defer wgWorkers.Done()
			worker(ctx, dscClient, discountAskCh, results)
		}()
	}

	//calling for discount
	shouldStop := false
	for i := 0; i < len(products) && !shouldStop; i++ {
		select {
		case <-ctx.Done():
			log.Println("ctx.done")
			shouldStop = true
			continue
		default:
			job := jobParams{userID, products[i]}
			discountAskCh <- job
		}

	}
	close(discountAskCh)
	log.Println("WaitGroup...")
	wgWorkers.Wait()
	log.Println("WaitGroup...done")
	close(results)
	log.Println("waitin AllResultsDone")
	<-allResultsDone
	log.Println("waitin AllResultsDone done.")

	//if timeout happens responds with zero discount
	for _, prod := range respMap {
		resp = append(resp, prod)
	}
	return resp
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	innerCtx, innCanel := context.WithTimeout(r.Context(), 250*time.Millisecond)
	defer innCanel()
	dscClt := dc.NewDiscountClient(discountAddress, askDiscountTimeout)
	askDiscount(innerCtx, dscClt, productsFixtures, "1")
}

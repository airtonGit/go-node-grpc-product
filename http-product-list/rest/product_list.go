package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/airtonGit/go-node-grpc-product/dc"
)

const (
	askDiscountTimeout = time.Millisecond * 5000
	resquestTimeout    = time.Millisecond * 5000
)

//jsonEncode jsonfy the response product list
func jsonEncode(products []product) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(products); err != nil {
		log.Println("jsonEncode fail err:", err.Error())
		return nil, err
	}
	return buf.Bytes(), nil
}

type jobEnvelope struct {
	UserID  string
	Product product
}

func worker(ctx context.Context, clt discountClient, jobs <-chan jobEnvelope, results chan<- product) error {
	for j := range jobs {
		select {
		case <-ctx.Done():
			return errors.New("ctx.done")

		default:

		}
		pResp := j.Product
		var err error
		pResp.Discount.Percent, pResp.Discount.ValueInCents, err = clt.DiscountAsk(ctx, j.UserID, j.Product.ID)
		if err != nil {
			log.Println("worker askDiscount err:", err.Error())
			continue
		}
		results <- pResp
	}
	return nil
}

type discountClient interface {
	DiscountAsk(context.Context, string, string) (float32, int32, error)
}

//what I have done wihd
func discountsReceiver(ctx context.Context, productsMap map[string]product, results <-chan product, done chan<- struct{}) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("results agregation gr ctx.done err:", ctx.Err())
			close(done)
			return errors.New("ctx.done")
		case prod, ok := <-results:
			if !ok {
				log.Println("results agregation ch closed")
				close(done)
				return errors.New("results ch closed")
			}
			productsMap[prod.ID] = prod
		default:
			log.Println("workser result aggregation sleep waiting results")
			time.Sleep(time.Second)
		}
	}
}

//askDiscount makes a grpc call to discount service for each product.
//Each grpc request respect the provided timeout, in case of timeout the discount will be zero.
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
	discountAskCh := make(chan jobEnvelope, 50)
	//workers results collection
	go discountsReceiver(ctx, respMap, results, allResultsDone)

	//4 workers launch
	wgWorkers.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
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
			job := jobEnvelope{userID, products[i]}
			discountAskCh <- job
		}

	}
	close(discountAskCh)

	//waiting the 4 workers finish
	wgWorkers.Wait()

	//No more results to receive
	close(results)

	//waiting results maping
	<-allResultsDone

	//if timeout happens responds with zero discount
	for _, prod := range respMap {
		resp = append(resp, prod)
	}
	return resp
}

func makeProductsHandler(params AppParams) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		innerCtx, innCancel := context.WithTimeout(r.Context(), resquestTimeout)
		defer innCancel()

		dscClt := dc.NewDiscountClient(params.DiscountAddr, askDiscountTimeout)

		closeFunc, err := dscClt.Dial()
		if err != nil {
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
		defer closeFunc()
		products := askDiscount(innerCtx, dscClt, productsFixtures, "1") //"1" is the 'loggedin' user

		buf, err := jsonEncode(products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(buf)
	}
}

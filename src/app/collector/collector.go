package collector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mjedari/vgang-project/src/app/auth"
	"github.com/mjedari/vgang-project/src/app/configs"
	"github.com/mjedari/vgang-project/src/domain/contracts"
	"github.com/mjedari/vgang-project/src/domain/products"
	"github.com/mjedari/vgang-project/src/infra/utils"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Collector struct {
	client  *auth.Client
	storage contracts.IStorage
	config  configs.Collector
}

func NewCollector(client *auth.Client, storage contracts.IStorage, config configs.Collector) *Collector {
	return &Collector{client: client, storage: storage, config: config}
}

func (c *Collector) FetchProducts(ctx context.Context, pagination PaginationQuery) (*products.ProductsResponse, error) {
	path := fmt.Sprintf("%v%v", "/retailers/products?search=shirt&sort=Latest&dont_show_out_of_stock=1", pagination.toString())

	productsRequest := auth.GetRequest{Path: path}
	res, err := c.client.Get(ctx, productsRequest)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		logrus.Error("status code: ", res.StatusCode)
		return nil, errors.New("request was not successful")
	}

	body, err := io.ReadAll(res.Body)

	var productsResponse products.ProductsResponse
	if jsonErr := json.Unmarshal(body, &productsResponse); jsonErr != nil {
		logrus.Fatalf("can not unmarshalling response: %v", jsonErr)
	}

	return &productsResponse, nil
}

func (c *Collector) ConcurrentFetchProducts(ctx context.Context, pagination PaginationQuery) (<-chan []byte, error) {
	path := fmt.Sprintf("%v%v", "/retailers/products?search=shirt&sort=Latest&dont_show_out_of_stock=1", pagination.toString())
	output := make(chan []byte) // be careful fo close the channel

	productsRequest := auth.GetRequest{Path: path}
	res, err := c.client.Get(ctx, productsRequest)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		logrus.Error("status code: ", res.StatusCode)
		return nil, errors.New("request was not successful")
	}

	go func() {
		defer close(output)
		body, _ := io.ReadAll(res.Body)
		output <- body
	}()

	return output, nil
}

func (c *Collector) Start(ctx context.Context) {
	switch c.config.Interval {
	case 0:
		c.run(ctx)
	default:
		c.runConcurrent(ctx)
	}
}

func (c *Collector) run(ctx context.Context) {
	// determine the exact total numbers
	totalItems, err := c.getTotalProductsCount(ctx)
	if err != nil {
		logrus.Errorf("could not make a requet: %v", err)
	}

	query := PaginationQuery{
		Count:  int(totalItems),
		Offset: 0,
	}

	productResponse, err := c.FetchProducts(ctx, query)
	if err != nil {
		logrus.Errorf("could not to fatch products: %v", err)
	}

	productList := products.NewProductList(productResponse.Products)
	productList.AddHash()

	// todo: List is ready here to store in redis: productList

	if storingErr := c.store(ctx, productList); storingErr != nil {
		logrus.Error(err)
	}

}

func (c *Collector) runConcurrent(ctx context.Context) {
	sources := make([]<-chan []byte, 0)

	// determine the exact total numbers
	totalItems, err := c.getTotalProductsCount(ctx)
	if err != nil {
		logrus.Errorf("could not make a requet: %v", err)
	}

	queryList := initiateList(c.config.SplitFactor, int(totalItems))

	fmt.Println(queryList)

	var sensorOutputCh <-chan []byte

	for _, query := range queryList {
		fmt.Println("run query...")
		sensorOutputCh, _ = c.ConcurrentFetchProducts(ctx, query)
		sources = append(sources, sensorOutputCh)
	}

	outputCh := utils.Funnel(sources...)

	output := <-outputCh

	var productsResponse products.ProductsResponse
	if jsonErr := json.Unmarshal(output, &productsResponse); jsonErr != nil {
		logrus.Fatalf("can not unmarshalling response: %v", jsonErr)
	}

	productList := products.NewProductList(productsResponse.Products)
	productList.AddHash()

	// todo: List is ready here to store in redis: productList

	if storingErr := c.store(ctx, productList); storingErr != nil {
		logrus.Error(err)
	}

}

func (c *Collector) store(ctx context.Context, list *products.ProductList) error {
	myMap := make(map[string]string)
	//prepare date
	for _, item := range list.List {
		str, err := item.ToString()
		if err != nil {
			return err
		}

		myMap[item.GetKey()] = str
	}

	err := c.storage.BatchStore(ctx, myMap, 0)
	if err != nil {
		return err
	}

	return nil
}

func (c *Collector) getTotalProductsCount(ctx context.Context) (uint64, error) {
	query := PaginationQuery{
		Count:  1,
		Offset: 0,
	}

	path := fmt.Sprintf("%v%v", "/retailers/products?search=shirt&sort=Latest&dont_show_out_of_stock=1", query.toString())

	productsRequest := auth.GetRequest{Path: path}
	res, err := c.client.Get(ctx, productsRequest)
	if err != nil {
		return 0, err
	}

	if res.StatusCode != http.StatusOK {
		logrus.Error("status code: ", res.StatusCode)
		return 0, errors.New("request was not successful")
	}

	body, err := io.ReadAll(res.Body)

	var productsResponse products.ProductsResponse
	if jsonErr := json.Unmarshal(body, &productsResponse); jsonErr != nil {
		logrus.Fatalf("can not unmarshalling response: %v", jsonErr)
	}

	return productsResponse.Total, nil
}
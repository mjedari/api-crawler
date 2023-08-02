package products

import (
	"encoding/json"
	"fmt"
	"github.com/mjedari/vgang-project/src/app/configs"
	"net"
)

const (
	// Define the character set for the hash (62 characters - a-z, A-Z, 0-9)
	characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	base       = 62
	hashLength = 5
)

type Product struct {
	Id             uint64  `json:"id"`
	Title          string  `json:"title"`
	SellerName     string  `json:"sellerName"`
	SellerCurrency string  `json:"sellerCurrency"`
	SellerID       uint64  `json:"sellerID"`
	MinPrice       float32 `json:"minPrice"`
	MaxPrice       float32 `json:"maxPrice"`
	MinRetailPrice float32 `json:"minRetailPrice"`
	MaxRetailPrice float32 `json:"maxRetailPrice"`
	Stock          uint64  `json:"stock"`
	Link           string  `json:"link"`
	Hash           string  `json:"hash"`
}

func (p *Product) AddLink() {
	address := net.JoinHostPort(configs.Config.Server.Host, configs.Config.Server.Port)
	p.Link = fmt.Sprintf("http://%v/%v", address, p.Hash)
}

func (p *Product) GetKey() string {
	return fmt.Sprintf("%v:%v", "products", p.Hash)
}

func (p *Product) ToString() (string, error) {
	jsonStr, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return string(jsonStr), nil
}

type ProductList struct {
	List []*Product
}

func NewProductList(list []*Product) *ProductList {
	return &ProductList{List: list}
}

func (l *ProductList) AddHash() {
	for i := 0; i < len(l.List); i++ {
		_ = l.List[i].AddShortKey()
	}
}

func (p *Product) generateHash() string {
	hash := make([]byte, hashLength)
	product := *p

	for i := 0; i < hashLength; i++ {
		// Get the remainder when dividing the unique ID by the base
		remainder := product.Id % base
		// Use the remainder to get the corresponding character from the character set
		hash[i] = characters[remainder]
		// Update the unique ID to continue the iteration
		product.Id /= base
	}

	return string(hash)
}

func (p *Product) AddShortKey() error {
	p.Hash = p.generateHash()
	return nil
}

type ProductsResponse struct {
	Products []*Product `json:"products"`
	Total    uint64     `json:"totalCount"`
}

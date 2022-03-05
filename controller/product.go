package controller

import (
	"context"
	"errors"
	"log"
	"micro-logistic/communicate"
	model "micro-logistic/models"
)

type ProductServer struct{}

func (s *ProductServer) ProductListAll(ctx context.Context, request *communicate.ProductListAllRequest) (*communicate.ProductListAllResponse, error) {
	res := &communicate.ProductListAllResponse{}

	var product model.Product

	products, total, err := product.GetProductByNamePaginate(request.Name, request.Page, request.Limit)

	if err != nil {
		return res, err
	}

	data := &communicate.DataProduct{}
	for _, c := range products {
		product := &communicate.Product{}
		product.Id = c.Id
		product.Name = c.Name
		product.Price = c.Price
		product.Nfe = c.Nfe
		product.IdClient = c.Client.Id
		data.Product = append(data.Product, product)
	}

	res.Total = total
	res.Page = request.Page
	res.Limit = request.Limit

	res.Data = data
	return res, nil
}

func (s *ProductServer) ListOneProductById(ctx context.Context, request *communicate.ListOneProductByIdRequest) (*communicate.ListOneProductByIdResponse, error) {
	res := &communicate.ListOneProductByIdResponse{}

	product := model.Product{}

	if err := product.GetById(request.Id); err != nil || product.Id == 0 {
		return res, errors.New("Product not found!")
	}

	res.Product = &communicate.Product{
		Id:       product.Id,
		Name:     product.Name,
		Price:    product.Price,
		Nfe:      product.Nfe,
		IdClient: product.Client.Id,
	}

	return res, nil
}

func (s *ProductServer) CreateProduct(ctx context.Context, request *communicate.CreateProductRequest) (*communicate.CreateProductResponse, error) {
	res := &communicate.CreateProductResponse{}

	product := model.Product{
		Name:   request.Name,
		Price:  request.Price,
		Nfe:    request.Nfe,
		Client: model.Client{Id: request.IdClient},
	}

	if err := product.GetByNameAndAndIdClient(request.Name, request.IdClient); err != nil || product.Id != 0 {
		return res, errors.New("product not duplicated!")
	}

	if err := product.CreateProduct(); err != nil {
		log.Println("err", err)
		return res, errors.New("Error creating product!")
	}

	res.Created = true

	return res, nil
}

func (s *ProductServer) UpdateProductById(ctx context.Context, request *communicate.UpdateProductByIdRequest) (*communicate.UpdateProductByIdResponse, error) {
	res := &communicate.UpdateProductByIdResponse{}

	product := model.Product{}

	if err := product.GetById(request.Id); err != nil || product.Id == 0 {
		return res, errors.New("Product not found!")
	}

	product = model.Product{
		Id:     request.Id,
		Name:   request.Name,
		Price:  request.Price,
		Nfe:    request.Nfe,
		Client: model.Client{Id: request.IdClient},
	}

	if err := product.UpdateProductById(); err != nil {
		return res, errors.New("Erro updating product!")
	}

	res.Updated = true

	return res, nil
}

func (s *ProductServer) DeleteProductById(ctx context.Context, request *communicate.DeleteProductByIdRequest) (*communicate.DeleteProductByIdResponse, error) {
	res := &communicate.DeleteProductByIdResponse{}

	product := model.Product{}

	if err := product.GetById(request.Id); err != nil || product.Id == 0 {
		return res, errors.New("Product not found!")
	}

	product.Id = request.Id

	if err := product.DeleteById(); err != nil {
		return res, errors.New("Erro deleting product!")
	}

	res.Deleted = true

	return res, nil
}

func (s *ProductServer) ValidateProductById(ctx context.Context, request *communicate.ValidateProductByIdRequest) (*communicate.ValidateProductByIdResponse, error) {

	res := &communicate.ValidateProductByIdResponse{}

	product := model.Product{}

	if err := product.GetById(request.IdProduct); err != nil {
		return res, errors.New("Product Id invalid!")
	}

	res.Valid = true

	return res, nil
}

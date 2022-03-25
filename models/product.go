package model

import (
	"fmt"
	"micro-client/db"
	"micro-client/helpers"
	"strconv"
)

type Product struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Nfe    string `json:"nfe"`
	Price  string `json:"price"`
	Client Client `json:"Client"`
}

func (product *Product) GetById(id int64) error {

	sql := db.ConnectDatabase()

	query := `select id, name, price, nfe, id_client from products where id = ? limit 1;`

	requestConfig, err := sql.Query(query, id)

	if err != nil {
		return err
	}

	for requestConfig.Next() {
		var name, price, nfe string
		var id, idClient int64
		_ = requestConfig.Scan(&id, &name, &price, &nfe, &idClient)
		if id != 0 {
			product.Id = id
			product.Name = name
			product.Price = price
			product.Nfe = nfe
			product.Client.Id = idClient
		}
	}

	return nil
}

func (product *Product) GetByNameAndAndIdClient(name string, idClient int64) error {

	sql := db.ConnectDatabase()

	query := `select id, name, price, id_client from products where name = ? and id_client = ? limit 1;`

	requestConfig, err := sql.Query(query, name, idClient)

	if err != nil {
		return err
	}

	for requestConfig.Next() {
		var id, name, price string
		var idCarry, idClient int64
		_ = requestConfig.Scan(&id, &name, &price, &idClient, &idCarry)
		i64, _ := strconv.ParseInt(id, 10, 64)
		product.Id = i64
		product.Name = name
		product.Price = price
		product.Client.Id = idClient
	}
	return nil
}

func (product *Product) DeleteById() error {
	sql := db.ConnectDatabase()

	query := "delete from products where id = ?"

	deleteDeposit, err := sql.Prepare(query)

	if err != nil {
		return err
	}

	_, e := deleteDeposit.Exec(product.Id)

	if e != nil {
		return e
	}

	return nil
}

func (product *Product) UpdateProductById() error {
	sql := db.ConnectDatabase()

	query := "update products set name = ?, price = ?, nfe = ?, id_client = ? where id = ?"

	destinationUpdate, err := sql.Prepare(query)

	if err != nil {
		return err
	}

	_, e := destinationUpdate.Exec(product.Name, product.Price, product.Nfe, product.Client.Id, product.Id)

	if e != nil {
		return e
	}

	return nil
}

func (product *Product) CreateProduct() error {
	sql := db.ConnectDatabase()

	query := "insert into products (`name`, price, nfe, id_client) values (?, ?, ?, ?);"

	createDestination, err := sql.Prepare(query)

	if err != nil {
		return err
	}

	_, e := createDestination.Exec(product.Name, product.Price, product.Nfe, product.Client.Id)

	if e != nil {
		return e
	}

	return nil
}

func (product *Product) GetProductByNamePaginate(name string, page, limit int64) ([]Product, int64, error) {
	var productArray []Product
	var total int64

	name = "%" + name + "%"

	sql := db.ConnectDatabase()

	paginate := helpers.Paginate{
		Page:  page,
		Limit: limit,
	}

	paginate.PaginateMounted()
	paginate.MountedQuery("products")

	query := fmt.Sprintf("select id, name, price, nfe, id_client, %v from products where name like ? LIMIT ? OFFSET ?;", paginate.Query)

	requestConfig, err := sql.Query(query, name, paginate.Limit, paginate.Page)

	if err != nil {
		return productArray, total, err
	}

	for requestConfig.Next() {
		productGet := Product{}
		var name, price, nfe string
		var id, idClient int64
		_ = requestConfig.Scan(&id, &name, &price, &nfe, &idClient, &total)
		if id != 0 {
			productGet.Id = id
			productGet.Name = name
			productGet.Price = price
			productGet.Nfe = nfe
			productGet.Client.Id = idClient
			productArray = append(productArray, productGet)
		}

	}

	return productArray, total, nil
}

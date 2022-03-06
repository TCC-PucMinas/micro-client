package model

import (
	"encoding/json"
	"fmt"
	"micro-client/db"
	"micro-client/helpers"
	"strconv"
	"time"
)

type Product struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Nfe    string `json:"nfe"`
	Price  string `json:"price"`
	Client Client `json:"Client"`
}

var (
	keyProductRedisGetById          = "key-product-get-by-id"
	keyProductRedisGetByName        = "key-product-get-by-name"
	keyProductRedisGetByNameAndPage = "key-product-get-by-name-and-page"
)

func setRedisCacheProductgGetById(product *Product) error {
	db, err := db.ConnectDatabaseRedis()

	if err != nil {
		return err
	}

	json, err := json.Marshal(product)

	if err != nil {
		return err
	}
	key := fmt.Sprintf("%v - %v", keyProductRedisGetById, product.Id)

	return db.Set(key, json, 1*time.Minute).Err()
}

func getProductRedisCacheGetOneById(id int64) (Product, error) {
	product := Product{}

	redis, err := db.ConnectDatabaseRedis()

	if err != nil {
		return product, err
	}

	key := fmt.Sprintf("%v - %v", keyProductRedisGetById, id)

	value, err := redis.Get(key).Result()

	if err != nil {
		return product, err
	}

	if err := json.Unmarshal([]byte(value), &product); err != nil {
		return product, err
	}

	return product, nil
}

func getProductRedisCacheGetOneByNameAndIdClient(name string, idClient int64) (Product, error) {
	product := Product{}

	redis, err := db.ConnectDatabaseRedis()

	if err != nil {
		return product, err
	}

	key := fmt.Sprintf("%v - %v - %v", keyProductRedisGetByName, name, idClient)

	value, err := redis.Get(key).Result()

	if err != nil {
		return product, err
	}

	if err := json.Unmarshal([]byte(value), &product); err != nil {
		return product, err
	}

	return product, nil
}

func setRedisCacheProductGetByName(product *Product) error {
	redis, err := db.ConnectDatabaseRedis()

	if err != nil {
		return err
	}

	marshal, err := json.Marshal(product)

	if err != nil {
		return err
	}
	key := fmt.Sprintf("%v - %v - %v", keyProductRedisGetByName, product.Name, product.Client.Id)

	return redis.Set(key, marshal, 1*time.Minute).Err()
}

func getProductRedisCacheGetOneByNamePaginate(name string, page, limit int64) ([]Product, error) {
	var product []Product

	redis, err := db.ConnectDatabaseRedis()

	if err != nil {
		return product, err
	}

	key := fmt.Sprintf("%v - %v - %v - %v", keyProductRedisGetByNameAndPage, name, page, limit)

	value, err := redis.Get(key).Result()

	if err != nil {
		return product, err
	}

	if err := json.Unmarshal([]byte(value), &product); err != nil {
		return product, err
	}

	return product, nil
}

func setRedisCacheProductGetByPaginateByName(name string, page, limit int64, product []Product) error {
	redis, err := db.ConnectDatabaseRedis()

	if err != nil {
		return err
	}

	marshal, err := json.Marshal(product)

	if err != nil {
		return err
	}

	key := fmt.Sprintf("%v - %v - %v - %v", keyProductRedisGetByNameAndPage, name, page, limit)

	return redis.Set(key, marshal, 1*time.Minute).Err()
}

func (product *Product) GetById(id int64) error {

	if c, err := getProductRedisCacheGetOneById(id); err == nil {
		product = &c
		return nil
	}

	sql := db.ConnectDatabase()

	query := `select id, name, price, nfe, id_client from products where id = ? limit 1;`

	requestConfig, err := sql.Query(query, id)

	if err != nil {
		return err
	}

	for requestConfig.Next() {
		var name, price, nfe string
		var id int64
		var idCarry, idClient int64
		_ = requestConfig.Scan(&id, &name, &price, &nfe, &idClient, &idCarry)
		product.Id = id
		product.Name = name
		product.Price = price
		product.Nfe = nfe
		product.Client.Id = idClient
	}

	_ = setRedisCacheProductgGetById(product)

	return nil
}

func (product *Product) GetByNameAndAndIdClient(name string, idClient int64) error {

	if c, err := getProductRedisCacheGetOneByNameAndIdClient(name, idClient); err == nil {
		product = &c
		return nil
	}

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

	_ = setRedisCacheProductGetByName(product)

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

	if c, err := getProductRedisCacheGetOneByNamePaginate(name, page, limit); err == nil {
		productArray = c
		return productArray, total, nil
	}

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

	_ = setRedisCacheProductGetByPaginateByName(name, page, limit, productArray)

	return productArray, total, nil
}

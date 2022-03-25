package model

import (
	"fmt"
	"micro-client/helpers"
	"strconv"

	"micro-client/db"
)

type Client struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func (client *Client) GetById(idClient int64) error {

	sql := db.ConnectDatabase()

	query := `select id, name, email, phone from clients where id = ? limit 1;`

	requestConfig, err := sql.Query(query, idClient)

	if err != nil {
		return err
	}

	for requestConfig.Next() {
		var name, email, phone string
		var id int64
		_ = requestConfig.Scan(&id, &name, &email, &phone)
		client.Id = id
		client.Name = name
		client.Email = email
		client.Phone = phone
	}

	return nil
}

func (client *Client) GetByNameAndEmail() error {

	sql := db.ConnectDatabase()

	query := `select id, name, email, phone from clients where name = ? and email = ? limit 1;`

	requestConfig, err := sql.Query(query, client.Name, client.Email)

	if err != nil {
		return err
	}

	for requestConfig.Next() {
		var id, name, email, phone string
		_ = requestConfig.Scan(&id, &name, &email, &phone)
		i64, _ := strconv.ParseInt(id, 10, 64)
		client.Id = i64
		client.Name = name
		client.Email = email
		client.Phone = phone
	}

	return nil
}

func (client *Client) GetByNameLike(name string, page, limit int64) ([]Client, int64, error) {
	var clientArray []Client
	var total int64

	sql := db.ConnectDatabase()

	name = "%" + name + "%"

	paginate := helpers.Paginate{
		Page:  page,
		Limit: limit,
	}

	paginate.PaginateMounted()
	paginate.MountedQuery("clients")

	query := fmt.Sprintf("select id, name, email, phone, %v from clients where name like ? LIMIT ? OFFSET ?;", paginate.Query)

	requestConfig, err := sql.Query(query, name, paginate.Limit, paginate.Page)

	if err != nil {
		return clientArray, total, err
	}

	for requestConfig.Next() {
		var clientGet Client
		var id, name, email, phone string
		err = requestConfig.Scan(&id, &name, &email, &phone, &total)
		i64, _ := strconv.ParseInt(id, 10, 64)
		clientGet.Id = i64
		clientGet.Name = name
		clientGet.Email = email
		clientGet.Phone = phone

		clientArray = append(clientArray, clientGet)
	}

	return clientArray, total, nil
}

func (client *Client) CreateClient() error {
	sql := db.ConnectDatabase()

	query := "insert into clients (`name`, `email`, `phone`) values (?, ?, ?)"

	createClient, err := sql.Prepare(query)

	if err != nil {
		return err
	}

	_, e := createClient.Exec(client.Name, client.Email, client.Phone)

	if e != nil {
		return e
	}

	return nil
}

func (client *Client) UpdateClient() error {
	sql := db.ConnectDatabase()

	query := `update clients set name = ?, email = ?, phone = ? where id = ?`

	createClient, err := sql.Prepare(query)

	if err != nil {
		return err
	}

	_, e := createClient.Exec(client.Name, client.Email, client.Phone, client.Id)

	if e != nil {
		return e
	}

	return nil
}

func (client *Client) DeleteClientById() error {
	sql := db.ConnectDatabase()

	query := "delete from clients where id = ?"

	createClient, err := sql.Prepare(query)

	if err != nil {
		return err
	}

	_, e := createClient.Exec(client.Id)

	if e != nil {
		return e
	}

	return nil
}

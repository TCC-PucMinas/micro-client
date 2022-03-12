package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"micro-client/helpers"
	"time"

	"micro-client/db"
)

type Destination struct {
	Id       int64  `json:"id"`
	Street   string `json:"street"`
	District string `json:"district"`
	City     string `json:"city"`
	Country  string `json:"country"`
	State    string `json:"state"`
	Number   string `json:"number"`
	Lat      string `json:"lat"`
	Lng      string `json:"lng"`
	ZipCode  string `json:"zipCode"`
	Client   Client `json:"client"`
}

var (
	keyDestinationRedisGetOneByClientId      = "key-destination-get-by-client-id"
	keyDestinationRedisGetPaginateByClientId = "key-destination-get-paginate-client-id"
	keyDestinationRedisGetOneById            = "key-destination-get-by-id"
)

func setRedisCacheDestinationGetByClientId(destination *Destination) error {
	db, err := db.ConnectDatabaseRedis()

	if err != nil {
		return err
	}

	json, err := json.Marshal(destination)

	if err != nil {
		return err
	}
	key := fmt.Sprintf("%v - %v", keyDestinationRedisGetOneByClientId, json)

	return db.Set(key, json, 1*time.Minute).Err()
}

func getRedisCacheDestinationGetByClientId(id int64) (Destination, error) {
	destination := Destination{}

	redis, err := db.ConnectDatabaseRedis()

	if err != nil {
		return destination, err
	}

	key := fmt.Sprintf("%v - %v", keyDestinationRedisGetOneByClientId, id)

	value, err := redis.Get(key).Result()

	if err != nil {
		return destination, err
	}

	if err := json.Unmarshal([]byte(value), &destination); err != nil {
		return destination, err
	}

	return destination, nil
}

func setRedisCacheDestinationGetById(destination *Destination) error {
	db, err := db.ConnectDatabaseRedis()

	if err != nil {
		return err
	}

	json, err := json.Marshal(destination)

	if err != nil {
		return err
	}
	key := fmt.Sprintf("%v - %v", keyDestinationRedisGetOneById, destination.Id)

	return db.Set(key, json, 1*time.Minute).Err()
}

func getRedisCacheDestinationGetById(id int64) (Destination, error) {
	destination := Destination{}

	redis, err := db.ConnectDatabaseRedis()

	if err != nil {
		return destination, err
	}

	key := fmt.Sprintf("%v - %v", keyDestinationRedisGetOneById, id)

	value, err := redis.Get(key).Result()

	if err != nil {
		return destination, err
	}

	if err := json.Unmarshal([]byte(value), &destination); err != nil {
		return destination, err
	}

	return destination, nil
}

func (destination *Destination) DestinationGetByClientId(idClient int64) error {

	if c, err := getRedisCacheDestinationGetByClientId(idClient); err == nil {
		destination.Id = c.Id
		destination.Street = c.Street
		destination.District = c.District
		destination.Lng = c.Lng
		destination.Lat = c.Lat
		destination.District = c.District
		destination.State = c.State
		destination.Country = c.Country
		destination.Number = c.Number
		destination.City = c.City
		destination.ZipCode = c.ZipCode
		destination.Client.Id = c.Client.Id
		return nil
	}

	sql := db.ConnectDatabase()

	query := `select id, street, district, city, country, state, number, lat, lng, zipCode from destinations where id_client = ? limit 1;`

	requestConfig, err := sql.Query(query, idClient)

	if err != nil {
		return err
	}

	for requestConfig.Next() {
		var street, district, city, country, state, number, lat, lng, zipCode string
		var id int64
		_ = requestConfig.Scan(&id, &street, &district, &city, &country, &state, &number, &lat, &lng, &zipCode)
		destination.Id = id
		destination.Street = street
		destination.District = district
		destination.City = city
		destination.Country = country
		destination.State = state
		destination.Number = number
		destination.Lat = lat
		destination.ZipCode = zipCode
		destination.Lng = lng
	}

	if destination.Id == 0 {
		return errors.New("Not found key")
	}

	_ = setRedisCacheDestinationGetByClientId(destination)

	return nil
}

func (destination *Destination) GetDestinationById(idDestination int64) error {

	if c, err := getRedisCacheDestinationGetById(idDestination); err == nil && c.Id != 0 {
		destination.Id = c.Id
		destination.Street = c.Street
		destination.District = c.District
		destination.Lng = c.Lng
		destination.Lat = c.Lat
		destination.District = c.District
		destination.State = c.State
		destination.Country = c.Country
		destination.Number = c.Number
		destination.City = c.City
		destination.ZipCode = c.ZipCode
		destination.Client.Id = c.Client.Id
		return nil
	}

	sql := db.ConnectDatabase()

	query := `select id, street, district, city, country, state, number, lat, lng, zipCode, id_client from destinations where id = ? limit 1;`

	requestConfig, err := sql.Query(query, idDestination)

	if err != nil {
		return err
	}

	for requestConfig.Next() {
		var street, district, city, country, state, number, lat, lng, zipCode string
		var id, idClient int64
		_ = requestConfig.Scan(&id, &street, &district, &city, &country, &state, &number, &lat, &lng, &zipCode, &idClient)
		destination.Id = id
		destination.Street = street
		destination.District = district
		destination.City = city
		destination.Country = country
		destination.State = state
		destination.Number = number
		destination.Lat = lat
		destination.Lng = lng
		destination.ZipCode = zipCode
		destination.Client.Id = idClient
	}

	if destination.Id == 0 {
		return errors.New("Not found key")
	}

	_ = setRedisCacheDestinationGetById(destination)

	return nil
}

func (destination *Destination) CreateDestination() error {
	sql := db.ConnectDatabase()

	query := "insert into destinations (street, district, city, country, `state`, `number`, lat, lng, zipCode, id_client) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	createDestination, err := sql.Prepare(query)

	if err != nil {
		return err
	}

	_, e := createDestination.Exec(destination.Street, destination.District, destination.City, destination.Country, destination.State, destination.Number, destination.Lat, destination.Lng, destination.ZipCode, destination.Client.Id)

	if e != nil {
		return e
	}

	return nil
}

func (destination *Destination) UpdateDestination() error {
	sql := db.ConnectDatabase()

	query := "update destinations set street = ?, district = ?, city = ?, country = ?, `state` = ?, `number` = ?, lat = ?, lng = ?, zipCode = ? id_client = ? where id = ? "

	destinationUpdate, err := sql.Prepare(query)

	if err != nil {
		return err
	}

	_, e := destinationUpdate.Exec(destination.Street, destination.District, destination.City, destination.Country, destination.State, destination.Number, destination.Lat, destination.Lng, destination.ZipCode, destination.Client.Id, destination.Id)

	if e != nil {
		return e
	}

	return nil
}

func (destination *Destination) DeleteDestinationById() error {
	sql := db.ConnectDatabase()

	query := "delete from destinations where id = ?"

	deleteDestination, err := sql.Prepare(query)

	if err != nil {
		return err
	}

	_, e := deleteDestination.Exec(destination.Id)

	if e != nil {
		return e
	}

	return nil
}

func getClientRedisCacheGetOneByClientIdPaginate(clientId int64, page, limit int64) ([]Destination, error) {
	var destinations []Destination

	redis, err := db.ConnectDatabaseRedis()

	if err != nil {
		return destinations, err
	}

	key := fmt.Sprintf("%v - %v -%v -%v", keyDestinationRedisGetPaginateByClientId, clientId, page, limit)

	value, err := redis.Get(key).Result()

	if err != nil {
		return destinations, err
	}

	if err := json.Unmarshal([]byte(value), &destinations); err != nil {
		return destinations, err
	}

	return destinations, nil
}

func setRedisCacheClientGetByClientIdPaginate(clientId int64, page, limit int64, clients []Destination) error {
	redis, err := db.ConnectDatabaseRedis()

	if err != nil {
		return err
	}

	marshal, err := json.Marshal(clients)

	if err != nil {
		return err
	}

	key := fmt.Sprintf("%v - %v -%v -%v", keyDestinationRedisGetPaginateByClientId, clientId, page, limit)

	return redis.Set(key, marshal, 1*time.Minute).Err()
}

func (destination *Destination) GetDestinationsByClientIdPaginate(clientId int64, page, limit int64) ([]Destination, int64, error) {
	var destinationArray []Destination
	var total int64

	if c, err := getClientRedisCacheGetOneByClientIdPaginate(clientId, page, limit); err == nil {
		destinationArray = c
		return destinationArray, total, nil
	}

	sql := db.ConnectDatabase()

	paginate := helpers.Paginate{
		Page:  page,
		Limit: limit,
	}

	paginate.PaginateMounted()

	paginate.MountedQuery("destinations")

	query := fmt.Sprintf("select id, street, district, city, country, state, number, lat, lng, zipCode, id_client, %v from destinations where id_client = ? LIMIT ? OFFSET ?;", paginate.Query)

	requestConfig, err := sql.Query(query, clientId, paginate.Limit, paginate.Page)

	if err != nil {
		return destinationArray, total, err
	}

	for requestConfig.Next() {
		destinationGet := Destination{}
		var street, district, city, country, state, number, lat, lng, zipCode string
		var id, idClient int64
		_ = requestConfig.Scan(&id, &street, &district, &city, &country, &state, &number, &lat, &lng, &zipCode, &idClient, &total)

		if id != 0 {
			destinationGet.Id = id
			destinationGet.Street = street
			destinationGet.District = district
			destinationGet.City = city
			destinationGet.Country = country
			destinationGet.State = state
			destinationGet.Number = number
			destinationGet.Lat = lat
			destinationGet.Lng = lng
			destinationGet.ZipCode = zipCode
			destinationGet.Client.Id = idClient
			destinationArray = append(destinationArray, destinationGet)
		}

	}
	if len(destinationArray) > 0 {
		_ = setRedisCacheClientGetByClientIdPaginate(clientId, page, limit, destinationArray)
	}

	return destinationArray, total, nil
}

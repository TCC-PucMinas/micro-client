package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"micro-logistic/helpers"
	"strconv"
	"time"

	"micro-logistic/db"
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
	Client   Client `json:"client"`
}

var (
	keyDestinationRedisGetOneByClientId      = "key-client-get-by-client-id"
	keyDestinationRedisGetPaginateByClientId = "key-client-get-paginate-client-id"
	keyDestinationRedisGetOneById            = "key-client-get-by-id"
)

func setRedisCacheDestinationGetByClientId(destination *Destination) error {
	db := db.ConnectDatabaseRedis()

	json, err := json.Marshal(destination)

	if err != nil {
		return err
	}
	key := fmt.Sprintf("%v - %v", keyDestinationRedisGetOneByClientId, json)

	return db.Set(key, json, 1*time.Minute).Err()
}

func getRedisCacheDestinationGetByClientId(id int64) (Destination, error) {
	destination := Destination{}

	redis := db.ConnectDatabaseRedis()

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
	db := db.ConnectDatabaseRedis()

	json, err := json.Marshal(destination)

	if err != nil {
		return err
	}
	key := fmt.Sprintf("%v - %v", keyDestinationRedisGetOneById, destination.Id)

	return db.Set(key, json, 1*time.Minute).Err()
}

func getRedisCacheDestinationGetById(id int64) (Destination, error) {
	destination := Destination{}

	redis := db.ConnectDatabaseRedis()

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
		destination = &c
		return nil
	}

	sql := db.ConnectDatabase()

	query := `select id, street, district, city, country, state, number, lat, lng from destinations where id_client = ? limit 1;`

	requestConfig, err := sql.Query(query, idClient)

	if err != nil {
		return err
	}

	for requestConfig.Next() {
		var id, street, district, city, country, state, number, lat, lng string
		_ = requestConfig.Scan(&id, &street, &district, &city, &country, &state, &number, &lat, &lng)
		i64, _ := strconv.ParseInt(id, 10, 64)
		destination.Id = i64
		destination.Street = street
		destination.District = district
		destination.City = city
		destination.Country = country
		destination.State = state
		destination.Number = number
		destination.Lat = lat
		destination.Lng = lng
	}

	if destination.Id == 0 {
		return errors.New("Not found key")
	}

	_ = setRedisCacheDestinationGetByClientId(destination)

	return nil
}

func (destination *Destination) GetDestinationById(idDestination int64) error {

	if c, err := getRedisCacheDestinationGetById(idDestination); err == nil {
		destination = &c
		return nil
	}

	sql := db.ConnectDatabase()

	query := `select id, street, district, city, country, state, number, lat, lng, id_client from destinations where id = ? limit 1;`

	requestConfig, err := sql.Query(query, idDestination)

	if err != nil {
		return err
	}

	for requestConfig.Next() {
		var street, district, city, country, state, number, lat, lng string
		var id, idClient int64
		_ = requestConfig.Scan(&id, &street, &district, &city, &country, &state, &number, &lat, &lng, &idClient)
		destination.Id = id
		destination.Street = street
		destination.District = district
		destination.City = city
		destination.Country = country
		destination.State = state
		destination.Number = number
		destination.Lat = lat
		destination.Lng = lng
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

	query := "insert into destinations (street, district, city, country, `state`, `number`, lat, lng, id_client) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	createDestination, err := sql.Prepare(query)

	if err != nil {
		return err
	}

	_, e := createDestination.Exec(destination.Street, destination.District, destination.City, destination.Country, destination.State, destination.Number, destination.Lat, destination.Lng, destination.Client.Id)

	if e != nil {
		return e
	}

	return nil
}

func (destination *Destination) UpdateDestination() error {
	sql := db.ConnectDatabase()

	query := "update destinations set street = ?, district = ?, city = ?, country = ?, `state` = ?, `number` = ?, lat = ?, lng = ?, id_client = ? where id = ? "

	destinationUpdate, err := sql.Prepare(query)

	if err != nil {
		return err
	}

	_, e := destinationUpdate.Exec(destination.Street, destination.District, destination.City, destination.Country, destination.State, destination.Number, destination.Lat, destination.Lng, destination.Client.Id, destination.Id)

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

	redis := db.ConnectDatabaseRedis()

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
	redis := db.ConnectDatabaseRedis()

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

	query := fmt.Sprintf("select id, street, district, city, country, state, number, lat, lng, id_client, %v from destinations where id_client = ? LIMIT ? OFFSET ?;", paginate.Query)

	requestConfig, err := sql.Query(query, clientId, paginate.Limit, paginate.Page)

	if err != nil {
		return destinationArray, total, err
	}

	for requestConfig.Next() {
		destinationGet := Destination{}
		var street, district, city, country, state, number, lat, lng string
		var id, idClient int64
		_ = requestConfig.Scan(&id, &street, &district, &city, &country, &state, &number, &lat, &lng, &idClient, &total)

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
			destinationGet.Client.Id = idClient
			destinationArray = append(destinationArray, destinationGet)
		}

	}

	_ = setRedisCacheClientGetByClientIdPaginate(clientId, page, limit, destinationArray)

	return destinationArray, total, nil
}

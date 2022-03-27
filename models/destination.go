package model

import (
	"fmt"
	"micro-client/helpers"

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

func (destination *Destination) DestinationGetByClientId(idClient int64) error {

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

	return nil
}

func (destination *Destination) GetDestinationByIdProduct(idProduct int64) error {

	sql := db.ConnectDatabase()

	query := ` select d.id, d.street, d.district, d.city, d.country, d.state, d.number, d.lat, d.lng, d.zipCode, d.id_client from products p
				inner join clients c on c.id = p.id_client
				inner join destinations d on d.id_client = c.id
			where p.id = ? LIMIT 1;`

	requestConfig, err := sql.Query(query, idProduct)

	if err != nil {
		return err
	}

	for requestConfig.Next() {
		var street, district, city, country, state, number, lat, lng, zipCode string
		var id, idClient int64
		_ = requestConfig.Scan(&id, &street, &district, &city, &country, &state, &number, &lat, &lng, &zipCode, &idClient)
		if id != 0 {
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
	}

	return nil
}

func (destination *Destination) GetDestinationById(idDestination int64) error {

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

	query := "update destinations set street = ?, district = ?, city = ?, country = ?, state = ?, number = ?, lat = ?, lng = ?, zipCode = ?, id_client = ? where id = ?"

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

func (destination *Destination) GetDestinationsByClientIdPaginate(clientId int64, page, limit int64) ([]Destination, int64, error) {
	var destinationArray []Destination
	var total int64

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

	return destinationArray, total, nil
}

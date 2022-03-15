package controller

import (
	"context"
	"errors"
	"micro-client/communicate"
	model "micro-client/models"
	"micro-client/service"
)

type DestinationServer struct{}

func (s *DestinationServer) ValidateDestinationById(ctx context.Context, request *communicate.ValidateDestinationByIdRequest) (*communicate.ValidateDestinationByIdResponse, error) {

	res := &communicate.ValidateDestinationByIdResponse{}

	destination := model.Destination{}

	if err := destination.GetDestinationById(request.Id); err != nil {
		return res, errors.New("Destination Id invalid!")
	}

	res.Valid = true

	return res, nil
}

func (s *DestinationServer) DestinationListAll(ctx context.Context, request *communicate.DestinationListAllRequest) (*communicate.DestinationListAllResponse, error) {
	res := &communicate.DestinationListAllResponse{}

	var destination model.Destination

	destinations, total, err := destination.GetDestinationsByClientIdPaginate(request.IdClient, request.Page, request.Limit)

	if err != nil {
		return res, err
	}

	data := &communicate.DataDestination{}
	for _, c := range destinations {
		destination := &communicate.Destination{}
		destination.Id = c.Id
		destination.Street = c.Street
		destination.District = c.District
		destination.City = c.City
		destination.Country = c.Country
		destination.State = c.State
		destination.Number = c.Number
		destination.Lat = c.Lat
		destination.Lng = c.Lng
		destination.ZipCode = c.ZipCode
		destination.IdClient = c.Client.Id
		data.Destination = append(data.Destination, destination)
	}

	res.Total = total
	res.Page = request.Page
	res.Limit = request.Limit

	res.Data = data
	return res, nil
}

func (s *DestinationServer) CreateDestination(ctx context.Context, request *communicate.CreateDestinationRequest) (*communicate.CreateDestinationResponse, error) {
	res := &communicate.CreateDestinationResponse{}

	destination := model.Destination{
		Street:   request.Street,
		District: request.District,
		City:     request.City,
		Country:  request.Country,
		State:    request.State,
		ZipCode:  request.ZipCode,
		Number:   request.Number,
		Client:   model.Client{Id: request.IdClient},
	}

	lat, lng, err := service.GetLocation(destination)

	if err != nil {
		return res, err
	}

	destination.Lat = lat
	destination.Lng = lng

	if err := destination.CreateDestination(); err != nil {
		return res, errors.New("Error creating destination!")
	}

	res.Created = true

	return res, nil
}

func (s *DestinationServer) ListOneDestinationById(ctx context.Context, request *communicate.ListOneDestinationByIdRequest) (*communicate.ListOneDestinationByIdResponse, error) {
	res := &communicate.ListOneDestinationByIdResponse{}

	var c model.Destination

	if err := c.GetDestinationById(request.Id); err != nil {
		return res, err
	}

	destination := &communicate.Destination{}
	destination.Id = c.Id
	destination.Street = c.Street
	destination.District = c.District
	destination.City = c.City
	destination.Country = c.Country
	destination.State = c.State
	destination.Number = c.Number
	destination.Lat = c.Lat
	destination.Lng = c.Lng
	destination.ZipCode = c.ZipCode
	destination.IdClient = c.Client.Id

	res.Destination = destination

	return res, nil
}

func (s *DestinationServer) UpdateDestinationById(ctx context.Context, request *communicate.UpdateDestinationByIdRequest) (*communicate.UpdateDestinationByIdResponse, error) {
	res := &communicate.UpdateDestinationByIdResponse{}

	destination := model.Destination{
		Id:       request.Id,
		City:     request.City,
		Country:  request.Country,
		State:    request.State,
		Street:   request.Street,
		District: request.District,
		ZipCode:  request.ZipCode,
		Number:   request.Number,
		Client:   model.Client{Id: request.IdClient},
	}

	if err := destination.GetDestinationById(request.Id); err != nil || destination.Id == 0 {
		return res, errors.New("Destination not found!")
	}

	lat, lng, err := service.GetLocation(destination)

	if err != nil {
		return res, err
	}

	destination.Lat = lat
	destination.Lng = lng

	if err := destination.UpdateDestination(); err != nil {
		return res, errors.New("Erro updating destination!")
	}

	res.Updated = true

	return res, nil
}

func (s *DestinationServer) DeleteDestinationById(ctx context.Context, request *communicate.DeleteDestinationByIdRequest) (*communicate.DeleteDestinationByIdResponse, error) {
	res := &communicate.DeleteDestinationByIdResponse{}

	destination := model.Destination{}

	if err := destination.GetDestinationById(request.Id); err != nil || destination.Id == 0 {
		return res, errors.New("Destination not found!")
	}

	if err := destination.DeleteDestinationById(); err != nil {
		return res, errors.New("Erro deleting destination!")
	}

	res.Deleted = true

	return res, nil
}

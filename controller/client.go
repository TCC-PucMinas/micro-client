package controller

import (
	"context"
	"errors"
	"micro-client/communicate"
	model "micro-client/models"
)

type ClientServer struct{}

func (s *ClientServer) ClientListAll(ctx context.Context, request *communicate.ClientListAllRequest) (*communicate.ClientListAllResponse, error) {
	res := &communicate.ClientListAllResponse{}

	var client model.Client

	clients, total, err := client.GetByNameLike(request.Name, request.Page, request.Limit)

	if err != nil {
		return res, err
	}

	data := &communicate.DataClient{}
	for _, c := range clients {
		user := &communicate.Client{}
		user.Id = c.Id
		user.Name = c.Name
		user.Email = c.Email
		user.Phone = c.Phone
		data.Client = append(data.Client, user)
	}

	res.Total = total
	res.Page = request.Page
	res.Limit = request.Limit

	res.Data = data
	return res, nil
}

func (s *ClientServer) ListOneClientById(ctx context.Context, request *communicate.ListOneClientByIdRequest) (*communicate.ListOneClientByIdResponse, error) {
	res := &communicate.ListOneClientByIdResponse{}

	client := model.Client{}

	if err := client.GetById(request.Id); err != nil {
		return res, err
	}

	clientGet := &communicate.Client{
		Id:    client.Id,
		Email: client.Email,
		Name:  client.Name,
		Phone: client.Phone,
	}

	res.Client = clientGet

	return res, nil
}

func (s *ClientServer) CreateClient(ctx context.Context, request *communicate.CreateClientRequest) (*communicate.CreateClientResponse, error) {
	res := &communicate.CreateClientResponse{}

	client := model.Client{
		Name:  request.Name,
		Email: request.Email,
		Phone: request.Phone,
	}

	if err := client.GetByNameAndEmail(); err != nil {
		return res, errors.New("Client not duplicated!")
	}

	if client.Id != 0 {
		return res, errors.New("Client not duplicated!")
	}

	if err := client.CreateClient(); err != nil {
		return res, errors.New("Error creating client!")
	}

	res.Created = true

	return res, nil
}

func (s *ClientServer) UpdateClientById(ctx context.Context, request *communicate.UpdateClientByIdRequest) (*communicate.UpdateClientByIdResponse, error) {
	res := &communicate.UpdateClientByIdResponse{}

	client := model.Client{
		Id:    request.Id,
		Name:  request.Name,
		Email: request.Email,
	}

	if err := client.GetById(request.Id); err != nil {
		return res, errors.New("Client not found!")
	}

	if err := client.UpdateClient(); err != nil {
		return res, errors.New("Erro updating client!")
	}

	res.Updated = true

	return res, nil
}

func (s *ClientServer) DeleteClientById(ctx context.Context, request *communicate.DeleteClientByIdRequest) (*communicate.DeleteClientByIdResponse, error) {
	res := &communicate.DeleteClientByIdResponse{}

	client := model.Client{}

	if err := client.GetById(request.Id); err != nil {
		return res, errors.New("Client not found!")
	}

	if client.Id == 0 {
		return res, errors.New("Client not found!")
	}

	if err := client.DeleteClientById(); err != nil {
		return res, errors.New("Erro deleting client!")
	}

	res.Deleted = true

	return res, nil
}

func (s *ClientServer) ValidateClientExist(ctx context.Context, request *communicate.ValidateClientCreateRequest) (*communicate.ValidateClientCreateResponse, error) {
	res := &communicate.ValidateClientCreateResponse{}

	client := model.Client{
		Name:  request.Name,
		Email: request.Email,
	}

	if err := client.GetByNameAndEmail(); err != nil || client.Id != 0 {
		return res, errors.New("Client duplicated!")
	}

	res.Valid = true

	return res, nil
}

func (s *ClientServer) ValidateClientById(ctx context.Context, request *communicate.ValidateClientByIdRequest) (*communicate.ValidateClientByIdResponse, error) {

	res := &communicate.ValidateClientByIdResponse{}

	client := model.Client{}

	if err := client.GetById(request.IdClient); err != nil {
		return res, errors.New("Client Id invalid!")
	}

	res.Valid = true

	return res, nil
}

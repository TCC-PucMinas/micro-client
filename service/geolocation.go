package service

import (
	"context"
	"google.golang.org/grpc"
	"micro-client/communicate"
	model "micro-client/models"
	"time"
)

const attemptRetry = 20

func integrationGeolocation(c *communicate.GelocationRequest, retry int) (*communicate.GelocationResponse, error) {
	ctx := context.Background()
	connGeolocation, err := grpc.Dial(":7000", grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	defer connGeolocation.Close()

	serviceLocation := communicate.NewGelocationCommunicateClient(connGeolocation)

	location, err := serviceLocation.GetLocation(ctx, c)

	if err != nil {
		return attempRetryLatency(retry, location, err, c)
	}

	return location, nil
}

func attempRetryLatency(retry int, location *communicate.GelocationResponse, err error, c *communicate.GelocationRequest) (*communicate.GelocationResponse, error) {
	retry += 1
	if retry <= attemptRetry {
		time.Sleep(1 * time.Second)
		return integrationGeolocation(c, retry)
	}
	return location, err
}

func GetLocation(request model.Destination) (string, string, error) {
	requestLocation := &communicate.GelocationRequest{
		Street:   request.Street,
		District: request.District,
		City:     request.City,
		Country:  request.Country,
		ZipCode:  request.ZipCode,
		State:    request.State,
		Number:   request.Number,
	}

	location, err := integrationGeolocation(requestLocation, 1)

	if err != nil {
		return "", "", err
	}

	return location.Lat, location.Lng, nil
}

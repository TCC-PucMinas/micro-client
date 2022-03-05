package service

import (
	"context"
	"fmt"
	"time"

	"googlemaps.github.io/maps"
)

var keyGoogle = "AIzaSyD5Qni6QsI9nC4GGtwr1kpBss24Zo9KIN8"

// var keyGoogle = "AIzaSyAWl46SC_v1R3MrePMvJqUhhwrX28KWMxU"

type LatAndLng struct {
	Lat string
	Lng string
}

type Calculate struct {
	Origin        LatAndLng
	Destiny       LatAndLng
	HumanReadable string
	Meters        int
	Duration      time.Duration
}

func (calc *Calculate) CalculateRoute() error {

	c, err := maps.NewClient(maps.WithAPIKey(keyGoogle))

	if err != nil {
		return err
	}

	r := &maps.DirectionsRequest{
		Origin:      fmt.Sprintf("%v, %v", calc.Origin.Lat, calc.Origin.Lng),
		Destination: fmt.Sprintf("%v, %v", calc.Destiny.Lat, calc.Destiny.Lng),
	}

	rout, _, err := c.Directions(context.Background(), r)

	if err != nil {
		return err
	}

	for _, v := range rout {
		for _, a := range v.Legs {
			calc.Meters = a.Meters
			calc.HumanReadable = a.HumanReadable
			calc.Duration = a.Duration
		}
	}
	return nil
}

func (l *LatAndLng) GetLatAndLngByAddress(address string) error {

	c, err := maps.NewClient(maps.WithAPIKey(keyGoogle))

	if err != nil {
		return err
	}

	r := &maps.GeocodingRequest{
		Address: address,
	}

	rout, err := c.Geocode(context.Background(), r)

	if err != nil {
		return err
	}

	for _, v := range rout {
		l.Lat = fmt.Sprintf("%f", v.Geometry.Location.Lat)
		l.Lng = fmt.Sprintf("%f", v.Geometry.Location.Lng)
	}

	return nil
}

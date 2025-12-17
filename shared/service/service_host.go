package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceHost interface {
	// Core
	GetContext() context.Context
	GetServices() []Service
	AddService(svc Service) error
	Quit() error

	// DB
	GetMongoClient() *mongo.Client
}

func GetService[T Service](host ServiceHost) (T, error) {
	var ret T

	if host == nil {
		return ret, fmt.Errorf("host is nil")
	}

	for _, svc := range host.GetServices() {
		if casted, ok := svc.(T); ok {
			return casted, nil
		}
	}

	return ret, fmt.Errorf("service not found: %T", ret)
}

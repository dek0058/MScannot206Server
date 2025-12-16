package serverinfo

import (
	"MScannot206/shared/service"
	"errors"
)

func NewServerInfoService(host service.ServiceHost, serverName, dbName string) (*ServerInfoService, error) {
	if host == nil {
		return nil, errors.New("service host is nil")
	}

	return &ServerInfoService{
		host:       host,
		serverName: serverName,
		dbName:     dbName,
	}, nil
}

type ServerInfoService struct {
	host       service.ServiceHost
	serverName string
	dbName     string

	serverInfoRepo *ServerInfoRepository
}

func (s *ServerInfoService) Init() error {
	var err error

	s.serverInfoRepo, err = NewServerInfoRepository(s.host.GetMongoClient(), s.dbName)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServerInfoService) Start() error {
	return nil
}

func (s *ServerInfoService) Stop() error {
	return nil
}

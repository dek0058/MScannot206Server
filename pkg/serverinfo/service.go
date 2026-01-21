package serverinfo

import (
	"MScannot206/shared/service"
	"context"
	"errors"
)

func NewServerInfoService(host service.ServiceHost, serverName, dbName string) (*ServerInfoService, error) {
	var err error

	if host == nil {
		return nil, errors.New("service host is nil")
	}

	s := &ServerInfoService{
		host:       host,
		serverName: serverName,
		dbName:     dbName,
	}

	s.serverInfoRepo, err = NewServerInfoRepository(s.host.GetMongoClient(), s.dbName)
	if err != nil {
		return nil, err
	}

	info, err := s.serverInfoRepo.GetInfo(s.host.GetContext(), s.serverName)
	if err != nil {
		return nil, err
	} else if info == nil {
		info = &ServerInfo{
			Name: s.serverName,

			GameDBName: "MSgame",
			LogDBName:  "MSlog",

			Status: StatusActive,

			Description: "자동 생성 된 서버 정보",
		}

		if err := s.serverInfoRepo.SetInfo(s.host.GetContext(), info); err != nil {
			return nil, err
		}
	}

	return s, nil
}

type ServerInfoService struct {
	host       service.ServiceHost
	serverName string
	dbName     string

	serverInfoRepo *ServerInfoRepository
}

func (s *ServerInfoService) Start(ctx context.Context) error {
	return nil
}

func (s *ServerInfoService) Stop(ctx context.Context) error {
	return nil
}

func (s *ServerInfoService) GetInfo() (*ServerInfo, error) {
	return s.serverInfoRepo.GetInfo(s.host.GetContext(), s.serverName)
}

func (s *ServerInfoService) GetGameDBName() (string, error) {
	info, err := s.GetInfo()
	if err != nil {
		return "", err
	}

	return info.GameDBName, nil
}

func (s *ServerInfoService) GetLogDBName() (string, error) {
	info, err := s.GetInfo()

	if err != nil {
		return "", err
	}

	return info.LogDBName, nil
}

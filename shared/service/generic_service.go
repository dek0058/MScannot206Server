package service

type GenericService interface {
	Init() error
	Start() error
	Stop() error
}

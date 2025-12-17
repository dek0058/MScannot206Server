package service

type Service interface {
	// 서비스 초기화 우선순위 (높을 수록 먼저 초기화 됨)
	GetPriority() int

	Init() error
	Start() error
	Stop() error
}

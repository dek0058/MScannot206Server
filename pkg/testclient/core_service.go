package testclient

import (
	"MScannot206/shared/service"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

func NewCoreService(host service.ServiceHost) (*CoreService, error) {
	if host == nil {
		return nil, errors.New("host is nil")
	}

	s := &CoreService{
		host: host,

		inputChan: make(chan string),
		commands:  make(map[string]ClientCommand, 100),
		doneChan:  make(chan struct{}),
	}

	return s, nil
}

type CoreService struct {
	host service.ServiceHost

	inputChan chan string
	doneChan  chan struct{}

	commands map[string]ClientCommand
}

func (s *CoreService) addCommand(cmd ClientCommand) {
	for _, c := range cmd.Commands() {
		absC := strings.ToLower(c)
		if _, exists := s.commands[absC]; exists {
			log.Warn().Msgf("명령어 중복 등록 시도: %s", c)
			continue
		}
		s.commands[absC] = cmd
	}
}

func (s *CoreService) Init() error {
	s.addCommand(NewLoginCommand(s.host))

	return nil
}

func (s *CoreService) Start() error {
	go s.taskCore()

	return nil
}

func (s *CoreService) Stop() error {
	return nil
}

func (s *CoreService) taskCore() {
	go s.taskInput()

	for {
		select {
		case input := <-s.inputChan:
			s.handleInput(input)

		case <-s.host.GetContext().Done():
			return
		}
	}
}

func (s *CoreService) taskInput() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		tInput := strings.TrimSpace(input)
		if tInput == "" {
			continue
		}

		s.inputChan <- tInput
		<-s.doneChan
	}
}

func (s *CoreService) handleInput(input string) {
	if strings.EqualFold(input, "-help") || strings.EqualFold(input, "-?") || strings.EqualFold(input, "-h") {
		s.printHelp()
		println()
		s.doneChan <- struct{}{}
	} else if strings.EqualFold(input, "-exit") || strings.EqualFold(input, "-quit") || strings.EqualFold(input, "-q") {
		fmt.Println("클라이언트 종료")
		s.host.Quit()
	} else {
		if parts := strings.Fields(input); len(parts) > 0 {
			// parts[0]가 -로 시작하는지 확인
			if !strings.HasPrefix(parts[0], "-") {
				log.Error().Msgf("명령어는 '-'로 시작해야 합니다: %s", parts[0])
				fmt.Println("Usage: -help, -?, -h")
				println()
				s.doneChan <- struct{}{}
				return
			}

			cmdName := strings.ToLower(parts[0][1:])
			if cmd, exists := s.commands[cmdName]; exists {
				if err := cmd.Execute(parts[1:]); err != nil {
					log.Error().Msgf("명령어 실행 오류: %v", err)
					fmt.Println("Usage: -help, -?, -h")
					println()
				}
			} else {
				log.Error().Msgf("알 수 없는 명령어: %s", cmdName)
				fmt.Println("Usage: -help, -?, -h")
				println()
			}
		}
		s.doneChan <- struct{}{}
	}
}

func (s *CoreService) printHelp() {
	fmt.Println("사용 가능한 명령어 목록:")
	println("-exit, -quit, -q: 프로그램 종료")
	for _, cmd := range s.commands {
		fmt.Println()
		fmt.Println(cmd.Description())
	}
}

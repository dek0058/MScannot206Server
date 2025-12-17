package client

import (
	"MScannot206/pkg/testclient/config"
	"MScannot206/pkg/testclient/framework"
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const HTTP_TIMEOUT = 30 * time.Second

func NewClient(ctx context.Context, cfg *config.ClientConfig) (*Client, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}

	webClientCfg := cfg
	if webClientCfg == nil {
		webClientCfg = &config.ClientConfig{
			Url:  "http://localhost",
			Port: 8080,
		}
	}

	client := http.Client{
		Timeout: HTTP_TIMEOUT,
	}

	url := webClientCfg.Url + ":" + fmt.Sprintf("%v", webClientCfg.Port)
	if url == "" {
		return nil, errors.New("웹 클라이언트 URL이 비어있습니다")
	}

	ctxWithCancel, cancel := context.WithCancel(ctx)

	self := &Client{
		ctx:        ctxWithCancel,
		cancelFunc: cancel,

		cfg: webClientCfg,
		url: url,

		client: &client,

		logics: make([]framework.Logic, 0, 8),

		inputChan: make(chan string),
		commands:  make(map[string]framework.ClientCommand, 100),
		doneChan:  make(chan struct{}),
	}

	return self, nil
}

type Client struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	// Config
	cfg *config.ClientConfig
	url string

	client *http.Client

	logics []framework.Logic

	inputChan chan string
	doneChan  chan struct{}

	commands map[string]framework.ClientCommand
}

func (c *Client) GetContext() context.Context {
	return c.ctx
}

func (c *Client) Init() error {
	var errs error

	for _, l := range c.logics {
		if err := l.Init(); err != nil {
			errs = errors.Join(errs, err)
			log.Err(err)
		}
	}

	if errs != nil {
		return errs
	}

	return nil
}

func (c *Client) Start() error {
	for _, l := range c.logics {
		if err := l.Start(); err != nil {
			return err
		}
	}

	go c.taskCore()

	<-c.ctx.Done()
	return nil
}

func (c *Client) Quit() error {
	for _, l := range c.logics {
		if err := l.Stop(); err != nil {
			return err
		}
	}

	if c.cancelFunc != nil {
		c.cancelFunc()
	}

	return nil
}

func (c *Client) AddLogic(logic framework.Logic) error {
	if logic == nil {
		return errors.New("logic is nil")
	}

	c.logics = append(c.logics, logic)
	return nil
}

func (c *Client) GetLogics() []framework.Logic {
	return c.logics
}

func (c *Client) AddCommand(cmd framework.ClientCommand) error {
	if cmd == nil {
		return errors.New("cmd is nil")
	}

	for _, command := range cmd.Commands() {
		absC := strings.ToLower(command)
		if _, exists := c.commands[absC]; exists {
			log.Warn().Msgf("명령어 중복 등록 시도: %s", command)
			continue
		}
		c.commands[absC] = cmd
	}

	return nil
}

func (c *Client) GetUrl() string {
	return c.url
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *Client) taskCore() {
	go c.taskInput()

	for {
		select {
		case input := <-c.inputChan:
			c.handleInput(input)

		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) taskInput() {
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

		c.inputChan <- tInput
		<-c.doneChan
	}
}

func (c *Client) handleInput(input string) {
	if strings.EqualFold(input, "-help") || strings.EqualFold(input, "-?") || strings.EqualFold(input, "-h") {
		c.printHelp()
		println()
		c.doneChan <- struct{}{}
	} else if strings.EqualFold(input, "-exit") || strings.EqualFold(input, "-quit") || strings.EqualFold(input, "-q") {
		fmt.Println("클라이언트 종료")
		c.Quit()
	} else {
		if parts := strings.Fields(input); len(parts) > 0 {
			// parts[0]가 -로 시작하는지 확인
			if !strings.HasPrefix(parts[0], "-") {
				log.Error().Msgf("명령어는 '-'로 시작해야 합니다: %s", parts[0])
				fmt.Println("Usage: -help, -?, -h")
				println()
				c.doneChan <- struct{}{}
				return
			}

			cmdName := strings.ToLower(parts[0][1:])
			if cmd, exists := c.commands[cmdName]; exists {
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
		c.doneChan <- struct{}{}
	}
}

func (c *Client) printHelp() {
	fmt.Println("사용 가능한 명령어 목록:")
	println("-exit, -quit, -q: 프로그램 종료")
	for _, cmd := range c.commands {
		fmt.Println()
		fmt.Println(cmd.Description())
	}
}

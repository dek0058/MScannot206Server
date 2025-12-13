package manager

import (
	"MScannot206/shared/config"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// 기본 설정값
const (
	DefaultMaxFileSize = 10 * 1024 * 1024 // 10MB
	DateFormat         = "2006-01-02"
)

type LogManager struct {
	writer *RollingFileWriter
}

var instance *LogManager
var once sync.Once

func GetLogManager() *LogManager {
	once.Do(func() {
		instance = &LogManager{}
	})
	return instance
}

func (m *LogManager) Init(cfg config.LogConfig) error {
	if cfg.AppName == "" {
		cfg.AppName = "App"
	}
	if cfg.LogDir == "" {
		// 실행 경로 가져오기
		ex, err := os.Executable()
		if err != nil {
			return err
		}
		cfg.LogDir = filepath.Dir(ex)
	} else {
		// 풀경로인지 확인
		if !filepath.IsAbs(cfg.LogDir) {
			// 실행 경로 기준으로 절대경로 변환
			ex, err := os.Executable()
			if err != nil {
				return err
			}
			cfg.LogDir = filepath.Join(filepath.Dir(ex), cfg.LogDir)
		}
	}
	if cfg.MaxFileSize <= 0 {
		cfg.MaxFileSize = DefaultMaxFileSize
	}

	// RollingFileWriter 생성
	rfw := &RollingFileWriter{
		appName:     cfg.AppName,
		logDir:      cfg.LogDir,
		maxFileSize: cfg.MaxFileSize,
	}

	// 초기 파일 설정 (이어쓰기 등을 위해)
	if err := rfw.rotateOrOpen(); err != nil {
		return fmt.Errorf("failed to initialize log file: %w", err)
	}

	m.writer = rfw

	// Zerolog 설정
	var writers []io.Writer
	writers = append(writers, rfw)

	if cfg.DebugMode {
		// 콘솔 출력 시 보기 좋게 포맷팅
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	multi := zerolog.MultiLevelWriter(writers...)

	// 전역 로거 설정
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	return nil
}

func (m *LogManager) Close() error {
	if m.writer != nil {
		return m.writer.Close()
	}
	return nil
}

// RollingFileWriter는 지정된 크기를 초과하면 파일을 회전시키는 io.Writer 구현체입니다.
type RollingFileWriter struct {
	mu sync.Mutex

	appName     string
	logDir      string
	maxFileSize int64

	file         *os.File
	currentSize  int64
	currentDate  string
	currentIndex int
}

func (w *RollingFileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 날짜 변경 체크
	nowDate := time.Now().Format(DateFormat)
	if nowDate != w.currentDate {
		if err := w.rotateOrOpen(); err != nil {
			return 0, err
		}
	}

	// 용량 체크
	writeLen := int64(len(p))
	if w.file == nil || (w.currentSize+writeLen > w.maxFileSize) {
		// 파일이 없거나 꽉 찼으면 로테이션(다음 번호)
		if w.file != nil && w.currentSize+writeLen > w.maxFileSize {
			w.currentIndex++
			if err := w.openNewFile(); err != nil {
				return 0, err
			}
		} else {
			// 파일이 nil인 경우 (rotateOrOpen 실패 후 재시도 등)
			if err := w.rotateOrOpen(); err != nil {
				return 0, err
			}
			// 열었는데도 꽉 찼으면 (rotateOrOpen 내부에서 처리하지만 안전장치)
			if w.currentSize+writeLen > w.maxFileSize {
				w.currentIndex++
				if err := w.openNewFile(); err != nil {
					return 0, err
				}
			}
		}
	}

	n, err = w.file.Write(p)
	w.currentSize += int64(n)
	return n, err
}

// rotateOrOpen은 현재 날짜에 맞는 적절한 로그 파일을 찾아 엽니다.
func (w *RollingFileWriter) rotateOrOpen() error {
	nowDate := time.Now().Format(DateFormat)

	// 날짜가 바뀌었으면 인덱스 리셋
	if w.currentDate != nowDate {
		w.currentDate = nowDate
		w.currentIndex = 0
	}

	// 디렉토리 생성
	if err := os.MkdirAll(w.logDir, 0755); err != nil {
		return err
	}

	// 해당 날짜의 마지막 인덱스 파일 찾기 (프로그램 재시작 시 이어쓰기 위해)
	// 이미 파일이 열려있고 날짜가 같다면 굳이 찾을 필요 없음
	if w.file == nil {
		lastIndex := 0

		files, err := os.ReadDir(w.logDir)
		if err == nil {
			prefix := fmt.Sprintf("%s_%s_", w.appName, w.currentDate)
			for _, f := range files {
				if f.IsDir() {
					continue
				}
				name := f.Name()
				if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, ".log") {
					// 인덱스 파싱
					idxStr := strings.TrimSuffix(strings.TrimPrefix(name, prefix), ".log")
					idx, err := strconv.Atoi(idxStr)
					if err == nil {
						if idx >= lastIndex {
							lastIndex = idx
						}
					}
				}
			}
		}
		w.currentIndex = lastIndex
	}

	// 파일 열기 (존재하면 Append, 없으면 Create)
	filePath := w.getFilePath()
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	w.currentSize = stat.Size()

	// 만약 열었는데 이미 꽉 차있다면? -> 다음 인덱스로
	if w.currentSize >= w.maxFileSize {
		file.Close()
		w.currentIndex++
		return w.openNewFile()
	}

	// 기존 파일 닫기
	if w.file != nil {
		w.file.Close()
	}
	w.file = file

	return nil
}

func (w *RollingFileWriter) openNewFile() error {
	filePath := w.getFilePath()

	// 기존 파일 닫기
	if w.file != nil {
		w.file.Close()
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w.file = file
	w.currentSize = 0
	return nil
}

func (w *RollingFileWriter) getFilePath() string {
	fileName := fmt.Sprintf("%s_%s_%d.log", w.appName, w.currentDate, w.currentIndex)
	return filepath.Join(w.logDir, fileName)
}

// Close 파일 닫기
func (w *RollingFileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

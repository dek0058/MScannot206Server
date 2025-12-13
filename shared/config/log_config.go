package config

type LogConfig struct {
	AppName     string `yaml:"app_name"`      // 애플리케이션 이름
	LogDir      string `yaml:"log_dir"`       // 빈 문자열이면 실행 경로
	MaxFileSize int64  `yaml:"max_file_size"` // 0이면 기본값 사용
	DebugMode   bool   `yaml:"debug_mode"`    // 콘솔 출력 여부 및 레벨 설정
}

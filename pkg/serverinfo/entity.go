package serverinfo

import "time"

type EnvironmentType string
type ServerStatus string

const (
	StatusActive      ServerStatus = "Active"      // 운영 중
	StatusMaintenance ServerStatus = "Maintenance" // 점검 중
	StatusHidden      ServerStatus = "Hidden"      // 일반 유저에게 안 보임
)

type ServerInfo struct {
	Name string `bson:"_id" json:"name"`

	GameDBName string `bson:"game_db_name" json:"game_db_name"`
	LogDBName  string `bson:"log_db_name" json:"log_db_name"`

	Status ServerStatus `bson:"status" json:"status"`

	Description string    `bson:"description" json:"description"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

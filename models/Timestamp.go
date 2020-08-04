package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Timestamp struct {
	gorm.Model
	Category  string
	StartTime time.Time
	EndTime   time.Time
	Comment   string
}

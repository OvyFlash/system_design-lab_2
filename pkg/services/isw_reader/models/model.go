package models

import (
	"time"

	"github.com/lib/pq"
)

type RawISWPage struct {
	Date    time.Time `gorm:"type:DATE;not null;primaryKey"`
	RawPage string    `gorm:"not null"`
}

type ParsedISWPage struct {
	Date       time.Time      `gorm:"type:DATE;not null;primaryKey"`
	ParsedPage pq.StringArray `gorm:"not null;type:text[]"`
}

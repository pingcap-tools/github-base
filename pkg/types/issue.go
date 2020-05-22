package types

import "time"

// Issue struct
type Issue struct {
	ID          int       `gorm:"column:id"`
	Owner       string    `gorm:"column:owner"`
	Repo        string    `gorm:"column:repo"`
	Number      int       `gorm:"column:issue_number"`
	Title       string    `gorm:"column:title"`
	Body        string    `gorm:"column:body"`
	User        string    `gorm:"column:user"`
	Association string    `gorm:"column:association"`
	Relation    string    `gorm:"column:relation"`
	Label       string    `gorm:"column:label"`
	Status      string    `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	ClosedAt    time.Time `gorm:"column:closed_at"`
}

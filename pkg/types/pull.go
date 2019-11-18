package types

import (
	"github.com/pingcap/github-base/util"
	"time"
)

// Pull struct
type Pull struct {
	ID          int `gorm:"column:id"`
	Owner       string `gorm:"column:owner"`
	Repo        string `gorm:"column:repo"`
	Number      int `gorm:"column:pull_number"`
	Title       string `gorm:"column:title"`
	Body        string `gorm:"column:body"`
	User        string `gorm:"column:user"`
	Association string `gorm:"column:association"`
	Relation    string `gorm:"column:relation"`
	Label       string `gorm:"column:label"`
	Status      string `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	ClosedAt    time.Time `gorm:"column:closed_at"`
	MergedAt    time.Time `gorm:"column:merged_at"`
}

// GetID get pull's ID
func (p *Pull)GetID() int {
	return p.ID
}

// GetOwner get pull's owner
func (p *Pull)GetOwner() string {
	return p.Owner
}

// GetRepo get pull's repo
func (p *Pull)GetRepo() string {
	return p.Repo
}

// GetNumber get pull's number
func (p *Pull)GetNumber() int {
	return p.Number
}

// GetTitle get pull's title
func (p *Pull)GetTitle() string {
	return p.Title
}

// GetBody get pull's body
func (p *Pull)GetBody() string {
	return p.Body
}

// GetUser get pull's user
func (p *Pull)GetUser() string {
	return p.User
}

// GetAssiciation get pull's association
func (p *Pull)GetAssociation() string {
	return p.Association
}

// GetRelation get pull's relation
func (p *Pull)GetRelation() string {
	return p.Relation
}

// GetLabel get pull's label slice
func (p *Pull)GetLabel() []string {
	return util.ParseStringSlice(p.Label)
}

// GetStatus get pull's status
func (p *Pull)GetStatus() string {
	return p.Status
}

// GetCreatedAt get pull's created time
func (p *Pull)GetCreatedAt() time.Time {
	return p.CreatedAt
}

// GetUpdatedAt get pull's updated time
func (p *Pull)GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

// GetClosedAt get pull's closed time
func (p *Pull)GetClosedAt() time.Time {
	return p.ClosedAt
}

// GetMergedAt get pull's merged time
func (p *Pull)GetMergedAt() time.Time {
	return p.MergedAt
}

// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
	"github.com/bryant-rh/srew/pkg/sqlx"
)

const TableNameDetail = "detail"

// Detail mapped from table <detail>
type Detail struct {
	ID               int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	PluginID         string    `gorm:"column:plugin_id;not null;default:0" json:"plugin_id"`
	PluginName       string    `gorm:"column:plugin_name;not null" json:"plugin_name"`
	Version          string    `gorm:"column:version;not null" json:"version"`
	Homepage         string    `gorm:"column:homepage;not null" json:"homepage"`
	ShortDescription string    `gorm:"column:shortDescription;not null" json:"shortDescription"`
	Description      string    `gorm:"column:description;not null" json:"description"`
	Caveats          string    `gorm:"column:caveats;not null" json:"caveats"`
	Platforms        sqlx.StringSlice    `gorm:"column:platforms;not null" json:"platforms"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName Detail's table name
func (*Detail) TableName() string {
	return TableNameDetail
}
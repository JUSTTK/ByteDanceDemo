// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameLike = "likes"

// Like mapped from table <likes>
type Like struct {
	ID        int64          `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement:true;comment:主键" json:"id"` // 主键
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime(3);comment:记录创建时间" json:"created_at"`                   // 记录创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime(3);comment:记录更新时间" json:"updated_at"`                   // 记录更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3);comment:软删除时间" json:"deleted_at"`                    // 软删除时间
	UserID    int64          `gorm:"column:user_id;type:bigint(20) unsigned;not null;comment:点赞用户id" json:"user_id"`        // 点赞用户id
	VideoID   int64          `gorm:"column:video_id;type:bigint(20) unsigned;not null;comment:点赞视频id" json:"video_id"`      // 点赞视频id
	Liked     int64          `gorm:"column:liked;type:bigint(20);not null;default:1;comment:默认1表示已点赞，0表示未点赞" json:"liked"`  // 默认1表示已点赞，0表示未点赞
}

// TableName Like's table name
func (*Like) TableName() string {
	return TableNameLike
}

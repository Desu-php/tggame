package models

import "time"

type ReferralUser struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`                                                           // Primary key with auto increment
	UserID         uint      `json:"user_id" gorm:"not null"`                                                                      // Foreign key to User
	ReferredUserID uint      `json:"referred_user_id" gorm:"not null"`                                                             // Foreign key to referred User
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`                                                             // Automatically managed creation timestamp
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`                                                             // Automatically managed update timestamp
	ReferredUser   *User     `json:"referred_user" gorm:"foreignKey:ReferredUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // Adds foreign key constraint
}

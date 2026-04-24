package model

type User struct {
	ID                 uint   `gorm:"primaryKey" json:"id"`
	Username           string `gorm:"unique;not null" json:"username"`
	Password           string `gorm:"not null" json:"-"`
	Role               string `gorm:"type:text" json:"role"`
	RevokeTokensBefore int64  `gorm:"type:integer" json:"revoke_tokens_before"`
}

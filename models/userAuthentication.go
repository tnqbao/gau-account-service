package models

type UserAuthentication struct {
	UserId     uint    `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Permission string  `json:"permission"`
	Username   *string `gorm:"unique" json:"username"`
	Password   *string `gorm:"size:255" json:"password"`
}

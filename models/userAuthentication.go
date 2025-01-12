package models

type UserAuthentication struct {
	UserId     uint    `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Permission string  `gorm:"index:idx_username_permission" json:"permission"`
	Username   *string `gorm:"unique;index:idx_username_permission" json:"username"`
	Password   *string `gorm:"size:255" json:"password"`
}

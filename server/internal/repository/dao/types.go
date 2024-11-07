package dao

type model struct {
	Id uint `gorm:"primaryKey,autoIncrement"`
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}

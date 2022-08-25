package permission

type Permission struct {
	Name string `gorm:"column:name;primaryKey"`
}

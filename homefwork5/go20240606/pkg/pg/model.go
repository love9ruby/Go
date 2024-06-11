package pg

type Event struct {
	Password string `json:"password" gorm:"foreignKey:users(password)"`
	Time     string `json:"time"`
	Url      string `json:"url"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password" gorm:"primaryKey"`
}

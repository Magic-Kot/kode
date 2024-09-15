package models

type User struct {
	ID       int    `json:"id"           validate:"required,min=1"`
	Age      int    `json:"age"          validate:"lte=120"` //gte=14,
	Username string `json:"login"        validate:"max=20"`  //min=4,
	Name     string `json:"name"         validate:"min=1,max=20"`
	Surname  string `json:"surname"      validate:"min=1,max=20"`
	Email    string `json:"email"        validate:"email"`
	Avatar   string `json:"avatar"`
}

type UserLogin struct {
	ID       int    `json:"id"`
	Username string `json:"login"        validate:"required,min=4,max=20"`
	Password string `json:"password"     validate:"required,min=6,max=20"`
}

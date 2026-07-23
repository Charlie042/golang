package models 

type User struct {
	ID string `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Created_at string `json:"created_at" db:"created_at"`
	Updated_at string `json:"updated_at" db:"updated_at"`

}
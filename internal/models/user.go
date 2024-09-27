package models

type User struct {
    ID    uint   `json:"id" gorm:"primaryKey"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Password  string `json:"password"`
}

type UserIDs struct {
    IDS   string `json:"ids"`
}

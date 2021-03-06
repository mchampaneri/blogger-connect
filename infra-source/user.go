package main

import (
	"encoding/json"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (*User) HashPassword(password string) string {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashBytes)
}

func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		Dob string `json:"dob"`
		*Alias
	}{
		Dob:   u.Dob.Format("02-01-2006"),
		Alias: (*Alias)(u),
	})
}

type User struct {
	Id        int       `json:"id" xorm:"int(12) not null unique pk 'id'"   `
	Role      int       `json:"role" xorm:"int(12) not null  'role'"`
	Name      string    `json:"name"  xorm:"varchar(25) not null  'name' " valid:"required"`
	NickName  string    `json:"nickname" xorm:"varchar(25) 'nickname' "`
	Gender    int       `json:"gender" xorm:"int(12) 'gender' " `
	Dob       time.Time `json:"dobinput" xorm:"timestamp 'dob' "`
	Domain    string    `json:"domain"  xorm:"varchar(25) 'domain' "`
	DobString string    `json:"dob" xorm:"-"`
	Email     string    `json:"email" xorm:"varchar(25) not null unique 'email' "  valid:"required,email"`
	Quote     string    `json:"quote"  xorm:"varchar(25)  'quote' "`
	Password  string    `json:"password" xorm:"varchar(25) not null 'password' "  valid:"required"`
	Status    int       `json:"id" xorm:"int(2) not null default false 'status' " `

	ProfilePic string `json:"profile_pic" xorm:"varchar(100) 'profile_pic' "`
	CoverPic   string `json:"cover_pic" xorm:"varchar(100) 'cover_pic' "`

	Nationality string `json:"nationality" xorm:"varchar(60) 'nationality' "`
	Language    string `json:"language" xorm:"varchar(60) 'language' "`

	Interest      string `json:"interest" xorm:"varchar(60) 'interest' "`
	FacebookToken string `json:"facebook_token" xorm:"varchar(200) 'facebook_token'"`
	GoogleToken   string `json:"google_token" xorm:"varchar(200) 'google_token'"`
	Profession    string `json:"profession" xorm:"varchar(60) 'profession' "`

	Slug  string `json:"slug" xorm:"varchar(100) 'slug' "`
	Phone string `json:"phone" xorm:"varchar(15)  'phone' "`

	CreatedAt time.Time `json:"created_at" xorm:"timestamp 'created_at'" `
	UpdatedAt time.Time `json:"updated_at" xorm:"timestamp 'updated_at'" `
	DeletedAt time.Time `json:"deleted_at" xorm:"timestamp 'deleted_at'" `
}

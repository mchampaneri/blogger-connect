package main

import (
	"time"
)

type BloggerUserSession struct {
	AccessToken  string
	RefreshToken string
	Id           string
	AssigendAt   time.Time
	ExpeireIn    float64
}

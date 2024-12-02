// Package domain contains the business logic for user/auth entities
package user

type User struct {
	ID   int
	Guid string `example:"guid"`
	Hash string `example:"$2a$10$OD6gRRUd0O8cTxdGQjPzqOuk1cmMkX/FON.1jkfVpz0I.AuQXqvMa"`
	IP   string
}

type NewUser struct {
	Guid string `example:"guid"`
	Hash string `example:"$2a$10$OD6gRRUd0O8cTxdGQjPzqOuk1cmMkX/FON.1jkfVpz0I.AuQXqvMa"`
	IP   string
}

package epinio

import "time"

type User struct {
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	Namespaces        []string
	Role              string
	Roles             []string
	CreationTimestamp time.Time

	secret string
}

func NewUser(secretName string) User {
	return User{
		secret: secretName,
	}
}

func (u User) SecretName() string {
	return u.secret
}

func (u User) GetID() string {
	return u.Username
}

func (u User) IsAdmin() bool {
	if u.Role == "admin" {
		return true
	}

	for _, r := range u.Roles {
		if r == "admin" {
			return true
		}
	}

	return false
}

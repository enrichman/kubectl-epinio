package epinio

type User struct {
	Username   string `yaml:"username"`
	Password   string
	Namespaces []string
	Role       string
	Roles      []string

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

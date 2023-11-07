package epinio

type User struct {
	Username string `yaml:"username"`
	Password string
	//TODO: add namespaces and roles
}

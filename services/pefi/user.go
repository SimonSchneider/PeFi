package pefi

type (
	UserService interface {
		Create(name string) (*User, error)
		Update(name string, new interface{}) error
		Delete(name string) error
		Get(name string) (*User, error)
	}

	User struct {
		Name string `json:name`
	}
)

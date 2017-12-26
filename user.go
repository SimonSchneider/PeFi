package pefi

type (
	UserService interface {
		Create(name ID) (*User, error)
		Update(name ID, new interface{}) error
		Delete(name ID) error
		Get(name ID) (*User, error)
	}

	User struct {
		ID   ID     `json:id`
		Name string `json:name`
	}
)

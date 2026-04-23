package user

type UserRepo interface {
	Create(User) (User, error)
	// Login(email string, pw string) (User, error)
	EmailExists(email string) (bool, error)
	GetByEmail(email string) (User, error)
}

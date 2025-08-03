package application

import "github.com/florian-renfer/b0red/internal/domain"

type UserRegistry struct {
	users []*domain.User
}

func (ur *UserRegistry) AddUser(user *domain.User) {
	if ur.users == nil {
		ur.users = make([]*domain.User, 0)
	}
	ur.users = append(ur.users, user)
}

func (ur *UserRegistry) RemoveUser(user *domain.User) {
	for i, u := range ur.users {
		if u.Name == user.Name {
			ur.users = append(ur.users[:i], ur.users[i+1:]...)
			return
		}
	}
}

func (ur *UserRegistry) GetUsers() []*domain.User {
	if ur.users == nil {
		return []*domain.User{}
	}
	return ur.users
}

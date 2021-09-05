package app

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

type Status uint8

const (
	_ = iota
	Active
	InActive
	Banned
)

const limit = 5

var (
	ErrNotFound              = errors.New("not found")
	ErrContactAlreadyExists  = errors.New("contact already exists")
	ErrContactsLimitExceeded = errors.New("contacts list is full")
)

type ErrContactNotFound struct {
	ID int
}

func NewErrContactNotFound(ID int) *ErrContactNotFound {
	return &ErrContactNotFound{ID: ID}
}

func (e *ErrContactNotFound) Error() string {
	return fmt.Sprintf("contact with ID: %v not found", e.ID)
}

type Contact struct {
	UserID int
}

type Contacts []*Contact

func NewContact(userID int) *Contact {
	return &Contact{UserID: userID}
}

type User struct {
	ID        int       `firestore:"id"`
	Contacts  Contacts  `firestore:"contacts"`
	Status    Status    `firestore:"status"`
	CreatedAt time.Time `firestore:"created_at"`
}

func NewUser(ID int) *User {
	return &User{ID: ID, Status: Active, CreatedAt: time.Now()}
}

func (u *User) AddContact(c *Contact) error {
	if len(u.Contacts) >= limit {
		return ErrContactsLimitExceeded
	}

	if u.Contacts.contains(c.UserID) {
		return ErrContactAlreadyExists
	}

	u.Contacts = append(u.Contacts, c)

	return nil
}

func (u *User) RemoveContact(c *Contact) error {
	u.Contacts.remove(c.UserID)

	return nil
}

func (u *User) GetContacts() Contacts {
	return u.Contacts
}

func (u *User) Ban() error {
	u.Status = Banned

	return nil
}

type UserStorage interface {
	Add(ctx context.Context, u *User) error
	Get(ctx context.Context, ID int) (*User, error)
	Exists(ctx context.Context, ID int) (bool, error)
	Remove(ctx context.Context, ID int) error
}

func (c Contacts) contains(id int) bool {
	for _, contact := range c {
		if contact.UserID == id {
			return true
		}
	}

	return false
}

func (c Contacts) remove(id int) bool {
	for _, contact := range c {
		if contact.UserID == id {
			return true
		}
	}

	return false
}

package storage

import (
	"cloud.google.com/go/firestore"
	"context"
	app "github.com/VMAnalytic/alarm-bot/internal"
	"github.com/VMAnalytic/alarm-bot/pkg/convertor"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const collectionUsers = "users"

type UserFirestoreStorage struct {
	client *firestore.Client
}

func NewUserFirestoreStorage(client *firestore.Client) *UserFirestoreStorage {
	return &UserFirestoreStorage{client: client}
}

func (s UserFirestoreStorage) Add(ctx context.Context, u *app.User) error {
	_, err := s.client.Collection(collectionUsers).Doc(convertor.ToString(u.ID)).Set(ctx, u)

	return errors.WithStack(err)
}

func (s UserFirestoreStorage) Get(ctx context.Context, ID int) (*app.User, error) {
	doc, err := s.client.Collection(collectionUsers).Doc(convertor.ToString(ID)).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, app.ErrNotFound
		}

		return nil, errors.WithStack(err)
	}

	if !doc.Exists() {
		return nil, app.ErrNotFound
	}

	var user app.User
	err = doc.DataTo(&user)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &user, nil
}

func (s UserFirestoreStorage) Exists(ctx context.Context, ID int) (bool, error) {
	doc, err := s.client.Collection(collectionUsers).Doc(convertor.ToString(ID)).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}

		return false, errors.WithStack(err)
	}

	if !doc.Exists() {
		return false, nil
	}

	return true, nil
}

func (s UserFirestoreStorage) Remove(ctx context.Context, ID int) error {
	_, err := s.client.Collection(collectionUsers).Doc(convertor.ToString(ID)).Delete(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	//query := s.client.Collection(collectionUsers).Where("contacts", "array-contains", ID).Documents(ctx)

	return nil
}

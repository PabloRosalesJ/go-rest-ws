package repository

import (
	"context"

	"github.com/PabloRosalesJ/go/rest-ws/models"
	_ "github.com/lib/pq"
)

var implementation Repository

type Repository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	/* ============ POST ============ */

	InsertPost(ctx context.Context, post *models.Post) error

	Close() error
}

func SetRepository(repo Repository) {
	implementation = repo
}

func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}

func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return implementation.GetUserById(ctx, id)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, email)
}

/* ============ POST ============ */

func InsertPost(ctx context.Context, post *models.Post) error {
	return implementation.InsertPost(ctx, post)
}

/* ============ Global ============ */

func Close() error {
	return implementation.Close()
}
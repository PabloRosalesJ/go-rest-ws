package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/PabloRosalesJ/go/rest-ws/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES($1, $2, $3)", user.Id, user.Email, user.Password)
	return err
}

func (repo *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user = models.User{}

	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "Select id, email, password from users where email = $1", email)
	if err != nil {
		log.Fatal("GetUserByEmail.rows", err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal("GetUserByEmail.close", err)
		}
	}()

	var user = models.User{}

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email, &user.Password); err != nil {
			log.Fatal("GetUserByEmail.next", err)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

/* ============ POST ============ */

func (repo *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO posts (id, user_id, post_content) values ($1, $2, $3)", post.Id, post.UserId, post.PostContent)
	return err
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}

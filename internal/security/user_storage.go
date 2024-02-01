package security

import (
	"database/sql"
	"log"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
)

type UserStorage struct {
	database *sql.DB
}

func NewUserStorage(database *sql.DB) *UserStorage {
	return &UserStorage{database}
}

func (this *UserStorage) FindUserByName(name string) (entity.User, error) {
	var user entity.User
	query := `SELECT * FROM "Users" WHERE "Name" = $1`
	row := this.database.QueryRow(query, name)
	if err := row.Scan(&user.Id, &user.Name, &user.Key); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err.Error())
		}

		return user, err
	}

	return user, nil
}

func (this *UserStorage) FindUserById(id int) (entity.User, error) {
	var user entity.User
	query := `SELECT * FROM "Users" WHERE "Id" = $1`
	row := this.database.QueryRow(query, id)
	if err := row.Scan(&user.Id, &user.Name, &user.Key); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err.Error())
		}

		return user, err
	}

	return user, nil
}

func (this *UserStorage) InsertUserIfNotExists(user entity.User) (int, error) {
	transaction, err := this.database.Begin()
	if err != nil {
		return 0, err
	}

	defer transaction.Rollback()

	exists := false
	findQuery := `SELECT EXISTS(SELECT 1 FROM "Users" WHERE "Name" = $1 )`
	row := this.database.QueryRow(findQuery, user.Name)
	err = row.Scan(&exists)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err.Error())
		}

		exists = false
	}

	if exists {
		return 0, ErrUserAlreadyExists
	}

	userQuery := `INSERT INTO "Users" ("Name", "Password") VALUES ($1, $2) RETURNING "Id"`
	id := 0
	row = transaction.QueryRow(userQuery, &user.Name, &user.Key)
	if err := row.Scan(&id); err != nil {
		log.Println(err.Error())
		return 0, err
	}

	err = transaction.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (this *UserStorage) UpdateUserKey(user entity.User) error {
	query := `UPDATE "Users" SET "Password" = $2 WHERE "Id" = $1`
	if _, err := this.database.Exec(query, &user.Id, &user.Key); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

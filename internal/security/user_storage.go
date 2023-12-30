package security

import (
	"database/sql"
	"log"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
)

type userStorage struct {
	database *sql.DB
}

func newUserStorage(database *sql.DB) *userStorage {
	return &userStorage{database}
}

func (this *userStorage) userExists(name string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM "Users" WHERE "Name" = $1 )`
	row := this.database.QueryRow(query, name)
	err := row.Scan(&exists)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err.Error())
		}

		return false
	}

	return exists
}

func (this *userStorage) findUserByName(name string) (entity.User, error) {
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

func (this *userStorage) findUserById(id int) (entity.User, error) {
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

func (this *userStorage) insertUser(user entity.User) (int, error) {
	var id int
	query := `INSERT INTO "Users" ("Name", "Password") VALUES ($1, $2) RETURNING "Id"`
	row := this.database.QueryRow(query, &user.Name, &user.Key)
	if err := row.Scan(&id); err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return id, nil
}

func (this *userStorage) updateUserKey(user entity.User) error {
	query := `UPDATE "Users" SET "Password" = $2 WHERE "Id" = $1`
	if _, err := this.database.Exec(query, &user.Id, &user.Key); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

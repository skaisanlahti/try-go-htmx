package psql

import (
	"database/sql"
	"log"

	"github.com/skaisanlahti/try-go-htmx/todoapp"
)

type UserStorage struct {
	Database *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db}
}

func (storage *UserStorage) FindUsers() []todoapp.User {
	var users []todoapp.User
	query := `SELECT * FROM "Users"`
	rows, err := storage.Database.Query(query)
	if err != nil {
		log.Println(err.Error())
		return users
	}

	defer rows.Close()
	for rows.Next() {
		var user todoapp.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Key); err != nil {
			log.Println(err.Error())
			return users
		}

		users = append(users, user)
	}

	return users
}

func (storage *UserStorage) FindUserByName(name string) (todoapp.User, error) {
	var user todoapp.User
	query := `SELECT * FROM "Users" WHERE "Name" = $1`
	row := storage.Database.QueryRow(query, name)
	if err := row.Scan(&user.Id, &user.Name, &user.Key); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err.Error())
		}

		return user, err
	}

	return user, nil
}

func (storage *UserStorage) AddUser(user todoapp.User) (int, error) {
	var id int
	query := `INSERT INTO "Users" ("Name", "Password") VALUES ($1, $2) RETURNING "Id"`
	row := storage.Database.QueryRow(query, &user.Name, &user.Key)
	if err := row.Scan(&id); err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return id, nil
}

func (storage *UserStorage) UpdateUserKey(user todoapp.User) error {
	query := `UPDATE "Users" SET "Password" = $2 WHERE "Id" = $1`
	if _, err := storage.Database.Exec(query, &user.Id, &user.Key); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (storage *UserStorage) RemoveUser(id int) error {
	query := `DELETE FROM "Users" WHERE "Id" = $1`
	if _, err := storage.Database.Exec(query, id); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

package auth

import (
	"database/sql"
	"log"
)

type userStorage struct {
	database *sql.DB
}

func newUserStorage(db *sql.DB) *userStorage {
	return &userStorage{db}
}

func (storage *userStorage) findUsers() []user {
	var users []user
	query := `SELECT * FROM "Users"`
	rows, err := storage.database.Query(query)
	if err != nil {
		log.Println(err.Error())
		return users
	}

	defer rows.Close()
	for rows.Next() {
		var user user
		if err := rows.Scan(&user.Id, &user.Name, &user.Key); err != nil {
			log.Println(err.Error())
			return users
		}

		users = append(users, user)
	}

	return users
}

func (storage *userStorage) userExists(name string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM "Users" WHERE "Name" = $1 )`
	row := storage.database.QueryRow(query, name)
	err := row.Scan(&exists)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err.Error())
		}

		return false
	}

	return exists
}

func (storage *userStorage) findUserByName(name string) (user, error) {
	var user user
	query := `SELECT * FROM "Users" WHERE "Name" = $1`
	row := storage.database.QueryRow(query, name)
	if err := row.Scan(&user.Id, &user.Name, &user.Key); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err.Error())
		}

		return user, err
	}

	return user, nil
}

func (storage *userStorage) addUser(user user) (int, error) {
	var id int
	query := `INSERT INTO "Users" ("Name", "Password") VALUES ($1, $2) RETURNING "Id"`
	row := storage.database.QueryRow(query, &user.Name, &user.Key)
	if err := row.Scan(&id); err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return id, nil
}

func (storage *userStorage) updateUserKey(user user) error {
	query := `UPDATE "Users" SET "Password" = $2 WHERE "Id" = $1`
	if _, err := storage.database.Exec(query, &user.Id, &user.Key); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (storage *userStorage) removeUser(id int) error {
	query := `DELETE FROM "Users" WHERE "Id" = $1`
	if _, err := storage.database.Exec(query, id); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

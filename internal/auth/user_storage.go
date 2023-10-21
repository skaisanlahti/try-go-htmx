package auth

import (
	"database/sql"
	"log"
)

type userStorage struct {
	Database *sql.DB
}

func NewUserStorage(db *sql.DB) *userStorage {
	return &userStorage{db}
}

func (storage *userStorage) findUsers() []user {
	var users []user
	query := `SELECT * FROM "Users"`
	rows, err := storage.Database.Query(query)
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

func (storage *userStorage) findUserByName(name string) (user, error) {
	var user user
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

func (storage *userStorage) addUser(user user) (int, error) {
	var id int
	query := `INSERT INTO "Users" ("Name", "Password") VALUES ($1, $2) RETURNING "Id"`
	row := storage.Database.QueryRow(query, &user.Name, &user.Key)
	if err := row.Scan(&id); err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return id, nil
}

func (storage *userStorage) updateUserKey(user user) error {
	query := `UPDATE "Users" SET "Password" = $2 WHERE "Id" = $1`
	if _, err := storage.Database.Exec(query, &user.Id, &user.Key); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (storage *userStorage) removeUser(id int) error {
	query := `DELETE FROM "Users" WHERE "Id" = $1`
	if _, err := storage.Database.Exec(query, id); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

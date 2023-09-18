package repositories

import (
	"database/sql"
	"log"

	"github.com/skaisanlahti/try-go-htmx/users/domain"
)

type UserRepository struct {
	Database *sql.DB
}

func NewUserRepository(database *sql.DB) *UserRepository {
	return &UserRepository{database}
}

const selectUsers string = `SELECT * FROM "Users"`

func (repository *UserRepository) GetUsers() []domain.User {
	var users []domain.User
	rows, err := repository.Database.Query(selectUsers)
	if err != nil {
		log.Println(err.Error())
		return users
	}

	defer rows.Close()
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Password); err != nil {
			log.Println(err.Error())
			return users
		}

		users = append(users, user)
	}

	return users
}

const selectUserByName string = `SELECT * FROM "Users" WHERE "Name" = $1`

func (repository *UserRepository) GetUserByName(name string) (domain.User, error) {
	var user domain.User
	row := repository.Database.QueryRow(selectUserByName, name)
	if err := row.Scan(&user.Id, &user.Name, &user.Password); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err.Error())
		}

		return user, err
	}

	return user, nil
}

const insertUser string = `INSERT INTO "Users" ("Name", "Password") VALUES ($1, $2)`

func (repository *UserRepository) AddUser(user domain.User) error {
	if _, err := repository.Database.Exec(insertUser, &user.Name, &user.Password); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

const updateUserPassword string = `Update "Users" Set "Password" = $2 WHERE "Id" = $1`

func (repository *UserRepository) UpdateUserPassword(user domain.User) error {
	if _, err := repository.Database.Exec(updateUserPassword, &user.Id, &user.Password); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

const deleteUser string = `DELETE FROM "Users" WHERE "Id" = $1`

func (repository *UserRepository) RemoveUser(id int) error {
	if _, err := repository.Database.Exec(deleteUser, id); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

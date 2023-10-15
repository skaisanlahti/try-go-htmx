package psql

import (
	"database/sql"
	"log"

	"github.com/skaisanlahti/try-go-htmx/todoapp"
)

type UserAccessor struct {
	Database *sql.DB
}

func NewUserAccessor(db *sql.DB) *UserAccessor {
	return &UserAccessor{db}
}

func (accessor *UserAccessor) FindUsers() []todoapp.User {
	var users []todoapp.User
	query := `SELECT * FROM "Users"`
	rows, err := accessor.Database.Query(query)
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

func (accessor *UserAccessor) FindUserByName(name string) (todoapp.User, error) {
	var user todoapp.User
	query := `SELECT * FROM "Users" WHERE "Name" = $1`
	row := accessor.Database.QueryRow(query, name)
	if err := row.Scan(&user.Id, &user.Name, &user.Key); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err.Error())
		}

		return user, err
	}

	return user, nil
}

func (accessor *UserAccessor) AddUser(user todoapp.User) (int, error) {
	var id int
	query := `INSERT INTO "Users" ("Name", "Password") VALUES ($1, $2) RETURNING "Id"`
	row := accessor.Database.QueryRow(query, &user.Name, &user.Key)
	if err := row.Scan(&id); err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return id, nil
}

func (accessor *UserAccessor) UpdateUserKey(user todoapp.User) error {
	query := `UPDATE "Users" SET "Password" = $2 WHERE "Id" = $1`
	if _, err := accessor.Database.Exec(query, &user.Id, &user.Key); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (accessor *UserAccessor) RemoveUser(id int) error {
	query := `DELETE FROM "Users" WHERE "Id" = $1`
	if _, err := accessor.Database.Exec(query, id); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

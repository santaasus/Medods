package db

import (
	"database/sql"
	"fmt"

	domain "Medods/auth_service/inner_layer/domain"
	dbCore "Medods/db_service"
)

const (
	SQLSELECT = "SELECT * FROM userss"
	SQLUPDATE = "UPDATE userss SET "
	SQLDELETE = "DELETE FROM userss"
	SQLINSERT = "INSERT INTO userss"
)

func CreateUser(newUser *domain.NewUser) (*domain.User, error) {
	db, err := dbCore.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	var userId int
	query := SQLINSERT + ` (user_guid, refresh_hash, ip)
			  VALUES ($1, $2, $3) RETURNING id;`

	err = db.QueryRow(
		query,
		newUser.Guid, newUser.Hash, newUser.IP,
	).Scan(&userId)

	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:   userId,
		Guid: newUser.Guid,
		Hash: newUser.Hash,
		IP:   newUser.IP}, nil
}

func GetUserByGuid(guid string) (*domain.User, error) {
	db, err := dbCore.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	query := SQLSELECT + " WHERE user_guid = $1;"

	row := db.QueryRow(query, guid)
	err = row.Err()
	if err != nil {
		return nil, err
	}

	var user domain.User

	err = scanToUser(row, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUserByHash(hash string) error {
	db, err := dbCore.Connect()
	if err != nil {
		return nil
	}

	defer db.Close()

	query := SQLDELETE + " WHERE refresh_hash=$1;"
	_, err = db.Exec(query, hash)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserByParams(params map[string]any, userId int) (*domain.User, error) {
	db, err := dbCore.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	query, queryArgs := configureQueryBy(SQLUPDATE, params)
	query += fmt.Sprintf(" WHERE id=%d", userId) + " RETURNING *;"

	row := db.QueryRow(query, queryArgs...)
	err = row.Err()
	if err != nil {
		return nil, err
	}

	var user domain.User

	if err := scanToUser(row, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func configureQueryBy(sqlAction string, params map[string]any) (string, []any) {
	query := sqlAction
	queryArgs := []any{}
	argumentCounter := 1

	// Assemble query, args by map arams. It's looks like:
	// SELECT * FROM users WHERE field1 = $1 AND field1 = $1
	for key, value := range params {
		if value == nil {
			continue
		}

		query += fmt.Sprintf("%s = $%d", key, argumentCounter)

		if len(params) != argumentCounter {
			if sqlAction == SQLUPDATE {
				query += ","
			} else {
				query += " AND "
			}
		}

		queryArgs = append(queryArgs, value)

		argumentCounter++
	}

	return query, queryArgs
}

func scanToUser(row *sql.Row, user *domain.User) error {
	if err := row.Scan(&user.ID, &user.Guid, &user.Hash, &user.IP); err != nil {
		return err
	}

	return nil
}

package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rohanvsuri/minecraft-dashboard/internal/graph/model"
)

func UserExistsByEmail(pool *pgxpool.Pool, email string) (bool, error) {
	// rows, err := pool.QueryRow(context.Background(),
	// 	"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email)

	users, err := AllUsers(pool)
	fmt.Println(users, err)

	if err != nil {
		return false, err
	}
	for _, user := range users {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
	// return exists, err
}

func AllUsers(pool *pgxpool.Pool) ([]*model.User, error) {
	rows, err := pool.Query(context.Background(),
		"SELECT id, email, name FROM users")
	print(rows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Email, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

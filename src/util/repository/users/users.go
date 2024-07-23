package users

import (
	"database/sql"
	"fmt"
	"golang-shopeekuy/src/util/repository/model/users"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{
		db: db,
	}
}

func (s *store) RegisterUser(bReq users.User) (*uuid.UUID, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	var userID uuid.UUID

	queryCreate := `
	INSERT INTO users (email, username, role, address, category_preferences, created_at, updated_at
	) VALUES(
	 $1,
	 $2,
     $3,
     $4,
     $5,
     NOW()
	 ) RETURNING id
	`
	if err := tx.QueryRow(
		queryCreate,
		bReq.Email,
		bReq.Username,
		bReq.Role,
		bReq.Address,
		pq.Array(bReq.CategoryPreferences),
	).Scan(&userID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}
	return &userID, nil

}

func (s *store) GetUserDetails(bReq users.User) (*users.User, error) {
	querySelect := `
	SELECT
		* 
	FROM
		users 
	`
	var queryConditions []string
	if bReq.Email != "" {
		queryConditions = append(queryConditions, fmt.Sprintf("email='%s'", bReq.Email))
	}

	if bReq.ID != uuid.Nil {
		queryConditions = append(queryConditions, fmt.Sprintf("id='%v'", bReq.ID))
	}

	if len(queryConditions) > 0 {
		querySelect += " WHERE " + strings.Join(queryConditions, " AND ")
	}
	var response users.User

	rows, err := s.db.Query(querySelect)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&response.ID,
			&response.Email,
			&response.Username,
			&response.Role,
			&response.Address,
			pq.Array(&response.CategoryPreferences),
			&response.CreatedAt,
			&response.UpdatedAt,
			&response.DeletedAt,
		); err != nil {
			if err != sql.ErrNoRows {
				return nil, fmt.Errorf("no users found")
			}
			return nil, fmt.Errorf("failed to fetch data users")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed iterate users: %v", err)

	}
	return &response, nil
}

package basicauthenticator

import (
	"bufio"
	"io"
	"log"
	"strings"
)

// User represents a single user.
type User struct {
	Userid   string `json:"sub"`
	Password string `json:"-"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// UserDB represents the user database.
type UserDB map[string]User

// NewUserDB reads a TSV stream and creates the user database from that. Each
// line in the stream represents a user entry.
func NewUserDB(r io.Reader) UserDB {
	db := make(UserDB)
	var line string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line = scanner.Text()
		parts := strings.Split(line, "\t")
		len := len(parts)
		if len < 2 {
			log.Fatalf("each user entry must have at least 2 fields - userid and password")
		}
		u := User{}
		if len > 3 {
			u.Email = parts[3]
		}
		if len > 2 {
			u.Name = parts[2]
		}
		u.Userid = parts[0]
		u.Password = parts[1]
		db[u.Userid] = u
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error while reading user database: %v", err)
	}

	return db
}

// Authenticate returns the User struct if successful and nil otherwise.
func (db UserDB) Authenticate(userid, password string) *User {
	u, ok := db[userid]
	if !ok {
		return nil
	}
	if u.Password != password {
		return nil
	}
	return &u
}

// Size returns the number of users in the database,
func (db UserDB) Size() int {
	return len(db)
}

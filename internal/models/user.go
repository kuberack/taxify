package models

type User struct {
	Id        int
	PhoneNum  string
	VerifySid string
}

// Create a new user, save user info into the database
func (user *User) Create() (err error) {
	// Postgres does not automatically return the last insert id, because it would be wrong to assume
	// you're always using a sequence.You need to use the RETURNING keyword in your insert to get this
	// information from postgres.
	statement := "insert into users (phone_number, verify_sid) values (?, ?)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	// Exec the statement
	result, err := stmt.Exec(user.PhoneNum, user.VerifySid)
	if err != nil {
		return
	}

	user_id, err := result.LastInsertId()
	if err != nil {
		return
	}

	// use QueryRow to return a row and scan the returned id into the User struct
	err = Db.QueryRow("SELECT user_id, phone_number, verify_sid FROM users WHERE user_id = ?", user_id).
		Scan(&user.Id, &user.PhoneNum, &user.VerifySid)
	return
}

// Get a single user given the UUID
func UserByID(id int) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT user_id, phone_number, verify_sid FROM users WHERE user_id = ?", user.Id).
		Scan(&user.Id, &user.PhoneNum, &user.VerifySid)
	return
}

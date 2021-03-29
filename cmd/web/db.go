package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// połączenie z bazą danych
func (ldb *DictionaryDatabase) openDB(dsn string) error {
	hDb, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if err = hDb.Ping(); err != nil {
		return err
	}
	ldb.db = hDb
	return nil
}

// wyszukiwanie tłumaczenia na angielski (source=polish) lub polski (source=english)
func (ldb *DictionaryDatabase) recordFind(source, word string) []Word {
	var words []Word

	query := fmt.Sprintf("SELECT id, english, polish from engpol where %s=?", source)
	result, err := ldb.db.Query(query, word)
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	for result.Next() {
		var word Word
		err := result.Scan(&word.ID, &word.English, &word.Polish)
		if err != nil {
			log.Fatal(err)
		}
		words = append(words, word)
	}
	return words
}

// dołączanie rekordu do słownika
func (ldb *DictionaryDatabase) recordAdd(english, polish string) (bool, error) {
	stmt, err := ldb.db.Prepare("INSERT INTO engpol(english, polish) VALUES(?, ?)")
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(english, polish)
	if err != nil {
		return false, err
	}
	return true, nil
}

// aktualizacja danych w słowniku
func (ldb *DictionaryDatabase) recordUpdate(id int, english, polish string) (bool, error) {
	stmt, err := ldb.db.Prepare("UPDATE engpol SET english=?, polish=? WHERE id=?")
	if err != nil {
		return false, err
	}

	result, err := stmt.Exec(english, polish, id)
	if err != nil {
		return false, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, errors.New("update failed")
	}
	return true, nil
}

// funkcja usuwa zapis z bazy danych
func (ldb *DictionaryDatabase) recordDelete(id int) (bool, error) {
	stmt, err := ldb.db.Prepare("DELETE FROM engpol WHERE id=?")
	if err != nil {
		return false, err
	}

	result, err := stmt.Exec(id)
	if err != nil {
		return false, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, errors.New("delete failed`")
	}
	return true, nil
}

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	ErrorUpdateFailed error = errors.New("database - update record failed")
	ErrorDeleteFailed error = errors.New("database - delete record failed")
	ErrorServerFailed error = errors.New("unable to connect to mysql server")
)

// połączenie z bazą danych
func (ldb *DictionaryDatabase) openDB(dsn string) error {
	var hDb *sql.DB
	var err error
	var count int = 0

	for {
		count += 1

		hDb, err = sql.Open("mysql", dsn)
		if err != nil {
			return err
		}

		err = hDb.Ping()
		if err == nil {
			fmt.Println("mySQL server is up.")
			break
		}

		if count <= *cfg.wait {
			if count%2 == 0 {
				fmt.Println("Waiting for the mySQL server...", err.Error())
			}
			time.Sleep(1 * time.Second)
			continue
		} else {
			return ErrorServerFailed
		}

	}

	ldb.db = hDb
	return nil
}

// wyszukiwanie tłumaczenia na angielski (source=polish) lub polski (source=english)
// jeżeli source=english przeszukiwana jest najpierw szybsza baza Redis
func (ldb *DictionaryDatabase) recordFind(source, word string) []Word {
	var words []Word
	var conn redis.Conn
	var err error

	if source == "english" {
		conn, err = redis.Dial("tcp", "redis:6379")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		reply, err := redis.StringMap(conn.Do("HGETALL", word))
		if err != nil {
			log.Fatal(err)
		}

		if len(reply) != 0 {
			fmt.Println(word, "- Redis")
			rWord := Word{}
			rWord.ID = reply["id"]
			rWord.Polish = reply["polish"]
			rWord.English = reply["english"]
			words = append(words, rWord)

			return words
		}
	}

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

	// save (only first) translations if exists
	if source == "english" && len(words) > 0 {
		fmt.Println(word, "- MySQL")
		_, err = conn.Do("HSET", word,
			"english", words[0].English,
			"polish", words[0].Polish,
			"id", words[0].ID,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	return words
}

// dołączanie rekordu do słownika
func (ldb *DictionaryDatabase) recordAdd(english, polish string) (bool, int64, error) {
	stmt, err := ldb.db.Prepare("INSERT INTO engpol(english, polish) VALUES(?, ?)")
	if err != nil {
		return false, 0, err
	}

	result, err := stmt.Exec(english, polish)
	if err != nil {
		return false, 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return false, 0, err
	}

	return true, id, nil
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
		return false, ErrorUpdateFailed
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
		return false, ErrorDeleteFailed
	}
	return true, nil
}

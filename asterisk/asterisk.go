package asterisk

import (
	"database/sql"
	"fmt"
	"log"

	"time"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

// Missed структра содержащая информацию о звонке
type Missed struct {
	uid   string `db:"uniqueid"`
	Src   string `db:"src"`
	dst   string `db:"dst"`
	Did   string `db:"did"`
	staus string `db:"disposition"`
}

// Load подключается к базе и загружает информацию по
// пропущенным звонкам
func Load(conf string, sec int64) []Missed {
	log.Println("Load", conf)
	result := make([]Missed, 0)
	db, err := sql.Open("mysql", conf)
	if err != nil {
		log.Println(err.Error())
		return result
	}
	now := time.Now()
	log.Println("now =", now)
	minuteago := now.Add(time.Duration(sec * 1000000000))
	log.Println("minuteago =", minuteago)
	query := fmt.Sprintf("SELECT uniqueid,src,dst,did,disposition FROM cdr WHERE calldate > '%s' AND did != ''", minuteago.Format("2006-01-02 15:04:05"))
	log.Println(query)
	answeredid := []string{}
	addedid := []string{}
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err.Error())
		return []Missed{}
	}
	for rows.Next() {
		var (
			uid, src, dst, did, status string
		)
		err := rows.Scan(&uid, &src, &dst, &did, &status)
		if err != nil {
			log.Println(err.Error())
		}
		log.Println(src, dst, did, status)
		if status == "ANSWERED" {
			answeredid = append(answeredid, uid)
		} else {
			if notcontain(answeredid, uid) && notcontain(addedid, uid) && did != "" {
				addedid = append(addedid, uid)
				result = append(result, Missed{uid: uid, Did: did, dst: dst, staus: status, Src: src})
			}
		}
	}
	db.Close()
	log.Printf("%#v\n", result)
	return result
}

// noncontain не содержится в слайсе
func notcontain(lst []string, item string) bool {
	for _, it := range lst {
		if it == item {
			return false
		}
	}
	return true
}

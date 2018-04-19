package asterisk

import (
	"log"

	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Missed struct {
	uid   string `db:"uniqueid"`
	src   string `db:"src"`
	dst   string `db:"dst"`
	did   string `db:"did"`
	staus string `db:"disposition"`
}

// Load подключается к базе и загружает информацию по
// пропущенным звонкам звонкам
func Load(conf string) []Missed {
	result := []Missed{}
	tmp := []Missed{}
	db, err := sqlx.Open("mysql", conf)
	if err != nil {
		log.Println(err.Error())
		return result
	}
	now := time.Now()
	minuteago := now.Add(time.Duration(-60))
	db.Select(&tmp, "SELECT src,dst,did,uniqueid,disposition FROM cdr WHERE calldate > $1 AND did != '' order by disposition DESC", minuteago)
	answeredid := []string{}
	addedid := []string{}
	for _, call := range tmp {
		if call.staus == "ANSWERED" {
			answeredid = append(answeredid, call.uid)
		}
	}
	for _, call := range tmp {
		if notcontain(answeredid, call.uid) && notcontain(addedid, call.uid) {
			addedid = append(addedid, call.uid)
			result = append(result, call)
		}
	}
	db.Close()
	return result
}

func notcontain(lst []string, item string) bool {
	for _, it := range lst {
		if it == item {
			return false
		}
		return true
	}
}

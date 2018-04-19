package asterisk

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Missed struct {
	src string
	dst string
	did string
}

// Load подключается к базе и загружает информацию по
// пропущенным звонкам звонкам
func Load(conf string) []Missed {
	result := make([]Missed, 0)
	tmp := &Missed{src: "", dst: "", did: ""}
	db, err := sqlx.Open("mysql", conf)
	if err != nil {
		log.Println(err.Error())
		return result
	}

}

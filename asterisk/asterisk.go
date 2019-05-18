package asterisk

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql" // mysql driver
)

// Missed структра содержащая информацию о звонке
type Missed struct {
	Uid    string `db:"uniqueid"`
	Src    string `db:"src"`
	dst    string `db:"dst"`
	Did    string `db:"did"`
	status string `db:"disposition"`
	file   string `db:"recordingfile"`
}

// Load подключается к базе и загружает информацию по
// пропущенным звонкам
func Load(conf string, sec int64) []Missed {
	result := make([]Missed, 0)

	config, err := mysql.ParseDSN(conf)
	if err != nil {
		log.Println(err.Error())
		return result
	}
	config.AllowNativePasswords = true
	connector, err := mysql.NewConnector(config)
	if err != nil {
		log.Println(err.Error())
		return result
	}

	db := sql.OpenDB(connector)
	query := `SELECT uniqueid, src, dst, did, disposition, recordingfile FROM cdr WHERE 
	calldate BETWEEN NOW() - INTERVAL 2 MINUTE AND NOW() - INTERVAL 1 MINUTE AND did != ''`
	addedid := []string{}
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err.Error())
		return result
	}
	for rows.Next() {
		var call Missed
		err := rows.Scan(&call.Uid, &call.Src, &call.dst, &call.Did, &call.status, &call.file)
		if err != nil {
			log.Println(err.Error())
		}
		if contain(addedid, call.Uid) {
			continue
		}
		if call.IsVoiceMail() {
			addedid = append(addedid, call.Uid)
			result = append(result, call)
		}
		if !call.IsAnswered() && call.Did != "" {
			addedid = append(addedid, call.Uid)
			result = append(result, call)
		}
	}
	db.Close()
	log.Printf("calls result:\n%#v\n", result)
	return result
}

// contain содержится в слайсе
func contain(lst []string, item string) bool {
	for _, it := range lst {
		if it == item {
			return true
		}
	}
	return false
}

// IsAnswered проверка статуса вызова
func (c *Missed) IsAnswered() bool {
	return c.status == "ANSWERED"
}

// IsVoiceMail вызов с голосовым сообщением?
func (c *Missed) IsVoiceMail() bool {
	return strings.HasPrefix(c.dst, "vm")
}

// GetFilePath получение имени файла
func (c *Missed) GetFilePath() string {
	path := time.Now().Format("/var/spool/asterisk/monitor/2006/01/02/")
	return path + c.file
}

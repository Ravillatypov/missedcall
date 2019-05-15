package asterisk

import (
	"database/sql"
	"fmt"
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
	log.Println("Load", conf)
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

	now := time.Now()
	minuteago := now.Add(time.Duration(-sec * 1000000000))
	query := fmt.Sprintf(`SELECT uniqueid, src, dst, did, disposition, recordingfile 
						  FROM cdr WHERE calldate > '%s' AND did != ''
						  AND is_notify IS NOT NULL`, minuteago.Format("2006-01-02 15:04:05"))
	log.Println(query)
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
	_, err = db.Query(fmt.Sprintf(`UPDATE cdr SET is_notify=1 WHERE uniqueid in (%s)`, strings.Join(addedid, ",")))
	if err != nil {
		log.Println(err.Error())
		return result
	}
	db.Close()
	log.Printf("%#v\n", result)
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

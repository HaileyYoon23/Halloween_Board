package db

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

var DB *sql.DB

var (
	ErrDuplicateGift = errors.New("already exist gift")
)

type Chat struct {
	Talk   string
	Size   int
	Length int
	Likes  int
}

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:315931@tcp(127.0.0.1:3306)/halloween")
	//DB, err = sql.Open("mysql", "doadmin:VmKcD8m0tWnKOv3E@tcp(xmas-do-user-10452914-0.b.db.ondigitalocean.com:25060)/xmas")

	if err != nil {
		panic(err)
	}

	maxOpen, _ := strconv.Atoi(os.Getenv("MAX_OPEN_CONN"))
	if maxOpen <= 0 {
		maxOpen = runtime.NumCPU() * 2
	}

	DB.SetConnMaxIdleTime(10 * time.Minute)
	DB.SetMaxIdleConns(maxOpen)
	DB.SetMaxOpenConns(maxOpen)

	createTableQuery := `
		create table IF NOT EXISTS chat ( 
		id BIGINT NOT NULL AUTO_INCREMENT,
		talk VARCHAR(200),
		likes	INT,
		PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		panic(err)
	}
}

func AddTalk(db *sql.DB, talk string) (err error) {
	var id int64

	err = db.QueryRow("select id from chat where talk=?", talk).Scan(&id)
	if id > 0 {
		err = ErrDuplicateGift
		return
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into chat (talk, likes) values (?, ?)")
	_, err = stmt.Exec(talk, 0)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	tx.Commit()
	return
}

func GetTalks(db *sql.DB) (talkList []Chat, err error) {
	var talk string
	var likes int

	rows, err := db.Query("select talk, likes from chat order by likes desc")

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&talk, &likes)
		if err != nil {
			return
		}
		size := likes*5 + 20
		talkList = append(talkList, Chat{
			Talk:   talk,
			Size:   size,
			Length: len(talk) * size / 2,
			Likes:  likes,
		})
	}

	return
}

func UpdateLike(talk string) (err error) {
	var likes int

	err = DB.QueryRow("select likes from chat where talk=?", talk).Scan(&likes)
	if err != nil {
		return
	}

	_, err = DB.Exec("update chat set likes=? where talk=?", likes+1, talk)
	if err != nil {
		return
	}

	return
}

func UpdateDislike(present string) (err error) {
	var dislikes int

	err = DB.QueryRow("select dislikes from presents where present=?", present).Scan(&dislikes)
	if err != nil {
		return
	}

	_, err = DB.Exec("update presents set dislikes=? where present=?", dislikes+1, present)
	if err != nil {
		return
	}

	return
}

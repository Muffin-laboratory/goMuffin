package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var wg sync.WaitGroup

// 이 스크립트는 MariaDB -> MongoDB로의 전환을 위해 만들었음.
func main() {
	mariaURL := os.Getenv("PREVIOUS_DATABASE_URL")
	mongoURL := configs.Config.DatabaseURL

	dbConnectionQuery := "?parseTime=true"

	wg.Add(3)

	// statement -> text
	go func() {
		defer wg.Done()

		var text, persona string
		var createdAt time.Time
		mariaDB, err := sql.Open("mysql", mariaURL+dbConnectionQuery)
		if err != nil {
			panic(err)
		}

		mongoDB, err := mongo.Connect(options.Client().ApplyURI(mongoURL))
		if err != nil {
			panic(err)
		}

		defer mongoDB.Disconnect(context.TODO())
		defer mariaDB.Close()
		rows, err := mariaDB.Query("select text, persona, created_at from statement;")
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		i := 1

		for rows.Next() {
			fmt.Printf("statement %d\n", i)
			err = rows.Scan(&text, &persona, &createdAt)
			if err != nil {
				panic(err)
			}

			if text == "" {
				text = "살ㄹ려주세요"
			}
			databases.Texts.InsertOne(context.TODO(), databases.InsertText{
				Text:      text,
				Persona:   persona,
				CreatedAt: createdAt,
			})
			i++
		}
	}()

	// nsfw_content -> text
	go func() {
		defer wg.Done()

		var text, persona string
		var createdAt time.Time
		mariaDB, err := sql.Open("mysql", mariaURL+dbConnectionQuery)
		if err != nil {
			panic(err)
		}

		mongoDB, err := mongo.Connect(options.Client().ApplyURI(mongoURL))
		if err != nil {
			panic(err)
		}

		defer mongoDB.Disconnect(context.TODO())
		defer mariaDB.Close()
		rows, err := mariaDB.Query("select text, persona, created_at from nsfw_content;")
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		i := 1

		for rows.Next() {
			fmt.Printf("nsfw_content %d\n", i)
			err = rows.Scan(&text, &persona, &createdAt)
			if err != nil {
				panic(err)
			}

			if text == "" {
				text = "살ㄹ려주세요"
			}
			databases.Texts.InsertOne(context.TODO(), databases.InsertText{
				Text:      text,
				Persona:   persona,
				CreatedAt: createdAt,
			})
			i++
		}
	}()

	// learn -> learn
	go func() {
		defer wg.Done()

		var command, result, userId string
		var createdAt time.Time
		mariaDB, err := sql.Open("mysql", mariaURL+dbConnectionQuery)
		if err != nil {
			panic(err)
		}

		mongoDB, err := mongo.Connect(options.Client().ApplyURI(mongoURL))
		if err != nil {
			panic(err)
		}

		defer mongoDB.Disconnect(context.TODO())
		defer mariaDB.Close()
		rows, err := mariaDB.Query("select command, result, user_id, created_at from learn;")
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		i := 1

		for rows.Next() {
			fmt.Printf("learn %d\n", i)
			err = rows.Scan(&command, &result, &userId, &createdAt)
			if err != nil {
				panic(err)
			}

			databases.Learns.InsertOne(context.TODO(), databases.InsertLearn{
				Command:   command,
				Result:    result,
				UserId:    userId,
				CreatedAt: createdAt,
			})
			i++
		}
	}()

	// 모든 고루틴이 끝날 떄 까지 대기
	wg.Wait()
}

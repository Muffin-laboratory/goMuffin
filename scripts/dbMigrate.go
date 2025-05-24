package scripts

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"git.wh64.net/muffin/goMuffin/configs"

	"github.com/devproje/commando"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var wg sync.WaitGroup

// 이 스크립트는 MariaDB -> MongoDB로의 전환을 위해 만들었음.
func DBMigrate(n *commando.Node) error {
	mariaURL := os.Getenv("PREVIOUS_DATABASE_URL")
	mongoURL := configs.Config.DatabaseURL
	dbName := configs.Config.DatabaseName

	dbConnectionQuery := "?parseTime=true"

	wg.Add(3)

	fmt.Println("[경고] 해당 명령어는 다음 버전에서 사라져요.")

	// statement -> text
	go func() {
		defer wg.Done()

		newDataList := []any{}
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
			var text, persona string
			var createdAt time.Time

			fmt.Printf("statement %d\n", i)
			err = rows.Scan(&text, &persona, &createdAt)
			if err != nil {
				panic(err)
			}

			if text == "" {
				text = "살ㄹ려주세요"
			}

			newDataList = append(newDataList, bson.M{
				"text":       text,
				"persona":    persona,
				"created_at": createdAt,
			})
			i++
		}
		_, err = mongoDB.Database(dbName).Collection("text").InsertMany(context.TODO(), newDataList)
		if err != nil {
			panic(err)
		}
	}()

	// nsfw_content -> text
	go func() {
		defer wg.Done()

		newDataList := []any{}
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
			var text, persona string
			var createdAt time.Time

			fmt.Printf("nsfw_content %d\n", i)
			err = rows.Scan(&text, &persona, &createdAt)
			if err != nil {
				panic(err)
			}

			if text == "" {
				text = "살ㄹ려주세요"
			}

			newDataList = append(newDataList, bson.M{
				"text":       text,
				"persona":    persona,
				"created_at": createdAt,
			})
			i++
		}

		_, err = mongoDB.Database(dbName).Collection("text").InsertMany(context.TODO(), newDataList)
		if err != nil {
			panic(err)
		}
	}()

	// learn -> learn
	go func() {
		defer wg.Done()

		newDataList := []any{}
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
			var command, result, userId string
			var createdAt time.Time

			fmt.Printf("learn %d\n", i)
			err = rows.Scan(&command, &result, &userId, &createdAt)
			if err != nil {
				panic(err)
			}

			newDataList = append(newDataList, bson.M{
				"command":    command,
				"result":     result,
				"user_id":    userId,
				"created_at": createdAt,
			})
			i++
		}

		_, err = mongoDB.Database(dbName).Collection("learn").InsertMany(context.TODO(), newDataList)
		if err != nil {
			panic(err)
		}
	}()

	// 모든 고루틴이 끝날 떄 까지 대기
	wg.Wait()
	fmt.Println("데이터 마이그레이션이 끝났어요.")
	return nil
}

package scripts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/devproje/commando"
	"github.com/devproje/commando/option"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var date time.Time = time.Now()

type textJSONLData struct {
	Text    string `json:"text"`
	Persona string `json:"persona,omitempty"`
}

type learnJSONLData struct {
	Command string `json:"command"`
	Result  string `json:"result"`
}

func getDate() string {
	year := strconv.Itoa(date.Year())
	month := strconv.Itoa(int(date.Month()))
	day := strconv.Itoa(date.Day())
	hour := strconv.Itoa(date.Hour())
	minute := strconv.Itoa(date.Minute())
	sec := strconv.Itoa(date.Second())

	if len(month) < 2 {
		month = "0" + month
	}

	if len(day) < 2 {
		day = "0" + day
	}

	if len(hour) < 2 {
		hour = "0" + hour
	}

	if len(minute) < 2 {
		minute = "0" + minute
	}

	if len(sec) < 2 {
		sec = "0" + sec
	}
	return year + month + day + hour + minute + sec
}

func checkDir(path string) error {
	_, err := os.ReadDir(path)
	if err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveFileToJSON(path, name string, data any) error {
	bytes, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/%s.json", path, name))
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func saveFileToJSONL(path, name string, data any) error {
	var content string

	switch data := data.(type) {
	case []textJSONLData:
		for _, data := range data {
			bytes, err := json.Marshal(data)
			if err != nil {
				return err
			}
			content += string(bytes) + "\n"
		}
	case []learnJSONLData:
		for _, data := range data {
			bytes, err := json.Marshal(data)
			if err != nil {
				return err
			}
			content += string(bytes) + "\n"
		}
	}

	f, err := os.Create(fmt.Sprintf("%s/%s.jsonl", path, name))
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func ExportData(n *commando.Node) error {
	defer databases.Database.Client.Disconnect(context.TODO())

	var wg sync.WaitGroup
	ch := make(chan error, 3)

	fileType, err := option.ParseString(*n.MustGetOpt("type"), n)
	if err != nil {
		return err
	}

	if fileType != "json" && fileType != "jsonl" {
		return fmt.Errorf("파일 형식은 txt또는 json또는 jsonl이여야 해요")
	}

	refined, err := option.ParseBool(*n.MustGetOpt("refined"), n)
	if err != nil {
		return err
	}

	path, err := option.ParseString(*n.MustGetOpt("export-path"), n)
	if err != nil {
		return err
	}

	path += "/" + getDate()

	err = checkDir(path)
	if err != nil {
		return err
	}

	wg.Add(3)

	if fileType == "jsonl" {
		fmt.Println("NOTE: 파일 형식이 'jsonl'인 경우 일부데이터만 추출 됩니다.")
	}

	// 머핀 데이터 추출
	go func() {
		defer wg.Done()

		var data []databases.Text

		cur, err := databases.Database.Texts.Find(context.TODO(), bson.D{{Key: "persona", Value: "muffin"}})
		if err != nil {
			ch <- err
			return
		}

		defer cur.Close(context.TODO())

		err = cur.All(context.TODO(), &data)
		if err != nil {
			ch <- err
			return
		}

		if refined {
			for i, text := range data {
				if utils.EmojiRegexp.Match([]byte(text.Text)) {
					data = append(data[:i], data[i+1:]...)
					return
				}

				text.Text = strings.TrimPrefix(text.Text, "머핀아 ")
			}
		}

		if fileType == "json" {
			err = saveFileToJSON(path, "muffin", data)
			if err != nil {
				ch <- err
				return
			}
		} else if fileType == "jsonl" {
			var newData []textJSONLData

			for _, data := range data {
				newData = append(newData, textJSONLData{data.Text, ""})
			}

			err = saveFileToJSONL(path, "muffin", newData)
			if err != nil {
				ch <- err
				return
			}
		}

		fmt.Println("머핀 데이터 추출 완료")
	}()

	// nsfw 데이터 추출
	go func() {
		defer wg.Done()

		var data []databases.Text

		cur, err := databases.Database.Texts.Find(context.TODO(), bson.D{
			{
				Key: "persona",
				Value: bson.M{
					"$regex": "^user",
				},
			},
		})
		if err != nil {
			ch <- err
			return
		}

		defer cur.Close(context.TODO())

		err = cur.All(context.TODO(), &data)
		if err != nil {
			ch <- err
			return
		}

		if fileType == "json" {
			err = saveFileToJSON(path, "nsfw", data)
			if err != nil {
				ch <- err
				return
			}
		} else if fileType == "jsonl" {
			var newData []textJSONLData

			for _, data := range data {
				newData = append(newData, textJSONLData{data.Text, data.Persona})
			}

			err = saveFileToJSONL(path, "nsfw", newData)
			if err != nil {
				ch <- err
				return
			}
		}

		fmt.Println("nsfw 데이터 추출 완료")
	}()

	// 지식 데이터 추출
	go func() {
		defer wg.Done()

		var data []databases.Learn

		cur, err := databases.Database.Learns.Find(context.TODO(), bson.D{{}})
		if err != nil {
			ch <- err
			return
		}

		defer cur.Close(context.TODO())

		err = cur.All(context.TODO(), &data)
		if err != nil {
			ch <- err
			return
		}

		if fileType == "json" {
			err = saveFileToJSON(path, "learn", data)
			if err != nil {
				ch <- err
				return
			}
		} else if fileType == "jsonl" {
			var newData []learnJSONLData

			for _, data := range data {
				newData = append(newData, learnJSONLData{data.Command, data.Result})
			}

			err = saveFileToJSONL(path, "learn", newData)
			if err != nil {
				ch <- err
				return
			}
		}

		fmt.Println("지식 데이터 추출 완료")
	}()

	wg.Wait()
	close(ch)

	for err = range ch {
		fmt.Println(err)
	}
	return nil
}

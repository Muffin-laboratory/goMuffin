package scripts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/devproje/commando"
	"github.com/devproje/commando/option"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var date time.Time = time.Now()

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

	path += getDate()
	err = checkDir(path)
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

func saveFileToTXT(path, name string, data []databases.Text) error {
	var content string

	for _, data := range data {
		content += data.Text + "\n"
	}

	path += getDate()
	err := checkDir(path)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/%s.txt", path, name))
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
	ch := make(chan error)
	defer databases.Client.Disconnect(context.TODO())
	fileType, err := option.ParseString(*n.MustGetOpt("type"), n)
	if err != nil {
		return err
	}

	if fileType != "txt" && fileType != "json" {
		return fmt.Errorf("파일 형식은 txt또는 json이여야 해요")
	}

	refined, err := option.ParseBool(*n.MustGetOpt("refined"), n)
	if err != nil {
		return err
	}

	path, err := option.ParseString(*n.MustGetOpt("export-path"), n)
	if err != nil {
		return err
	}

	texts := databases.Texts
	// learns := databases.Learns

	// 머핀 데이터 추출
	go func() {
		var data []databases.Text

		cur, err := texts.Find(context.TODO(), bson.D{{Key: "persona", Value: "muffin"}})
		if err != nil {
			ch <- err
		}

		cur.All(context.TODO(), &data)

		if refined {
			for i, text := range data {
				if utils.EmojiRegexp.Match([]byte(text.Text)) {
					data = append(data[:i], data[i+1:]...)
				}

				text.Text = strings.TrimPrefix(text.Text, "머핀아 ")
			}
		}

		if fileType == "json" {
			ch <- saveFileToJSON(path, "muffin", data)
		} else {
			fmt.Println("NOTE: 파일 형식이 'txt'인 경우 머핀 데이터만 txt형식으로 저장되고, 나머지는 json으로 저장됩니다.")
			ch <- saveFileToTXT(path, "muffin", data)
		}
	}()
	return <-ch
}

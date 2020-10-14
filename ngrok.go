package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
)

//const ngrokarm = "https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-2.2.8-linux-arm.zip"
const ngrokarm = "https://bin.equinox.io/a/nmkK3DkqZEB/ngrok-2.2.8-linux-arm64.zip"
const ngrokarm64 = "https://bin.equinox.io/a/nmkK3DkqZEB/ngrok-2.2.8-linux-arm64.zip"
const ngrokamd64 = "https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-amd64.zip"
const ngrok386 = "https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-386.zip"

func unzip(src, dest string) ([]string, error) { // 3.party

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		fmt.Println("1", err)
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		/*if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}*/

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func checkNgrok() bool {
	if _, err := os.Stat("./ngrok"); os.IsNotExist(err) {
		return false
	}
	return true
}

func dandunzipngrok(url string) error {
	req, err := http.Get(url)
	if err != err {
		fmt.Println("req err : ", err)
		return err
	}
	file, err := os.Create("ngrok.zip")
	if err != err {
		fmt.Println("create ngrok.zip err : ", err)
		return err
	}
	_, err = io.Copy(file, req.Body)
	if err != nil {
		fmt.Println("io copy err : ", err)
		return err
	}
	_, err = unzip("ngrok.zip", ".")
	if err != nil {
		fmt.Println("unzip err : ", err)
		return err
	}
	err = os.Chmod("ngrok", 0111)
	if err != nil {
		fmt.Println("chmod err : ", err)
		return err
	}
	os.Remove("ngrok.zip")
	return nil

}

func getNgrok() error {
	switch runtime.GOARCH {
	case "arm", "ARM":
		return dandunzipngrok(ngrokarm)
	case "arm64", "ARM64":
		return dandunzipngrok(ngrokarm64)
	case "amd64", "AMD64":
		return dandunzipngrok(ngrokamd64)
	case "386", "x86":
		return dandunzipngrok(ngrok386)
	default:
		return fmt.Errorf("UnSupported Arch %s", runtime.GOARCH)
	}
}

func runNgrok(port string) (<-chan cmd.Status, error) {
	if !checkNgrok() {
		return nil, fmt.Errorf("Ngrok Not Found")
	}
	ngrokCommand := cmd.NewCmdOptions(cmd.Options{Buffered: false, Streaming: false}, "./ngrok", "http", port)
	ngrokcmd := ngrokCommand.Start()
	defer ngrokCommand.Stop()

	/*cmd := exec.Command("./ngrok", "http", port)
	err := cmd.Start()
	if err != nil {
		return err
	}*/
	return ngrokcmd, nil
}

func getNgrokLink() (string, error) {
	req, err := http.Get("http://127.0.0.1:4040/api/tunnels")

	if err != nil {
		log.Println(err)
		return "", err

	}
	ham, _ := ioutil.ReadAll(req.Body)
	html := string(ham)
	html = strings.Replace(html, "&#34;", "", -1)
	html = strings.Replace(html, "\"", "", -1)
	publicURLIndex := strings.Index(html, "public_url") + 11
	ngrokURL := strings.Split(html[publicURLIndex:], ",")[0]
	return ngrokURL, nil
	//return html, nil
}

func getNgrokLinkStable() (string, error) {
	var errlog error
	for i := 0; i < 6; i++ {
		log.Println("Ngrok Waiting : ", i)
		time.Sleep(time.Second * 2)
		url, err := getNgrokLink()
		if err != nil {
			errlog = err
			log.Println(err.Error())

		} else {
			fmt.Println("NGROK : ", url)
			if strings.Contains(url, "ngrok.io") {
				return url, nil
			}
		}

	}
	if errlog != nil {
		return "", errlog
	}
	return "", fmt.Errorf("Ngrok cannot start!\nYou must open manuel;\n{ 127.0.0.1:8089 }")
}

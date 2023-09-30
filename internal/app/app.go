package app

import (
	"url_checker/internal/config"
	"url_checker/pkg/logger"

	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger
	config *config.Config
}

type Urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []Url    `xml:"url"`
}

type Url struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

func New() *App {
	app := &App{}

	app.config = config.Setup(app.logger)
	app.logger = logger.New(app.config.LogPath)

	return app
}

func (a *App) Run() {
	file, _ := os.Open(a.config.XmlPath)
	defer file.Close()

	var sitemap Urlset
	decoder := xml.NewDecoder(file)
	err := decoder.Decode(&sitemap)

	if err != nil {
		a.logger.Fatal("Failure parsing xml data to structure")
	}

	runtime.GOMAXPROCS(int(a.config.MaxProc) + 1)

	var statusCh = make(chan string)
	var errCh = make(chan error)

	var wg sync.WaitGroup
	wg.Add(int(a.config.MaxProc) + 1)

	urls1, urls2 := a.splitSlice(sitemap.Urls)

	timeout := time.After(time.Duration(a.config.Timeout) * time.Minute)

	go func() {
		defer wg.Done()
		a.asyncCheck(urls1, statusCh, errCh)
	}()

	go func() {
		defer wg.Done()
		a.asyncCheck(urls2, statusCh, errCh)
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case err := <-errCh:
				a.logger.Error("Error during request error message: ", zap.Error(err))
				return
			case status := <-statusCh:
				a.logger.Info(status)
				a.logger.Info("--------------")
			case <-timeout:
				a.logger.Info("Ending timeout")
				return
			}
		}
	}()

	wg.Wait()

	defer close(statusCh)
	defer close(errCh)
	a.logger.Info("Get all data- success")
}

func (a *App) asyncCheck(urls []Url, statusCh chan<- string, errCh chan<- error) {
	for _, page := range urls {
		statusText, err := a.check(page.Loc)

		if err != nil {
			errCh <- err
			continue
		}

		statusCh <- fmt.Sprintf("Request http get url: %s -> %s", page.Loc, statusText)
	}
}

func (a *App) check(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	res, err1 := http.DefaultClient.Do(req)

	if err1 != nil {
		return "", err1
	}

	defer res.Body.Close()

	return res.Status, nil
}

func (a *App) splitSlice(urls []Url) ([]Url, []Url) {
	var res []Url
	var res1 []Url

	for i := 0; i < len(urls); i++ {
		if i%2 == 0 {
			res = append(res, urls[i])
		} else {
			res1 = append(res1, urls[i])
		}
	}

	return res, res1
}

package service

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/service/api"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type (
	// Transcriber returns words IPA transcriptiops array
	Transcriber interface {
		Process(string, bool) ([]*api.ResultWord, error)
	}

	// TranscriberOne returns possible transcription wariants for one word
	WordTranscriber interface {
		Process(string, bool) (*api.WordInfo, error)
	}

	//Data is service operation data
	Data struct {
		Transcriber     Transcriber
		WordTranscriber WordTranscriber
		Port            int
	}
)

//StartWebServer starts the HTTP service and listens for the requests
func StartWebServer(data *Data) error {
	goapp.Log.Infof("Starting HTTP service at %d", data.Port)
	portStr := strconv.Itoa(data.Port)

	e := initRoutes(data)

	e.Server.Addr = ":" + portStr
	e.Server.ReadHeaderTimeout = 5 * time.Second
	e.Server.ReadTimeout = 15 * time.Second
	e.Server.WriteTimeout = 15 * time.Second

	w := goapp.Log.Writer()
	defer w.Close()
	l := log.New(w, "", 0)
	gracehttp.SetLogger(l)

	return gracehttp.Serve(e.Server)
}

var p *prometheus.Prometheus

func initRoutes(data *Data) *echo.Echo {
	e := echo.New()
	if p == nil {
		p = prometheus.NewPrometheus("tag", nil)
		p.Use(e)
	}

	e.POST("/ipa", handleText(data))
	e.GET("/ipa/:word", handleWord(data))
	e.GET("/live", live(data))

	goapp.Log.Info("Routes:")
	for _, r := range e.Routes() {
		goapp.Log.Infof("  %s %s", r.Method, r.Path)
	}
	return e
}

type textBinder struct{}

func (cb *textBinder) Bind(s *string, c echo.Context) error {
	bodyBytes, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't get data").SetInternal(err)
	}
	*s = string(bodyBytes)
	*s = strings.TrimSpace(string(bodyBytes))
	if *s == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "No input", err)
	}
	return nil
}

func handleText(data *Data) func(echo.Context) error {
	return func(c echo.Context) error {
		defer goapp.Estimate("Service method: ipa")()
		tb := &textBinder{}
		var text string
		if err := tb.Bind(&text, c); err != nil {
			goapp.Log.Error(err)
			return err
		}

		res, err := data.Transcriber.Process(text, c.QueryParam("sampa") == "1")
		if err != nil {
			goapp.Log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Can't segment")
		}

		return c.JSON(http.StatusOK, res)
	}
}

func live(data *Data) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.JSONBlob(http.StatusOK, []byte(`{"service":"OK"}`))
	}
}

func handleWord(data *Data) func(echo.Context) error {
	return func(c echo.Context) error {
		defer goapp.Estimate("Service method: word")()

		word := c.Param("word")
		if word == "" {
			goapp.Log.Error("No word")
			return echo.NewHTTPError(http.StatusBadRequest, "No word")
		}

		res, err := data.WordTranscriber.Process(word, c.QueryParam("sampa") == "1")
		if err != nil {
			goapp.Log.Error(errors.Wrap(err, "Cannot process "+word))
			return echo.NewHTTPError(http.StatusInternalServerError, "Cannot process "+word)
		}
		return c.JSON(http.StatusOK, res)
	}
}

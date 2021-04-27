package main

import (
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/oneword"
	oneworker "github.com/airenas/lt-text-to-ipa/internal/pkg/oneword/worker"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/process"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/process/worker"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/service"
	"github.com/labstack/gommon/color"
	"github.com/spf13/viper"

	"github.com/pkg/errors"
)

func main() {
	goapp.StartWithDefault()

	data := service.Data{}
	data.Port = goapp.Config.GetInt("port")
	var err error

	printBanner()

	mw := &process.MainWorker{}
	err = addProcessors(mw, goapp.Config)
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "Can't init processors"))
	}

	data.Transcriber = mw

	wordW := &oneword.MainWorker{}
	err = addWordProcessors(wordW, goapp.Config)
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "Can't init word processors"))
	}

	data.WordTranscriber = wordW

	err = service.StartWebServer(&data)
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "Can't start the service"))
	}
}

func addProcessors(mw *process.MainWorker, cfg *viper.Viper) error {
	mw.Add(worker.NewCleaner())
	pr, err := worker.NewTagger(cfg.GetString("tagger.url"))
	if err != nil {
		return errors.Wrap(err, "can't init tagger")
	}
	mw.Add(pr)
	pr, err = worker.NewAccentuator(cfg.GetString("accenter.url"))
	if err != nil {
		return errors.Wrap(err, "can't init accenter")
	}
	mw.Add(pr)
	pr, err = worker.NewClitics(cfg.GetString("cliticsDetector.url"))
	if err != nil {
		return errors.Wrap(err, "can't init clitics detector")
	}
	mw.Add(pr)
	pr, err = worker.NewTranscriber(cfg.GetString("transcriber.url"))
	if err != nil {
		return errors.Wrap(err, "can't init transcriber")
	}
	mw.Add(pr)
	pr, err = worker.NewToIPA(cfg.GetString("ipaConverter.url"))
	if err != nil {
		return errors.Wrap(err, "can't init transcriber")
	}
	mw.Add(pr)
	mw.Add(worker.NewResultMaker())
	return nil
}

func addWordProcessors(mw *oneword.MainWorker, cfg *viper.Viper) error {
	pr, err := oneworker.NewAccentuator(cfg.GetString("accenter.url"))
	if err != nil {
		return errors.Wrap(err, "can't init accenter")
	}
	mw.Add(pr)
	pr, err = oneworker.NewTranscriber(cfg.GetString("transcriber.url"))
	if err != nil {
		return errors.Wrap(err, "can't init transcriber")
	}
	mw.Add(pr)
	pr, err = oneworker.NewToIPA(cfg.GetString("ipaConverter.url"))
	if err != nil {
		return errors.Wrap(err, "can't init transcriber")
	}
	mw.Add(pr)
	mw.Add(oneworker.NewResultMaker())
	return nil
}

var (
	version string
)

func printBanner() {
	banner := `
    __  ______   __            __ 
   / / /_  __/  / /____  _  __/ /_
  / /   / /    / __/ _ \| |/_/ __/
 / /___/ /    / /_/  __/>  </ /_  
/_____/_/     \__/\___/_/|_|\__/  
                                  
       __           ________  ___ 
      / /_____     /  _/ __ \/   |
     / __/ __ \    / // /_/ / /| |
    / /_/ /_/ /  _/ // ____/ ___ |
    \__/\____/  /___/_/   /_/  |_|  v: %s 

	%s
________________________________________________________                                                 

`
	cl := color.New()
	cl.Printf(banner, cl.Red(version), cl.Green("https://github.com/airenas/lt-text-to-ipa"))
}

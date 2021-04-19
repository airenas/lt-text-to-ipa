package main

import (
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/lt-text-to-ipa/internal/pkg/service"
	"github.com/labstack/gommon/color"

	"github.com/pkg/errors"
)

func main() {
	goapp.StartWithDefault()

	data := service.Data{}
	data.Port = goapp.Config.GetInt("port")
	var err error

	printBanner()

	err = service.StartWebServer(&data)
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "Can't start the service"))
	}
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

package main

import (
	"encoding/json"
	"github.com/moehandi/go-webapp-base/route"
	"github.com/moehandi/go-webapp-base/helper/database"
	"github.com/moehandi/go-webapp-base/helper/email"
	"github.com/moehandi/go-webapp-base/helper/recaptcha"
	"github.com/moehandi/go-webapp-base/helper/server"
	"github.com/moehandi/go-webapp-base/helper/session"
	"github.com/moehandi/go-webapp-base/helper/view"
	"github.com/moehandi/go-webapp-base/helper/view/plugin"
	"log"
	"os"
	"runtime"
		"github.com/moehandi/go-webapi-base/helper/tomlconfig"
		"github.com/BurntSushi/toml"
)

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)
	// Use all CPU Cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var configPath = "config"+string(os.PathSeparator)+"config.toml"

func main() {
	// Load the configuration file
	tomlconfig.Load(configPath, config)

	// Configure the session cookie store
	session.Configure(config.Session)

	// Connect to database
	database.Connect(config.Database)

	// Configure the Google reCAPTCHA prior to loading view plugins
	recaptcha.Configure(config.Recaptcha)

	// Setup the views
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	view.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		plugin.PrettyTime(),
		recaptcha.Plugin())

	// Start the listener
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)
}

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
		Database  database.Info   `toml:"Database"`
		Email     email.SMTPInfo  `toml:"Email"`
		Recaptcha recaptcha.Info  `toml:"Recaptcha"`
		Server    server.Server   `toml:"Server"`
		Session   session.Session `toml:"Session"`
		Template  view.Template   `toml:"Template"`
		View      view.View       `toml:"View"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

// ParseTOML unmarshals bytes to structs
func (c *configuration) ParseTOML(b []byte) error {
		return toml.Unmarshal(b, &config)
}

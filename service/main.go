package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kardianos/service"
	"github.com/rs/zerolog"
)

var (
	re        = regexp.MustCompile(`to do id ([A-Za-z0-9]+)`)
	log       zerolog.Logger
	srv       *http.Server
	procGroup *sync.WaitGroup
	plist     = `<?xml version="1.0" encoding="UTF-8"?>
    <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
    <plist version="1.0">
        <dict>
            <key>Label</key>
            <string>com.voilaweb.fusion.CraftyThings</string>
            <key>ProgramArguments</key>
            <array>
              <string>/usr/local/bin/craftythingshelper</string>
            </array>
            <key>RunAtLoad</key>
            <true/>
            <key>StandardErrorPath</key>
            <string>/tmp/craftythings.err</string>
        </dict>
    </plist>`
)

type addquery struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type program struct{}

func runCmd(name string, arg ...string) *bytes.Buffer {
	buf := new(bytes.Buffer)
	cmd := exec.Command(name, arg...)
	cmd.Stdout = buf
	cmd.Stderr = buf
	if err := cmd.Run(); err != nil {
		log.Error().Err(err).Msg("Exec error")
		return nil
	}
	return buf
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		log.Error().Msg("Error. Form is not json-encoded.")
		return
	}
	var q addquery
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&q); err != nil {
		log.Error().Err(err).Msg("Decoder error")
		return
	}

	log.Debug().Str("title", q.Title).Msg("Debug")
	log.Debug().Str("content", q.Content).Msg("Debug")

	buf := runCmd("/usr/bin/osascript",
		"-e", "tell application \"Things3\"",
		"-e", "set newToDo to make new to do with properties {name:\""+strings.ReplaceAll(q.Title, `"`, `\"`)+"\", notes:\""+strings.ReplaceAll(q.Content, `"`, `\"`)+"\"}",
		"-e", "end tell")
	if buf == nil {
		return
	}
	gp := re.FindAllStringSubmatch(buf.String(), -1)
	if len(gp) > 0 {
		if len(gp[0]) > 1 {
			w.Write([]byte(gp[0][1]))
		}
	}
}

func (p *program) Start(s service.Service) error {
	log.Info().Msg("Starting service.")

	procGroup = &sync.WaitGroup{}
	procGroup.Add(1)

	go p.run()
	return nil
}

func (p *program) run() {
	http.HandleFunc("/form", formHandler)

	log.Info().Msg("Starting server at port 48484")

	srv = &http.Server{Addr: ":48484"}
	go func() {
		defer procGroup.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Could not start API endpoint.")
		}
	}()
}

func (p *program) Stop(s service.Service) error {
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Error().Err(err).Msg("Could not stop service cleanly.")
		return nil
	}
	procGroup.Wait()

	log.Info().Msg("Service stopped.")
	return nil
}

func selfInstall() (exit bool) {
	buf := runCmd("/bin/sh", "-c", "/usr/bin/pgrep craftythingshelper || true")
	if buf == nil {
		return true
	}
	if len(buf.Bytes()) > 0 {
		log.Error().Msg("I am already running, likely as a service. Bailing.")
		runCmd("/usr/bin/osascript",
			"-e", "display alert \"Crafty Things Helper is already running, likely as a service. Bailing.\"")
		return true
	}

	// So, not running as a service... yet.
	buf = runCmd("/usr/bin/id", "-u")
	if buf == nil {
		return true
	}
	whoami := strings.TrimSpace(buf.String())

	// Look for existing service
	buf = runCmd("/bin/sh", "-c", fmt.Sprintf("/bin/launchctl print gui/%s | grep com.voilaweb.fusion.CraftyThings || true", whoami))
	if buf == nil {
		return true
	}
	if len(buf.Bytes()) > 0 {
		return false
	}
	log.Info().Msg("Installing myself as a service.")
	runCmd("/usr/bin/osascript",
		"-e", "display alert \"Installing Crafty Things Helper as a service.\"")

	// Install executable
	exePath, err := os.Executable()
	if err != nil {
		log.Error().Err(err).Msg("Unable to find location of executable.")
		return true
	}
	buf = runCmd("/usr/bin/install", exePath, "/usr/local/bin/craftythingshelper")
	if buf == nil {
		return true
	}

	err = ioutil.WriteFile("/tmp/com.voilaweb.fusion.craftythings.plist", []byte(plist), 0600)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create plist.")
		return true
	}
	// Bootstrap using plist
	buf = runCmd("/bin/launchctl", "bootstrap", fmt.Sprintf("gui/%s", whoami), "/tmp/com.voilaweb.fusion.craftythings.plist")
	if buf == nil {
		return true
	}
	return true
}

func main() {
	log = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC1123Z}).Level(
		zerolog.InfoLevel).With().Timestamp().Logger()
	log.Info().Str("Version", "0.0.1").Str("Author", "Chris F Ravenscroft").Msg("Crafty Things Craft-to-Things Helper.")

	if selfInstall() {
		return
	}

	svcConfig := &service.Config{
		Name:        "CraftyThings",
		DisplayName: "Crafty Things",
		Description: "Crafty Things helper.",
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Error().Err(err).Msg("Could not create service.")
		return
	}
	err = s.Run()
	if err != nil {
		log.Error().Err(err).Msg("Could not start service.")
		return
	}
}

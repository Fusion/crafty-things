package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kardianos/service"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	re        = regexp.MustCompile(`to do id ([A-Za-z0-9]+)`)
	log       zerolog.Logger
	srv       *http.Server
	procGroup *sync.WaitGroup
)

type addquery struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type program struct{}

func formHandler(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
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

	buf := new(bytes.Buffer)
	cmd := exec.Command("/usr/bin/osascript",
		"-e", "tell application \"Things3\"",
		"-e", "set newToDo to make new to do with properties {name:\""+strings.ReplaceAll(q.Title, `"`, `\"`)+"\", notes:\""+strings.ReplaceAll(q.Content, `"`, `\"`)+"\"}",
		"-e", "end tell")
	cmd.Stdout = buf
	cmd.Stderr = buf
	if err := cmd.Run(); err != nil {
		log.Error().Err(err).Msg("Exec error")
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

func main() {
	log = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC1123Z}).Level(
		zerolog.InfoLevel).With().Timestamp().Logger()

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

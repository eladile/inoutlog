package server

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	telegram "inoutlog/telegram/client"
	"inoutlog/timelogger"
)

const timeLayout = "15:04"

var (
	inRegexp    = mustCompile("/in.*")
	outRegexp   = mustCompile("/out.*")
	aliveRegexp = mustCompile("/alive.*")
)

func mustCompile(s string) *regexp.Regexp {
	r, err := regexp.Compile(s)
	if err != nil {
		panic(err)
	}
	return r
}

type handler func(string, string) error

type Server struct {
	RegToHandler   map[*regexp.Regexp]handler
	TelegramClient telegram.Client
	TimeLogger     timelogger.Logger
}

func NewServer(telegram telegram.Client, logger timelogger.Logger) *Server {
	s := Server{
		TelegramClient: telegram,
		TimeLogger:     logger,
	}
	s.RegToHandler = map[*regexp.Regexp]handler{
		inRegexp:    s.handleIn,
		outRegexp:   s.handleOut,
		aliveRegexp: s.handleAlive,
	}
	return &s
}

func (s *Server) HandleUpdates(updates []telegram.Update) error {
	for _, update := range updates {
		if update.Message == nil || update.Message.Text == nil {
			continue
		}
		for regx, handler := range s.RegToHandler {
			if text := strings.ToLower(*update.Message.Text); regx.MatchString(text) {
				if err := handler(fmt.Sprintf("%d", update.Message.Chat.Id), text); err != nil {
					log.Printf("Can't handle text %s that matched %s, got error:%s\n",
						text, regx.String(), err.Error())
				}
			}
		}
	}
	return nil
}

func (s *Server) handleOut(id, text string) error {
	t, err := s.getDate(text, id)
	if err != nil {
		return err
	}
	pay, err := s.TimeLogger.Out(t)
	if err != nil {
		_ = s.TelegramClient.SendMessage(id, fmt.Sprintf("Failed to log out: %s", err.Error()))
		return err
	}
	return s.TelegramClient.SendMessage(id, fmt.Sprintf("Out time logged successfully: %s, pay %d", t.Format(timeLayout), pay))
}

func (s *Server) handleIn(id, text string) error {
	t, err := s.getDate(text, id)
	if err != nil {
		return err
	}
	err = s.TimeLogger.In(t)
	if err != nil {
		_ = s.TelegramClient.SendMessage(id, fmt.Sprintf("Failed to log in: %s", err.Error()))
		return err
	}
	return s.TelegramClient.SendMessage(id, fmt.Sprintf("In time logged successfully: %s", t.Format(timeLayout)))
}

// getDate demands a single word prefix!
func (s *Server) getDate(text string, id string) (time.Time, error) {
	words := strings.Fields(text)
	if len(words) <= 1 {
		// /in or /out only
		return time.Now(), nil
	}

	if len(words) > 2 {
		errText := fmt.Sprintf("too long text: %s. should be [HH:MM(optional)]", text)
		log.Println(errText)
		_ = s.TelegramClient.SendMessage(id, errText)
		return time.Time{}, errors.New(errText)
	}

	t, err := time.Parse(timeLayout, words[1])
	if err != nil {
		errText := fmt.Sprintf("date extract time from text should be [HH:MM(optional)]: %s", text)
		log.Println(errText)
		_ = s.TelegramClient.SendMessage(id, errText)
		return time.Time{}, errors.New(errText)
	}
	hour, min, _ := t.Clock()
	year, month, day := time.Now().Date()
	ret := time.Date(year, month, day, hour, min, 0, 0, time.Now().Location())
	return ret, nil
}

func (s *Server) handleAlive(id, _ string) error {
	return s.TelegramClient.SendMessage(id, "yeah I'm fine, thanks!")
}

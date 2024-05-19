package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"messagequeue/internal/config"
	"messagequeue/internal/logging"
	"messagequeue/internal/pkg/entities"
	"net/http"
	"net/url"
)

type Sender struct {
	SaveURL string
	Token   string
	logger  logging.ILogger
	client  *http.Client
}

func NewFactSender(config *config.AppConfig, logger logging.ILogger) *Sender {
	return &Sender{
		SaveURL: config.SaveURL,
		Token:   config.Token,
		client:  &http.Client{},
		logger:  logger,
	}
}

func (s *Sender) Process(facts []entities.Fact) {
	amount := len(facts)
	sem := make(chan struct{}, 1)
	requests := make(chan *http.Request, amount)
	for i, fact := range facts {
		go s.formRequest(fact, requests, i, amount)
	}
	for request := range requests {
		sem <- struct{}{}
		resp := entities.SaveResponse{}
		response, err := s.client.Do(request)
		if err != nil {
			s.logger.Error(err.Error())
		}
		err = json.NewDecoder(response.Body).Decode(&resp)
		if err != nil {
			s.logger.Error(err.Error())
		}
		s.logger.Debug(fmt.Sprintf("Response fact ID %d", resp.Data["indicator_to_mo_fact_id"]))
		<-sem
	}
}

func (s *Sender) formRequest(fact entities.Fact, out chan *http.Request, i, total int) {
	s.logger.Debug(fmt.Sprintf("Processing fact %d out of %d", i+1, total))
	params := url.Values{}
	params.Add("period_start", fact.PeriodStart)
	params.Add("period_end", fact.PeriodEnd)
	params.Add("period_key", fact.PeriodKey)
	params.Add("indicator_to_mo_id", fact.MoID)
	params.Add("indicator_to_mo_fact_id", fact.MoFactID)
	params.Add("value", fact.Value)
	params.Add("fact_time", fact.FactTime)
	params.Add("is_plan", fact.IsPlan)
	params.Add("auth_user_id", fact.AuthUserID)
	params.Add("comment", fact.Comment)
	payload := bytes.NewReader([]byte(params.Encode()))
	request, err := http.NewRequest("POST", s.SaveURL, payload)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error processing fact %d: %s", i+1, err.Error()))
	}
	out <- request
	s.logger.Debug(fmt.Sprintf("Done processing fact %d", i+1))
}

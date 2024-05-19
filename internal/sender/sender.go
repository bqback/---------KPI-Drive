package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"messagequeue/internal/config"
	"messagequeue/internal/logging"
	"messagequeue/internal/pkg/entities"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
)

type Sender struct {
	saveURL     string
	token       string
	logger      logging.ILogger
	client      *http.Client
	maxRequests int
}

func NewFactSender(config *config.AppConfig, logger logging.ILogger) *Sender {
	return &Sender{
		saveURL:     config.SaveURL,
		token:       config.Token,
		client:      &http.Client{},
		logger:      logger,
		maxRequests: config.MaxRequests,
	}
}

// Process преобразовывает список фактов в http-запросы и отправляет их на saveURL
func (s *Sender) Process(facts []entities.Fact) {
	amount := len(facts)

	requests := make(chan *http.Request, amount)
	// Очерёдность реализована простым семафором
	semaphore := make(chan struct{}, s.maxRequests)
	wg := &sync.WaitGroup{}

	wg.Add(amount)
	go s.doRequest(wg, requests, semaphore)

	// Факты преобразовываются в запросы асинхронно
	// Для 10 запросов это роли играть не будет, но для 10-20к разница может наблюдаться
	for i, fact := range facts {
		wg.Add(1)
		go s.formRequest(wg, fact, requests, i, amount)
	}
	wg.Wait()
	close(requests)
}

// doRequest вытаскивает запросы из канала, отправляет их на saveURL и читает ответ
func (s *Sender) doRequest(wg *sync.WaitGroup, requests chan *http.Request, sem chan struct{}) {
	for request := range requests {
		sem <- struct{}{}
		go func(request *http.Request) {
			defer wg.Done()

			resp := entities.SaveResponse{}
			response, err := s.client.Do(request)
			if err != nil {
				s.logger.Error("Error performing request: " + err.Error())
				<-sem
				return
			}
			err = json.NewDecoder(response.Body).Decode(&resp)
			if err != nil {
				s.logger.Error("Error decoding JSON: " + err.Error())
				<-sem
				return
			}
			if resp.Status == "OK" {
				s.logger.Debug(fmt.Sprintf("Response fact ID %d", resp.Data["indicator_to_mo_fact_id"]))
			} else {
				s.logger.Debug(fmt.Sprintf("Server returned error: %v", resp.Messages.Error))
			}

			<-sem
		}(request)
	}
	close(sem)
}

// formRequest преобразовывает entities.Fact в http.Request и кладёт его в канал для doRequest
func (s *Sender) formRequest(wg *sync.WaitGroup, fact entities.Fact, out chan *http.Request, i, total int) {
	defer wg.Done()

	s.logger.Debug(fmt.Sprintf("Processing fact %d out of %d", i+1, total))
	var err error
	values := map[string]io.Reader{
		"period_start":            strings.NewReader(fact.PeriodStart),
		"period_end":              strings.NewReader(fact.PeriodEnd),
		"period_key":              strings.NewReader(fact.PeriodKey),
		"indicator_to_mo_id":      strings.NewReader(fact.MoID),
		"indicator_to_mo_fact_id": strings.NewReader(fact.MoFactID),
		"value":                   strings.NewReader(fact.Value),
		"fact_time":               strings.NewReader(fact.FactTime),
		"is_plan":                 strings.NewReader(fact.IsPlan),
		"auth_user_id":            strings.NewReader(fact.AuthUserID),
		"comment":                 strings.NewReader(fact.Comment),
	}
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for k, v := range values {
		var fw io.Writer
		if x, ok := v.(io.Closer); ok {
			defer x.Close()
		}
		if fw, err = w.CreateFormField(k); err != nil {
			return
		}
		if _, err = io.Copy(fw, v); err != nil {
			return
		}
	}
	w.Close()

	request, err := http.NewRequest("POST", s.saveURL, &body)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error processing fact %d: %s", i+1, err.Error()))
	}
	// При заголовке "application/x-www-form-urlencoded" сервер отдаёт 500, в Postman та же проблема
	// Запрос при этом проходит и факт отдаётся get_facts
	request.Header.Add("Content-Type", w.FormDataContentType())
	request.Header.Add("Authorization", "Bearer "+s.token)
	out <- request
	s.logger.Debug(fmt.Sprintf("Done processing fact %d", i+1))
}

package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jafarsirojov/bank-cards/pkg/core/auth"
	"github.com/jafarsirojov/bank-cards/pkg/core/cards"
	"github.com/jafarsirojov/mux/pkg/mux"
	"github.com/jafarsirojov/mux/pkg/mux/middleware/jwt"
	"github.com/jafarsirojov/rest/pkg/rest"
	"log"
	"net/http"
	"strconv"
	"time"
)

const addrHistorySvc = "http://localhost:9010"

type MainServer struct {
	exactMux *mux.ExactMux
	cardsSvc *cards.Service
}

func NewMainServer(exactMux *mux.ExactMux, cardsSvc *cards.Service) *MainServer {
	return &MainServer{exactMux: exactMux, cardsSvc: cardsSvc}
}

func (m *MainServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// delegation
	m.exactMux.ServeHTTP(writer, request)
}

func (m *MainServer) HandleGetAllCards(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwt.FromContext(request.Context()).(*auth.Auth)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	log.Print(authentication)
	if authentication.Id == 0 {
		log.Print("admin")
		log.Print("all cards")
		models, err := m.cardsSvc.All()
		if err != nil {
			log.Print("can't get all cards")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Print(models)
		writer.Header().Set("Content-Type", "application/json")
		err = rest.WriteJSONBody(writer, models)
		if err != nil {
			log.Print("can't write json get all cards")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	id := authentication.Id
	log.Printf("user by id: %d", id)
	models, err := m.cardsSvc.ViewCardsByOwnerId(id)
	if err != nil {
		log.Print("can't get all cards")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print(models)
	err = rest.WriteJSONBody(writer, models)
	if err != nil {
		log.Print("can't write json get all cards")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func (m *MainServer) HandleGetCardById(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwt.FromContext(request.Context()).(*auth.Auth)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	log.Print(authentication)
	log.Print("cards by id")
	value, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		log.Print("can't get all cards")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		log.Print("can't strconv atoi to show card by id")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if authentication.Id == 0 {
		models, err := m.cardsSvc.ById(id)
		if err != nil {
			log.Print("can't get all cards")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Print(models)
		err = rest.WriteJSONBody(writer, models)
		if err != nil {
			log.Print("can't write json get all cards")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	models, err := m.cardsSvc.ByIdUserCard(id, authentication.Id)
	if err != nil {
		log.Print("can't get all cards")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Print(models)
	err = rest.WriteJSONBody(writer, models)
	if err != nil {
		log.Print("can't write json get all cards")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (m *MainServer) HandlePostCard(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwt.FromContext(request.Context()).(*auth.Auth)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	log.Print(authentication)
	if authentication.Id != 0 {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't is not admin post cards")
		return
	}
	log.Print("post card")
	model := cards.Cards{}
	err := rest.ReadJSONBody(request, &model)
	if err != nil {
		log.Printf("can't READ json POST model: %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(model)
	if model.Id != 0 {
		log.Print("bad request")
		log.Print("post body id not 0!!!")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = m.cardsSvc.AddCard(model)
	if err != nil {
		log.Printf("can't add card %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (m *MainServer) HandleBlockById(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwt.FromContext(request.Context()).(*auth.Auth)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	log.Print(authentication)
	model := cards.ModelBlockCard{}
	err := rest.ReadJSONBody(request, &model)
	if err != nil {
		log.Printf("can't read json body: %d", err)
		return
	}
	if authentication.Id == 0 {
		err = m.cardsSvc.BlockById(model.Id)
		if err != nil {
			log.Print("can't to blocked card by id")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err = m.cardsSvc.UserBlockCardById(authentication.Id, model)
	if err != nil {
		log.Printf("can't user block card by id: %d", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (m *MainServer) HandleUnBlockById(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwt.FromContext(request.Context()).(*auth.Auth)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	log.Print(authentication)
	if authentication.Id != 0 {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't is not admin unblock card")
		return
	}
	model := cards.ModelBlockCard{}
	err := rest.ReadJSONBody(request, &model)
	if err != nil {
		log.Printf("can't read json body: %d", err)
		return
	}
	err = m.cardsSvc.UnBlockedById(model.Id)
	if err != nil {
		log.Print("can't to blocked card by id")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (m *MainServer) HandleTransferMoneyCardToCard(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwt.FromContext(request.Context()).(*auth.Auth)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	log.Print(authentication)

	transfer := cards.ModelTransferMoneyCardToCard{}
	err := rest.ReadJSONBody(request, &transfer)
	if err != nil {
		log.Printf("can't READ json transfer money: %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(transfer)
	err, senderData, recipientData := m.cardsSvc.TransferMoneyCardToCard(authentication.Id, transfer)
	if err != nil {
		log.Printf("can't transfer money: %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print("transfer ok to handler")

	token:= request.Header.Get("Authorization")
	log.Printf("Bearer token %s",token)
	token = token[7:]
	log.Printf("token %s",token)

	err = saveHistoryToSvcHistory(senderData,token)
	if err != nil {
		log.Printf("can't save to db svc history tranth sender: %s", err)
		//save with local db
	}
	err = saveHistoryToSvcHistory(recipientData,token)
	if err != nil {
		log.Printf("can't save to db svc history tranth recipient: %s", err)
		//save with local db
	}
	value := request.Context()
	log.Printf("value: %s", value)

}

func saveHistoryToSvcHistory(data cards.ModelOperationsLog, token string) (err error) {
	log.Print("starting sender request to history Svc")
	requestBody, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("can't encode requestBody %v: %w", data, err)
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/history", addrHistorySvc),
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return fmt.Errorf("can't create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s",token))
	log.Print("started sender request to history Svc")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("can't send request: %w", err)
	}
	defer response.Body.Close()

	log.Print("finish sender request to history Svc")

	switch response.StatusCode {
	case 200:
		log.Print("200 request to history Svc")
		return nil
	case 400:
		log.Print("400 request to history Svc")
		return fmt.Errorf("bad request is server: %s", addrHistorySvc)
	case 401:
		log.Print("401 unauthorized to history Svc")
		return fmt.Errorf("unauthorized is server: %s", addrHistorySvc)
	case 500:
		log.Print("500 request to history Svc")
		return fmt.Errorf("internel server error is server: %s", addrHistorySvc)
	default:
		log.Printf("response status code: %s", err)
		return fmt.Errorf("err: %s", addrHistorySvc)
	}
}

package app

import (
	"bank-cards/pkg/core/auth"
	"bank-cards/pkg/core/cards"
	"github.com/jafarsirojov/mux/pkg/mux"
	"github.com/jafarsirojov/mux/pkg/mux/middleware/jwt"
	"github.com/jafarsirojov/rest/pkg/rest"
	"log"
	"net/http"
	"strconv"
)

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
	if authentication.Id != 0 {
		log.Printf("can't authentication is not admin, this id user = %d", authentication.Id)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Print("all cards")
	models, err := m.cardsSvc.All()
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
	log.Print("user by id")
	value, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		log.Print("can't get all cards")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(value)
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
	log.Print("post user")
	model := cards.Cards{}
	err := rest.ReadJSONBody(request, &model)
	if err != nil {
		log.Printf("can't READ json POST model: %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(model)
	m.cardsSvc.AddCard(model)

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
	value, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		log.Print("can't read id FromContext")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		log.Print("can't strconv atoi to blocked by id card")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = m.cardsSvc.BlockById(id)
	if err != nil {
		log.Print("can't to blocked card by id")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (m *MainServer) HandleUnBlockedById(writer http.ResponseWriter, request *http.Request) {
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
	value, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		log.Print("can't read id FromContext")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		log.Print("can't strconv atoi to blocked by id card")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = m.cardsSvc.UnBlockedById(id)
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
	value, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		log.Print("can't read id FromContext")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		log.Print("can't strconv atoi to blocked by id card")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	transfer := cards.ModelTransferMoneyCardToCard{}
	err = rest.ReadJSONBody(request, &transfer)
	if err != nil {
		log.Printf("can't READ json transfer money: %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(transfer)
	err = m.cardsSvc.TransferMoneyCardToCard(id /*transfer.IdCardSender*/, transfer.NumberCardRecipient, transfer.Count)
	log.Print("transfer ok to handler")

}

func (m *MainServer) HandleGetCardsByOwnerId(writer http.ResponseWriter, request *http.Request) {
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
	log.Print("user by id")
	id := authentication.Id
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
}

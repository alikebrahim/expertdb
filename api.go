package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/account", makeHTTPHandleFunc(s.HandleAccount))
	mux.HandleFunc("/account/", makeHTTPHandleFunc(s.HandleGetAccountByID))

	log.Println("Server running on port: ", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, mux))
}

// Handlers
func (s *APIServer) HandleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.HandleGetAccount(w, r)
	case "POST":
		return s.HandleCreateAccount(w, r)
	case "PUT":
		return s.HandleEditAccount(w, r)
	case "DELETE":
		return s.HandleDeleteAccount(w, r)
	default:
		return WriteJson(w, http.StatusMethodNotAllowed, ApiError{Error: fmt.Sprintf("error: %v, method %v not supported", http.StatusMethodNotAllowed, r.Method)})

	}
	return nil
}

func getID(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func (s *APIServer) HandleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, accounts)
}
func (s *APIServer) HandleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id := getID(r)
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal("error converting id to int", err)
	}
	expert := &Account{ID: idInt}
	fmt.Println(id)
	return WriteJson(w, http.StatusOK, expert)
}

func (s *APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accReq := new(CreateAccountRequest) // CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(accReq); err != nil {
		return err
	}
	account := NewAccount(accReq.Prefix,
		accReq.Name,
		accReq.Affiliation,
		accReq.BH,
		accReq.Available,
		accReq.Rating,
		accReq.Role,
		accReq.Type,
		accReq.GeneralArea,
		accReq.SpecialisedArea,
		accReq.Trained,
		accReq.PrimaryContact,
		accReq.SecondaryContact,
		accReq.Email,
		accReq.Published)
	fmt.Println("About to create account")
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, account)
}

func (s *APIServer) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) HandleEditAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

// json
func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// server API
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpClient *http.Client
	router     chi.Router
}

func NewServer() *Server {
	s := &Server{
		httpClient: &http.Client{},
	}
	r := chi.NewRouter()
	s.router = r

	// Routes
	r.Use(s.withPolicyAuthorization())
	r.Get("/", s.handleHome())
	r.Get("/health", s.handleHealth())
	r.Get("/payments/{user}", s.handleListPaymentsForUser())

	return s
}

func (s *Server) withPolicyAuthorization() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowed, err := s.checkPolicy(r)
			if err != nil {
				log.Printf("error: check policy: %v", err)
				http.Error(w, "Policy check failed", http.StatusInternalServerError)
				return
			}
			if !allowed {
				http.Error(w, "Not allowed", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type policyInput struct {
	Method string   `json:"method"`
	Path   []string `json:"path"`
	User   string   `json:"user"`
}

type policyOutput struct {
	Result bool `json:"result"`
}

func (s *Server) checkPolicy(r *http.Request) (bool, error) {
	reqPath := strings.Split(r.URL.Path, "/")
	input := map[string]any{
		"input": policyInput{
			Method: r.Method,
			Path:   reqPath[1:],
			User:   r.Header.Get("X-User"),
		},
	}

	data := &bytes.Buffer{}
	if err := json.NewEncoder(io.MultiWriter(data, os.Stderr)).Encode(input); err != nil {
		return false, err
	}
	opaURL := "http://localhost:8081/v1/data/demo/api/allow"
	opaReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, opaURL, data)
	if err != nil {
		return false, err
	}
	opaResp, err := s.httpClient.Do(opaReq)
	if err != nil {
		return false, err
	}
	body := opaResp.Body
	defer func() {
		_, _ = io.ReadAll(body)
		body.Close()
	}()
	if opaResp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("received status %d from opa", opaResp.StatusCode)
	}

	var result policyOutput
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return false, err
	}
	return result.Result, nil
}

func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "healthy")
	}
}

func (s *Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "home")
	}
}

func (s *Server) handleListPaymentsForUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "payments for user %q\n", chi.URLParam(r, "user"))
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

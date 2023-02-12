package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"server/auth"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func routes() *mux.Router {
	r := mux.NewRouter()
	// r.Use(middleware.ExternalOriginMiddleware)
	// r.Use(middleware.AddResponseHeaders)
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/signup", signup).Methods("POST")
	r.HandleFunc("/forgotpassword", forgotPassword).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(invalidEndpoint)
	return r
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("endpoint %v with method %v\n", r.URL.Path, r.Method)
	// email/phone, password
	// loggedin: bool, jwt token
	rawResp, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = auth.VerifyCredentials(rawResp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("endpoint %v with method %v\n", r.URL.Path, r.Method)
	// name, email, password, phone
	// loggedin: bool, jwt token
	rawResp, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = auth.AddNewUser(rawResp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func forgotPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("endpoint %v with method %v\n", r.URL.Path, r.Method)
	// email/phone
	// resp true/false
	// await otp
	rawResp, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = auth.ResetPassword(rawResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	// err = auth.SendOTP(rawResp)
	// if err != nil {
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
}

func invalidEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Endpoint Hit: %v with %v method\n", r.URL.Path, r.Method)
	http.Error(w, "Endpoint does not exist", http.StatusNotFound)
	w.WriteHeader(http.StatusNoContent)
}

func handleRequests() {
	router := routes()

	credentials := handlers.AllowCredentials()
	origins := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"POST"})
	//ttl := handlers.MaxAge(3600)

	fmt.Println("Starting server on port: ", 10500)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%v", 10500), handlers.CORS(credentials, headers, methods, origins)(router)))
}

func StartServer() {
	handleRequests()
}

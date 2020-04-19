package test

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/julienschmidt/httprouter"

	"golang.org/x/oauth2"
)

// OAUTHGoogle .
func OAUTHGoogle(router *httprouter.Router) {
	clientID, _ := os.LookupEnv("CLIENT_ID_GOOGLE")
	clientSecret, _ := os.LookupEnv("CLIENT_SECRET_GOOGLE")
	// httpClient := &http.Client{}

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/drive.metadata.readonly", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
		RedirectURL: "http://localhost:8080/api/google/oauth",
	}

	ctx := context.Background()

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth: %s\n", url)
	router.GET("/api/google/oauth", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Println("HIT API/GOOGLE/OAUTH")

		// fmt.Println(r.URL.Query().Get("client_id"))
		tok, err := conf.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(tok)
		w.Write([]byte(tok.AccessToken))
		// client := conf.Client(ctx, tok)
		// client.Get("https://openidconnect.googleapis.com/v1/userinfo")
		// todo: request to https://oauth2.googleapis.com/tokeninfo?access_token=token

		// fmt.Print("m: ")
	})

	// router.POST("/api/google/token", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 	fmt.Println("HIT API/GOOGLE/TOKEN")
	// 	fmt.Println(r.Header.Get("Authorization"))
	// 	if r.Header.Get("Authorization") == "" {
	// 		w.Write([]byte("x"))
	// 		return
	// 	}
	// 	// decoder := json.NewDecoder(r.Body)
	// 	// b := map[string]interface{}{}
	// 	// err := decoder.Decode(&b)
	// 	// fmt.Println(b)
	// 	// if err != nil {
	// 	// 	log.Fatal(err)
	// 	// }
	// 	// e, _ := json.Marshal(b)

	// 	w.Write([]byte(`{ "access_token": "` + strings.Split(r.Header.Get("Authorization"), " ")[1] + `"}`))
	// })

	router.GET("/google", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Println("HIT /GOOGLE")
		w.Write([]byte("{}"))
	})
}

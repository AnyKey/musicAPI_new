package handlers

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/streadway/amqp"
	"log"
	"musicAPI/model"
	"musicAPI/repository"
	"net/http"
	"strings"
	"time"
)

type TokenHandler struct {
	Repo  repository.Repository
	Chann *amqp.Channel
}

func (th TokenHandler) AuthUser(next http.Handler) http.Handler {
	q := newQueue(th.Chann)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/api/refresh") || strings.HasPrefix(r.RequestURI, "/api/login") || strings.HasPrefix(r.RequestURI, "/index") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		token := r.Header.Get("token")
		if token != "" {
			sd := th.CheckToken(r.Context(), token, q, r)
			if sd == true {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
		_ = WriteJsonToResponse(w, "Unauthorized")
	})

}

func NewTokenA(user string) (*jwt.Token, error) {
	return jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		jwt.MapClaims{
			"name": user,
			"exp":  time.Now().Add(time.Hour * 1).Unix(),
			"root": true,
		}), nil
}
func NewTokenR(user string) (*jwt.Token, error) {
	return jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		jwt.MapClaims{
			"name": user,
			"exp":  time.Now().Add(time.Hour * 3).Unix(),
		}), nil
}

func (th TokenHandler) CheckToken(ctx context.Context, myToken string, queue amqp.Queue, r *http.Request) bool {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(myToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})

	if err == nil && token.Valid {
		user := (claims["name"]).(string)
		res := th.Repo.GetTokens(ctx, user)
		if res.Access == myToken && res.Valid == true {
			body := model.LogBody{
				Name:   user,
				Action: r.RequestURI,
				Time:   time.Now(),
			}
			bytes, err := json.Marshal(body)
			if err != nil {
				log.Println(err)
			}
			th.pushToChan(bytes, queue)
			return true
		}
	}
	return false
}
func newQueue(ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(
		"main_queue", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Println(err)
	}
	return q
}
func (th TokenHandler) pushToChan(body []byte, q amqp.Queue) error {
	err := th.Chann.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Println(err)
	}
	return nil
}

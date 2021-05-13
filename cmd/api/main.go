package main

import (
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
	"musicAPI/client"
	"musicAPI/elastic"
	elasticM "musicAPI/elastic/delivery"
	elasticRep "musicAPI/elastic/repository/elastic"
	elasticUseCase "musicAPI/elastic/usecase"
	"musicAPI/logs"
	logsM "musicAPI/logs/delivery"
	rabbitMqLogs "musicAPI/logs/repository/rabbitmq"
	logsUseCase "musicAPI/logs/usecase"
	"musicAPI/music"
	"musicAPI/music/delivery/api"
	musicH "musicAPI/music/delivery/http"
	esMusicRep "musicAPI/music/repository/elastic"
	dbMusicRep "musicAPI/music/repository/postgres"
	redisMusicRep "musicAPI/music/repository/redis"
	musicUseCase "musicAPI/music/usecase"
	"musicAPI/user"
	userM "musicAPI/user/delivery"
	userH "musicAPI/user/delivery/http"
	redisUserRep "musicAPI/user/repository/redis"
	userUseCase "musicAPI/user/usecase"
	"net/http"
	"time"
)

type App struct {
	httpAddress string
	userUC      user.UseCase
	musicUC     music.UseCase
	logsUC      logs.UseCase
	elasticUC   elastic.UseCase
}

func NewApp() *App {
	var reg = NewReg()
	userDB := redisUserRep.New(reg.rConn)
	musicRedis := redisMusicRep.New(reg.rConn)
	musicDB := dbMusicRep.New(reg.dbConn)
	musicA := api.New()
	musicEs := esMusicRep.New(reg.esConn)
	logsRab := rabbitMqLogs.New(reg.qConn)
	elasticEs := elasticRep.New(reg.esConn)
	return &App{
		httpAddress: reg.address,
		userUC: userUseCase.New(
			userDB,
		),
		musicUC: musicUseCase.New(
			musicRedis,
			musicDB,
			musicA,
			musicEs,
		),
		logsUC: logsUseCase.New(
			logsRab,
		),
		elasticUC: elasticUseCase.New(
			elasticEs,
		),
	}
}
func (a *App) Run() {

	router := mux.NewRouter()
	router.Use(userM.NewUserHandler(a.userUC).UserMiddleware)
	router.Use(logsM.NewLogHandler(a.logsUC).LogMiddleware)
	router.Use(elasticM.NewTrackHandler(a.elasticUC).WsHandler)
	router.Use(mux.CORSMethodMiddleware(router))
	userH.UserHandlers(router, a.userUC)
	musicH.MusicHandlers(router, a.musicUC)
	client.Template(router)

	srv := &http.Server{
		Handler:      router,
		Addr:         a.httpAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Serve http ON", a.httpAddress)
	log.Fatal(srv.ListenAndServe())

}

func main() {
	app := NewApp()
	app.Run()
}

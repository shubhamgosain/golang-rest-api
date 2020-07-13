package routes

import (
	"encoding/json"
	"fmt"
	"golang-rest-api/dboperations"
	customlogger "golang-rest-api/logger"
	"golang-rest-api/middlewares/prometheus"
	"golang-rest-api/middlewares/reqid"
	"golang-rest-api/response"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type (
	configPath struct {
		path string
	}
	dbConfig dboperations.DbConfig
)

var (
	logger *logrus.Logger
	logRID *logrus.Logger
)

func readData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	render.Render(w, r, response.SuccessResponse(dboperations.ReadRecords()))
}

func insertData(w http.ResponseWriter, r *http.Request) {
	var data dboperations.InputData
	err := render.DecodeJSON(r.Body, &data)
	if err != nil {
		log.Println(err)
		render.Render(w, r, response.ErrBadRequest(err))
		return
	}
	if err = dboperations.AddRecord(data); err != nil {
		render.Render(w, r, response.ErrBadRequest(err))
		customlogger.GetLoggerWithRID(logger, r).Infoln(err)
	} else {
		render.Render(w, r, response.SuccessResponse("Successfully Added record"))
	}
}

func deleteData(w http.ResponseWriter, r *http.Request) {
	var data dboperations.InputData
	err := render.DecodeJSON(r.Body, &data)
	if err != nil {
		render.Render(w, r, response.ErrBadRequest(err))
		return
	}
	if err := dboperations.DeleteRecord(data); err != nil {
		render.Render(w, r, response.ErrBadRequest(err))
	} else {
		render.Render(w, r, response.SuccessResponse("Succesfully Deleted Record"))
	}
}

func (config configPath) readConfig() (dbconfig dbConfig) {
	rawConfig, err := ioutil.ReadFile(string(config.path))
	if err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(rawConfig, &dbconfig); err != nil {
		log.Fatal(err)
	}
	return
}

func (dbconfig dbConfig) initDB() {
	dboperations.CreateConnection(dbconfig.PostgresDB)
}

//RestHandler A handler to handle all the rest api. This is holding a chi router which is been developed to handle GET,POST and DELETE methods
func RestHandler(configFile string) {
	configpath := &configPath{configFile}
	config := configpath.readConfig()
	config.initDB()
	r := chi.NewRouter()
	logger = logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	r.Route("/orders", func(r chi.Router) {
		r.Use(reqid.RequestID)
		r.Use(customlogger.Logger(logger))
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(60 * time.Second))
		r.Use(prometheus.Metrics)
		r.Get("/", readData)
		r.Post("/", insertData)
		r.Delete("/", deleteData)
	})
	r.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%v", config.App.Port), r)
}

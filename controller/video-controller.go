package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"company/system/microservices/entity"
	"company/system/microservices/service"
	customValidators "company/system/microservices/validators"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	kafka "github.com/segmentio/kafka-go"
)

type VideoController interface {
	FindAll() []entity.Video
	Save(ctx *gin.Context) error
	ShowAll(ctx *gin.Context)
	Validate(ctx *gin.Context) string
}

type videoController struct {
	service service.VideoService
}

var validate *validator.Validate

const (
	searchDocumentEvented = "searchDocumentEvented"
	foundDocumentEvented  = "foundDocumentEvented"
)

func NewVideoController(service service.VideoService) VideoController {
	validate = validator.New()
	validate.RegisterValidation("is-gmail", customValidators.ValidateIsGmail)
	return &videoController{
		service: service,
	}
}

func (controller *videoController) FindAll() []entity.Video {
	return controller.service.FindAll()
}

func (controller *videoController) Save(context *gin.Context) error {
	var video entity.Video
	err := context.ShouldBindJSON(&video)
	if err != nil {
		return err
	}
	err = validate.Struct(video)
	if err != nil {
		return err
	}
	controller.service.Save(video)
	return nil
}

func (controller *videoController) ShowAll(context *gin.Context) {
	videos := controller.service.FindAll()
	data := gin.H{
		"title":  "Video Page",
		"videos": videos,
	}
	context.HTML(http.StatusOK, "index.html", data)
}

func (controller *videoController) Validate(context *gin.Context) string {
	var client entity.Client
	err := context.ShouldBindJSON(&client)
	if err != nil {
		return err.Error()
	}

	body, _ := json.Marshal(client)

	kafkaURL := "localhost:9092"

	kafkaWriter := getKafkaWriter(kafkaURL, searchDocumentEvented)
	defer kafkaWriter.Close()
	msg := kafka.Message{
		Key:   []byte("A"),
		Value: body,
	}

	err = kafkaWriter.WriteMessages(context, msg)
	if err != nil {
		fmt.Println(err)
	} else {
		log.Println("produced - ", string(msg.Value))
	}

	response := <-consumeFoundDocumentEvented(context, client.Document)

	return response
}

func consumeFoundDocumentEvented(context *gin.Context, document string) <-chan string {
	response := make(chan string)

	log.Println("Antes de Go Func")
	go func() {
		defer close(response)

		kafkaURL := "localhost:9092"
		log.Println("Antes de Iniciar Reader")
		reader := getKafkaReader(kafkaURL, foundDocumentEvented, "logger-group")
		log.Println("Despues de Iniciar Reader")
		foundData := ""

		for {
			log.Println("Antes de Leer Reader")
			m, err := reader.ReadMessage(context)
			log.Println("Despues de Leer Reader")
			if err != nil {
				fmt.Println(err)
			}
			foundData = string(m.Value)
			log.Println("consumed at topic:", m.Topic, " - ", foundData)

			if strings.Contains(foundData, document) {
				break
			}
		}

		defer reader.Close()

		response <- foundData
	}()

	return response
}

func longRunningTask() <-chan int32 {
	r := make(chan int32)

	go func() {
		defer close(r)

		// simulate a workload
		time.Sleep(time.Second * 3)
		r <- rand.Int31n(100)
	}()

	return r
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

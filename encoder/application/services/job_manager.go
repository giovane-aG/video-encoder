package services

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/giovane-aG/video-encoder/encoder/application/repositories"
	"github.com/giovane-aG/video-encoder/encoder/domain"
	"github.com/giovane-aG/video-encoder/encoder/infrastructure/queue"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type JobManager struct {
	Db               *gorm.DB
	Job              domain.Job
	MessageChannel   chan amqp.Delivery
	JobReturnChannel chan JobWorkerResult
	RabbitMQ         *queue.RabbitMQ
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManager(
	db *gorm.DB,
	messageChannel chan amqp.Delivery,
	jobReturnChannel chan JobWorkerResult,
	rabbitMQ *queue.RabbitMQ,
) *JobManager {
	return &JobManager{
		Db:               db,
		Job:              domain.Job{},
		MessageChannel:   messageChannel,
		JobReturnChannel: jobReturnChannel,
		RabbitMQ:         rabbitMQ,
	}
}

func (j *JobManager) Start(ch *amqp.Channel) {
	videoService := NewVideoService()
	videoService.VideoRepository = repositories.NewVideoRepository(j.Db)

	jobService := JobService{
		JobRepository: *repositories.NewJobRepository(j.Db),
		VideoService:  videoService,
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))

	if err != nil {
		log.Fatalf("error loading var: CONCURRENCY_WORKERS")
	}

	for numProcesses := 0; numProcesses < concurrency; numProcesses++ {
		go JobWorker(j.MessageChannel, j.JobReturnChannel, jobService, j.Job, numProcesses)
	}

	for jobResult := range j.JobReturnChannel {
		if jobResult.Error != nil {
			err = j.checkParseErrors(jobResult)
		} else {
			err = j.notifySuccess(jobResult, ch)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}
}

func (j *JobManager) checkParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID #{jobResult.Message.DeliveryTag}. Error parsing job:  #{jobResult.Job.ID}")
	} else {
		log.Printf("MessageID #{jobResult.Message.DeliveryTag}. Error parsing message: #{jobResult.Error}")
	}

	errorMsg := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error:   jobResult.Job.Error,
	}

	jobJson, err := json.Marshal(errorMsg)
	if err != nil {
		return err
	}

	err = j.notify(jobJson)
	if err != nil {
		return err
	}

	err = jobResult.Message.Reject(false)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobManager) notify(jobJson []byte) error {
	err := j.RabbitMQ.Notify(
		string(jobJson),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)

	if err != nil {
		return err
	}

	return nil
}

func (j *JobManager) notifySuccess(jobResult JobWorkerResult, ch *amqp.Channel) error {
	jobJson, err := json.Marshal(jobResult.Job)
	if err != nil {
		return err
	}

	err = j.notify(jobJson)
	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)
	if err != nil {
		return err
	}

	return nil
}

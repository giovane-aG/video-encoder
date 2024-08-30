package services

import (
	"encoding/json"
	"os"
	"time"

	"github.com/giovane-aG/video-encoder/encoder/domain"
	"github.com/giovane-aG/video-encoder/encoder/infrastructure/utils"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(
	messageChannel chan amqp.Delivery,
	returnChan chan JobWorkerResult,
	jobService JobService,
	job domain.Job,
	workerID int,
) {
	for message := range messageChannel {

		if err := utils.IsJson(string(message.Body)); err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err := json.Unmarshal(message.Body, &jobService.VideoService.Video)
		jobService.VideoService.Video.ID = uuid.NewString()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoService.Video.Validate()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoService.InsertVideo()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		job.Video = jobService.VideoService.Video
		job.OutputBucketPath = os.Getenv("outputBucketPath")
		job.ID = uuid.NewString()
		job.Status = "STARTING"
		job.CreatedAt = time.Now()

		_, err = jobService.JobRepository.Insert(&job)

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.Job = &job

		err = jobService.Start()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		returnChan <- returnJobResult(job, message, nil)
	}
}

func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWorkerResult {
	return JobWorkerResult{
		Job:     job,
		Message: &message,
		Error:   err,
	}
}

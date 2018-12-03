package logger

import (
	"github.com/gocraft/health"
	"os"
)

var databaseHealthStream *health.Stream
var jobsHealthStream *health.Stream

func DatabaseStream() *health.Stream {
	if databaseHealthStream == nil {
		databaseHealthStream = health.NewStream()
		databaseHealthStream.AddSink(&health.WriterSink{os.Stdout})
	}
	return databaseHealthStream
}

func JobsStream() *health.Stream {
	if jobsHealthStream == nil {
		jobsHealthStream = health.NewStream()
		jobsHealthStream.AddSink(&health.WriterSink{os.Stdout})
	}
	return jobsHealthStream
}

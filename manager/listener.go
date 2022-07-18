package manager

import "log"

//will use for websocket based notification delivery later
type ConversionListener interface {
	OnConversionStarted(string)
	OnConversionCompleted(string)
	OnConversionFailed(string, error)
}

type LoggerListener struct {
}

func (l *LoggerListener) OnConversionStarted(id string) {
	log.Printf("[OnConversionStarted]: %s\n", id)
}
func (l *LoggerListener) OnConversionCompleted(id string) {
	log.Printf("[OnConversionCompleted]: %s\n", id)
}
func (l *LoggerListener) OnConversionFailed(id string, err error) {
	log.Printf("[OnConversionFailed]: %s | %s\n", id, err.Error())
}

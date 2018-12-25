package logschan

import "fmt"

type Logger interface {
	AddLog(params ...string) error
	TakeLogs(params ...string) (string, error)
}
type Log struct {
	Logger
	cAddLog chan string
}

func (l *Log) AddLog(params ...string) error {
	fmt.Println(params)
	return nil
}
func (l *Log) TakeLogs(params ...string) (string, error) {
	return "", nil
}
func init() {
	l := Log{}
	go ListenPortLog(l)
}
func ListenPortLog(l Log) {
	for {
		select {
		case s := <-l.cAddLog:
			l.AddLog(s)
			if s == "Exit Log" {
				return
			}
		}
	}
}

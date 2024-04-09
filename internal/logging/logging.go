package logging

type Logger interface {
	Infof(template string, args ...any)
	Info(args ...any)
	Debugf(template string, args ...any)
	Debug(args ...any)
	Errorf(template string, args ...any)
	Error(args ...any)
	Fatalf(template string, args ...any)
	Fatal(args ...any)
}

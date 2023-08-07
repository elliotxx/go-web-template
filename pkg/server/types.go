package server

type Server interface {
	PreRun() error
	Run() error
}

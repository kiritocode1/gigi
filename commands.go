package main

type Gigi interface {
	Init(path string) error
	Add(path string) error
	Commit(message string) error
	Push(remote string) error
	Pull(remote string) error
	Clone(remote string) error
	Log() error
}

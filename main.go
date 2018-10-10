package main

import (
	"github.com/nbaztec/npm-beacon/handler"
)

func main() {
	conf := handler.LoadConfiguration()
	handler.Process(conf.Repositories, conf.GithubToken, conf.MinDaysNewRelease)
}

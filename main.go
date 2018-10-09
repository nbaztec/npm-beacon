package main

import (
	"github.com/nbaztec/npm-beacon/config"
	"github.com/nbaztec/npm-beacon/handler"
)

func main() {
	conf := config.Load()
	handler.Process(conf.Repositories, conf.GithubToken)
}

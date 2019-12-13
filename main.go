// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"net/http"
	"strings"

	"drone-env-merge/plugin"
	"github.com/drone/drone-go/plugin/environ"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type spec struct {
	Bind       string `envconfig:"DRONE_BIND"`
	Debug      bool   `envconfig:"DRONE_DEBUG"`
	Secret     string `envconfig:"DRONE_SECRET"`
	SkipVerify bool   `envconfig:"DRONE_SKIPVERIFY"`
	Upstreams  string `envconfig:"DRONE_UPSTREAMS"`
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Bind == "" {
		spec.Bind = ":80"
	}

	handler := environ.Handler(
		spec.Secret,
		plugin.New(
			spec.Secret,
			strings.Split(spec.Upstreams, ","),
			spec.SkipVerify,
		),
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}

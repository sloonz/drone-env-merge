// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"sync"

	"github.com/drone/drone-go/plugin/environ"
	"github.com/sirupsen/logrus"
)

func New(secret string, upstreams []string, skipVerify bool) environ.Plugin {
	logrus.Debugf("upstreams: %v", upstreams)
	return &plugin{
		secret:     secret,
		upstreams:  upstreams,
		skipVerify: skipVerify,
	}
}

type plugin struct {
	secret     string
	upstreams  []string
	skipVerify bool
}

func (p *plugin) List(ctx context.Context, req *environ.Request) ([]*environ.Variable, error) {
	var wg sync.WaitGroup
	var m sync.Mutex
	var env []*environ.Variable

	for _, upstream := range p.upstreams {
		wg.Add(1)
		go (func() {
			defer wg.Done()

			client := environ.Client(upstream, p.secret, p.skipVerify)
			upstreamEnv, err := client.List(ctx, req)
			if err != nil {
				logrus.Warning(err)
				return
			}

			logrus.Debugf("%s data: %v", upstream, env)

			m.Lock()
			env = append(env, upstreamEnv...)
			m.Unlock()
		})()
	}

	wg.Wait()

	logrus.Debugf("merged: %v", env)

	return env, nil
}

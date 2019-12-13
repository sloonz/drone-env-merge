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
		secret: secret,
		upstreams: upstreams,
		skipVerify: skipVerify,
	}
}

type plugin struct {
	secret string
	upstreams []string
	skipVerify bool
}

func unwrap(syncMap *sync.Map) map[string]string {
	m := make(map[string]string)
	syncMap.Range(func (key, val interface {}) bool {
		m[key.(string)] = val.(string)
		return true
	})
	return m
}

func (p *plugin) fetchUpstreamEnviron(upstream string, ctx context.Context, req *environ.Request, result *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()

	logrus.Debugf("requesting from %s", upstream)

	client := environ.Client(upstream, p.secret, p.skipVerify)
	env, err := client.List(ctx, req)
	if err != nil {
		logrus.Warning(err)
		return
	}

	for key, val := range env {
		result.Store(key, val)
	}

	logrus.Debugf("%s data: %v", upstream, env)
}

func (p *plugin) List(ctx context.Context, req *environ.Request) (map[string]string, error) {
	var result sync.Map
	var wg sync.WaitGroup

	for _, upstream := range p.upstreams {
		wg.Add(1)
		go p.fetchUpstreamEnviron(upstream, ctx, req, &result, &wg)
	}

	wg.Wait()

	env := unwrap(&result)
	logrus.Debugf("merged: %v", env)

	return env, nil
}

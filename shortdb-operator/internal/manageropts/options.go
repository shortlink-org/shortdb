/*
Copyright 2022 Viktor Login.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package manageropts builds sigs.k8s.io/controller-runtime manager.Options via functional options.
package manageropts

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

type config struct {
	scheme                 *runtime.Scheme
	metrics                metricsserver.Options
	webhook                webhook.Options
	healthProbeBindAddress string
	leaderElection         bool
	leaderElectionID       string
}

// Option configures the manager ctrl.Options produced by Build.
type Option func(*config) error

// WithScheme sets the runtime.Scheme (required).
func WithScheme(s *runtime.Scheme) Option {
	return func(c *config) error {
		if s == nil {
			return errors.New("scheme must not be nil")
		}
		c.scheme = s
		return nil
	}
}

// WithMetrics sets metricsserver.Options (e.g. BindAddress, TLS, extra handlers).
func WithMetrics(m metricsserver.Options) Option {
	return func(c *config) error {
		c.metrics = m
		return nil
	}
}

// WithWebhook sets webhook.Options passed to webhook.NewServer.
func WithWebhook(o webhook.Options) Option {
	return func(c *config) error {
		c.webhook = o
		return nil
	}
}

// WithHealthProbeBindAddress sets HealthProbeBindAddress on ctrl.Options (empty disables probes server binding).
func WithHealthProbeBindAddress(addr string) Option {
	return func(c *config) error {
		c.healthProbeBindAddress = addr
		return nil
	}
}

// WithLeaderElection toggles leader election on the manager.
func WithLeaderElection(on bool) Option {
	return func(c *config) error {
		c.leaderElection = on
		return nil
	}
}

// WithLeaderElectionID sets LeaderElectionID (required when leader election is enabled).
func WithLeaderElectionID(id string) Option {
	return func(c *config) error {
		c.leaderElectionID = id
		return nil
	}
}

// Build returns ctrl.Options for ctrl.NewManager after applying opts.
func Build(opts ...Option) (ctrl.Options, error) {
	var c config
	for i, opt := range opts {
		if opt == nil {
			return ctrl.Options{}, fmt.Errorf("option %d is nil", i)
		}
		if err := opt(&c); err != nil {
			return ctrl.Options{}, fmt.Errorf("option %d: %w", i, err)
		}
	}
	if c.scheme == nil {
		return ctrl.Options{}, errors.New("WithScheme is required")
	}
	if c.leaderElection && c.leaderElectionID == "" {
		return ctrl.Options{}, errors.New("WithLeaderElectionID is required when leader election is enabled")
	}

	return ctrl.Options{
		Scheme:                 c.scheme,
		Metrics:                c.metrics,
		WebhookServer:          webhook.NewServer(c.webhook),
		HealthProbeBindAddress: c.healthProbeBindAddress,
		LeaderElection:         c.leaderElection,
		LeaderElectionID:       c.leaderElectionID,
	}, nil
}

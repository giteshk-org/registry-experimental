// Copyright 2021 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"

	longrunning "cloud.google.com/go/longrunning/autogen/longrunningpb"
	experimental_rpc "github.com/apigee/registry-experimental/rpc"
	"github.com/apigee/registry-experimental/server/search"
	"github.com/apigee/registry/pkg/log"
	"github.com/apigee/registry/pkg/log/interceptor"
	registry_rpc "github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"
)

// ServerConfig is the top-level configuration structure.
type ServerConfig struct {
	// Server port. If unset or zero, an open port will be assigned.
	Port     int            `yaml:"port"`
	Database DatabaseConfig `yaml:"database"`
	Logging  LoggingConfig  `yaml:"logging"`
	Pubsub   PubsubConfig   `yaml:"pubsub"`
}

// DatabaseConfig holds database configuration.
type DatabaseConfig struct {
	// Driver for the database connection.
	// Values: [ sqlite3, postgres, cloudsqlpostgres ]
	Driver string `yaml:"driver"`
	// Config for the database connection. The format is a data source name (DSN).
	// PostgreSQL Reference: See "Connection Strings" at https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
	// SQLite Reference: See "URI filename examples" at https://www.sqlite.org/c3ref/open.html
	Config string `yaml:"config"`
}

// LoggingConfig holds logging configuration.
type LoggingConfig struct {
	// Level of logging to print to standard output.
	// Values: [ debug, info, warn, error, fatal ]
	Level string `yaml:"level"`
	// Format of log entries.
	// Options: [ json, text ]
	Format string `yaml:"format"`
}

// PubsubConfig holds pubsub (notification) configuration.
type PubsubConfig struct {
	// Enable Pub/Sub for event notification publishing.
	// Values: [ true, false ]
	Enable bool `yaml:"enable"`
	// Project ID of the Google Cloud project to use for Pub/Sub.
	// Reference: https://cloud.google.com/resource-manager/docs/creating-managing-projects
	Project string `yaml:"project"`
}

// default configuration
var config = ServerConfig{
	Port: 8080,
	Database: DatabaseConfig{
		Driver: "sqlite3",
		Config: "file:/tmp/registry.db",
	},
	Logging: LoggingConfig{
		Level:  "info",
		Format: "text",
	},
	Pubsub: PubsubConfig{
		Enable:  false,
		Project: "",
	},
}

func main() {
	var configPath string
	pflag.StringVarP(&configPath, "configuration", "c", "", "The server configuration file to load.")
	pflag.Parse()

	// Use a default logger configuration until we load the server config.
	bootLogger := log.NewLogger()
	if configPath != "" {
		bootLogger.Infof("Loading configuration from %s", configPath)
		raw, err := ioutil.ReadFile(configPath)
		if err != nil {
			bootLogger.WithError(err).Fatal("Failed to open config file")
		}
		// Expand environment variables before unmarshaling.
		expanded := []byte(os.ExpandEnv(string(raw)))
		err = yaml.Unmarshal(expanded, &config)
		if err != nil {
			bootLogger.WithError(err).Fatalf("Failed to read config file")
		}
	}

	if err := validateConfig(); err != nil {
		bootLogger.WithError(err).Fatalf("Invalid configuration")
	}

	// Use logging options from the server config.
	var (
		logOpts        = loggerOptions(config.Logging)
		logger         = log.NewLogger(logOpts...)
		logInterceptor = interceptor.CallLogger(logOpts...)
	)

	logger.Infof("Configured port %d", config.Port)
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: config.Port,
	})
	if err != nil {
		logger.WithError(err).Fatalf("Failed to create TCP listener")
	}
	defer listener.Close()

	registryServer, err := registry.New(registry.Config{
		Database:  config.Database.Driver,
		DBConfig:  config.Database.Config,
		LogLevel:  config.Logging.Level,
		LogFormat: config.Logging.Format,
		Notify:    config.Pubsub.Enable,
		ProjectID: config.Pubsub.Project,
	})
	if err != nil {
		logger.WithError(err).Fatalf("Failed to create registry server")
	}

	searchServer := search.New(search.Config{
		Database: config.Database.Driver,
		DBConfig: config.Database.Config,
	}, registryServer)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(logInterceptor))
	reflection.Register(grpcServer)
	registry_rpc.RegisterRegistryServer(grpcServer, registryServer)
	registry_rpc.RegisterAdminServer(grpcServer, registryServer)
	experimental_rpc.RegisterSearchServer(grpcServer, searchServer)
	longrunning.RegisterOperationsServer(grpcServer, searchServer)

	go func() {
		_ = grpcServer.Serve(listener)
	}()
	logger.Infof("Listening on %s", listener.Addr())

	// Wait for an interruption signal.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done
}

func validateConfig() error {
	if config.Port < 0 {
		return fmt.Errorf("invalid port %q: must be non-negative", config.Port)
	}

	switch driver := config.Database.Driver; driver {
	case "sqlite3", "postgres", "cloudsqlpostgres":
	default:
		return fmt.Errorf("invalid database.driver %q: must be one of [sqlite3, postgres, cloudsqlpostgres]", driver)
	}

	switch level := config.Logging.Level; level {
	case "fatal", "error", "warn", "info", "debug":
	default:
		return fmt.Errorf("invalid logging.level %q: must be one of [fatal, error, warn, info, debug]", level)
	}

	switch format := config.Logging.Format; format {
	case "json", "text":
	default:
		return fmt.Errorf("invalid logging format %q: must be one of [json, text]", format)
	}

	if project := config.Pubsub.Project; config.Pubsub.Enable && project == "" {
		return fmt.Errorf("invalid pubsub.project %q: pubsub cannot be enabled without GCP project ID", project)
	}

	return nil
}

func loggerOptions(conf LoggingConfig) []log.Option {
	opts := make([]log.Option, 0, 2)
	switch conf.Level {
	case "debug":
		opts = append(opts, log.DebugLevel)
	case "info":
		opts = append(opts, log.InfoLevel)
	case "warn":
		opts = append(opts, log.WarnLevel)
	case "error":
		opts = append(opts, log.ErrorLevel)
	case "fatal":
		opts = append(opts, log.FatalLevel)
	}

	switch conf.Format {
	case "json":
		opts = append(opts, log.JSONFormat(os.Stderr))
	case "text":
		opts = append(opts, log.TextFormat(os.Stderr))
	}

	return opts
}

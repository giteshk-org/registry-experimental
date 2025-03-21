// Copyright 2023 Google LLC. All Rights Reserved.
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

package count

import (
	"context"
	"fmt"

	"github.com/apigee/registry/cmd/registry/tasks"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/log"
	"github.com/apigee/registry/pkg/names"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/protobuf/field_mask"
)

func deploymentsCommand() *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   "deployments",
		Short: "Count the number of deployments of specified APIs",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			c, err := connection.ActiveConfig()
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get config")
			}
			args[0] = c.FQName(args[0])

			client, err := connection.NewRegistryClientWithSettings(ctx, c)
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get client")
			}
			// Initialize task queue.
			jobs, err := cmd.Flags().GetInt("jobs")
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to get jobs from flags")
			}
			taskQueue, wait := tasks.WorkerPoolIgnoreError(ctx, jobs)
			defer wait()

			api, err := names.ParseApi(args[0])
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed parse")
			}

			// Iterate through a collection of APIs and count the number of deployments of each.
			err = visitor.ListAPIs(ctx, client, api, filter, func(ctx context.Context, api *rpc.Api) error {
				taskQueue <- &countApiDeploymentsTask{
					client: client,
					api:    api,
				}
				return nil
			})
			if err != nil {
				log.FromContext(ctx).WithError(err).Fatal("Failed to list APIs")
			}
		},
	}

	cmd.Flags().StringVar(&filter, "filter", "", "filter selected resources")
	return cmd
}

type countApiDeploymentsTask struct {
	client connection.RegistryClient
	api    *rpc.Api
}

func (task *countApiDeploymentsTask) String() string {
	return "count deployments " + task.api.Name
}

func (task *countApiDeploymentsTask) Run(ctx context.Context) error {
	count := 0
	request := &rpc.ListApiDeploymentsRequest{
		Parent: task.api.Name,
	}
	it := task.client.ListApiDeployments(ctx, request)
	for {
		_, err := it.Next()
		if err == iterator.Done {
			break
		} else if err == nil {
			count++
		} else {
			return err
		}
	}
	log.Debugf(ctx, "%d\t%s", count, task.api.Name)
	if task.api.Labels == nil {
		task.api.Labels = make(map[string]string, 0)
	}
	task.api.Labels["deployments"] = fmt.Sprintf("%d", count)
	_, err := task.client.UpdateApi(ctx,
		&rpc.UpdateApiRequest{
			Api: task.api,
			UpdateMask: &field_mask.FieldMask{
				Paths: []string{"labels"},
			},
		})
	return err
}

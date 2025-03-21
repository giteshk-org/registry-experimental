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

package extract

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/apigee/registry-experimental/pkg/yamlquery"
	"github.com/apigee/registry/cmd/registry/compress"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/mime"
	"github.com/apigee/registry/pkg/names"
	"github.com/apigee/registry/pkg/visitor"
	"github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract PATTERN",
		Short: "Extract properties from specs and artifacts stored in the registry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c, err := connection.ActiveConfig()
			if err != nil {
				return err
			}
			pattern := c.FQName(args[0])
			filter, err := cmd.Flags().GetString("filter")
			if err != nil {
				return err
			}
			registryClient, err := connection.NewRegistryClientWithSettings(ctx, c)
			if err != nil {
				return err
			}
			v := &extractVisitor{
				registryClient: registryClient,
			}
			return visitor.Visit(ctx, v, visitor.VisitorOptions{
				RegistryClient: registryClient,
				Pattern:        pattern,
				Filter:         filter,
			})
		},
	}
	cmd.Flags().String("filter", "", "Filter selected resources")
	cmd.Flags().Int("jobs", 10, "Number of actions to perform concurrently")
	return cmd
}

type extractVisitor struct {
	visitor.Unsupported
	registryClient connection.RegistryClient
}

var empty = ""

func (v *extractVisitor) SpecHandler() visitor.SpecHandler {
	return func(ctx context.Context, spec *rpc.ApiSpec) error {
		fmt.Printf("%s\n", spec.Name)
		err := visitor.FetchSpecContents(ctx, v.registryClient, spec)
		if err != nil {
			return err
		}
		bytes := spec.Contents
		if mime.IsGZipCompressed(spec.MimeType) {
			bytes, err = compress.GUnzippedBytes(bytes)
			if err != nil {
				return err
			}
		}
		if mime.IsOpenAPIv2(spec.MimeType) || mime.IsOpenAPIv3(spec.MimeType) {
			var node yaml.Node
			if err := yaml.Unmarshal(bytes, &node); err != nil {
				return err
			}

			openapi := yamlquery.QueryString(&node, "openapi")

			swagger := yamlquery.QueryString(&node, "swagger")

			description := yamlquery.QueryString(&node, "info.description")
			if description == nil {
				description = &empty
			}
			*description = markdownify(*description)

			title := yamlquery.QueryString(&node, "info.title")
			if title == nil {
				title = &empty
			}

			provider := yamlquery.QueryString(&node, "info.x-providerName")

			categories := yamlquery.QueryNode(&node, "info.x-apisguru-categories")

			// Set API (displayName, description) from (title, description).
			specName, _ := names.ParseSpec(spec.Name)
			apiName := specName.Api()
			api, err := v.registryClient.GetApi(ctx,
				&rpc.GetApiRequest{
					Name: apiName.String(),
				},
			)
			if err != nil {
				return err
			}
			labels := api.Labels
			if labels == nil {
				labels = make(map[string]string)
			}
			labels["openapi"] = "true"
			delete(labels, "style-openapi")
			labels["categories"] = strings.Join(yamlquery.QueryStringArray(categories), ",")
			if provider != nil {
				labels["provider"] = *provider
			}
			_, err = v.registryClient.UpdateApi(ctx,
				&rpc.UpdateApiRequest{
					Api: &rpc.Api{
						Name:        apiName.String(),
						DisplayName: *title,
						Description: *description,
						Labels:      labels,
					},
				},
			)
			if err != nil {
				return err
			}

			// Set the spec mimetype (this should not bump the revision!).
			if openapi != nil || swagger != nil {
				var compression string
				if mime.IsGZipCompressed(spec.MimeType) {
					compression = "+gzip"
				}
				var mimeType string
				if openapi != nil {
					mimeType = mime.OpenAPIMimeType(compression, *openapi)
				} else if swagger != nil {
					mimeType = mime.OpenAPIMimeType(compression, *swagger)
				}
				specName, _ := names.ParseSpec(spec.Name)
				_, err := v.registryClient.UpdateApiSpec(ctx,
					&rpc.UpdateApiSpecRequest{
						ApiSpec: &rpc.ApiSpec{
							Name:     specName.String(),
							MimeType: mimeType,
						},
					},
				)
				if err != nil {
					return err
				}
			}
		}
		if mime.IsDiscovery(spec.MimeType) {
			var node yaml.Node
			if err := yaml.Unmarshal(bytes, &node); err != nil {
				return err
			}
			styleForYAML(&node)
			//fmt.Printf("discovery:\n%s\n", yamlquery.Describe(&node))

			description := yamlquery.QueryString(&node, "description")
			if description == nil {
				description = &empty
			}

			title := yamlquery.QueryString(&node, "canonicalName")
			if title == nil {
				title = &empty
			}

			provider := yamlquery.QueryString(&node, "ownerDomain")

			// Set API (displayName, description) from (title, description).
			specName, _ := names.ParseSpec(spec.Name)
			apiName := specName.Api()
			api, err := v.registryClient.GetApi(ctx,
				&rpc.GetApiRequest{
					Name: apiName.String(),
				},
			)
			if err != nil {
				return err
			}
			labels := api.Labels
			if labels == nil {
				labels = make(map[string]string)
			}
			labels["discovery"] = "true"
			delete(labels, "style-discovery")
			if provider != nil {
				labels["provider"] = *provider
			}
			_, err = v.registryClient.UpdateApi(ctx,
				&rpc.UpdateApiRequest{
					Api: &rpc.Api{
						Name:        apiName.String(),
						DisplayName: *title,
						Description: *description,
						Labels:      labels,
					},
				},
			)
			if err != nil {
				return err
			}

		}
		if mime.IsProto(spec.MimeType) {
			// create a tmp directory
			root, err := os.MkdirTemp("", "extract-protos-")
			if err != nil {
				return err
			}
			// whenever we finish, delete the tmp directory
			defer os.RemoveAll(root)
			// unzip the protos to the temp directory
			_, err = compress.UnzipArchiveToPath(spec.Contents, root)
			if err != nil {
				return err
			}

			var displayName string
			var description string

			if err = filepath.Walk(root, func(filepath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// Skip everything that's not a YAML file.
				if info.IsDir() || !strings.HasSuffix(filepath, ".yaml") {
					return nil
				}

				bytes, err := os.ReadFile(filepath)
				if err != nil {
					return err
				}

				sc := &ServiceConfig{}
				if err := yaml.Unmarshal(bytes, sc); err != nil {
					return err
				}

				// Skip invalid API service configurations.
				if sc.Type != "google.api.Service" || sc.Title == "" || sc.Name == "" {
					return nil
				}

				displayName = sc.Title
				description = strings.ReplaceAll(sc.Documentation.Summary, "\n", " ")

				// Skip the directory after we find an API service configuration.
				return fs.SkipDir
			}); err != nil {
				return err
			}

			specName, _ := names.ParseSpec(spec.Name)
			apiName := specName.Api()
			api, err := v.registryClient.GetApi(ctx,
				&rpc.GetApiRequest{
					Name: apiName.String(),
				},
			)
			if err != nil {
				return err
			}
			labels := api.Labels
			if labels == nil {
				labels = make(map[string]string)
			}
			labels["grpc"] = "true"
			delete(labels, "style-grpc")
			labels["provider"] = "google.com"
			_, err = v.registryClient.UpdateApi(ctx,
				&rpc.UpdateApiRequest{
					Api: &rpc.Api{
						Name:        apiName.String(),
						DisplayName: displayName,
						Description: description,
						Labels:      labels,
					},
				},
			)
			if err != nil {
				return err
			}

		}
		return nil
	}
}

// The API Service Configuration contains important API properties.
type ServiceConfig struct {
	Type          string `yaml:"type"`
	Name          string `yaml:"name"`
	Title         string `yaml:"title"`
	Documentation struct {
		Summary string `yaml:"summary"`
	} `yaml:"documentation"`
}

// styleForYAML sets the style field on a tree of yaml.Nodes for YAML export.
func styleForYAML(node *yaml.Node) {
	node.Style = 0
	for _, n := range node.Content {
		styleForYAML(n)
	}
}

func markdownify(text string) string {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(text)
	if err != nil {
		return text
	}
	return markdown
}

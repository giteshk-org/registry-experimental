// Copyright 2020 Google LLC. All Rights Reserved.
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
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count quantities in the API Registry",
	}

	cmd.AddCommand(deploymentsCommand())
	cmd.AddCommand(revisionsCommand())
	cmd.AddCommand(versionsCommand())

	cmd.PersistentFlags().Int("jobs", 10, "Number of actions to perform concurrently")
	return cmd
}

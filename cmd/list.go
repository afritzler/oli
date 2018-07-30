// Copyright © 2018 NAME HERE <andreas.fritzler@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/afritzler/oli/pkg/client"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List everything LBaaS specific in your tenant",
	Long:  `List everything LBaaS specific in your tenant.`,
	Run: func(cmd *cobra.Command, args []string) {
		listEverything()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listEverything() {
	osClient, err := client.NewOpenStackProvider()
	if err != nil {
		panic(fmt.Errorf("failed to create os client %s", err))
	}
	lbs, err := osClient.ListLBaaS()
	if err != nil {
		panic(fmt.Errorf("failed to list lb ids %s", err))
	}
	fmt.Printf("-------------------------------------------------------------------------\n")
	fmt.Printf("LoadBalancer IDs                     | Name\n")
	fmt.Printf("-------------------------------------------------------------------------\n")
	for _, lb := range lbs {
		fmt.Printf("%s | %s\n", lb.ID, lb.Name)
	}
	listenerIDs, err := osClient.ListListenerIDsForCurrentTenant()
	if err != nil {
		panic(fmt.Errorf("failed to list listener ids %s", err))
	}
	fmt.Printf("-------------------------------------------------------------------------\n")
	fmt.Printf("Listener IDs\n")
	fmt.Printf("-------------------------------------------------------------------------\n")
	for _, id := range listenerIDs {
		fmt.Printf("%s\n", id)
	}
	monitorIDs, err := osClient.ListMonitorIDsForCurrentTenant()
	if err != nil {
		panic(fmt.Errorf("failed to list healthmonitor ids %s", err))
	}
	fmt.Printf("-------------------------------------------------------------------------\n")
	fmt.Printf("Healthmonitor IDs\n")
	fmt.Printf("-------------------------------------------------------------------------\n")
	for _, id := range monitorIDs {
		fmt.Printf("%s\n", id)
	}
	poolIDs, err := osClient.GetPoolIDsForCurrentTenant()
	if err != nil {
		panic(fmt.Errorf("failed to list healthmonitor ids %s", err))
	}
	fmt.Printf("-------------------------------------------------------------------------\n")
	fmt.Printf("Pool IDs\n")
	fmt.Printf("-------------------------------------------------------------------------\n")
	for _, id := range poolIDs {
		fmt.Printf("%s\n", id)
	}
}

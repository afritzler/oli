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
	"github.com/afritzler/oli/pkg/renderer"
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
	r := renderer.NewTreeRenderer()

	osClient, err := client.NewDefaultOpenStackProvider()
	if err != nil {
		panic(fmt.Errorf("failed to create os client %s", err))
	}
	lbs, err := osClient.ListLBaaS()
	if err != nil {
		panic(fmt.Errorf("failed to list lb ids %s", err))
	}
	for _, lb := range lbs {
		r.AddLoadBalancer(lb)
	}
	listeners, err := osClient.ListListenersForCurrentTenant()
	if err != nil {
		panic(fmt.Errorf("failed to list listener %s", err))
	}
	for _, listener := range listeners {
		r.AddListener(listener)
	}
	pools, err := osClient.GetPoolsForCurrentTenant()
	if err != nil {
		panic(fmt.Errorf("failed to list pools %s", err))
	}
	for _, pool := range pools {
		r.AddPool(pool)
		members, err := osClient.GetMembersForPoolID(pool.ID)
		if err != nil {
			panic(fmt.Errorf("failed to list members for pool %s, %s", pool.ID, err))
		}
		for _, member := range members {
			r.AddMember(pool.ID, member)
		}
	}
	monitors, err := osClient.ListMonitorsForCurrentTenant()
	if err != nil {
		panic(fmt.Errorf("failed to list healthmonitor ids %s", err))
	}
	for _, monitor := range monitors {
		r.AddMonitor(monitor)
	}
	fmt.Println(r.GetTreeStringWithLegend())
}

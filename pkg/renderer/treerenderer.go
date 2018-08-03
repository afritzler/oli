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

package renderer

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/listeners"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/monitors"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/pools"
	"github.com/xlab/treeprint"
)

const (
	legend = "[LB] LoadBalancer, [L] Listener, [P] Pool, [M] Member, [HM] HealthMonitor"
)

type TreeRenderer interface {
	AddLoadBalancer(loadbalancer loadbalancers.LoadBalancer) treeprint.Tree
	AddListener(listener listeners.Listener) treeprint.Tree
	AddPool(pool pools.Pool) treeprint.Tree
	AddMonitor(monitor monitors.Monitor) treeprint.Tree
	AddMember(poolid string, member pools.Member) treeprint.Tree
	GetTreeString() string
	GetTreeStringWithLegend() string
}

type treerenderer struct {
	tree treeprint.Tree
}

func NewTreeRenderer() TreeRenderer {
	tree := treeprint.New()
	return &treerenderer{tree: tree}
}

func (t *treerenderer) AddLoadBalancer(loadbalancer loadbalancers.LoadBalancer) treeprint.Tree {
	t.tree.AddMetaBranch(loadbalancer.ID, t.renderName("LB", loadbalancer.Name, loadbalancer.AdminStateUp))
	return t.tree
}

func (t *treerenderer) AddListener(listener listeners.Listener) treeprint.Tree {
	for _, lb := range listener.Loadbalancers {
		lbNode := t.tree.FindByMeta(lb.ID)
		if lbNode == nil {
			t.addOrphan(listener.ID, t.renderName("L", listener.Name, listener.AdminStateUp))
		} else {
			lbNode.AddMetaNode(listener.ID, t.renderName("L", listener.Name, listener.AdminStateUp))
		}
	}
	return t.tree
}

func (t *treerenderer) AddPool(pool pools.Pool) treeprint.Tree {
	for _, listener := range pool.Listeners {
		lbNode := t.tree.FindByMeta(listener.ID)
		if lbNode == nil {
			t.addOrphan(pool.ID, t.renderName("P", pool.Name, pool.AdminStateUp))
		} else {
			lbNode.AddMetaNode(pool.ID, t.renderName("P", pool.Name, pool.AdminStateUp))
		}
	}
	return t.tree
}

func (t *treerenderer) AddMonitor(monitor monitors.Monitor) treeprint.Tree {
	for _, pool := range monitor.Pools {
		lbNode := t.tree.FindByMeta(pool.ID)
		if lbNode == nil {
			t.addOrphan(monitor.ID, t.renderName("HM", monitor.Name, monitor.AdminStateUp))
		} else {
			lbNode.AddMetaNode(monitor.ID, t.renderName("HM", monitor.Name, monitor.AdminStateUp))
		}
	}
	return t.tree
}

// pool id needed as seperate arg since the member.PoolID field is empty (bug?)
func (t *treerenderer) AddMember(poolid string, member pools.Member) treeprint.Tree {
	lbNode := t.tree.FindByMeta(poolid)
	if lbNode == nil {
		t.addOrphan(member.ID, t.renderName("M", member.Name, member.AdminStateUp))
	} else {
		lbNode.AddMetaNode(member.ID, t.renderName("M", member.Name, member.AdminStateUp))
	}
	return t.tree
}

func (t *treerenderer) addOrphan(meta string, name string) treeprint.Tree {
	orphan := t.tree.FindByMeta("42")
	if orphan == nil {
		orphan := t.tree.AddMetaBranch("42", "Orphan Objects")
		orphan.AddMetaNode(meta, name)
	} else {
		orphan.AddMetaNode(meta, name)
	}
	return t.tree
}

func (t *treerenderer) GetTreeString() string {
	return t.tree.String()
}

func (t *treerenderer) GetTreeStringWithLegend() string {
	return t.tree.String() + "\n" + legend
}

func (t *treerenderer) renderName(kind string, name string, state bool) string {
	return fmt.Sprintf("[%s] %s Up: %t", kind, name, state)
}

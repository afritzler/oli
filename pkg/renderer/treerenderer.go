package renderer

import (
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
	t.tree.AddMetaBranch(loadbalancer.ID, "[LB] "+loadbalancer.Name)
	return t.tree
}

func (t *treerenderer) AddListener(listener listeners.Listener) treeprint.Tree {
	for _, lb := range listener.Loadbalancers {
		lbNode := t.tree.FindByMeta(lb.ID)
		if lbNode == nil {
			t.addOrphan(listener.ID, "[M] "+listener.Name)
		} else {
			lbNode.AddMetaNode(listener.ID, "[L] "+listener.Name)
		}
	}
	return t.tree
}

func (t *treerenderer) AddPool(pool pools.Pool) treeprint.Tree {
	for _, listener := range pool.Listeners {
		lbNode := t.tree.FindByMeta(listener.ID)
		if lbNode == nil {
			t.addOrphan(pool.ID, "[M] "+pool.Name)
		} else {
			lbNode.AddMetaNode(pool.ID, "[P] "+pool.Name)
		}
	}
	return t.tree
}

func (t *treerenderer) AddMonitor(monitor monitors.Monitor) treeprint.Tree {
	for _, pool := range monitor.Pools {
		lbNode := t.tree.FindByMeta(pool.ID)
		if lbNode == nil {
			t.addOrphan(monitor.ID, "[M] "+monitor.Name)
		} else {
			lbNode.AddMetaNode(monitor.ID, "[HM] "+monitor.Name)
		}
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

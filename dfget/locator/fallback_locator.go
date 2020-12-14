/*
 * Copyright The Dragonfly Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package locator

import (
	"fmt"
	"sync/atomic"

	"github.com/dragonflyoss/Dragonfly/dfget/config"
	"github.com/dragonflyoss/Dragonfly/pkg/algorithm"
	"github.com/dragonflyoss/Dragonfly/pkg/netutils"
)

// FallbackLocator uses the nodes passed from configuration or CLI.
// Unlike the static locator it will always connect to the first node in list, using other nodes as fallback, shuffling them to distribute the load.
// This is useful when we want primary supernode to serve all local datacenter traffic and redirect to other when local fails
type FallbackLocator struct {
	idx   int32
	Group *SupernodeGroup
}

// ----------------------------------------------------------------------------
// constructors

// NewStaticLocator constructs FallbackLocator which uses the nodes passed from
// configuration or CLI.
func NewFallbackLocator(groupName string, nodes []*config.NodeWeight) *FallbackLocator {
	locator := &FallbackLocator{
		idx: -1,
	}
	if len(nodes) == 0 {
		return locator
	}
	group := &SupernodeGroup{
		Name: groupName,
	}
	for _, node := range nodes {
		ip, port := netutils.GetIPAndPortFromNode(node.Node, config.DefaultSupernodePort)
		if ip == "" {
			continue
		}
		supernode := &Supernode{
			Schema:    config.DefaultSupernodeSchema,
			IP:        ip,
			Port:      port,
			Weight:    node.Weight,
			GroupName: groupName,
		}

		// fallback locator uses same weight for all nodes
		group.Nodes = append(group.Nodes, supernode)
	}
	shuffleFallbackNodes(group.Nodes)
	locator.Group = group
	return locator
}

// NewStaticLocatorFromStr constructs FallbackLocator from string list.
// The format of nodes is: ip:port=weight
func NewFallbackLocatorFromStr(groupName string, nodes []string) (*FallbackLocator, error) {
	nodeWeight, err := config.ParseNodesSlice(nodes)
	if err != nil {
		return nil, err
	}
	return NewFallbackLocator(groupName, nodeWeight), nil
}

// ----------------------------------------------------------------------------
// implement api methods

// Get returns the current selected supernode, it should be idempotent.
// It should return nil before first calling the Next method.
func (s *FallbackLocator) Get() *Supernode {
	if s.Group == nil {
		return nil
	}
	return s.Group.GetNode(s.load())
}

// Next chooses the next available supernode for retrying or other
// purpose. The current supernode should be set as this result.
func (s *FallbackLocator) Next() *Supernode {
	if s.Group == nil || s.load() >= len(s.Group.Nodes) {
		return nil
	}
	return s.Group.GetNode(s.inc())
}

// Select chooses a supernode based on the giving key.
// It should not affect the result of method 'Get()'.
func (s *FallbackLocator) Select(key interface{}) *Supernode {
	// unnecessary to implement this method
	return nil
}

// GetGroup returns the group with the giving name.
func (s *FallbackLocator) GetGroup(name string) *SupernodeGroup {
	if s.Group == nil || s.Group.Name != name {
		return nil
	}
	return s.Group
}

// All returns all the supernodes.
func (s *FallbackLocator) All() []*SupernodeGroup {
	if s.Group == nil {
		return nil
	}
	return []*SupernodeGroup{s.Group}
}

// Size returns the number of all supernodes.
func (s *FallbackLocator) Size() int {
	if s.Group == nil {
		return 0
	}
	return len(s.Group.Nodes)
}

// Report records the metrics of the current supernode in order to choose a
// more appropriate supernode for the next time if necessary.
func (s *FallbackLocator) Report(node string, metrics *SupernodeMetrics) {
	// unnecessary to implement this method
	return
}

// Refresh refreshes all the supernodes.
func (s *FallbackLocator) Refresh() bool {
	atomic.StoreInt32(&s.idx, -1)
	return true
}

func (s *FallbackLocator) String() string {
	idx := s.load()
	if s.Group == nil || idx >= len(s.Group.Nodes) {
		return "empty"
	}

	nodes := make([]string, len(s.Group.Nodes)-idx-1)
	for i := idx + 1; i < len(s.Group.Nodes); i++ {
		n := s.Group.GetNode(i)
		nodes[i-idx-1] = fmt.Sprintf("%s:%d=%d", n.IP, n.Port, n.Weight)
	}
	return s.Group.Name + ":" + fmt.Sprintf("%v", nodes)
}

// ----------------------------------------------------------------------------
// private methods of FallbackLocator

func (s *FallbackLocator) load() int {
	return int(atomic.LoadInt32(&s.idx))
}

func (s *FallbackLocator) inc() int {
	return int(atomic.AddInt32(&s.idx, 1))
}

// ----------------------------------------------------------------------------
// helper functions

func shuffleFallbackNodes(nodes []*Supernode) []*Supernode {
	// shuffle all except the very first
	if length := len(nodes); length > 1 {
		algorithm.Shuffle(length - 1, func(i, j int) {
			nodes[i + 1], nodes[j + 1] = nodes[j + 1], nodes[i + 1]
		})
	}
	return nodes
}

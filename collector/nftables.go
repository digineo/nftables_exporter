// Copyright 2020 Intrinsec
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"log"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "nftables"

var (
	nftablesCounterBytesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "bytes_total"),
		"Total number of packets",
		[]string{"table", "chain", "comment"}, nil,
	)
	nftablesCounterPacketsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "packets_total"),
		"Total number of bytes",
		[]string{"table", "chain", "comment"}, nil,
	)
)

// Update collects the metrics
func (c *Collector) Update(ch chan<- prometheus.Metric) (err error) {
	nft := &nftables.Conn{}

	tables, err := nft.ListTables()
	if err != nil {
		log.Println(err)
		return err
	}

	chains, err := nft.ListChains()
	if err != nil {
		log.Println(err)
		return err
	}

	for _, table := range tables {
		for _, chain := range chains {
			rules, err := nft.GetRule(table, chain)
			if err != nil {
				log.Println(err)
				return err
			}

			for _, rule := range rules {
				c.addRule(ch, table, chain, rule)
			}
		}
	}
	return nil
}

func (c *Collector) addRule(ch chan<- prometheus.Metric, table *nftables.Table, chain *nftables.Chain, rule *nftables.Rule) {
	var comment string
	var bytes, packets uint64
	var hasCounter bool

	// extract comment
	if data := rule.UserData; len(data) > 2 {
		length := int(data[1])
		if data[0] == 0 && len(data) >= 2+length {
			comment = string(data[2 : length+1])
		}
	}

	// copy counter values
	for _, rawExpr := range rule.Exprs {
		switch ex := rawExpr.(type) {
		case *expr.Counter:
			hasCounter = true
			bytes = ex.Bytes
			packets = ex.Packets
		}
	}

	if hasCounter && comment != "" {
		labels := []string{table.Name, chain.Name, comment}
		ch <- prometheus.MustNewConstMetric(nftablesCounterBytesDesc, prometheus.CounterValue, float64(bytes), labels...)
		ch <- prometheus.MustNewConstMetric(nftablesCounterPacketsDesc, prometheus.CounterValue, float64(packets), labels...)
	}
}

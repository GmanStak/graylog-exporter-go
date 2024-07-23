package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// 定义自定义指标
var (
	graylog_index_count_num = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "graylog_index_count_num",
			Help: "Number of index sets in Graylog",
		},
		[]string{"id", "index_prefix"},
	)
	graylog_index_size_bytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "graylog_index_size_bytes",
			Help: "graylog 索引容量",
		},
		[]string{"id", "index_prefix"},
	)
	graylog_sidecar_count_num = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "graylog_sidecar_count_num",
			Help: "graylog sidecar 数量",
		},
		[]string{"host"},
	)
	graylog_cluster_info = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "graylog_cluster_info",
			Help: "graylog 集群参数",
		},
		[]string{"id", "index_prefix"},
	)
	graylog_journal_entries = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "graylog_journal_entries",
			Help: "uncommitted_journal_entries",
		},
		[]string{"name"},
	)
	graylog_sidecar_node_status = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "graylog_sidecar_node_status",
			Help: "graylog sidecar 节点运行状态",
		},
		[]string{"node_name", "operating_system", "node_ip", "message"},
	)
)

func init() {
	//手动删除go和process指标
	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.Unregister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	prometheus.MustRegister(graylog_index_size_bytes)
	prometheus.MustRegister(graylog_index_count_num)
	prometheus.MustRegister(graylog_sidecar_count_num)
	prometheus.MustRegister(graylog_cluster_info)
	prometheus.MustRegister(graylog_journal_entries)
	prometheus.MustRegister(graylog_sidecar_node_status)

	// 注册自定义指标，通过创建新的Registry只提供注册的指标
	//Registry := prometheus.NewRegistry()
	//Registry.MustRegister(graylog_index_size_bytes)
	//Registry.MustRegister(graylog_index_count_num)
	//Registry.MustRegister(graylog_sidecar_count_num)
	//Registry.MustRegister(graylog_cluster_info)
	//Registry.MustRegister(graylog_journal_entries)
	//Registry.MustRegister(graylog_sidecar_node_status)
}

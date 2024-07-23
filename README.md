# graylog-exporter-go
# 启动参数
-c /metircs 访问认证文件路径（auth_config.yaml）
-a graylog url 地址及认证路径（url_config.yaml）
-p /metircs 访问端口（:8080）
# 采集指标
graylog_index_count_num
graylog_index_size_bytes
graylog_sidecar_count_num
graylog_cluster_info
graylog_journal_entries
graylog_sidecar_node_status

package main

import "time"

type Journal_info struct {
	Enabled                   bool          `json:"enabled"`
	AppendEventsPerSecond     int           `json:"append_events_per_second"`
	ReadEventsPerSecond       int           `json:"read_events_per_second"`
	UncommittedJournalEntries int           `json:"uncommitted_journal_entries"`
	JournalSize               int           `json:"journal_size"`
	JournalSizeLimit          int64         `json:"journal_size_limit"`
	NumberOfSegments          int           `json:"number_of_segments"`
	OldestSegment             time.Time     `json:"oldest_segment"`
	JournalConfig             JournalConfig `json:"journal_config"`
}
type JournalConfig struct {
	Directory     string `json:"directory"`
	SegmentSize   int    `json:"segment_size"`
	SegmentAge    int    `json:"segment_age"`
	MaxSize       int64  `json:"max_size"`
	MaxAge        int    `json:"max_age"`
	FlushInterval int    `json:"flush_interval"`
	FlushAge      int    `json:"flush_age"`
}

type Data_d struct {
	Total     int        `json:"total"`
	IndexSets []IndexSet `json:"index_sets"`
	Stats     struct{}   `json:"stats"`
}

type IndexSet struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Index_prefix string `json:"index_prefix"`
	// Add more fields as needed
}

type ID_Index_sets struct {
	Indices   int `json:"indices"`
	Documents int `json:"documents"`
	Size      int `json:"size"`
}

type Sidecar_Count_num struct {
	Total int `json:"total"`
}

// MetricsConfig struct to hold configuration values for /metrics authentication
type MetricsConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// URLConfig struct to hold configuration values for API URL and authentication
type URLConfig struct {
	IP       string `yaml:"ip"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Sidecar_Total_Info struct {
	Query      string      `json:"query"`
	Total      int         `json:"total"`
	OnlyActive bool        `json:"only_active"`
	Sort       interface{} `json:"sort"`
	Order      interface{} `json:"order"`
	Sidecars   []Sidecars  `json:"sidecars"`
	Filters    interface{} `json:"filters"`
	Pagination Pagination  `json:"pagination"`
}
type Metrics struct {
	Disks75 []interface{} `json:"disks_75"`
	CPUIdle float64       `json:"cpu_idle"`
	Load1   float64       `json:"load_1"`
}
type Status struct {
	Status     int          `json:"status"`
	Message    string       `json:"message"`
	Collectors []Collectors `json:"collectors"`
}
type Collectors struct {
	CollectorID    string `json:"collector_id"`
	Status         int    `json:"status"`
	Message        string `json:"message"`
	VerboseMessage string `json:"verbose_message"`
}
type NodeDetails struct {
	OperatingSystem string      `json:"operating_system"`
	IP              string      `json:"ip"`
	Metrics         Metrics     `json:"metrics"`
	LogFileList     interface{} `json:"log_file_list"`
	Status          Status      `json:"status"`
}
type Sidecars struct {
	Active         bool          `json:"active"`
	NodeID         string        `json:"node_id"`
	NodeName       string        `json:"node_name"`
	NodeDetails    NodeDetails   `json:"node_details"`
	Assignments    []interface{} `json:"assignments"`
	LastSeen       time.Time     `json:"last_seen"`
	SidecarVersion string        `json:"sidecar_version"`
	Collectors     interface{}   `json:"collectors"`
}
type Pagination struct {
	Total   int `json:"total"`
	Count   int `json:"count"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

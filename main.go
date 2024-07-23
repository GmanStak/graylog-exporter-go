package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
)

func main() {
	auth_configFile := flag.String("c", "auth_config.yaml", "/metircs 访问认证文件路径（auth_config.yaml）")
	url_configFile := flag.String("a", "url_config.yaml", "graylog url 地址及认证路径（url_config.yaml）")
	port := flag.String("p", "8080", "/metircs 访问端口（:8080）")
	//port := flag.String("p", "8080", "Port number")
	flag.Parse()
	// Load /metrics configuration from config.yaml
	//metricsConfig, err := loadMetricsConfig("config.yaml")
	metricsConfig, err := loadMetricsConfig(*auth_configFile)
	if err != nil {
		log.Fatalf("Error loading /metrics config: %v", err)
	}

	// Load API URL configuration from url_config.yaml
	urlConfig, err := loadURLConfig(*url_configFile)
	if err != nil {
		log.Fatalf("Error loading URL config: %v", err)
	}

	// Expose Prometheus metrics endpoint with /metrics authentication
	http.Handle("/metrics", basicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Update metrics on each request
		err := updateMetrics(urlConfig)
		if err != nil {
			log.Printf("Error updating metrics: %v", err)
			http.Error(w, "Failed to fetch metrics", http.StatusInternalServerError)
			return
		}
		sidecar_num, err := Get_Sidecar_Num(urlConfig)
		if err != nil {
			log.Fatalf("Error getting data from API: %v", err)
		}
		graylog_sidecar_count_num.WithLabelValues(urlConfig.IP).Set(float64(sidecar_num.Total))

		journal_info, err := Get_Journal(urlConfig)
		journal := journal_info.UncommittedJournalEntries

		graylog_journal_entries.WithLabelValues(urlConfig.IP).Set(float64(journal))

		data, err := Get_Sidecar_Node_Status(urlConfig)
		if err != nil {
			log.Fatalf("Error getting data from API: %v", err)
		}
		//fmt.Println(data)
		//sidecar node status json info
		for _, sidecar := range data.Sidecars {
			node_name := sidecar.NodeName
			operating_system := sidecar.NodeDetails.OperatingSystem
			node_ip := sidecar.NodeID
			for _, collector := range sidecar.NodeDetails.Status.Collectors {
				message := collector.Message
				var runningStatus int
				if message == "Running" {
					runningStatus = 1
				} else {
					runningStatus = 0
				}
				graylog_sidecar_node_status.WithLabelValues(node_name, operating_system, node_ip, message).Set(float64(runningStatus))
			}
		}

		// Serve Prometheus metrics
		promhttp.Handler().ServeHTTP(w, r)
	}), metricsConfig.Username, metricsConfig.Password))

	log.Printf("Beginning to serve metrics on port :%s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

// basicAuth wraps a handler function and provides basic authentication
func basicAuth(handler http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

// loadMetricsConfig loads /metrics configuration from a YAML file
func loadMetricsConfig(filename string) (MetricsConfig, error) {
	var config MetricsConfig
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// loadURLConfig loads API URL configuration from a YAML file
func loadURLConfig(filename string) (URLConfig, error) {
	var config URLConfig
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// updateMetrics fetches data from Graylog API and updates Prometheus metrics
func updateMetrics(config URLConfig) error {
	indexURL := fmt.Sprintf("http://%s:%s/api/system/indices/index_sets", config.IP, config.Port)
	req, err := http.NewRequest("GET", indexURL, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(config.Username, config.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed: %s", resp.Status)
	}

	var data Data_d
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	// Reset Prometheus metrics before updating
	// Optionally, you can reset metrics here if needed

	for _, indexSet := range data.IndexSets {
		id := indexSet.ID
		indexPrefix := indexSet.Index_prefix
		I_data, err := Get_ID_Index_sets(config, id)
		if err != nil {
			log.Fatalf("Error getting data from API: %v", err)
		}
		graylog_index_count_num.WithLabelValues(id, indexPrefix).Set(float64(I_data.Documents))
		graylog_index_size_bytes.WithLabelValues(id, indexPrefix).Set(float64(I_data.Size))

	}

	log.Println("Metrics updated successfully")
	return nil
}

func Get_ID_Index_sets(config URLConfig, id string) (ID_Index_sets, error) {
	url := fmt.Sprintf("http://%s:%s/api/system/indices/index_sets/%s/stats", config.IP, config.Port, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ID_Index_sets{}, err
	}
	req.SetBasicAuth(config.Username, config.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ID_Index_sets{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ID_Index_sets{}, fmt.Errorf("HTTP request failed: %s", resp.Status)
	}

	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return "", err
	//}
	//
	//
	//data := string(body)
	//
	//return data, nil
	var data ID_Index_sets
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return ID_Index_sets{}, err
	}

	return data, nil
}

func Get_Sidecar_Num(config URLConfig) (Sidecar_Count_num, error) {
	url := fmt.Sprintf("http://%s:%s/api/sidecars?query=sssssss", config.IP, config.Port)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Sidecar_Count_num{}, err
	}
	req.SetBasicAuth(config.Username, config.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Sidecar_Count_num{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Sidecar_Count_num{}, fmt.Errorf("HTTP request failed: %s", resp.Status)
	}

	var data Sidecar_Count_num
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return Sidecar_Count_num{}, err
	}

	return data, nil

}

func Get_Sidecar_Node_Status(config URLConfig) (Sidecar_Total_Info, error) {
	url := fmt.Sprintf("http://%s:%s/api/sidecars/all", config.IP, config.Port)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Sidecar_Total_Info{}, err
	}
	req.SetBasicAuth(config.Username, config.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Sidecar_Total_Info{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Sidecar_Total_Info{}, fmt.Errorf("HTTP request failed: %s", resp.Status)
	}

	//返回请求string
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return Sidecar_Total_Info{}, err
	//}
	//
	//data := string(body)
	//
	//return data, nil

	//返回请求json
	var data Sidecar_Total_Info
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return Sidecar_Total_Info{}, err
	}

	return data, nil
}

func Get_Journal(config URLConfig) (Journal_info, error) {
	url := fmt.Sprintf("http://%s:%s/api/system/journal", config.IP, config.Port)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Journal_info{}, err
	}
	req.SetBasicAuth(config.Username, config.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Journal_info{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Journal_info{}, fmt.Errorf("HTTP request failed: %s", resp.Status)
	}

	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return "", err
	//}
	//
	//data := string(body)
	//
	//return data, nil
	var data Journal_info
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return Journal_info{}, err
	}

	return data, nil
}

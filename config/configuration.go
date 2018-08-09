package config

type Configuration struct {
	DatabaseConnectionURL     string
	StartMonitoring           string
	AlarmAt                   string
	AlarmAverageIntervals     []int
	SentenceAt                string
	PollingInterval           int
	MonitoringInterval        int
	MonitoringPollingInterval int
	GoEnv                     string
}

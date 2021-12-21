package types

// type ImageExporterConfig struct {
// 	RancherVersion      string   `yaml:"rancherVersion"`
// 	Targets             []string `yaml:"targets"`
// 	OsType              []string `yaml:"osType,omitempty"`
// 	ChartsRepoUrl       string   `yaml:"chartsRepoUrl"`
// 	ChartsBranch        string   `yaml:"chartsBranch"`
// 	SystemChartsRepoUrl string   `yaml:"systemChartsRepoUrl"`
// 	SystemChartsBranch  string   `yaml:"systemChartsBranch"`
// 	SystemChartsCommit  string   `yaml:"systemChartsCommit"`
// }
type ImageExporterConfig struct {
	RancherTag   string      `yaml:"rancherTag"`
	Targets      []string    `yaml:"targets"`
	Os           []string    `yaml:"os,omitempty"`
	Charts       ChartConfig `yaml:"charts"`
	Systemcharts ChartConfig `yaml:"systemcharts"`
	RKE2         RKE2Config  `yaml:"rke2"`
}

type ChartConfig struct {
	Repository `yaml:"repository"`
}

type Repository struct {
	Location string `yaml:"location"`
	Branch   string `yaml:"branch"`
	Commit   string `yaml:"commit"`
}

type RKE2Config struct {
	KDMData string `yaml:"kdmData"`
}

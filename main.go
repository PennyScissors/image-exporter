// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"

// 	"github.com/davecgh/go-spew/spew"
// 	"github.com/pennyscissors/image-exporter/pkg/charts"
// 	"github.com/pennyscissors/image-exporter/pkg/export"
// 	"github.com/pennyscissors/image-exporter/pkg/systemcharts"
// 	"github.com/pennyscissors/image-exporter/pkg/types"
// 	"github.com/sirupsen/logrus"
// 	cli "github.com/urfave/cli/v2"
// 	yaml "gopkg.in/yaml.v2"
// )

// const (
// 	Linux   = "Linux"
// 	Windows = "Windows"
// )

// // type OSType int

// // const (
// // 	Linux OSType = iota
// // 	Windows
// // )

// // func (o OSType) String() string {
// // 	return OsToString[o]
// // }

// // var OsToString = map[OSType]string{
// // 	Linux:   "Linux",
// // 	Windows: "Windows",
// // }

// const (
// 	Test                           = "Test"
// 	DefaultImageExporterConfigFile = "config.yaml"
// 	DefaultTarget                  = "all"
// 	DefaultRancherTag              = "v2.6.3"
// 	TargetEnvVar                   = "TARGET"
// 	RancherVersionEnvVar           = "TAG"
// )

// var (
// 	// Version represents the current version of the chart build scripts
// 	Version = "0.0.0-dev"
// 	// GitCommit represents the latest commit when building this script
// 	GitCommit = "HEAD"

// 	ImageExporterConfigFile string
// 	CurrentTarget           string
// 	CurrentRancherTag       string
// )

// func main() {
// 	app := cli.NewApp()
// 	app.Name = "image-exporter"
// 	app.Version = fmt.Sprintf("%s (%s)", Version, GitCommit)
// 	app.Usage = "Tool used to export airgap images of a given target component and Rancher version"
// 	configFlag := cli.StringFlag{
// 		Name:        "config",
// 		Aliases:     []string{"c"},
// 		Usage:       "Configuration file",
// 		TakesFile:   false,
// 		Destination: &ImageExporterConfigFile,
// 		Value:       DefaultImageExporterConfigFile,
// 		DefaultText: "config.yaml",
// 	}
// 	tagFlag := cli.StringFlag{
// 		Name:        "tag",
// 		Aliases:     []string{"t"},
// 		Usage:       "Rancher tag to export images",
// 		Required:    false,
// 		Destination: &CurrentRancherTag,
// 		Value:       DefaultRancherTag,
// 		// EnvVars:     []string{RancherVersionEnvVar},
// 		DefaultText: "DefaultRancherTag",
// 	}
// 	app.Commands = []*cli.Command{
// 		{
// 			Name:   "export",
// 			Usage:  "Export images for a given target",
// 			Action: exportAllImages,
// 			Flags: []cli.Flag{
// 				&tagFlag,
// 				&configFlag,
// 			},
// 			Subcommands: []*cli.Command{
// 				// {
// 				// 	Name:   "charts",
// 				// 	Usage:  "charts usage here",
// 				// 	Action: exportChartsImages,
// 				// 	Flags: []cli.Flag{
// 				// 		&tagFlag,
// 				// 		&configFlag,
// 				// 	},
// 				// },
// 				{
// 					Name:   "systemcharts",
// 					Usage:  "system charts usage here",
// 					Action: exportSystemChartsImages,
// 					Flags: []cli.Flag{
// 						&tagFlag,
// 						&configFlag,
// 					},
// 				},
// 				{
// 					Name:   "systemc",
// 					Usage:  "system charts usage here",
// 					Action: exportSystemChartsImages2,
// 					Flags: []cli.Flag{
// 						&tagFlag,
// 						&configFlag,
// 					},
// 				},
// 			},
// 		},
// 	}
// 	if err := app.Run(os.Args); err != nil {
// 		logrus.Fatal(err)
// 	}
// }

// // type ImageExporterConfig struct {
// // 	RancherVersion      string   `yaml:"rancherVersion"`
// // 	Targets             []string `yaml:"targets"`
// // 	OsType              []string `yaml:"osType,omitempty"`
// // 	ChartsRepoUrl       string   `yaml:"chartsRepoUrl"`
// // 	ChartsBranch        string   `yaml:"chartsBranch"`
// // 	SystemChartsRepoUrl string   `yaml:"systemChartsRepoUrl"`
// // 	SystemChartsBranch  string   `yaml:"systemChartsBranch"`
// // }

// // func exportImages(c *cli.Context) error {
// // 	fmt.Println("getImages()")
// // 	config := parseConfigFile()
// // 	rancherVersion, err := semver.NewVersion(config.RancherVersion)
// // 	if err != nil {
// // 		logrus.Fatalf("Unable to parse rancher version string: %s", err)
// // 	}
// // 	fmt.Printf("target: %v %v %v\n", config.Target, CurrentTarget, strings.EqualFold(config.Target, CurrentTarget))
// // 	fmt.Printf("rancher semver: %v\n", rancherVersion)
// // 	return nil
// // }

// func exportAllImages(c *cli.Context) error {
// 	config := parseConfigFile()
// 	export.GetAllImages(config)
// 	return nil
// }

// func exportChartsImages(c *cli.Context) error {
// 	parseConfigFile()
// 	charts.GetChartsImages()
// 	return nil
// }

// func exportSystemChartsImages(c *cli.Context) error {
// 	config := parseConfigFile()
// 	systemcharts.GetSystemChartsImages(config)
// 	return nil
// }

// type Repository struct {
// 	Url    string `yaml:"url"`
// 	Branch string `yaml:"branch"`
// 	Commit string `yaml:"commit"`
// }

// type ChartConfig struct {
// 	Repository `yaml:"repository"`
// }

// type Configg struct {
// 	Systemcharts ChartConfig `yaml:"systemcharts"`
// }

// func exportSystemChartsImages2(c *cli.Context) error {
// 	_ = parseConfigFile2()
// 	// systemcharts.GetSystemChartsImages(config)
// 	return nil
// }

// func parseConfigFile() types.ImageExporterConfig {
// 	configYaml, err := ioutil.ReadFile(ImageExporterConfigFile)
// 	if err != nil {
// 		logrus.Fatalf("Unable to find configuration file: %s", err)
// 	}
// 	config := types.ImageExporterConfig{}
// 	if err := yaml.UnmarshalStrict(configYaml, &config); err != nil {
// 		logrus.Fatalf("Unable to unmarshall configuration file: %s", err)
// 	}
// 	spew.Dump(config)
// 	return config
// }

// func parseConfigFile2() Configg {
// 	configYaml, err := ioutil.ReadFile(ImageExporterConfigFile)
// 	if err != nil {
// 		logrus.Fatalf("Unable to find configuration file: %s", err)
// 	}
// 	config := Configg{}
// 	if err := yaml.UnmarshalStrict(configYaml, &config); err != nil {
// 		logrus.Fatalf("Unable to unmarshall configuration file: %s", err)
// 	}
// 	spew.Dump(config)
// 	return config
// }

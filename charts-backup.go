<<<<<<< Updated upstream
=======
package image

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	libhelm "github.com/rancher/rancher/pkg/helm"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/repo"
)

const RancherVersionAnnotationKey = "catalog.cattle.io/rancher-version"

type Charts struct {
	config *ExportConfig
}

func (c Charts) fetchImages(imagesSet map[string]map[string]struct{}) error {
	if c.config.chartsPath == "" {
		return errors.New("repository path undefined")
	}
	// Load index.yaml file
	index, err := repo.LoadIndexFile(filepath.Join(c.config.chartsPath, "index.yaml"))
	if err != nil {
		return err
	}
	// Filter index entries based on their Rancher version constraint
	// Selecting the correct latest heavily relies on the charts-build-scripts `make standardize` command sorting versions in the index correctly.
	var filteredVersions repo.ChartVersions
	for _, versions := range index.Entries {
		if len(versions) >= 1 {
			// Always append the latest version of a chart
			latestVersion := versions[0]
			filteredVersions = append(filteredVersions, latestVersion)
		}
		if len(versions) > 1 {
			// If there is more than one version, append them if their Rancher version constraint
			// satisfies the given Rancher version or tag
			for _, version := range versions[1:] {
				// logrus.Infof("  version: %s:%s", version.Name, version.Version)
				isConstraintSatisfied, err := c.checkChartVersionConstraint(*version)
				if err != nil {
					return errors.Wrapf(err, "failed to check constraint of chart")
				}
				if isConstraintSatisfied {
					filteredVersions = append(filteredVersions, version)
				}
			}
		}
	}
	//
	for _, version := range filteredVersions {
		tgzPath := filepath.Join(c.config.chartsPath, version.URLs[0])
		// Find values.yaml files in tgz
		versionValues, err := decodeValuesFilesInTgz(tgzPath)
		if err != nil {
			logrus.Info(err)
			continue
		}
		chartNameAndVersion := fmt.Sprintf("%s:%s", version.Name, version.Version)
		for _, values := range versionValues {
			// Walk values.yaml and add images to set
			err = pickImagesFromValuesMap(imagesSet, values, chartNameAndVersion, c.config.osType)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c Charts) checkChartVersionConstraint(version repo.ChartVersion) (bool, error) {
	constraintStr, ok := version.Annotations[RancherVersionAnnotationKey]
	if !ok {
		// Log a warning when a chart doesn't have the rancher-version annotation, but return true so that images are exported.
		// logrus.Warnf("chart: %s:%s does not have a %s annotation defined", version.Name, version.Version, RancherVersionAnnotationKey)
		return true, nil
	}
	isConstraintSatisfied, err := compareRancherVersionToConstraint(c.config.rancherVersion, constraintStr)
	if err != nil {
		return false, err
	}
	return isConstraintSatisfied, nil
}

type SystemCharts struct {
	config *ExportConfig
}

type Questions struct {
	RancherMinVersion string `yaml:"rancher_min_version"`
	RancherMaxVersion string `yaml:"rancher_max_version"`
}

func (sc SystemCharts) fetchImages(imagesSet map[string]map[string]struct{}) error {
	if sc.config.systemChartsPath == "" {
		return errors.New("repository path undefined")
	}
	// Load system charts virtual index
	helm := libhelm.Helm{
		LocalPath: sc.config.systemChartsPath,
		IconPath:  sc.config.systemChartsPath,
		Hash:      "",
	}
	virtualIndex, err := helm.LoadIndex()
	if err != nil {
		return errors.Wrapf(err, "failed to load system charts index")
	}
	// Filter index entries based on their Rancher version constraint
	var filteredVersions libhelm.ChartVersions
	for _, versions := range virtualIndex.IndexFile.Entries {
		if len(versions) >= 1 {
			// Always append the latest version of a chart
			latestVersions := versions[0]
			filteredVersions = append(filteredVersions, latestVersions)
		}
		if len(versions) > 1 {
			// If there is more than one version, append them if their Rancher version constraint
			// satisfies the given Rancher version or tag
			for _, version := range versions[1:] {
				isConstraintSatisfied, err := sc.checkChartVersionConstraint(*version)
				if err != nil {
					return errors.Wrapf(err, "failed to filter chart versions")
				}
				if isConstraintSatisfied {
					filteredVersions = append(filteredVersions, version)
				}
			}
		}
	}
	// Iterate through versions and check for images in their values files
	for _, version := range filteredVersions {
		for _, file := range version.LocalFiles {
			if !isValuesFile(file) {
				continue
			}
			// logrus.Infof("file: %v", file)
			values, err := decodeValuesFile(file)
			if err != nil {
				return err
			}
			chartNameAndVersion := fmt.Sprintf("%s:%s", version.Name, version.Version)
			err = pickImagesFromValuesMap(imagesSet, values, chartNameAndVersion, sc.config.osType)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sc SystemCharts) checkChartVersionConstraint(version libhelm.ChartVersion) (bool, error) {
	questionsPath := filepath.Join(sc.config.systemChartsPath, version.Dir, "questions.yaml")
	questions, err := decodeQuestionsFile(questionsPath)
	if os.IsNotExist(err) {
		questionsPath = filepath.Join(sc.config.systemChartsPath, version.Dir, "questions.yml")
		questions, err = decodeQuestionsFile(questionsPath)
	}
	if err != nil {
		return false, err
	}
	constraintStr := minMaxToConstraintStr(questions.RancherMinVersion, questions.RancherMaxVersion)
	if constraintStr == "" {
		// Log a warning and export images when a chart doesn't have rancher version constraints in its questions file
		logrus.Warnf("system chart: %s does not have a rancher_min_version or rancher_max_version constraint defined in its questions file", version.Name)
		return true, nil
	}
	isConstraintSatisfied, err := compareRancherVersionToConstraint(sc.config.rancherVersion, constraintStr)
	if err != nil {
		return false, err
	}
	return isConstraintSatisfied, nil
}

// isRancherVersionInConstraintRange returns true if the rancher server version satisfies a given constraint (E.g ">=2.5.0 <=2.6"), false otherwise.
func compareRancherVersionToConstraint(rancherVersion, constraintStr string) (bool, error) {
	if constraintStr == "" {
		return false, errors.Errorf("Invalid constraint string: \"%s\"", constraintStr)
	}
	v, err := semver.NewVersion(rancherVersion)
	if err != nil {
		return false, err
	}
	// Decrease patch to 98 so that semver comparison for rancher-versin constraint annotations work
	patch := v.Patch()
	if patch == 99 {
		patch = 98
	}
	// Create rancher semver without pre-release
	// Remove the pre-release because the semver package will not consider a rancherVersion with a
	// pre-release unless the versions in the constraintStr has pre-releases as well.
	// For example: rancherVersion "2.6.4-rc1" and constraint "2.6.3 - 2.6.5" will return false because
	// there is no pre-release in the constraint "2.5.6 - 2.5.8" (This behavior is intentional).
	rancherSemVer, err := semver.NewVersion(fmt.Sprintf("%d.%d.%d", v.Major(), v.Minor(), patch))
	if err != nil {
		return false, err
	}
	constraint, err := semver.NewConstraint(constraintStr)
	if err != nil {
		return false, err
	}
	return constraint.Check(rancherSemVer), nil
}

// minMaxToConstraintStr converts min and max rancher version strings into a constraint string.
func minMaxToConstraintStr(min, max string) string {
	if min != "" && max != "" {
		return fmt.Sprintf("%s - %s", min, max)
	}
	if min != "" {
		return fmt.Sprintf(">= %s", min)
	}
	if max != "" {
		return fmt.Sprintf("<= %s", max)
	}
	return ""
}

// pickImagesFromValuesMap walks a values map to find images of a given OS type and adds them to imagesSet.
func pickImagesFromValuesMap(imagesSet map[string]map[string]struct{}, values map[interface{}]interface{}, chartNameAndVersion string, osType OSType) error {
	walkMap(values, func(inputMap map[interface{}]interface{}) {
		repository, ok := inputMap["repository"].(string)
		if !ok {
			return
		}
		tag, ok := inputMap["tag"].(string)
		if !ok {
			return
		}
		imageName := fmt.Sprintf("%s:%v", repository, tag)
		// By default, images are added to the generic images list ("linux"). For Windows and multi-OS
		// images to be considered, they must use a comma-delineated list (e.g. "os: windows",
		// "os: windows,linux", and "os: linux,windows").
		osList, ok := inputMap["os"].(string)
		if !ok {
			if inputMap["os"] != nil {
				errors.Errorf("field 'os:' for image %s contains neither a string nor nil", imageName)
			}
			if osType == Linux {
				addSourceToImage(imagesSet, imageName, chartNameAndVersion)
				return
			}
		}
		for _, os := range strings.Split(osList, ",") {
			os = strings.TrimSpace(os)
			if strings.EqualFold("windows", os) && osType == Windows {
				addSourceToImage(imagesSet, imageName, chartNameAndVersion)
				return
			}
			if strings.EqualFold("linux", os) && osType == Linux {
				addSourceToImage(imagesSet, imageName, chartNameAndVersion)
				return
			}
		}
	})
	return nil
}

// decodeValueFilesInTgz reads in given path to find values.yaml files and returns a slice of decoded values.yaml files.
func decodeValuesFilesInTgz(tgzPath string) ([]map[interface{}]interface{}, error) {
	tgz, err := os.Open(tgzPath)
	if err != nil {
		return nil, err
	}
	defer tgz.Close()
	gzr, err := gzip.NewReader(tgz)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	var valuesSlice []map[interface{}]interface{}
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return valuesSlice, nil
			// return any other error
		case err != nil:
			return nil, err
		case header.Typeflag == tar.TypeReg && isValuesFile(header.Name):
			var values map[interface{}]interface{}
			if err := decodeYAMLFile(tr, &values); err != nil {
				return nil, err
			}
			valuesSlice = append(valuesSlice, values)
		default:
			continue
		}
	}
}

// walkMap walks a map and executes the given walk function for each node.
func walkMap(data interface{}, walkFunc func(map[interface{}]interface{})) {
	if inputMap, isMap := data.(map[interface{}]interface{}); isMap {
		// Run the walkFunc on the root node and each child node
		walkFunc(inputMap)
		for _, value := range inputMap {
			walkMap(value, walkFunc)
		}
	} else if inputList, isList := data.([]interface{}); isList {
		// Run the walkFunc on each element in the root node, ignoring the root itself
		for _, elem := range inputList {
			walkMap(elem, walkFunc)
		}
	}
}

func decodeQuestionsFile(path string) (Questions, error) {
	var questions Questions
	file, err := os.Open(path)
	if err != nil {
		return Questions{}, err
	}
	defer file.Close()
	if err := decodeYAMLFile(file, &questions); err != nil {
		return Questions{}, err
	}
	return questions, nil
}

func decodeValuesFile(path string) (map[interface{}]interface{}, error) {
	var values map[interface{}]interface{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := decodeYAMLFile(file, &values); err != nil {
		return nil, err
	}
	return values, nil
}

func decodeYAMLFile(r io.Reader, target interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, target)
}

func isValuesFile(path string) bool {
	basename := filepath.Base(path)
	return basename == "values.yaml" || basename == "values.yml"
}
>>>>>>> Stashed changes

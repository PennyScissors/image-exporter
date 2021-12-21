


Commands:

Targets:
All - Default
charts
system-charts
core
rancher
rke/kdm/system
rke2
k3s


./image-exporter export -t charts -c config.yaml
./image-exporter export charts -c config.yaml

// getChartsImages(tag, arch, chartsRepoUrl, chartsRepoBranch)

remove the need to switch the branch as we can have a configuration specific for release/dev

kdmDataPath: 

os info:
charts - inside values.yaml



Images:
- os: can be windows or linux
- source: can come from a specific chart, system, core, rke2, k3s, etc


type Image struct {
    Repository string
    Tag string
    // Source of the image. E.g chart, system, rke2All, k3s upgrade, rancher
    Sources []string
    OS []string
}
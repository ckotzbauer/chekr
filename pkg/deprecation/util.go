package deprecation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	semver "github.com/Masterminds/semver/v3"
	"github.com/ckotzbauer/chekr/pkg/util"
	"github.com/sirupsen/logrus"
)

var ignoredDeprecatedKinds = []string{"Event"}

func fetchDeprecatedApis() []GroupVersion {
	file, err := ioutil.TempFile("", "k8s_deprecation")

	if err != nil {
		logrus.WithError(err).Fatalf("Could not create temp-file!")
	}

	err = util.DownloadFile(
		file.Name(),
		"https://raw.githubusercontent.com/ckotzbauer/chekr/main/data/k8s_deprecations_generated.json")

	if err != nil {
		logrus.WithError(err).Fatalf("Could not download deprecation-definition!")
	}

	buf, err := ioutil.ReadFile(file.Name())

	if err != nil {
		logrus.WithError(err).Fatalf("Could not read deprecation-definition!")
	}

	data := []GroupVersion{}
	err = json.Unmarshal([]byte(buf), &data)

	if err != nil {
		logrus.WithError(err).Fatalf("Could not unmarshal deprecation-definition!")
	}

	return data
}

func (d Deprecation) isVersionIgnored(deprecation string) bool {
	if d.K8sVersion == "" {
		info, err := d.KubeClient.Client.ServerVersion()

		if err == nil && info != nil {
			d.K8sVersion = fmt.Sprintf("%s.%s", info.Major, info.Minor)
		}
	}

	c, err := semver.NewConstraint("< " + deprecation)

	if err != nil {
		logrus.WithError(err).WithField("constraint", "< "+deprecation).Fatalf("Could not parse version-constraint!")
	}

	v, err := semver.NewVersion(d.K8sVersion)

	if err != nil {
		logrus.WithError(err).WithField("version1", deprecation).WithField("version2", d.K8sVersion).Fatalf("Could not compare versions!")
	}

	return c.Check(v)
}

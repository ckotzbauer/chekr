package deprecation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/ckotzbauer/chekr/pkg/util"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var kyvernoMessage = "{{ request.object.apiVersion }}/{{ request.object.kind }} is deprecated and will be removed in %s."

type Version struct {
	Version string
	Message string
	Name    string
	Kinds   []string
}

type Container struct {
	Versions                []Version
	Category                string
	Subject                 string
	ValidationFailureAction string
	Background              bool
}

func (d Deprecation) ExecuteKyvernoCreate() string {
	apis := fetchDeprecatedApis()

	rest.SetDefaultWarningHandler(rest.NoWarnings{})
	return d.createKyvernoPolicies(apis)
}

func (d Deprecation) createKyvernoPolicies(apis []GroupVersion) string {
	m := make(map[string][]string)
	ignoredDeprecatedKinds = append(ignoredDeprecatedKinds, d.IgnoredKinds...)

	for _, a := range apis {
		for _, r := range a.Resources {
			if util.Contains(ignoredDeprecatedKinds, r.Name) {
				continue
			}

			if x, found := m[r.Removed]; found {
				x = append(x, fmt.Sprintf("%s/%s/%s", a.Group, a.Version, r.Name))
				m[r.Removed] = x
			} else {
				m[r.Removed] = []string{fmt.Sprintf("%s/%s/%s", a.Group, a.Version, r.Name)}
			}
		}
	}

	container := Container{
		Versions:                []Version{},
		Category:                d.Category,
		Subject:                 d.Subject,
		ValidationFailureAction: d.ValidationFailureAction,
		Background:              d.Background,
	}

	for s, v := range m {
		if s == "" {
			s = "the future"
		}

		if s != "the future" {
			if d.isVersionIgnored(s) {
				continue
			}
		} else if d.K8sVersion != "" {
			continue
		}

		container.Versions = append(container.Versions, Version{
			Version: s,
			Message: fmt.Sprintf(kyvernoMessage, s),
			Name:    fmt.Sprintf("validate-%s-removals", strings.Replace(s, ".", "-", -1)),
			Kinds:   v,
		})
	}

	buf := new(bytes.Buffer)
	tpl := template.New("page")
	tpl, err := tpl.Parse(KyvernoTemplate)

	if err != nil {
		logrus.WithError(err).Fatalf("Could not parse policy-template!")
	}

	err = tpl.Execute(buf, container)

	if err != nil {
		logrus.WithError(err).Fatalf("Could not render policy-template!")
	}

	return buf.String()
}

func (d Deprecation) HandleKyvernoResult(stringOutput, output, outputFile string, dryRun bool) {
	if dryRun {
		if output == "json" {
			data, _ := parseYaml(stringOutput)
			printResult(string(data), outputFile)
			return
		}

		printResult(stringOutput, outputFile)
	} else {
		d.applyPolicy(stringOutput)
	}
}

func printResult(str, outputFile string) {
	if outputFile != "" {
		data := []byte(str)
		err := os.WriteFile(outputFile, data, 0640)

		if err != nil {
			logrus.WithError(err).Fatalf("Could write output to file!")
		}
	} else {
		fmt.Fprint(os.Stdout, str)
	}
}

func parseYaml(stringOutput string) ([]byte, *unstructured.Unstructured) {
	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, _, err := decUnstructured.Decode([]byte(stringOutput), nil, obj)
	if err != nil {
		logrus.WithError(err).Fatalf("Could not decode policy-string.")
	}

	data, err := json.Marshal(obj)
	if err != nil {
		logrus.WithError(err).Fatalf("Could not marshal object to json.")
	}

	return data, obj
}

func (d Deprecation) applyPolicy(stringOutput string) {
	data, obj := parseYaml(stringOutput)

	dyn, _ := dynamic.NewForConfig(d.KubeClient.Config)
	dr := dyn.Resource(schema.GroupVersionResource{Group: "kyverno.io", Version: "v1", Resource: "clusterpolicies"})

	_, err := dr.Patch(context.Background(), obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
		FieldManager: "chekr",
	})

	if err != nil {
		logrus.WithError(err).Fatalf("Could not apply policy.")
	} else {
		fmt.Printf("clusterpolicy/%s created/updated\n", obj.GetName())
	}
}

func (d Deprecation) DeletePolicy() {
	dyn, _ := dynamic.NewForConfig(d.KubeClient.Config)
	dr := dyn.Resource(schema.GroupVersionResource{Group: "kyverno.io", Version: "v1", Resource: "clusterpolicies"})

	name := "chekr-check-deprecated-apis"
	err := dr.Delete(context.Background(), name, metav1.DeleteOptions{})

	if err != nil && !apierrors.IsNotFound(err) {
		logrus.WithError(err).Fatalf("Could not delete policy.")
	} else {
		fmt.Printf("clusterpolicy/%s deleted\n", name)
	}
}

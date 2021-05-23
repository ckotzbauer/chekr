package deprecation

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	semver "github.com/Masterminds/semver/v3"
	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/ckotzbauer/chekr/pkg/util"
	"github.com/ddelizia/channelify"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var ignoredDeprecatedKinds = []string{"Event"}

func (d Deprecation) Execute() (printer.PrintableList, error) {
	apis, err := fetchDeprecatedApis()

	if err != nil {
		return nil, err
	}

	rest.SetDefaultWarningHandler(rest.NoWarnings{})
	return d.findDeprecatedVersions(apis)
}

func fetchDeprecatedApis() ([]GroupVersion, error) {
	file, err := ioutil.TempFile("", "k8s_deprecation")

	if err != nil {
		return nil, err
	}

	err = util.DownloadFile(
		file.Name(),
		"https://raw.githubusercontent.com/ckotzbauer/chekr/master/data/k8s_deprecations_generated.json")

	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadFile(file.Name())

	if err != nil {
		return nil, err
	}

	data := []GroupVersion{}
	err = json.Unmarshal([]byte(buf), &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (d Deprecation) findDeprecatedVersions(deprecatedGVs []GroupVersion) (DeprecatedResourceList, error) {
	kindVersions := d.KubeClient.DiscoverResourceNameAndPreferredGV()
	d.KubeClient.Config.Burst = d.ThrottleBurst
	dynamicClient, err := dynamic.NewForConfig(d.KubeClient.Config)
	ignoredDeprecatedKinds = append(ignoredDeprecatedKinds, d.IgnoredKinds...)

	if err != nil {
		return DeprecatedResourceList{}, err
	}

	fn1 := func(
		d Deprecation,
		deprecatedGV GroupVersion,
		deprecatedGVR Resource,
		client dynamic.Interface,
		kindVersions kubernetes.KindVersions) printer.PrintableResult {

		return d.analyzeDeprecatedResource(deprecatedGV, deprecatedGVR, dynamicClient, kindVersions)
	}

	ch1 := channelify.Channelify(fn1)
	var channels [](chan printer.PrintableResult)
	deprecatedList := DeprecatedResourceList{Items: []DeprecatedResource{}}

	for _, deprecatedGV := range deprecatedGVs {
		for _, deprecatedGVR := range deprecatedGV.Resources {
			ignored, err := d.isVersionIgnored(deprecatedGVR.Deprecated)

			if err != nil {
				return DeprecatedResourceList{}, err
			}

			if ignored || util.Contains(ignoredDeprecatedKinds, deprecatedGVR.Name) {
				continue
			}

			if _, ok := kindVersions[deprecatedGVR.Name]; !ok {
				// This deprecated kind does not exist anymore "Deleted"
				continue
			}

			ch := ch1.(func(Deprecation, GroupVersion, Resource, dynamic.Interface, kubernetes.KindVersions) chan printer.PrintableResult)(d, deprecatedGV, deprecatedGVR, dynamicClient, kindVersions)
			channels = append(channels, ch)
		}
	}

	for _, v := range channels {
		result := <-v

		if result.Error != nil {
			return DeprecatedResourceList{}, result.Error
		}

		if result.Item == nil {
			continue
		}

		deprecatedList.Items = append(deprecatedList.Items, result.Item.([]DeprecatedResource)...)
	}

	return deprecatedList, nil
}

func (d Deprecation) analyzeDeprecatedResource(
	deprecatedGV GroupVersion,
	deprecatedGVR Resource,
	client dynamic.Interface,
	kindVersions kubernetes.KindVersions) printer.PrintableResult {

	deprecated := make([]DeprecatedResource, 0)
	supportedVersions := kindVersions[deprecatedGVR.Name]
	supported := findSupportedKindVersion(supportedVersions, deprecatedGV.Version)

	if supported.Version == "" {
		// This deprecated version does not exist anymore "Deleted"
		return printer.PrintableResult{Item: deprecated}
	}

	gvr := schema.GroupVersionResource{Group: deprecatedGV.Group, Version: deprecatedGV.Version, Resource: supported.Name}
	deprecatedItems, err := client.Resource(gvr).List(context.TODO(), metav1.ListOptions{})

	if apierrors.IsNotFound(err) || apierrors.IsMethodNotSupported(err) {
		return printer.PrintableResult{Item: deprecated}
	}

	if apierrors.IsForbidden(err) {
		log.Fatalf("Failed to list objects in the cluster. Permission denied! Please check if you have the proper authorization")
		return printer.PrintableResult{Item: deprecated}
	}

	if err != nil {
		log.Fatalf("Failed communicating with k8s while listing objects [%v] %v. \nError: %v", gvr.String(), apierrors.ReasonForError(err), err)
	}

	if len(deprecatedItems.Items) > 0 {
		replacement := deprecatedGVR.Replacement

		if replacement.Group != "" {
			replacementGvr := schema.GroupVersionResource{Group: replacement.Group, Version: replacement.Version, Resource: supported.Name}
			replacementItems, err := client.Resource(replacementGvr).List(context.TODO(), metav1.ListOptions{})

			if apierrors.IsNotFound(err) {
				return printer.PrintableResult{Item: deprecated}
			}

			if apierrors.IsForbidden(err) {
				log.Fatalf("Failed to list objects in the cluster. Permission denied! Please check if you have the proper authorization")
			}

			if err != nil {
				log.Fatalf("Failed communicating with k8s while listing objects. \nError: %v", err)
			}

			for _, deprecatedItem := range deprecatedItems.Items {
				if !existsReplacementItem(deprecatedItem.GetNamespace(), deprecatedItem.GetName(), replacementItems) {
					deprecated = append(deprecated, createDeprecationItem(&deprecatedItem, deprecatedGVR))
				}
			}
		} else {
			for _, deprecatedItem := range deprecatedItems.Items {
				deprecated = append(deprecated, createDeprecationItem(&deprecatedItem, deprecatedGVR))
			}
		}
	}

	return printer.PrintableResult{Item: deprecated}
}

func (d Deprecation) isVersionIgnored(deprecation string) (bool, error) {
	if d.K8sVersion == "" {
		return false, nil
	}

	c, err := semver.NewConstraint("< " + deprecation)

	if err != nil {
		return false, err
	}

	v, err := semver.NewVersion(d.K8sVersion)

	if err != nil {
		return false, err
	}

	return c.Check(v), nil
}

func createDeprecationItem(deprecatedItem *unstructured.Unstructured, metadata Resource) DeprecatedResource {
	var deprecatedGv string
	var replacementGv string

	if deprecatedItem.GroupVersionKind().Group == "" {
		deprecatedGv = deprecatedItem.GroupVersionKind().Version
	} else {
		deprecatedGv = deprecatedItem.GroupVersionKind().Group + "/" + deprecatedItem.GroupVersionKind().Version
	}

	if metadata.Replacement.Version != "" {
		if metadata.Replacement.Group == "" {
			replacementGv = metadata.Replacement.Version
		} else {
			replacementGv = metadata.Replacement.Group + "/" + metadata.Replacement.Version
		}
	}

	return DeprecatedResource{
		Namespace:               deprecatedItem.GetNamespace(),
		Name:                    deprecatedItem.GetName(),
		DeprecatedGroupVersion:  deprecatedGv,
		DeprecatedKind:          deprecatedItem.GetKind(),
		ReplacementGroupVersion: replacementGv,
		ReplacementKind:         metadata.Replacement.Name,
		DeprecationVersion:      metadata.Deprecated + "+",
		RemovalVersion:          metadata.Removed + "+",
	}
}

func existsReplacementItem(namespace, name string, replacementItems *unstructured.UnstructuredList) bool {
	for _, replacement := range replacementItems.Items {
		if fmt.Sprintf("%s/%s", namespace, name) == fmt.Sprintf("%s/%s", replacement.GetNamespace(), replacement.GetName()) {
			return true
		}
	}

	return false
}

func findSupportedKindVersion(versions []kubernetes.KindVersion, deprecatedVersion string) kubernetes.KindVersion {
	for _, v := range versions {
		if v.Version == deprecatedVersion {
			return v
		}
	}

	return kubernetes.KindVersion{}
}

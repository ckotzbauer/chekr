package deprecation

import (
	"context"
	"fmt"

	"github.com/ckotzbauer/chekr/pkg/kubernetes"
	"github.com/ckotzbauer/chekr/pkg/printer"
	"github.com/ckotzbauer/chekr/pkg/util"
	"github.com/ddelizia/channelify"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

func (d Deprecation) ExecuteList() printer.PrintableList {
	apis := fetchDeprecatedApis()

	rest.SetDefaultWarningHandler(rest.NoWarnings{})
	return d.findDeprecatedVersions(apis)
}

func (d Deprecation) findDeprecatedVersions(deprecatedGVs []GroupVersion) DeprecatedResourceList {
	kindVersions := d.KubeClient.DiscoverResourceNameAndPreferredGV()
	d.KubeClient.Config.Burst = d.ThrottleBurst
	dynamicClient, err := dynamic.NewForConfig(d.KubeClient.Config)
	ignoredDeprecatedKinds = append(ignoredDeprecatedKinds, d.IgnoredKinds...)

	if err != nil {
		logrus.WithError(err).Fatalf("Could not create Kubernetes discovery client!")
	}

	fn1 := func(
		d Deprecation,
		deprecatedGV GroupVersion,
		deprecatedGVR Resource,
		client dynamic.Interface,
		kindVersions kubernetes.KindVersions) printer.Printable {

		return d.analyzeDeprecatedResource(deprecatedGV, deprecatedGVR, dynamicClient, kindVersions)
	}

	ch1 := channelify.Channelify(fn1)
	var channels [](chan printer.Printable)
	deprecatedList := DeprecatedResourceList{Items: []DeprecatedResource{}}

	for _, deprecatedGV := range deprecatedGVs {
		for _, deprecatedGVR := range deprecatedGV.Resources {
			ignored := d.isVersionIgnored(deprecatedGVR.Deprecated)

			if ignored || util.Contains(ignoredDeprecatedKinds, deprecatedGVR.Name) {
				continue
			}

			if _, ok := kindVersions[deprecatedGVR.Name]; !ok {
				// This deprecated kind does not exist anymore "Deleted"
				continue
			}

			ch := ch1.(func(Deprecation, GroupVersion, Resource, dynamic.Interface, kubernetes.KindVersions) chan printer.Printable)(d, deprecatedGV, deprecatedGVR, dynamicClient, kindVersions)
			channels = append(channels, ch)
		}
	}

	for _, v := range channels {
		result := <-v

		if result == nil {
			continue
		}

		deprecatedList.Items = append(deprecatedList.Items, result.([]DeprecatedResource)...)
	}

	return deprecatedList
}

func (d Deprecation) analyzeDeprecatedResource(
	deprecatedGV GroupVersion,
	deprecatedGVR Resource,
	client dynamic.Interface,
	kindVersions kubernetes.KindVersions) printer.Printable {

	deprecated := make([]DeprecatedResource, 0)
	supportedVersions := kindVersions[deprecatedGVR.Name]
	supported := findSupportedKindVersion(supportedVersions, deprecatedGV.Version)

	if supported.Version == "" {
		// This deprecated version does not exist anymore "Deleted"
		return deprecated
	}

	gvr := schema.GroupVersionResource{Group: deprecatedGV.Group, Version: deprecatedGV.Version, Resource: supported.Name}
	deprecatedItems, err := client.Resource(gvr).List(context.TODO(), metav1.ListOptions{})

	if apierrors.IsNotFound(err) || apierrors.IsMethodNotSupported(err) {
		return deprecated
	}

	if apierrors.IsForbidden(err) {
		logrus.WithError(err).Fatalf("Failed to list objects in the cluster. Permission denied! Please check if you have the proper authorization")
		return deprecated
	}

	if err != nil {
		logrus.WithError(err).WithField("reason", apierrors.ReasonForError(err)).Fatalf("Failed communicating with k8s while listing objects [%v]", gvr.String())
	}

	if len(deprecatedItems.Items) > 0 {
		replacement := deprecatedGVR.Replacement

		if replacement.Group != "" {
			replacementGvr := schema.GroupVersionResource{Group: replacement.Group, Version: replacement.Version, Resource: supported.Name}
			replacementItems, err := client.Resource(replacementGvr).List(context.TODO(), metav1.ListOptions{})

			if apierrors.IsNotFound(err) {
				return deprecated
			}

			if apierrors.IsForbidden(err) {
				logrus.WithError(err).Fatalf("Failed to list objects in the cluster. Permission denied! Please check if you have the proper authorization")
			}

			if err != nil {
				logrus.WithError(err).Fatalf("Failed communicating with k8s while listing objects.")
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

	return deprecated
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

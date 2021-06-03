package kubernetes

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {
	Client *kubernetes.Clientset
	Config *rest.Config
}

func IsInCluster() bool {
	_, err := rest.InClusterConfig()
	return err == nil
}

func BindFlags(flags *pflag.FlagSet) *clientcmd.ConfigOverrides {
	cmd := &clientcmd.ConfigOverrides{}
	overrides := clientcmd.RecommendedConfigOverrideFlags("")
	clientcmd.BindOverrideFlags(cmd, flags, overrides)
	return cmd
}

func NewClient(cmd *cobra.Command, configOverrides *clientcmd.ConfigOverrides) *KubeClient {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configFile, err := cmd.Root().PersistentFlags().GetString(clientcmd.RecommendedConfigPathFlag)

	if err != nil {
		fmt.Println(err) // TODO: logging
	}

	if configFile != "" {
		loadingRules.Precedence = []string{configFile}
	}

	// if you want to change the loading rules (which files in which order), you can do so here

	var config *rest.Config

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err = kubeConfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &KubeClient{Client: client, Config: config}
}

func (kubeClient *KubeClient) ListPods(namespace, labelSelector string) []corev1.Pod {
	list, err := kubeClient.Client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		panic(err.Error())
	}

	return list.Items
}

func (kubeClient *KubeClient) GetNamespacedPods(namespace string, names []string) []corev1.Pod {
	var pods []corev1.Pod

	if len(names) == 0 {
		// all pods of the given namespace
		list, err := kubeClient.Client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})

		if err == nil {
			pods = list.Items
		}
	} else {
		// the given pods in the given namespace
		for _, p := range names {
			pod, err := kubeClient.Client.CoreV1().Pods(namespace).Get(context.Background(), p, metav1.GetOptions{})

			if err == nil {
				pods = append(pods, *pod)
			}
		}
	}

	return pods
}

func (kubeClient *KubeClient) GetReplicaSet(pod corev1.Pod) (*appsv1.ReplicaSet, error) {
	var rsName string

	for _, ref := range pod.OwnerReferences {
		if ref.Kind == "ReplicaSet" {
			rsName = ref.Name
		}
	}

	if rsName == "" {
		return nil, nil
	}

	return kubeClient.Client.AppsV1().ReplicaSets(pod.Namespace).Get(context.Background(), rsName, metav1.GetOptions{})
}

func (kubeClient *KubeClient) GetDeployment(namespace string, name string) (*appsv1.Deployment, error) {
	return kubeClient.Client.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (kubeClient *KubeClient) GetDaemonSet(namespace string, name string) (*appsv1.DaemonSet, error) {
	return kubeClient.Client.AppsV1().DaemonSets(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (kubeClient *KubeClient) GetStatefulSet(namespace string, name string) (*appsv1.StatefulSet, error) {
	return kubeClient.Client.AppsV1().StatefulSets(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (kubeClient *KubeClient) GetPersistentVolumeClaim(namespace string, name string) (*corev1.PersistentVolumeClaim, error) {
	return kubeClient.Client.CoreV1().PersistentVolumeClaims(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (kubeClient *KubeClient) DiscoverResourceNameAndPreferredGV() KindVersions {
	kv := make(KindVersions)
	groups, resources, err := kubeClient.Client.DiscoveryClient.ServerGroupsAndResources()

	if err != nil {
		if apierrors.IsNotFound(err) {
			return kv
		}

		if apierrors.IsForbidden(err) {
			log.Fatalf("Failed to list objects for Name discovery. Permission denied! Please check if you have the proper authorization")
		}

		log.Fatalf("Failed communicating with k8s while discovering the object preferred name and gv. Error: %v", err)
	}

	for _, resourceGroup := range resources {
		parts := strings.Split(resourceGroup.GroupVersion, "/")
		var group, version string

		if len(parts) == 1 {
			group = ""
			version = parts[0]
		} else {
			group = parts[0]
			version = parts[1]
		}

		for _, resource := range resourceGroup.APIResources {
			if strings.Contains(resource.Name, "/") {
				continue
			}

			kindVersion := KindVersion{
				Group:     group,
				Version:   version,
				Name:      resource.Name,
				Preferred: isVersionPreferred(group, version, groups),
			}

			if _, ok := kv[resource.Kind]; !ok {
				kv[resource.Kind] = []KindVersion{kindVersion}
			} else {
				kv[resource.Kind] = append(kv[resource.Kind], kindVersion)
			}
		}
	}

	return kv
}

func isVersionPreferred(group, currentVersion string, allGroups []*metav1.APIGroup) bool {
	for _, supportedGroup := range allGroups {
		if supportedGroup.Name != group {
			continue
		}

		return supportedGroup.PreferredVersion.Version == currentVersion
	}

	return false
}

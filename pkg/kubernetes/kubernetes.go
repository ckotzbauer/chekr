package kubernetes

import (
	"context"

	"github.com/spf13/pflag"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {
	Client *kubernetes.Clientset
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

func NewClient(configOverrides *clientcmd.ConfigOverrides) *KubeClient {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	// if you want to change the loading rules (which files in which order), you can do so here

	var config *rest.Config
	var err error

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err = kubeConfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &KubeClient{Client: client}
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

package kubernetes

import (
	"context"

	"github.com/spf13/pflag"
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

func (kubeClient *KubeClient) GetNamedPods(namespace string, names []string) []corev1.Pod {
	var pods []corev1.Pod

	for _, p := range names {
		pod, err := kubeClient.Client.CoreV1().Pods(namespace).Get(context.Background(), p, metav1.GetOptions{})

		if err == nil {
			pods = append(pods, *pod)
		}
	}

	return pods
}

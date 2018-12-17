package kubernetes_test

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/topfreegames/pitaya-bot/cmd"
	pbKubernetes "github.com/topfreegames/pitaya-bot/kubernetes"
	"github.com/topfreegames/pitaya-bot/launcher"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreateManagerPod(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	specs, err := launcher.GetSpecs("../testing/specs/")
	assert.NoError(t, err)
	config := cmd.CreateConfig("../testing/config/config.yaml")
	pbKubernetes.CreateManagerPod(logrus.New(), clientset, config, specs, time.Minute)
	configMaps, err := clientset.CoreV1().ConfigMaps(corev1.NamespaceDefault).List(metav1.ListOptions{LabelSelector: "app=pitaya-bot-manager,game="})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(configMaps.Items))
	pods, err := clientset.CoreV1().Pods(corev1.NamespaceDefault).List(metav1.ListOptions{LabelSelector: "app=pitaya-bot-manager,game="})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(pods.Items))
}

func TestDeployJobs(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	specs, err := launcher.GetSpecs("../testing/specs/")
	assert.NoError(t, err)
	config := cmd.CreateConfig("../testing/config/config.yaml")
	pbKubernetes.DeployJobs(logrus.New(), clientset, config, specs, time.Minute)
	configMaps, err := clientset.CoreV1().ConfigMaps(corev1.NamespaceDefault).List(metav1.ListOptions{LabelSelector: "app=pitaya-bot,game="})
	assert.NoError(t, err)
	assert.Equal(t, len(specs)+1, len(configMaps.Items))
	jobs, err := clientset.BatchV1().Jobs(corev1.NamespaceDefault).List(metav1.ListOptions{LabelSelector: "app=pitaya-bot,game="})
	assert.NoError(t, err)
	assert.Equal(t, len(specs), len(jobs.Items))
}
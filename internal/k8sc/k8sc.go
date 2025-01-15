package k8sc

import (
	"context"
	"fmt"

	"github.com/daniarmas/notes/internal/clog"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sC interface {
	CreateJob(ctx context.Context, jobName, namespace, imageName string, args []string, envs []corev1.EnvFromSource) error
}

type k8sc struct {
	Clientset *kubernetes.Clientset
}

func NewClient() (K8sC, error) {
	var clientset *kubernetes.Clientset
	var clientsetError error

	// creates the in-cluster config
	config, configErr := rest.InClusterConfig()
	if configErr != nil {
		clog.Error(context.Background(), "error creating in-cluster config", configErr)
		return nil, configErr
	}
	// creates the clientset
	clientset, clientsetError = kubernetes.NewForConfig(config)
	if clientsetError != nil {
		clog.Error(context.Background(), "error creating kubernetes clientset", clientsetError)
		return nil, clientsetError
	}

	return &k8sc{
		Clientset: clientset,
	}, nil
}

// CreateJob creates a job in the k8s cluster
func (c *k8sc) CreateJob(ctx context.Context, jobName, namespace, imageName string, args []string, envs []corev1.EnvFromSource) error {
	// Create a job spec
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    jobName,
							Image:   imageName,
							Command: []string{"/app/notes"},
							Args:    args,
							EnvFrom: envs,
						},
					},
				},
			},
		},
	}

	// Create the job
	job, err := c.Clientset.BatchV1().Jobs(namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		clog.Error(context.Background(), "error creating job", err)
		return err
	}

	clog.Info(context.Background(), fmt.Sprintf("job created successful: %s", job.Name), nil)
	return nil
}

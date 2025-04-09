package k8sc

import (
	"context"
	"fmt"

	"github.com/daniarmas/clogg"
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
		clogg.Error(context.Background(), "error creating in-cluster config", clogg.String("error", configErr.Error()))
		return nil, configErr
	}
	// creates the clientset
	clientset, clientsetError = kubernetes.NewForConfig(config)
	if clientsetError != nil {
		clogg.Error(context.Background(), "error creating kubernetes clientset", clogg.String("error", clientsetError.Error()))
		return nil, clientsetError
	}

	return &k8sc{
		Clientset: clientset,
	}, nil
}

// CreateJob creates a job in the k8s cluster
func (c *k8sc) CreateJob(ctx context.Context, jobName, namespace, imageName string, args []string, envs []corev1.EnvFromSource) error {
	// Define the TTL duration in seconds
	ttlSecondsAfterFinished := int32(15) // 15 seconds
	// Define the backoff limit
	backoffLimit := int32(4) // Retry up to 4 times

	// Create a job spec
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			BackoffLimit:            &backoffLimit,
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
		clogg.Error(ctx, "error creating job", clogg.String("error", err.Error()))
		return err
	}

	clogg.Info(ctx, fmt.Sprintf("job created successful: %s", job.Name))
	return nil
}

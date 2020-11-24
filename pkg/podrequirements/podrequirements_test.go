package podrequirements

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestNoPods(t *testing.T) {
	pods := &corev1.PodList{
		Items: []corev1.Pod{},
	}

	reqs := GetPodsCPUAndMemoryRequirements(pods)

	assert.Equal(t, 0, len(reqs.Items), "Expected an empty list of pod requirements")
}

func TestOneEmptyPod(t *testing.T) {
	pods := &corev1.PodList{
		Items: []corev1.Pod{
			corev1.Pod{},
		},
	}

	reqs := GetPodsCPUAndMemoryRequirements(pods)

	assert.Equal(t, 1, len(reqs.Items), "Expected list with a single pod requirements")

	podReqs := reqs.Items[0]

	assertRequests(t, podReqs, "cpu", "0m")
	assertRequests(t, podReqs, "memory", "0Gi")
	assertLimits(t, podReqs, "cpu", "0m")
	assertLimits(t, podReqs, "memory", "0Gi")
}

func TestOnePodWithSingleContainerAndSomeMissingRequirements(t *testing.T) {
	pods := &corev1.PodList{
		Items: []corev1.Pod{
			corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu": resource.MustParse("50m"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("500m"),
									"memory": resource.MustParse("2Gi"),
								},
							},
						},
					},
				},
			},
		},
	}

	reqs := GetPodsCPUAndMemoryRequirements(pods)

	assert.Equal(t, 1, len(reqs.Items), "Expected a list with 1 pod requirements")

	podReqs := reqs.Items[0]

	assertRequests(t, podReqs, "cpu", "50m")
	assertRequests(t, podReqs, "memory", "0Gi")
	assertLimits(t, podReqs, "cpu", "500m")
	assertLimits(t, podReqs, "memory", "2Gi")
}

func TestOnePodWithSingleContainer(t *testing.T) {

	pods := &corev1.PodList{
		Items: []corev1.Pod{
			corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("50m"),
									"memory": resource.MustParse("1Gi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("500m"),
									"memory": resource.MustParse("2Gi"),
								},
							},
						},
					},
				},
			},
		},
	}

	reqs := GetPodsCPUAndMemoryRequirements(pods)

	assert.Equal(t, 1, len(reqs.Items), "Expected a list with 1 pod requirements")

	podReqs := reqs.Items[0]

	assertRequests(t, podReqs, "cpu", "50m")
	assertRequests(t, podReqs, "memory", "1Gi")
	assertLimits(t, podReqs, "cpu", "500m")
	assertLimits(t, podReqs, "memory", "2Gi")
}

func TestOnePodWithTwoStandarContainers(t *testing.T) {

	pods := &corev1.PodList{
		Items: []corev1.Pod{
			corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("50m"),
									"memory": resource.MustParse("1Gi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("500m"),
									"memory": resource.MustParse("2Gi"),
								},
							},
						},

						corev1.Container{
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("25m"),
									"memory": resource.MustParse("0.1Gi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("100m"),
									"memory": resource.MustParse("0.5Gi"),
								},
							},
						},
					},
				},
			},
		},
	}

	reqs := GetPodsCPUAndMemoryRequirements(pods)

	assert.Equal(t, 1, len(reqs.Items), "Expected a list with 1 pod requirements")

	podReqs := reqs.Items[0]

	assertRequests(t, podReqs, "cpu", "75m")
	assertRequests(t, podReqs, "memory", "1.1Gi")
	assertLimits(t, podReqs, "cpu", "600m")
	assertLimits(t, podReqs, "memory", "2.5Gi")
}

func TestOnePodWithOneInitContainerAndOneStandardContainer(t *testing.T) {

	pods := &corev1.PodList{
		Items: []corev1.Pod{
			corev1.Pod{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						corev1.Container{
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("50m"),
									"memory": resource.MustParse("1Gi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("500m"),
									"memory": resource.MustParse("2Gi"),
								},
							},
						},
					},
					Containers: []corev1.Container{
						corev1.Container{
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("25m"),
									"memory": resource.MustParse("0.1Gi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("100m"),
									"memory": resource.MustParse("0.5Gi"),
								},
							},
						},
					},
				},
			},
		},
	}

	reqs := GetPodsCPUAndMemoryRequirements(pods)

	assert.Equal(t, 1, len(reqs.Items), "Expected a list with 1 pod requirements")

	podReqs := reqs.Items[0]

	assertRequests(t, podReqs, "cpu", "75m")
	assertRequests(t, podReqs, "memory", "1.1Gi")
	assertLimits(t, podReqs, "cpu", "600m")
	assertLimits(t, podReqs, "memory", "2.5Gi")
}

func TestOnePodWithOneInitContainerAndOneEphemeralContainer(t *testing.T) {

	pods := &corev1.PodList{
		Items: []corev1.Pod{
			corev1.Pod{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						corev1.Container{
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("50m"),
									"memory": resource.MustParse("1Gi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("500m"),
									"memory": resource.MustParse("2Gi"),
								},
							},
						},
					},
					EphemeralContainers: []corev1.EphemeralContainer{
						corev1.EphemeralContainer{
							EphemeralContainerCommon: corev1.EphemeralContainerCommon{
								Resources: corev1.ResourceRequirements{
									Requests: map[corev1.ResourceName]resource.Quantity{
										"cpu":    resource.MustParse("25m"),
										"memory": resource.MustParse("0.1Gi"),
									},
									Limits: map[corev1.ResourceName]resource.Quantity{
										"cpu":    resource.MustParse("100m"),
										"memory": resource.MustParse("0.5Gi"),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	reqs := GetPodsCPUAndMemoryRequirements(pods)

	assert.Equal(t, 1, len(reqs.Items), "Expected a list with 1 pod requirements")

	podReqs := reqs.Items[0]

	assertRequests(t, podReqs, "cpu", "75m")
	assertRequests(t, podReqs, "memory", "1.1Gi")
	assertLimits(t, podReqs, "cpu", "600m")
	assertLimits(t, podReqs, "memory", "2.5Gi")
}

func TestTwoPodsWithSingleContainer(t *testing.T) {

	pods := &corev1.PodList{
		Items: []corev1.Pod{
			corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("50m"),
									"memory": resource.MustParse("1Gi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("500m"),
									"memory": resource.MustParse("2Gi"),
								},
							},
						},
					},
				},
			},

			corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("125m"),
									"memory": resource.MustParse("0.75Gi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("600m"),
									"memory": resource.MustParse("1.25Gi"),
								},
							},
						},
					},
				},
			},
		},
	}

	reqs := GetPodsCPUAndMemoryRequirements(pods)

	assert.Equal(t, 2, len(reqs.Items), "Expected a list with 2 pod requirements")

	podReqs := reqs.Items[0]

	assertRequests(t, podReqs, "cpu", "50m")
	assertRequests(t, podReqs, "memory", "1Gi")
	assertLimits(t, podReqs, "cpu", "500m")
	assertLimits(t, podReqs, "memory", "2Gi")

	podReqs = reqs.Items[1]

	assertRequests(t, podReqs, "cpu", "125m")
	assertRequests(t, podReqs, "memory", "0.75Gi")
	assertLimits(t, podReqs, "cpu", "600m")
	assertLimits(t, podReqs, "memory", "1.25Gi")
}

func assertRequests(t *testing.T, reqs PodRequirements, resourceName corev1.ResourceName, expected string) {
	actual := reqs.Requirements.Requests[resourceName]

	assert.True(t, resource.MustParse(expected).Equal(actual),
		fmt.Sprintf("Unexpected %s requests. Got: %v, Want: %v", resourceName, actual.String(), expected))
}

func assertLimits(t *testing.T, reqs PodRequirements, resourceName corev1.ResourceName, expected string) {
	actual := reqs.Requirements.Limits[resourceName]

	assert.True(t, resource.MustParse(expected).Equal(actual),
		fmt.Sprintf("Unexpected %s requests. Got: %v, Want: %v", resourceName, actual.String(), expected))
}

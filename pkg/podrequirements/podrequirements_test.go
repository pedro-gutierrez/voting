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

	assertRequests(t, "0m", podReqs, "cpu")
	assertRequests(t, "0Gi", podReqs, "memory")
	assertLimits(t, "0m", podReqs, "cpu")
	assertLimits(t, "0Gi", podReqs, "memory")
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

	assertRequests(t, "50m", podReqs, "cpu")
	assertRequests(t, "0Gi", podReqs, "memory")
	assertLimits(t, "500m", podReqs, "cpu")
	assertLimits(t, "2Gi", podReqs, "memory")

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

	assertRequests(t, "50m", podReqs, "cpu")
	assertRequests(t, "1Gi", podReqs, "memory")
	assertLimits(t, "500m", podReqs, "cpu")
	assertLimits(t, "2Gi", podReqs, "memory")
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

	assertRequests(t, "75m", podReqs, "cpu")
	assertRequests(t, "1.1Gi", podReqs, "memory")
	assertLimits(t, "600m", podReqs, "cpu")
	assertLimits(t, "2.5Gi", podReqs, "memory")
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

	assertRequests(t, "75m", podReqs, "cpu")
	assertRequests(t, "1.1Gi", podReqs, "memory")
	assertLimits(t, "600m", podReqs, "cpu")
	assertLimits(t, "2.5Gi", podReqs, "memory")
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

	assertRequests(t, "75m", podReqs, "cpu")
	assertRequests(t, "1.1Gi", podReqs, "memory")
	assertLimits(t, "600m", podReqs, "cpu")
	assertLimits(t, "2.5Gi", podReqs, "memory")
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

	assertRequests(t, "50m", podReqs, "cpu")
	assertRequests(t, "1Gi", podReqs, "memory")
	assertLimits(t, "500m", podReqs, "cpu")
	assertLimits(t, "2Gi", podReqs, "memory")

	podReqs = reqs.Items[1]

	assertRequests(t, "125m", podReqs, "cpu")
	assertRequests(t, "0.75Gi", podReqs, "memory")
	assertLimits(t, "600m", podReqs, "cpu")
	assertLimits(t, "1.25Gi", podReqs, "memory")
}

func assertRequests(t *testing.T, expected string, reqs PodRequirements, resourceName corev1.ResourceName) {
	actual := reqs.Requirements.Requests[resourceName]

	assert.True(t, resource.MustParse(expected).Equal(actual),
		fmt.Sprintf("Unexpected %s requests. Got: %v, Want: %v", resourceName, actual.String(), expected))
}

func assertLimits(t *testing.T, expected string, reqs PodRequirements, resourceName corev1.ResourceName) {
	actual := reqs.Requirements.Limits[resourceName]

	assert.True(t, resource.MustParse(expected).Equal(actual),
		fmt.Sprintf("Unexpected %s requests. Got: %v, Want: %v", resourceName, actual.String(), expected))
}

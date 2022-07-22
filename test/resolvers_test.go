/*
 Copyright 2022 The Tekton Authors

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

*/

package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/tektoncd/pipeline/pkg/pod"
	"github.com/tektoncd/pipeline/test/parse"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	knativetest "knative.dev/pkg/test"
	"knative.dev/pkg/test/helpers"
)

var hubFeatureFlags = requireAnyGate(map[string]string{
	"enable-hub-resolver": "true",
	"enable-api-fields":   "alpha",
})

var gitFeatureFlags = requireAnyGate(map[string]string{
	"enable-git-resolver": "true",
	"enable-api-fields":   "alpha",
})

func TestHubResolver(t *testing.T) {
	ctx := context.Background()
	c, namespace := setup(ctx, t, hubFeatureFlags)

	t.Parallel()

	knativetest.CleanupOnInterrupt(func() { tearDown(ctx, t, c, namespace) }, t.Logf)
	defer tearDown(ctx, t, c, namespace)

	prName := helpers.ObjectNameForTest(t)

	pipelineRun := parse.MustParsePipelineRun(t, fmt.Sprintf(`
metadata:
  name: %s
  namespace: %s
spec:
  workspaces:
    - name: output # this workspace name must be declared in the Pipeline
      volumeClaimTemplate:
        spec:
          accessModes:
            - ReadWriteOnce # access mode may affect how you can use this volume in parallel tasks
          resources:
            requests:
              storage: 1Gi
  pipelineSpec:
    workspaces:
      - name: output
    tasks:
      - name: task1
        workspaces:
          - name: output
        taskRef:
          resolver: hub
          resource:
          - name: kind
            value: task
          - name: name
            value: git-clone
          - name: version
            value: "0.3"
        params:
          - name: url
            value: https://github.com/tektoncd/pipeline
          - name: deleteExisting
            value: "true"
`, prName, namespace))

	_, err := c.PipelineRunClient.Create(ctx, pipelineRun, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create PipelineRun `%s`: %s", prName, err)
	}

	t.Logf("Waiting for PipelineRun %s in namespace %s to complete", prName, namespace)
	if err := WaitForPipelineRunState(ctx, c, prName, timeout, PipelineRunSucceed(prName), "PipelineRunSuccess"); err != nil {
		t.Fatalf("Error waiting for PipelineRun %s to finish: %s", prName, err)
	}

}

func TestHubResolver_Failure(t *testing.T) {
	ctx := context.Background()
	c, namespace := setup(ctx, t, hubFeatureFlags)

	t.Parallel()

	knativetest.CleanupOnInterrupt(func() { tearDown(ctx, t, c, namespace) }, t.Logf)
	defer tearDown(ctx, t, c, namespace)

	prName := helpers.ObjectNameForTest(t)

	pipelineRun := parse.MustParsePipelineRun(t, fmt.Sprintf(`
metadata:
  name: %s
  namespace: %s
spec:
  workspaces:
    - name: output # this workspace name must be declared in the Pipeline
      volumeClaimTemplate:
        spec:
          accessModes:
            - ReadWriteOnce # access mode may affect how you can use this volume in parallel tasks
          resources:
            requests:
              storage: 1Gi
  pipelineSpec:
    workspaces:
      - name: output
    tasks:
      - name: task1
        workspaces:
          - name: output
        taskRef:
          resolver: hub
          resource:
          - name: kind
            value: task
          - name: name
            value: git-clone-this-does-not-exist
          - name: version
            value: "0.3"
        params:
          - name: url
            value: https://github.com/tektoncd/pipeline
          - name: deleteExisting
            value: "true"
`, prName, namespace))

	_, err := c.PipelineRunClient.Create(ctx, pipelineRun, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create PipelineRun `%s`: %s", prName, err)
	}

	t.Logf("Waiting for PipelineRun %s in namespace %s to complete", prName, namespace)
	if err := WaitForPipelineRunState(ctx, c, prName, timeout,
		Chain(
			FailedWithReason(pod.ReasonCouldntGetTask, prName),
			FailedWithMessage("requested resource 'https://api.hub.tekton.dev/v1/resource/Tekton/task/git-clone-this-does-not-exist/0.3/yaml' not found on hub", prName),
		), "PipelineRunFailed"); err != nil {
		t.Fatalf("Error waiting for PipelineRun to finish with expected error: %s", err)
	}
}

func TestGitResolver(t *testing.T) {
	ctx := context.Background()
	c, namespace := setup(ctx, t, gitFeatureFlags)

	t.Parallel()

	knativetest.CleanupOnInterrupt(func() { tearDown(ctx, t, c, namespace) }, t.Logf)
	defer tearDown(ctx, t, c, namespace)

	prName := helpers.ObjectNameForTest(t)

	pipelineRun := parse.MustParsePipelineRun(t, fmt.Sprintf(`
metadata:
  name: %s
  namespace: %s
spec:
  workspaces:
    - name: output # this workspace name must be declared in the Pipeline
      volumeClaimTemplate:
        spec:
          accessModes:
            - ReadWriteOnce # access mode may affect how you can use this volume in parallel tasks
          resources:
            requests:
              storage: 1Gi
  pipelineSpec:
    workspaces:
      - name: output
    tasks:
      - name: task1
        workspaces:
          - name: output
        taskRef:
          resolver: git
          resource:
          - name: url
            value: https://github.com/tektoncd/catalog.git
          - name: pathInRepo
            value: /task/git-clone/0.7/git-clone.yaml
          - name: commit
            value: 783b4fe7d21148f3b1a93bfa49b0024d8c6c2955
        params:
          - name: url
            value: https://github.com/tektoncd/pipeline
          - name: deleteExisting
            value: "true"
`, prName, namespace))

	_, err := c.PipelineRunClient.Create(ctx, pipelineRun, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create PipelineRun `%s`: %s", prName, err)
	}

	t.Logf("Waiting for PipelineRun %s in namespace %s to complete", prName, namespace)
	if err := WaitForPipelineRunState(ctx, c, prName, timeout, PipelineRunSucceed(prName), "PipelineRunSuccess"); err != nil {
		t.Fatalf("Error waiting for PipelineRun %s to finish: %s", prName, err)
	}

}

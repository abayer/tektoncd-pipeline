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

package pipelinerun

import (
	"fmt"
	"regexp"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tektoncd/pipeline/pkg/apis/config"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tektoncd/pipeline/test"
	"github.com/tektoncd/pipeline/test/diff"
	"github.com/tektoncd/pipeline/test/names"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/apis"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
	logtesting "knative.dev/pkg/logging/testing"
)

func TestUpdatePipelineRunStatusFromInformer(t *testing.T) {
	testCases := []struct {
		name              string
		embeddedStatusVal string
	}{
		{
			name:              "default embedded status",
			embeddedStatusVal: config.DefaultEmbeddedStatus,
		},
		{
			name:              "full embedded status",
			embeddedStatusVal: config.FullEmbeddedStatus,
		},
		{
			name:              "both embedded status",
			embeddedStatusVal: config.BothEmbeddedStatus,
		},
		{
			name:              "minimal embedded status",
			embeddedStatusVal: config.MinimalEmbeddedStatus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			names.TestingSeed()

			pr := &v1beta1.PipelineRun{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pipeline-run",
					Namespace: "foo",
					Labels:    map[string]string{"mylabel": "myvale"},
				},
				Spec: v1beta1.PipelineRunSpec{
					PipelineSpec: &v1beta1.PipelineSpec{
						Tasks: []v1beta1.PipelineTask{
							{
								Name: "unit-test-task-spec",
								TaskSpec: &v1beta1.EmbeddedTask{
									TaskSpec: v1beta1.TaskSpec{
										Steps: []v1beta1.Step{{Container: corev1.Container{
											Name:  "mystep",
											Image: "myimage"}}},
									},
								},
							},
							{
								Name: "custom-task-ref",
								TaskRef: &v1beta1.TaskRef{
									APIVersion: "example.dev/v0",
									Kind:       "Example",
									Name:       "some-custom-task",
								},
							},
						},
					},
				},
			}

			cms := getConfigMapsWithEmbeddedStatus(tc.embeddedStatusVal)
			cms[0].Data["enable-custom-tasks"] = "true"
			d := test.Data{
				PipelineRuns: []*v1beta1.PipelineRun{pr},
				ConfigMaps:   cms,
			}
			prt := newPipelineRunTest(d, t)
			defer prt.Cancel()

			wantEvents := []string{
				"Normal Started",
				"Normal Running Tasks Completed: 0",
			}

			// Reconcile the PipelineRun.  This creates a Taskrun.
			reconciledRun, clients := prt.reconcileRun("foo", "test-pipeline-run", wantEvents, false)

			// Save the name of the TaskRun and Run that were created.
			taskRunName := ""
			runName := ""
			if shouldHaveFullEmbeddedStatus(tc.embeddedStatusVal) {
				if len(reconciledRun.Status.TaskRuns) != 1 {
					t.Fatalf("Expected 1 TaskRun but got %d", len(reconciledRun.Status.TaskRuns))
				}
				for k := range reconciledRun.Status.TaskRuns {
					taskRunName = k
					break
				}
				if len(reconciledRun.Status.Runs) != 1 {
					t.Fatalf("Expected 1 Run but got %d", len(reconciledRun.Status.Runs))
				}
				for k := range reconciledRun.Status.Runs {
					runName = k
					break
				}
			}
			if shouldHaveMinimalEmbeddedStatus(tc.embeddedStatusVal) {
				if len(reconciledRun.Status.ChildReferences) != 2 {
					t.Fatalf("Expected 2 ChildReferences but got %d", len(reconciledRun.Status.ChildReferences))
				}
				for _, cr := range reconciledRun.Status.ChildReferences {
					if cr.Kind == "TaskRun" {
						taskRunName = cr.Name
					}
					if cr.Kind == "Run" {
						runName = cr.Name
					}
				}
			}

			if taskRunName == "" {
				t.Fatal("expected to find a TaskRun name, but didn't")
			}
			if runName == "" {
				t.Fatal("expected to find a Run name, but didn't")
			}

			// Add a label to the PipelineRun.  This tests a scenario in issue 3126 which could prevent the reconciler
			// from finding TaskRuns that are missing from the status.
			reconciledRun.ObjectMeta.Labels["bah"] = "humbug"
			reconciledRun, err := clients.Pipeline.TektonV1beta1().PipelineRuns("foo").Update(prt.TestAssets.Ctx, reconciledRun, metav1.UpdateOptions{})
			if err != nil {
				t.Fatalf("unexpected error when updating status: %v", err)
			}

			// The label update triggers another reconcile.  Depending on timing, the PipelineRun passed to the reconcile may or may not
			// have the updated status with the name of the created TaskRun.  Clear the status because we want to test the case where the
			// status does not have the TaskRun.
			reconciledRun.Status = v1beta1.PipelineRunStatus{}
			if _, err := clients.Pipeline.TektonV1beta1().PipelineRuns("foo").UpdateStatus(prt.TestAssets.Ctx, reconciledRun, metav1.UpdateOptions{}); err != nil {
				t.Fatalf("unexpected error when updating status: %v", err)
			}

			reconciledRun, _ = prt.reconcileRun("foo", "test-pipeline-run", wantEvents, false)

			// Verify that the reconciler found the existing TaskRun instead of creating a new one.
			if shouldHaveFullEmbeddedStatus(tc.embeddedStatusVal) {
				if len(reconciledRun.Status.TaskRuns) != 1 {
					t.Fatalf("Expected 1 TaskRun after label change but got %d", len(reconciledRun.Status.TaskRuns))
				}
				for k := range reconciledRun.Status.TaskRuns {
					if k != taskRunName {
						t.Fatalf("Status has unexpected taskrun %s", k)
					}
				}
				if len(reconciledRun.Status.Runs) != 1 {
					t.Fatalf("Expected 1 Run after label change but got %d", len(reconciledRun.Status.Runs))
				}
				for k := range reconciledRun.Status.Runs {
					if k != runName {
						t.Fatalf("Status has unexpected Run %s", k)
					}
				}
			}
			if shouldHaveMinimalEmbeddedStatus(tc.embeddedStatusVal) {
				if len(reconciledRun.Status.ChildReferences) != 2 {
					t.Fatalf("Expected 2 ChildReferences after label change but got %d", len(reconciledRun.Status.ChildReferences))
				}
				for _, cr := range reconciledRun.Status.ChildReferences {
					if cr.Kind == "TaskRun" && cr.Name != taskRunName {
						t.Errorf("Status has unexpected taskrun %s", cr.Name)
					}
					if cr.Kind == "Run" && cr.Name != runName {
						t.Errorf("Status has unexpected Run %s", cr.Name)
					}
				}
			}
		})
	}
}

type updateStatusTaskRunsData struct {
	withConditions map[string]*v1beta1.PipelineRunTaskRunStatus
	missingTaskRun map[string]*v1beta1.PipelineRunTaskRunStatus
	foundTaskRun   map[string]*v1beta1.PipelineRunTaskRunStatus
	recovered      map[string]*v1beta1.PipelineRunTaskRunStatus
	simple         map[string]*v1beta1.PipelineRunTaskRunStatus
}

func getUpdateStatusTaskRunsData() updateStatusTaskRunsData {
	// PipelineRunConditionCheckStatus recovered by updatePipelineRunStatusFromTaskRuns
	// It does not include the status, which is then retrieved via the regular reconcile
	prccs2Recovered := map[string]*v1beta1.PipelineRunConditionCheckStatus{
		"pr-task-2-running-condition-check-xxyyy": {
			ConditionName: "running-condition-0",
		},
	}
	prccs3Recovered := map[string]*v1beta1.PipelineRunConditionCheckStatus{
		"pr-task-3-successful-condition-check-xxyyy": {
			ConditionName: "successful-condition-0",
		},
	}
	prccs4Recovered := map[string]*v1beta1.PipelineRunConditionCheckStatus{
		"pr-task-4-failed-condition-check-xxyyy": {
			ConditionName: "failed-condition-0",
		},
	}

	// PipelineRunConditionCheckStatus full is used to test the behaviour of updatePipelineRunStatusFromTaskRuns
	// when no orphan TaskRuns are found, to check we don't alter good ones
	prccs2Full := map[string]*v1beta1.PipelineRunConditionCheckStatus{
		"pr-task-2-running-condition-check-xxyyy": {
			ConditionName: "running-condition-0",
			Status: &v1beta1.ConditionCheckStatus{
				ConditionCheckStatusFields: v1beta1.ConditionCheckStatusFields{
					Check: corev1.ContainerState{
						Running: &corev1.ContainerStateRunning{},
					},
				},
				Status: duckv1beta1.Status{
					Conditions: []apis.Condition{{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown}},
				},
			},
		},
	}
	prccs3Full := map[string]*v1beta1.PipelineRunConditionCheckStatus{
		"pr-task-3-successful-condition-check-xxyyy": {
			ConditionName: "successful-condition-0",
			Status: &v1beta1.ConditionCheckStatus{
				ConditionCheckStatusFields: v1beta1.ConditionCheckStatusFields{
					Check: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{ExitCode: 0},
					},
				},
				Status: duckv1beta1.Status{
					Conditions: []apis.Condition{{Type: apis.ConditionSucceeded, Status: corev1.ConditionTrue}},
				},
			},
		},
	}
	prccs4Full := map[string]*v1beta1.PipelineRunConditionCheckStatus{
		"pr-task-4-failed-condition-check-xxyyy": {
			ConditionName: "failed-condition-0",
			Status: &v1beta1.ConditionCheckStatus{
				ConditionCheckStatusFields: v1beta1.ConditionCheckStatusFields{
					Check: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{ExitCode: 127},
					},
				},
				Status: duckv1beta1.Status{
					Conditions: []apis.Condition{{Type: apis.ConditionSucceeded, Status: corev1.ConditionFalse}},
				},
			},
		},
	}

	return updateStatusTaskRunsData{
		withConditions: map[string]*v1beta1.PipelineRunTaskRunStatus{
			"pr-task-1-xxyyy": {
				PipelineTaskName: "task-1",
				Status:           &v1beta1.TaskRunStatus{},
			},
			"pr-task-2-xxyyy": {
				PipelineTaskName: "task-2",
				Status:           nil,
				ConditionChecks:  prccs2Full,
			},
			"pr-task-3-xxyyy": {
				PipelineTaskName: "task-3",
				Status:           &v1beta1.TaskRunStatus{},
				ConditionChecks:  prccs3Full,
			},
			"pr-task-4-xxyyy": {
				PipelineTaskName: "task-4",
				Status:           nil,
				ConditionChecks:  prccs4Full,
			},
		},
		missingTaskRun: map[string]*v1beta1.PipelineRunTaskRunStatus{
			"pr-task-1-xxyyy": {
				PipelineTaskName: "task-1",
				Status:           &v1beta1.TaskRunStatus{},
			},
			"pr-task-2-xxyyy": {
				PipelineTaskName: "task-2",
				Status:           nil,
				ConditionChecks:  prccs2Full,
			},
			"pr-task-4-xxyyy": {
				PipelineTaskName: "task-4",
				Status:           nil,
				ConditionChecks:  prccs4Full,
			},
		},
		foundTaskRun: map[string]*v1beta1.PipelineRunTaskRunStatus{
			"pr-task-1-xxyyy": {
				PipelineTaskName: "task-1",
				Status:           &v1beta1.TaskRunStatus{},
			},
			"pr-task-2-xxyyy": {
				PipelineTaskName: "task-2",
				Status:           nil,
				ConditionChecks:  prccs2Full,
			},
			"pr-task-3-xxyyy": {
				PipelineTaskName: "task-3",
				Status:           &v1beta1.TaskRunStatus{},
				ConditionChecks:  prccs3Recovered,
			},
			"pr-task-4-xxyyy": {
				PipelineTaskName: "task-4",
				Status:           nil,
				ConditionChecks:  prccs4Full,
			},
		},
		recovered: map[string]*v1beta1.PipelineRunTaskRunStatus{
			"pr-task-1-xxyyy": {
				PipelineTaskName: "task-1",
				Status:           &v1beta1.TaskRunStatus{},
			},
			"orphaned-taskruns-pr-task-2-xxyyy": {
				PipelineTaskName: "task-2",
				Status:           nil,
				ConditionChecks:  prccs2Recovered,
			},
			"pr-task-3-xxyyy": {
				PipelineTaskName: "task-3",
				Status:           &v1beta1.TaskRunStatus{},
				ConditionChecks:  prccs3Recovered,
			},
			"orphaned-taskruns-pr-task-4-xxyyy": {
				PipelineTaskName: "task-4",
				Status:           nil,
				ConditionChecks:  prccs4Recovered,
			},
		},
		simple: map[string]*v1beta1.PipelineRunTaskRunStatus{
			"pr-task-1-xxyyy": {
				PipelineTaskName: "task-1",
				Status:           &v1beta1.TaskRunStatus{},
			},
		},
	}
}

func TestUpdatePipelineRunStatusFromTaskRuns(t *testing.T) {

	prUID := types.UID("11111111-1111-1111-1111-111111111111")
	otherPrUID := types.UID("22222222-2222-2222-2222-222222222222")

	taskRunsPRStatusData := getUpdateStatusTaskRunsData()

	prRunningStatus := duckv1beta1.Status{
		Conditions: []apis.Condition{
			{
				Type:    "Succeeded",
				Status:  "Unknown",
				Reason:  "Running",
				Message: "Not all Tasks in the Pipeline have finished executing",
			},
		},
	}

	prStatusWithCondition := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			TaskRuns: taskRunsPRStatusData.withConditions,
		},
	}

	prStatusMissingTaskRun := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			TaskRuns: taskRunsPRStatusData.missingTaskRun,
		},
	}

	prStatusFoundTaskRun := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			TaskRuns: taskRunsPRStatusData.foundTaskRun,
		},
	}

	prStatusWithEmptyTaskRuns := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			TaskRuns: map[string]*v1beta1.PipelineRunTaskRunStatus{},
		},
	}

	prStatusWithOrphans := v1beta1.PipelineRunStatus{
		Status: duckv1beta1.Status{
			Conditions: []apis.Condition{
				{
					Type:    "Succeeded",
					Status:  "Unknown",
					Reason:  "Running",
					Message: "Not all Tasks in the Pipeline have finished executing",
				},
			},
		},
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			TaskRuns: map[string]*v1beta1.PipelineRunTaskRunStatus{},
		},
	}

	prStatusRecovered := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			TaskRuns: taskRunsPRStatusData.recovered,
		},
	}

	prStatusRecoveredSimple := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			TaskRuns: taskRunsPRStatusData.simple,
		},
	}

	allTaskRuns := []*v1beta1.TaskRun{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-1-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-1",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-2-running-condition-check-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-2",
					pipeline.ConditionCheckKey:    "pr-task-2-running-condition-check-xxyyy",
					pipeline.ConditionNameKey:     "running-condition",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-3-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-3",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-3-successful-condition-check-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-3",
					pipeline.ConditionCheckKey:    "pr-task-3-successful-condition-check-xxyyy",
					pipeline.ConditionNameKey:     "successful-condition",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-4-failed-condition-check-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-4",
					pipeline.ConditionCheckKey:    "pr-task-4-failed-condition-check-xxyyy",
					pipeline.ConditionNameKey:     "failed-condition",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
	}

	taskRunsFromAnotherPR := []*v1beta1.TaskRun{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-1-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-1",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: otherPrUID}},
			},
		},
	}

	tcs := []struct {
		prName           string
		prStatus         v1beta1.PipelineRunStatus
		trs              []*v1beta1.TaskRun
		expectedPrStatus v1beta1.PipelineRunStatus
	}{
		{
			prName:           "no-status-no-taskruns",
			prStatus:         v1beta1.PipelineRunStatus{},
			trs:              nil,
			expectedPrStatus: v1beta1.PipelineRunStatus{},
		}, {
			prName:           "status-no-taskruns",
			prStatus:         prStatusWithCondition,
			trs:              nil,
			expectedPrStatus: prStatusWithCondition,
		}, {
			prName:   "status-nil-taskruns",
			prStatus: prStatusWithEmptyTaskRuns,
			trs: []*v1beta1.TaskRun{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pr-task-1-xxyyy",
						Labels: map[string]string{
							pipeline.PipelineTaskLabelKey: "task-1",
						},
						OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
					},
				},
			},
			expectedPrStatus: prStatusRecoveredSimple,
		}, {
			prName:   "status-missing-taskruns",
			prStatus: prStatusMissingTaskRun,
			trs: []*v1beta1.TaskRun{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pr-task-3-xxyyy",
						Labels: map[string]string{
							pipeline.PipelineTaskLabelKey: "task-3",
						},
						OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pr-task-3-successful-condition-check-xxyyy",
						Labels: map[string]string{
							pipeline.PipelineTaskLabelKey: "task-3",
							pipeline.ConditionCheckKey:    "pr-task-3-successful-condition-check-xxyyy",
							pipeline.ConditionNameKey:     "successful-condition",
						},
						OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
					},
				},
			},
			expectedPrStatus: prStatusFoundTaskRun,
		}, {
			prName:           "status-matching-taskruns-pr",
			prStatus:         prStatusWithCondition,
			trs:              allTaskRuns,
			expectedPrStatus: prStatusWithCondition,
		}, {
			prName:           "orphaned-taskruns-pr",
			prStatus:         prStatusWithOrphans,
			trs:              allTaskRuns,
			expectedPrStatus: prStatusRecovered,
		}, {
			prName:           "tr-from-another-pr",
			prStatus:         prStatusWithEmptyTaskRuns,
			trs:              taskRunsFromAnotherPR,
			expectedPrStatus: prStatusWithEmptyTaskRuns,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.prName, func(t *testing.T) {
			logger := logtesting.TestLogger(t)

			pr := &v1beta1.PipelineRun{
				ObjectMeta: metav1.ObjectMeta{Name: tc.prName, UID: prUID},
				Status:     tc.prStatus,
			}

			updatePipelineRunStatusFromTaskRuns(logger, pr, tc.trs, nil)
			actualPrStatus := pr.Status

			expectedPRStatus := tc.expectedPrStatus

			// The TaskRun keys for recovered taskruns will contain a new random key, appended to the
			// base name that we expect. Replace the random part so we can diff the whole structure
			actualTaskRuns := actualPrStatus.PipelineRunStatusFields.TaskRuns
			if actualTaskRuns != nil {
				fixedTaskRuns := make(map[string]*v1beta1.PipelineRunTaskRunStatus)
				re := regexp.MustCompile(`^[a-z\-]*?-task-[0-9]`)
				for k, v := range actualTaskRuns {
					newK := re.FindString(k)
					fixedTaskRuns[newK+"-xxyyy"] = v
				}
				actualPrStatus.PipelineRunStatusFields.TaskRuns = fixedTaskRuns
			}

			if d := cmp.Diff(expectedPRStatus, actualPrStatus); d != "" {
				t.Errorf("expected the PipelineRun status to match %#v. Diff %s", expectedPRStatus, diff.PrintWantGot(d))
			}
		})
	}
}

type updateStatusChildRefsData struct {
	withConditions []v1beta1.ChildStatusReference
	missingTaskRun []v1beta1.ChildStatusReference
	foundTaskRun   []v1beta1.ChildStatusReference
	missingRun     []v1beta1.ChildStatusReference
	recovered      []v1beta1.ChildStatusReference
	simple         []v1beta1.ChildStatusReference
}

func getUpdateStatusChildRefsData() updateStatusChildRefsData {
	// PipelineRunChildConditionCheckStatus recovered by updatePipelineRunStatusFromChildRefs
	// It does not include the status, which is then retrieved via the regular reconcile
	prccs2Recovered := []*v1beta1.PipelineRunChildConditionCheckStatus{{
		ConditionCheckName: "pr-task-2-running-condition-check-xxyyy",
		PipelineRunConditionCheckStatus: v1beta1.PipelineRunConditionCheckStatus{
			ConditionName: "running-condition-0",
		},
	}}
	prccs3Recovered := []*v1beta1.PipelineRunChildConditionCheckStatus{{
		ConditionCheckName: "pr-task-3-successful-condition-check-xxyyy",
		PipelineRunConditionCheckStatus: v1beta1.PipelineRunConditionCheckStatus{
			ConditionName: "successful-condition-0",
		},
	}}
	prccs4Recovered := []*v1beta1.PipelineRunChildConditionCheckStatus{{
		ConditionCheckName: "pr-task-4-failed-condition-check-xxyyy",
		PipelineRunConditionCheckStatus: v1beta1.PipelineRunConditionCheckStatus{
			ConditionName: "failed-condition-0",
		},
	}}

	// PipelineRunChildConditionCheckStatus full is used to test the behaviour of updatePipelineRunStatusFromChildRefs
	// when no orphan TaskRuns are found, to check we don't alter good ones
	prccs2Full := []*v1beta1.PipelineRunChildConditionCheckStatus{{
		ConditionCheckName: "pr-task-2-running-condition-check-xxyyy",
		PipelineRunConditionCheckStatus: v1beta1.PipelineRunConditionCheckStatus{
			ConditionName: "running-condition-0",
			Status: &v1beta1.ConditionCheckStatus{
				ConditionCheckStatusFields: v1beta1.ConditionCheckStatusFields{
					Check: corev1.ContainerState{
						Running: &corev1.ContainerStateRunning{},
					},
				},
				Status: duckv1beta1.Status{
					Conditions: []apis.Condition{{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown}},
				},
			},
		},
	}}
	prccs3Full := []*v1beta1.PipelineRunChildConditionCheckStatus{{
		ConditionCheckName: "pr-task-3-successful-condition-check-xxyyy",
		PipelineRunConditionCheckStatus: v1beta1.PipelineRunConditionCheckStatus{
			ConditionName: "successful-condition-0",
			Status: &v1beta1.ConditionCheckStatus{
				ConditionCheckStatusFields: v1beta1.ConditionCheckStatusFields{
					Check: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{ExitCode: 0},
					},
				},
				Status: duckv1beta1.Status{
					Conditions: []apis.Condition{{Type: apis.ConditionSucceeded, Status: corev1.ConditionTrue}},
				},
			},
		},
	}}
	prccs4Full := []*v1beta1.PipelineRunChildConditionCheckStatus{{
		ConditionCheckName: "pr-task-4-failed-condition-check-xxyyy",
		PipelineRunConditionCheckStatus: v1beta1.PipelineRunConditionCheckStatus{
			ConditionName: "failed-condition-0",
			Status: &v1beta1.ConditionCheckStatus{
				ConditionCheckStatusFields: v1beta1.ConditionCheckStatusFields{
					Check: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{ExitCode: 127},
					},
				},
				Status: duckv1beta1.Status{
					Conditions: []apis.Condition{{Type: apis.ConditionSucceeded, Status: corev1.ConditionFalse}},
				},
			},
		},
	}}

	return updateStatusChildRefsData{
		withConditions: []v1beta1.ChildStatusReference{
			childRefForPipelineTask("pr-task-1-xxyyy", "task-1", "TaskRun", "v1beta1", nil, nil),
			childRefForPipelineTask("pr-task-2-xxyyy", "task-2", "TaskRun", "v1beta1", nil, prccs2Full),
			childRefForPipelineTask("pr-task-3-xxyyy", "task-3", "TaskRun", "v1beta1", nil, prccs3Full),
			childRefForPipelineTask("pr-task-4-xxyyy", "task-4", "TaskRun", "v1beta1", nil, prccs4Full),
			childRefForPipelineTask("pr-run-6-xxyyy", "task-6", "Run", "v1alpha1", nil, nil),
		},
		missingTaskRun: []v1beta1.ChildStatusReference{
			childRefForPipelineTask("pr-task-1-xxyyy", "task-1", "TaskRun", "v1beta1", nil, nil),
			childRefForPipelineTask("pr-task-2-xxyyy", "task-2", "TaskRun", "v1beta1", nil, prccs2Full),
			childRefForPipelineTask("pr-task-4-xxyyy", "task-4", "TaskRun", "v1beta1", nil, prccs4Full),
			childRefForPipelineTask("pr-run-6-xxyyy", "task-6", "Run", "v1alpha1", nil, nil),
		},
		foundTaskRun: []v1beta1.ChildStatusReference{
			childRefForPipelineTask("pr-task-1-xxyyy", "task-1", "TaskRun", "v1beta1", nil, nil),
			childRefForPipelineTask("pr-task-2-xxyyy", "task-2", "TaskRun", "v1beta1", nil, prccs2Full),
			childRefForPipelineTask("pr-task-3-xxyyy", "task-3", "TaskRun", "v1beta1", nil, prccs3Recovered),
			childRefForPipelineTask("pr-task-4-xxyyy", "task-4", "TaskRun", "v1beta1", nil, prccs4Full),
			childRefForPipelineTask("pr-run-6-xxyyy", "task-6", "Run", "v1alpha1", nil, nil),
		},
		missingRun: []v1beta1.ChildStatusReference{
			childRefForPipelineTask("pr-task-1-xxyyy", "task-1", "TaskRun", "v1beta1", nil, nil),
			childRefForPipelineTask("pr-task-2-xxyyy", "task-2", "TaskRun", "v1beta1", nil, prccs2Full),
			childRefForPipelineTask("pr-task-3-xxyyy", "task-3", "TaskRun", "v1beta1", nil, prccs3Full),
			childRefForPipelineTask("pr-task-4-xxyyy", "task-4", "TaskRun", "v1beta1", nil, prccs4Full),
		},
		recovered: []v1beta1.ChildStatusReference{
			childRefForPipelineTask("pr-task-1-xxyyy", "task-1", "TaskRun", "v1beta1", nil, nil),
			childRefForPipelineTask("orphaned-taskruns-pr-task-2-xxyyy", "task-2", "TaskRun", "v1beta1", nil, prccs2Recovered),
			childRefForPipelineTask("pr-task-3-xxyyy", "task-3", "TaskRun", "v1beta1", nil, prccs3Recovered),
			childRefForPipelineTask("orphaned-taskruns-pr-task-4-xxyyy", "task-4", "TaskRun", "v1beta1", nil, prccs4Recovered),
			childRefForPipelineTask("pr-run-6-xxyyy", "task-6", "Run", "v1alpha1", nil, nil),
		},
		simple: []v1beta1.ChildStatusReference{
			childRefForPipelineTask("pr-task-1-xxyyy", "task-1", "TaskRun", "v1beta1", nil, nil),
		},
	}
}

func TestUpdatePipelineRunStatusFromChildRefs(t *testing.T) {
	prUID := types.UID("11111111-1111-1111-1111-111111111111")
	otherPrUID := types.UID("22222222-2222-2222-2222-222222222222")

	childRefsPRStatusData := getUpdateStatusChildRefsData()

	prRunningStatus := duckv1beta1.Status{
		Conditions: []apis.Condition{
			{
				Type:    "Succeeded",
				Status:  "Unknown",
				Reason:  "Running",
				Message: "Not all Tasks in the Pipeline have finished executing",
			},
		},
	}

	prStatusWithCondition := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: childRefsPRStatusData.withConditions,
		},
	}

	prStatusMissingTaskRun := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: childRefsPRStatusData.missingTaskRun,
		},
	}

	prStatusFoundTaskRun := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: childRefsPRStatusData.foundTaskRun,
		},
	}

	prStatusMissingRun := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: childRefsPRStatusData.missingRun,
		},
	}

	prStatusWithEmptyChildRefs := v1beta1.PipelineRunStatus{
		Status:                  prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{},
	}

	prStatusWithOrphans := v1beta1.PipelineRunStatus{
		Status: duckv1beta1.Status{
			Conditions: []apis.Condition{
				{
					Type:    "Succeeded",
					Status:  "Unknown",
					Reason:  "Running",
					Message: "Not all Tasks in the Pipeline have finished executing",
				},
			},
		},
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{},
	}

	prStatusRecovered := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: childRefsPRStatusData.recovered,
		},
	}

	prStatusRecoveredSimple := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: childRefsPRStatusData.simple,
		},
	}

	prStatusRecoveredSimpleWithRun := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: []v1beta1.ChildStatusReference{{
				TypeMeta: runtime.TypeMeta{
					APIVersion: "tekton.dev/v1alpha1",
					Kind:       "Run",
				},
				Name:             "pr-run-6-xxyyy",
				PipelineTaskName: "task-6",
			}},
		},
	}

	allTaskRuns := []*v1beta1.TaskRun{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-1-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-1",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-2-running-condition-check-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-2",
					pipeline.ConditionCheckKey:    "pr-task-2-running-condition-check-xxyyy",
					pipeline.ConditionNameKey:     "running-condition",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-3-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-3",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-3-successful-condition-check-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-3",
					pipeline.ConditionCheckKey:    "pr-task-3-successful-condition-check-xxyyy",
					pipeline.ConditionNameKey:     "successful-condition",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-4-failed-condition-check-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-4",
					pipeline.ConditionCheckKey:    "pr-task-4-failed-condition-check-xxyyy",
					pipeline.ConditionNameKey:     "failed-condition",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
	}

	taskRunsFromAnotherPR := []*v1beta1.TaskRun{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-1-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-1",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: otherPrUID}},
			},
		},
	}

	runsFromAnotherPR := []*v1alpha1.Run{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-run-1-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-1",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: otherPrUID}},
			},
		},
	}

	tcs := []struct {
		prName           string
		prStatus         v1beta1.PipelineRunStatus
		trs              []*v1beta1.TaskRun
		runs             []*v1alpha1.Run
		expectedPrStatus v1beta1.PipelineRunStatus
	}{
		{
			prName:           "no-status-no-taskruns-or-runs",
			prStatus:         v1beta1.PipelineRunStatus{},
			trs:              nil,
			runs:             nil,
			expectedPrStatus: v1beta1.PipelineRunStatus{},
		}, {
			prName:           "status-no-taskruns-or-runs",
			prStatus:         prStatusWithCondition,
			trs:              nil,
			runs:             nil,
			expectedPrStatus: prStatusWithCondition,
		}, {
			prName:   "status-nil-taskruns",
			prStatus: prStatusWithEmptyChildRefs,
			trs: []*v1beta1.TaskRun{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pr-task-1-xxyyy",
						Labels: map[string]string{
							pipeline.PipelineTaskLabelKey: "task-1",
						},
						OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
					},
				},
			},
			expectedPrStatus: prStatusRecoveredSimple,
		}, {
			prName:   "status-nil-runs",
			prStatus: prStatusWithEmptyChildRefs,
			runs: []*v1alpha1.Run{{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pr-run-6-xxyyy",
					Labels: map[string]string{
						pipeline.PipelineTaskLabelKey: "task-6",
					},
					OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
				},
			}},
			expectedPrStatus: prStatusRecoveredSimpleWithRun,
		}, {
			prName:   "status-missing-taskruns",
			prStatus: prStatusMissingTaskRun,
			trs: []*v1beta1.TaskRun{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pr-task-3-xxyyy",
						Labels: map[string]string{
							pipeline.PipelineTaskLabelKey: "task-3",
						},
						OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pr-task-3-successful-condition-check-xxyyy",
						Labels: map[string]string{
							pipeline.PipelineTaskLabelKey: "task-3",
							pipeline.ConditionCheckKey:    "pr-task-3-successful-condition-check-xxyyy",
							pipeline.ConditionNameKey:     "successful-condition",
						},
						OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
					},
				},
			},
			expectedPrStatus: prStatusFoundTaskRun,
		}, {
			prName:   "status-missing-runs",
			prStatus: prStatusMissingRun,
			runs: []*v1alpha1.Run{{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pr-run-6-xxyyy",
					Labels: map[string]string{
						pipeline.PipelineTaskLabelKey: "task-6",
					},
					OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
				},
			}},
			expectedPrStatus: prStatusWithCondition,
		}, {
			prName:           "status-matching-taskruns-pr",
			prStatus:         prStatusWithCondition,
			trs:              allTaskRuns,
			expectedPrStatus: prStatusWithCondition,
		}, {
			prName:   "orphaned-taskruns-pr",
			prStatus: prStatusWithOrphans,
			trs:      allTaskRuns,
			runs: []*v1alpha1.Run{{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pr-run-6-xxyyy",
					Labels: map[string]string{
						pipeline.PipelineTaskLabelKey: "task-6",
					},
					OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
				},
			}},
			expectedPrStatus: prStatusRecovered,
		}, {
			prName:           "tr-from-another-pr",
			prStatus:         prStatusWithEmptyChildRefs,
			trs:              taskRunsFromAnotherPR,
			expectedPrStatus: prStatusWithEmptyChildRefs,
		}, {
			prName:           "run-from-another-pr",
			prStatus:         prStatusWithEmptyChildRefs,
			runs:             runsFromAnotherPR,
			expectedPrStatus: prStatusWithEmptyChildRefs,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.prName, func(t *testing.T) {
			logger := logtesting.TestLogger(t)

			pr := &v1beta1.PipelineRun{
				ObjectMeta: metav1.ObjectMeta{Name: tc.prName, UID: prUID},
				Status:     tc.prStatus,
			}

			_ = updatePipelineRunStatusFromChildRefs(logger, pr, tc.trs, tc.runs)

			actualPrStatus := pr.Status

			actualChildRefs := actualPrStatus.ChildReferences
			if len(actualChildRefs) != 0 {
				var fixedChildRefs []v1beta1.ChildStatusReference
				re := regexp.MustCompile(`^[a-z\-]*?-(task|run)-[0-9]`)
				for _, cr := range actualChildRefs {
					cr.Name = fmt.Sprintf("%s-xxyyy", re.FindString(cr.Name))
					fixedChildRefs = append(fixedChildRefs, cr)
				}
				actualPrStatus.ChildReferences = fixedChildRefs
			}

			// Sort the ChildReferences to deal with annoying ordering issues.
			sort.Slice(actualPrStatus.ChildReferences, func(i, j int) bool {
				return actualPrStatus.ChildReferences[i].PipelineTaskName < actualPrStatus.ChildReferences[j].PipelineTaskName
			})

			if d := cmp.Diff(tc.expectedPrStatus, actualPrStatus); d != "" {
				t.Errorf("expected the PipelineRun status to match %#v. Diff %s", tc.expectedPrStatus, diff.PrintWantGot(d))
			}
		})
	}
}

func TestUpdatePipelineRunStatusFromSlices(t *testing.T) {
	prUID := types.UID("11111111-1111-1111-1111-111111111111")
	otherPrUID := types.UID("22222222-2222-2222-2222-222222222222")

	childRefsPRStatusData := getUpdateStatusChildRefsData()

	prRunningStatus := duckv1beta1.Status{
		Conditions: []apis.Condition{
			{
				Type:    "Succeeded",
				Status:  "Unknown",
				Reason:  "Running",
				Message: "Not all Tasks in the Pipeline have finished executing",
			},
		},
	}

	prStatusWithCondition := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: childRefsPRStatusData.withConditions,
		},
	}

	prStatusMissingTaskRun := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: childRefsPRStatusData.missingTaskRun,
		},
	}

	prStatusFoundTaskRun := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: childRefsPRStatusData.foundTaskRun,
		},
	}

	prStatusMissingRun := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: childRefsPRStatusData.missingRun,
		},
	}

	prStatusWithEmptyChildRefs := v1beta1.PipelineRunStatus{
		Status:                  prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{},
	}

	prStatusWithOrphans := v1beta1.PipelineRunStatus{
		Status: duckv1beta1.Status{
			Conditions: []apis.Condition{
				{
					Type:    "Succeeded",
					Status:  "Unknown",
					Reason:  "Running",
					Message: "Not all Tasks in the Pipeline have finished executing",
				},
			},
		},
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{},
	}

	prStatusRecovered := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: childRefsPRStatusData.recovered,
		},
	}

	prStatusRecoveredSimple := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: childRefsPRStatusData.simple,
		},
	}

	prStatusRecoveredSimpleWithRun := v1beta1.PipelineRunStatus{
		Status: prRunningStatus,
		PipelineRunStatusFields: v1beta1.PipelineRunStatusFields{
			ChildReferences: []v1beta1.ChildStatusReference{{
				TypeMeta: runtime.TypeMeta{
					APIVersion: "tekton.dev/v1alpha1",
					Kind:       "Run",
				},
				Name:             "pr-run-6-xxyyy",
				PipelineTaskName: "task-6",
			}},
		},
	}

	allTaskRuns := []*v1beta1.TaskRun{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-1-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-1",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-2-running-condition-check-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-2",
					pipeline.ConditionCheckKey:    "pr-task-2-running-condition-check-xxyyy",
					pipeline.ConditionNameKey:     "running-condition",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-3-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-3",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-3-successful-condition-check-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-3",
					pipeline.ConditionCheckKey:    "pr-task-3-successful-condition-check-xxyyy",
					pipeline.ConditionNameKey:     "successful-condition",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-4-failed-condition-check-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-4",
					pipeline.ConditionCheckKey:    "pr-task-4-failed-condition-check-xxyyy",
					pipeline.ConditionNameKey:     "failed-condition",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
			},
		},
	}

	taskRunsFromAnotherPR := []*v1beta1.TaskRun{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-task-1-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-1",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: otherPrUID}},
			},
		},
	}

	runsFromAnotherPR := []*v1alpha1.Run{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pr-run-1-xxyyy",
				Labels: map[string]string{
					pipeline.PipelineTaskLabelKey: "task-1",
				},
				OwnerReferences: []metav1.OwnerReference{{UID: otherPrUID}},
			},
		},
	}

	tcs := []struct {
		prName           string
		prStatus         v1beta1.PipelineRunStatus
		trs              []*v1beta1.TaskRun
		runs             []*v1alpha1.Run
		expectedPrStatus v1beta1.PipelineRunStatus
	}{
		{
			prName:           "no-status-no-taskruns-or-runs",
			prStatus:         v1beta1.PipelineRunStatus{},
			trs:              nil,
			runs:             nil,
			expectedPrStatus: v1beta1.PipelineRunStatus{},
		}, {
			prName:           "status-no-taskruns-or-runs",
			prStatus:         prStatusWithCondition,
			trs:              nil,
			runs:             nil,
			expectedPrStatus: prStatusWithCondition,
		}, {
			prName:   "status-nil-taskruns",
			prStatus: prStatusWithEmptyChildRefs,
			trs: []*v1beta1.TaskRun{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pr-task-1-xxyyy",
						Labels: map[string]string{
							pipeline.PipelineTaskLabelKey: "task-1",
						},
						OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
					},
				},
			},
			expectedPrStatus: prStatusRecoveredSimple,
		}, {
			prName:   "status-nil-runs",
			prStatus: prStatusWithEmptyChildRefs,
			runs: []*v1alpha1.Run{{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pr-run-6-xxyyy",
					Labels: map[string]string{
						pipeline.PipelineTaskLabelKey: "task-6",
					},
					OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
				},
			}},
			expectedPrStatus: prStatusRecoveredSimpleWithRun,
		}, {
			prName:   "status-missing-taskruns",
			prStatus: prStatusMissingTaskRun,
			trs: []*v1beta1.TaskRun{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pr-task-3-xxyyy",
						Labels: map[string]string{
							pipeline.PipelineTaskLabelKey: "task-3",
						},
						OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pr-task-3-successful-condition-check-xxyyy",
						Labels: map[string]string{
							pipeline.PipelineTaskLabelKey: "task-3",
							pipeline.ConditionCheckKey:    "pr-task-3-successful-condition-check-xxyyy",
							pipeline.ConditionNameKey:     "successful-condition",
						},
						OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
					},
				},
			},
			expectedPrStatus: prStatusFoundTaskRun,
		}, {
			prName:   "status-missing-runs",
			prStatus: prStatusMissingRun,
			runs: []*v1alpha1.Run{{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pr-run-6-xxyyy",
					Labels: map[string]string{
						pipeline.PipelineTaskLabelKey: "task-6",
					},
					OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
				},
			}},
			expectedPrStatus: prStatusWithCondition,
		}, {
			prName:           "status-matching-taskruns-pr",
			prStatus:         prStatusWithCondition,
			trs:              allTaskRuns,
			expectedPrStatus: prStatusWithCondition,
		}, {
			prName:   "orphaned-taskruns-pr",
			prStatus: prStatusWithOrphans,
			trs:      allTaskRuns,
			runs: []*v1alpha1.Run{{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pr-run-6-xxyyy",
					Labels: map[string]string{
						pipeline.PipelineTaskLabelKey: "task-6",
					},
					OwnerReferences: []metav1.OwnerReference{{UID: prUID}},
				},
			}},
			expectedPrStatus: prStatusRecovered,
		}, {
			prName:           "tr-from-another-pr",
			prStatus:         prStatusWithEmptyChildRefs,
			trs:              taskRunsFromAnotherPR,
			expectedPrStatus: prStatusWithEmptyChildRefs,
		}, {
			prName:           "run-from-another-pr",
			prStatus:         prStatusWithEmptyChildRefs,
			runs:             runsFromAnotherPR,
			expectedPrStatus: prStatusWithEmptyChildRefs,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.prName, func(t *testing.T) {
			logger := logtesting.TestLogger(t)

			pr := &v1beta1.PipelineRun{
				ObjectMeta: metav1.ObjectMeta{Name: tc.prName, UID: prUID},
				Status:     tc.prStatus,
			}

			_ = updatePipelineRunStatusFromChildRefs(logger, pr, tc.trs, tc.runs)

			actualPrStatus := pr.Status

			actualChildRefs := actualPrStatus.ChildReferences
			if len(actualChildRefs) != 0 {
				var fixedChildRefs []v1beta1.ChildStatusReference
				re := regexp.MustCompile(`^[a-z\-]*?-(task|run)-[0-9]`)
				for _, cr := range actualChildRefs {
					cr.Name = fmt.Sprintf("%s-xxyyy", re.FindString(cr.Name))
					fixedChildRefs = append(fixedChildRefs, cr)
				}
				actualPrStatus.ChildReferences = fixedChildRefs
			}

			// Sort the ChildReferences to deal with annoying ordering issues.
			sort.Slice(actualPrStatus.ChildReferences, func(i, j int) bool {
				return actualPrStatus.ChildReferences[i].PipelineTaskName < actualPrStatus.ChildReferences[j].PipelineTaskName
			})

			if d := cmp.Diff(tc.expectedPrStatus, actualPrStatus); d != "" {
				t.Errorf("expected the PipelineRun status to match %#v. Diff %s", tc.expectedPrStatus, diff.PrintWantGot(d))
			}
		})
	}
}

// TODO(abayer): Update from Runs
// TODO(abayer): Update from Slices

func childRefForPipelineTask(taskRunName, pipelineTaskName, kind, apiVersion string, whenExpressions []v1beta1.WhenExpression,
	conditionChecks []*v1beta1.PipelineRunChildConditionCheckStatus) v1beta1.ChildStatusReference {
	return v1beta1.ChildStatusReference{
		TypeMeta: runtime.TypeMeta{
			APIVersion: fmt.Sprintf("%s/%s", pipeline.GroupName, apiVersion),
			Kind:       kind,
		},
		Name:             taskRunName,
		PipelineTaskName: pipelineTaskName,
		WhenExpressions:  whenExpressions,
		ConditionChecks:  conditionChecks,
	}
}

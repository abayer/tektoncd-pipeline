/*
Copyright 2020 The Tekton Authors

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

package v1beta1

import (
	"context"
	"fmt"
	"time"

	"github.com/tektoncd/pipeline/pkg/apis/config"
	apisconfig "github.com/tektoncd/pipeline/pkg/apis/config"
	"github.com/tektoncd/pipeline/pkg/apis/validate"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
)

var _ apis.Validatable = (*PipelineRun)(nil)

// Validate pipelinerun
func (pr *PipelineRun) Validate(ctx context.Context) *apis.FieldError {
	errs := validate.ObjectMetadata(pr.GetObjectMeta()).ViaField("metadata")

	if apis.IsInDelete(ctx) {
		return nil
	}

	if pr.IsPending() && pr.HasStarted() {
		errs = errs.Also(apis.ErrInvalidValue("PipelineRun cannot be Pending after it is started", "spec.status"))
	}

	return errs.Also(pr.Spec.Validate(apis.WithinSpec(ctx)).ViaField("spec"))
}

// Validate pipelinerun spec
func (ps *PipelineRunSpec) Validate(ctx context.Context) (errs *apis.FieldError) {
	cfg := config.FromContextOrDefaults(ctx)

	// Must have exactly one of pipelineRef and pipelineSpec.
	if ps.PipelineRef == nil && ps.PipelineSpec == nil {
		errs = errs.Also(apis.ErrMissingOneOf("pipelineRef", "pipelineSpec"))
	}
	if ps.PipelineRef != nil && ps.PipelineSpec != nil {
		errs = errs.Also(apis.ErrMultipleOneOf("pipelineRef", "pipelineSpec"))
	}

	// Validate PipelineRef if it's present
	if ps.PipelineRef != nil {
		errs = errs.Also(ps.PipelineRef.Validate(ctx).ViaField("pipelineRef"))
	}

	// Validate PipelineSpec if it's present
	if ps.PipelineSpec != nil {
		errs = errs.Also(ps.PipelineSpec.Validate(ctx).ViaField("pipelineSpec"))
	}

	if ps.Timeout != nil {
		// timeout should be a valid duration of at least 0.
		if ps.Timeout.Duration < 0 {
			errs = errs.Also(apis.ErrInvalidValue(fmt.Sprintf("%s should be >= 0", ps.Timeout.Duration.String()), "timeout"))
		}
	}

	// This is an alpha feature and will fail validation if it's used in a pipelinerun spec
	// when the enable-api-fields feature gate is anything but "alpha".
	if ps.Timeouts != nil {
		if ps.Timeout != nil {
			// can't have both at the same time
			errs = errs.Also(apis.ErrDisallowedFields("timeout", "timeouts"))
		}
		errs = errs.Also(validateTimeoutFields(ctx, ps.Timeouts))
	}

	errs = errs.Also(validateSpecStatus(ctx, ps.Status))

	if ps.Workspaces != nil {
		wsNames := make(map[string]int)
		for idx, ws := range ps.Workspaces {
			errs = errs.Also(ws.Validate(ctx).ViaFieldIndex("workspaces", idx))
			if prevIdx, alreadyExists := wsNames[ws.Name]; alreadyExists {
				errs = errs.Also(apis.ErrGeneric(fmt.Sprintf("workspace %q provided by pipelinerun more than once, at index %d and %d", ws.Name, prevIdx, idx), "name").ViaFieldIndex("workspaces", idx))
			}
			wsNames[ws.Name] = idx
		}
	}

	for idx, trs := range ps.TaskRunSpecs {
		errs = errs.Also(validateTaskRunSpec(ctx, trs).ViaIndex(idx).ViaField("taskRunSpecs"))
	}

	if len(ps.PipelineRunSpecs) > 0 {
		if cfg.FeatureFlags.EnableAPIFields == config.AlphaAPIFields {
			for idx, prs := range ps.PipelineRunSpecs {
				errs = errs.Also(validatePipelineRunSpec(ctx, prs).ViaIndex(idx).ViaField("pipelineRunSpecs"))
			}
		} else {
			errs = errs.Also(apis.ErrDisallowedFields("pipelineRunSpecs"))
		}
	}

	return errs
}

func validateSpecStatus(ctx context.Context, status PipelineRunSpecStatus) *apis.FieldError {
	switch status {
	case "":
		return nil
	case PipelineRunSpecStatusPending,
		PipelineRunSpecStatusCancelledDeprecated:
		return nil
	case PipelineRunSpecStatusCancelled,
		PipelineRunSpecStatusCancelledRunFinally,
		PipelineRunSpecStatusStoppedRunFinally:
		return nil
	}

	return apis.ErrInvalidValue(fmt.Sprintf("%s should be %s, %s, %s or %s", status,
		PipelineRunSpecStatusCancelled,
		PipelineRunSpecStatusCancelledRunFinally,
		PipelineRunSpecStatusStoppedRunFinally,
		PipelineRunSpecStatusPending), "status")

}

func validateTimeoutDuration(field string, d *metav1.Duration) (errs *apis.FieldError) {
	if d != nil && d.Duration < 0 {
		fieldPath := fmt.Sprintf("timeouts.%s", field)
		return errs.Also(apis.ErrInvalidValue(fmt.Sprintf("%s should be >= 0", d.Duration.String()), fieldPath))
	}
	return nil
}

func validatePipelineTimeout(timeoutFields *TimeoutFields, timeout time.Duration, errorMsg string) (errs *apis.FieldError) {
	if timeoutFields.Tasks != nil {
		tasksTimeoutErr := false
		tasksTimeoutStr := timeoutFields.Tasks.Duration.String()
		if timeoutFields.Tasks.Duration > timeout {
			tasksTimeoutErr = true
		}
		if timeoutFields.Tasks.Duration == apisconfig.NoTimeoutDuration && timeout != apisconfig.NoTimeoutDuration {
			tasksTimeoutErr = true
			tasksTimeoutStr += " (no timeout)"
		}
		if tasksTimeoutErr {
			errs = errs.Also(apis.ErrInvalidValue(fmt.Sprintf("%s %s", tasksTimeoutStr, errorMsg), "timeouts.tasks"))
		}
	}

	if timeoutFields.Finally != nil {
		finallyTimeoutErr := false
		finallyTimeoutStr := timeoutFields.Finally.Duration.String()
		if timeoutFields.Finally.Duration > timeout {
			finallyTimeoutErr = true
		}
		if timeoutFields.Finally.Duration == apisconfig.NoTimeoutDuration && timeout != apisconfig.NoTimeoutDuration {
			finallyTimeoutErr = true
			finallyTimeoutStr += " (no timeout)"
		}
		if finallyTimeoutErr {
			errs = errs.Also(apis.ErrInvalidValue(fmt.Sprintf("%s %s", finallyTimeoutStr, errorMsg), "timeouts.finally"))
		}
	}

	if timeoutFields.Tasks != nil && timeoutFields.Finally != nil {
		if timeoutFields.Tasks.Duration+timeoutFields.Finally.Duration > timeout {
			errs = errs.Also(apis.ErrInvalidValue(fmt.Sprintf("%s + %s %s", timeoutFields.Tasks.Duration.String(), timeoutFields.Finally.Duration.String(), errorMsg), "timeouts.tasks"))
			errs = errs.Also(apis.ErrInvalidValue(fmt.Sprintf("%s + %s %s", timeoutFields.Tasks.Duration.String(), timeoutFields.Finally.Duration.String(), errorMsg), "timeouts.finally"))
		}
	}
	return errs
}

func validateTaskRunSpec(ctx context.Context, trs PipelineTaskRunSpec) (errs *apis.FieldError) {
	cfg := config.FromContextOrDefaults(ctx)
	if cfg.FeatureFlags.EnableAPIFields == config.AlphaAPIFields {
		if trs.StepOverrides != nil {
			errs = errs.Also(validateStepOverrides(trs.StepOverrides).ViaField("stepOverrides"))
		}
		if trs.SidecarOverrides != nil {
			errs = errs.Also(validateSidecarOverrides(trs.SidecarOverrides).ViaField("sidecarOverrides"))
		}
	} else {
		if trs.StepOverrides != nil {
			errs = errs.Also(apis.ErrDisallowedFields("stepOverrides"))
		}
		if trs.SidecarOverrides != nil {
			errs = errs.Also(apis.ErrDisallowedFields("sidecarOverrides"))
		}
	}
	return errs
}

func validateTimeoutFields(ctx context.Context, timeoutFields *TimeoutFields) (errs *apis.FieldError) {
	// tasks timeout should be a valid duration of at least 0.
	errs = errs.Also(validateTimeoutDuration("tasks", timeoutFields.Tasks))

	// finally timeout should be a valid duration of at least 0.
	errs = errs.Also(validateTimeoutDuration("finally", timeoutFields.Finally))

	// pipeline timeout should be a valid duration of at least 0.
	errs = errs.Also(validateTimeoutDuration("pipeline", timeoutFields.Pipeline))

	if timeoutFields.Pipeline != nil {
		errs = errs.Also(validatePipelineTimeout(timeoutFields, timeoutFields.Pipeline.Duration, "should be <= pipeline duration"))
	} else {
		defaultTimeout := time.Duration(config.FromContextOrDefaults(ctx).Defaults.DefaultTimeoutMinutes)
		errs = errs.Also(validatePipelineTimeout(timeoutFields, defaultTimeout, "should be <= default timeout duration"))
	}

	return errs
}

func validatePipelineRunSpec(ctx context.Context, prs PipelinePipelineRunSpec) (errs *apis.FieldError) {
	// This is an alpha feature and will fail validation if it's used in a pipelinerun spec
	// when the enable-api-fields feature gate is anything but "alpha".
	if prs.Timeouts != nil {
		errs = errs.Also(validateTimeoutFields(ctx, prs.Timeouts))
	}

	for idx, trs := range prs.TaskRunSpecs {
		errs = errs.Also(validateTaskRunSpec(ctx, trs).ViaIndex(idx).ViaField("taskRunSpecs"))
	}

	return errs
}

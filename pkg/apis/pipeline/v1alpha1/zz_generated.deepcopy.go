// +build !ignore_autogenerated

/*
Copyright 2018 The Knative Authors

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

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package v1alpha1

import (
	duck_v1alpha1 "github.com/knative/pkg/apis/duck/v1alpha1"
	core_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactBucket) DeepCopyInto(out *ArtifactBucket) {
	*out = *in
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]SecretParam, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactBucket.
func (in *ArtifactBucket) DeepCopy() *ArtifactBucket {
	if in == nil {
		return nil
	}
	out := new(ArtifactBucket)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactPVC) DeepCopyInto(out *ArtifactPVC) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactPVC.
func (in *ArtifactPVC) DeepCopy() *ArtifactPVC {
	if in == nil {
		return nil
	}
	out := new(ArtifactPVC)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BuildGCSResource) DeepCopyInto(out *BuildGCSResource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BuildGCSResource.
func (in *BuildGCSResource) DeepCopy() *BuildGCSResource {
	if in == nil {
		return nil
	}
	out := new(BuildGCSResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResource) DeepCopyInto(out *ClusterResource) {
	*out = *in
	if in.CAData != nil {
		in, out := &in.CAData, &out.CAData
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]SecretParam, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResource.
func (in *ClusterResource) DeepCopy() *ClusterResource {
	if in == nil {
		return nil
	}
	out := new(ClusterResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterTask) DeepCopyInto(out *ClusterTask) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterTask.
func (in *ClusterTask) DeepCopy() *ClusterTask {
	if in == nil {
		return nil
	}
	out := new(ClusterTask)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterTask) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterTaskList) DeepCopyInto(out *ClusterTaskList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterTask, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterTaskList.
func (in *ClusterTaskList) DeepCopy() *ClusterTaskList {
	if in == nil {
		return nil
	}
	out := new(ClusterTaskList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterTaskList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCSResource) DeepCopyInto(out *GCSResource) {
	*out = *in
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]SecretParam, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCSResource.
func (in *GCSResource) DeepCopy() *GCSResource {
	if in == nil {
		return nil
	}
	out := new(GCSResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitResource) DeepCopyInto(out *GitResource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitResource.
func (in *GitResource) DeepCopy() *GitResource {
	if in == nil {
		return nil
	}
	out := new(GitResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageResource) DeepCopyInto(out *ImageResource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageResource.
func (in *ImageResource) DeepCopy() *ImageResource {
	if in == nil {
		return nil
	}
	out := new(ImageResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Inputs) DeepCopyInto(out *Inputs) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]TaskResource, len(*in))
		copy(*out, *in)
	}
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make([]TaskParam, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Inputs.
func (in *Inputs) DeepCopy() *Inputs {
	if in == nil {
		return nil
	}
	out := new(Inputs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Outputs) DeepCopyInto(out *Outputs) {
	*out = *in
	if in.Results != nil {
		in, out := &in.Results, &out.Results
		*out = make([]TestResult, len(*in))
		copy(*out, *in)
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]TaskResource, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Outputs.
func (in *Outputs) DeepCopy() *Outputs {
	if in == nil {
		return nil
	}
	out := new(Outputs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Param) DeepCopyInto(out *Param) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Param.
func (in *Param) DeepCopy() *Param {
	if in == nil {
		return nil
	}
	out := new(Param)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Pipeline) DeepCopyInto(out *Pipeline) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Pipeline.
func (in *Pipeline) DeepCopy() *Pipeline {
	if in == nil {
		return nil
	}
	out := new(Pipeline)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Pipeline) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineDeclaredResource) DeepCopyInto(out *PipelineDeclaredResource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineDeclaredResource.
func (in *PipelineDeclaredResource) DeepCopy() *PipelineDeclaredResource {
	if in == nil {
		return nil
	}
	out := new(PipelineDeclaredResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineList) DeepCopyInto(out *PipelineList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Pipeline, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineList.
func (in *PipelineList) DeepCopy() *PipelineList {
	if in == nil {
		return nil
	}
	out := new(PipelineList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PipelineList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineParam) DeepCopyInto(out *PipelineParam) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineParam.
func (in *PipelineParam) DeepCopy() *PipelineParam {
	if in == nil {
		return nil
	}
	out := new(PipelineParam)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineRef) DeepCopyInto(out *PipelineRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineRef.
func (in *PipelineRef) DeepCopy() *PipelineRef {
	if in == nil {
		return nil
	}
	out := new(PipelineRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineResource) DeepCopyInto(out *PipelineResource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineResource.
func (in *PipelineResource) DeepCopy() *PipelineResource {
	if in == nil {
		return nil
	}
	out := new(PipelineResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PipelineResource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineResourceBinding) DeepCopyInto(out *PipelineResourceBinding) {
	*out = *in
	out.ResourceRef = in.ResourceRef
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineResourceBinding.
func (in *PipelineResourceBinding) DeepCopy() *PipelineResourceBinding {
	if in == nil {
		return nil
	}
	out := new(PipelineResourceBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineResourceList) DeepCopyInto(out *PipelineResourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PipelineResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineResourceList.
func (in *PipelineResourceList) DeepCopy() *PipelineResourceList {
	if in == nil {
		return nil
	}
	out := new(PipelineResourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PipelineResourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineResourceRef) DeepCopyInto(out *PipelineResourceRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineResourceRef.
func (in *PipelineResourceRef) DeepCopy() *PipelineResourceRef {
	if in == nil {
		return nil
	}
	out := new(PipelineResourceRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineResourceSpec) DeepCopyInto(out *PipelineResourceSpec) {
	*out = *in
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make([]Param, len(*in))
		copy(*out, *in)
	}
	if in.SecretParams != nil {
		in, out := &in.SecretParams, &out.SecretParams
		*out = make([]SecretParam, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineResourceSpec.
func (in *PipelineResourceSpec) DeepCopy() *PipelineResourceSpec {
	if in == nil {
		return nil
	}
	out := new(PipelineResourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineResourceStatus) DeepCopyInto(out *PipelineResourceStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineResourceStatus.
func (in *PipelineResourceStatus) DeepCopy() *PipelineResourceStatus {
	if in == nil {
		return nil
	}
	out := new(PipelineResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineRun) DeepCopyInto(out *PipelineRun) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineRun.
func (in *PipelineRun) DeepCopy() *PipelineRun {
	if in == nil {
		return nil
	}
	out := new(PipelineRun)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PipelineRun) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineRunList) DeepCopyInto(out *PipelineRunList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PipelineRun, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineRunList.
func (in *PipelineRunList) DeepCopy() *PipelineRunList {
	if in == nil {
		return nil
	}
	out := new(PipelineRunList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PipelineRunList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineRunSpec) DeepCopyInto(out *PipelineRunSpec) {
	*out = *in
	out.PipelineRef = in.PipelineRef
	out.Trigger = in.Trigger
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]PipelineResourceBinding, len(*in))
		copy(*out, *in)
	}
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make([]Param, len(*in))
		copy(*out, *in)
	}
	if in.Results != nil {
		in, out := &in.Results, &out.Results
		if *in == nil {
			*out = nil
		} else {
			*out = new(Results)
			**out = **in
		}
	}
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		if *in == nil {
			*out = nil
		} else {
			*out = new(v1.Duration)
			**out = **in
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		if *in == nil {
			*out = nil
		} else {
			*out = new(core_v1.Affinity)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineRunSpec.
func (in *PipelineRunSpec) DeepCopy() *PipelineRunSpec {
	if in == nil {
		return nil
	}
	out := new(PipelineRunSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineRunStatus) DeepCopyInto(out *PipelineRunStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(duck_v1alpha1.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Results != nil {
		in, out := &in.Results, &out.Results
		if *in == nil {
			*out = nil
		} else {
			*out = new(Results)
			**out = **in
		}
	}
	if in.StartTime != nil {
		in, out := &in.StartTime, &out.StartTime
		if *in == nil {
			*out = nil
		} else {
			*out = new(v1.Time)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.CompletionTime != nil {
		in, out := &in.CompletionTime, &out.CompletionTime
		if *in == nil {
			*out = nil
		} else {
			*out = new(v1.Time)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.TaskRuns != nil {
		in, out := &in.TaskRuns, &out.TaskRuns
		*out = make(map[string]*PipelineRunTaskRunStatus, len(*in))
		for key, val := range *in {
			if val == nil {
				(*out)[key] = nil
			} else {
				(*out)[key] = new(PipelineRunTaskRunStatus)
				val.DeepCopyInto((*out)[key])
			}
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineRunStatus.
func (in *PipelineRunStatus) DeepCopy() *PipelineRunStatus {
	if in == nil {
		return nil
	}
	out := new(PipelineRunStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineRunTaskRunStatus) DeepCopyInto(out *PipelineRunTaskRunStatus) {
	*out = *in
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		if *in == nil {
			*out = nil
		} else {
			*out = new(TaskRunStatus)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineRunTaskRunStatus.
func (in *PipelineRunTaskRunStatus) DeepCopy() *PipelineRunTaskRunStatus {
	if in == nil {
		return nil
	}
	out := new(PipelineRunTaskRunStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineSpec) DeepCopyInto(out *PipelineSpec) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]PipelineDeclaredResource, len(*in))
		copy(*out, *in)
	}
	if in.Tasks != nil {
		in, out := &in.Tasks, &out.Tasks
		*out = make([]PipelineTask, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make([]PipelineParam, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineSpec.
func (in *PipelineSpec) DeepCopy() *PipelineSpec {
	if in == nil {
		return nil
	}
	out := new(PipelineSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineStatus) DeepCopyInto(out *PipelineStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineStatus.
func (in *PipelineStatus) DeepCopy() *PipelineStatus {
	if in == nil {
		return nil
	}
	out := new(PipelineStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineTask) DeepCopyInto(out *PipelineTask) {
	*out = *in
	out.TaskRef = in.TaskRef
	if in.RunAfter != nil {
		in, out := &in.RunAfter, &out.RunAfter
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		if *in == nil {
			*out = nil
		} else {
			*out = new(PipelineTaskResources)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make([]Param, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineTask.
func (in *PipelineTask) DeepCopy() *PipelineTask {
	if in == nil {
		return nil
	}
	out := new(PipelineTask)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineTaskInputResource) DeepCopyInto(out *PipelineTaskInputResource) {
	*out = *in
	if in.From != nil {
		in, out := &in.From, &out.From
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineTaskInputResource.
func (in *PipelineTaskInputResource) DeepCopy() *PipelineTaskInputResource {
	if in == nil {
		return nil
	}
	out := new(PipelineTaskInputResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineTaskOutputResource) DeepCopyInto(out *PipelineTaskOutputResource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineTaskOutputResource.
func (in *PipelineTaskOutputResource) DeepCopy() *PipelineTaskOutputResource {
	if in == nil {
		return nil
	}
	out := new(PipelineTaskOutputResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineTaskParam) DeepCopyInto(out *PipelineTaskParam) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineTaskParam.
func (in *PipelineTaskParam) DeepCopy() *PipelineTaskParam {
	if in == nil {
		return nil
	}
	out := new(PipelineTaskParam)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineTaskResources) DeepCopyInto(out *PipelineTaskResources) {
	*out = *in
	if in.Inputs != nil {
		in, out := &in.Inputs, &out.Inputs
		*out = make([]PipelineTaskInputResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Outputs != nil {
		in, out := &in.Outputs, &out.Outputs
		*out = make([]PipelineTaskOutputResource, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineTaskResources.
func (in *PipelineTaskResources) DeepCopy() *PipelineTaskResources {
	if in == nil {
		return nil
	}
	out := new(PipelineTaskResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineTaskRun) DeepCopyInto(out *PipelineTaskRun) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineTaskRun.
func (in *PipelineTaskRun) DeepCopy() *PipelineTaskRun {
	if in == nil {
		return nil
	}
	out := new(PipelineTaskRun)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineTrigger) DeepCopyInto(out *PipelineTrigger) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineTrigger.
func (in *PipelineTrigger) DeepCopy() *PipelineTrigger {
	if in == nil {
		return nil
	}
	out := new(PipelineTrigger)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Results) DeepCopyInto(out *Results) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Results.
func (in *Results) DeepCopy() *Results {
	if in == nil {
		return nil
	}
	out := new(Results)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretParam) DeepCopyInto(out *SecretParam) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretParam.
func (in *SecretParam) DeepCopy() *SecretParam {
	if in == nil {
		return nil
	}
	out := new(SecretParam)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StepState) DeepCopyInto(out *StepState) {
	*out = *in
	in.ContainerState.DeepCopyInto(&out.ContainerState)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StepState.
func (in *StepState) DeepCopy() *StepState {
	if in == nil {
		return nil
	}
	out := new(StepState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Task) DeepCopyInto(out *Task) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Task.
func (in *Task) DeepCopy() *Task {
	if in == nil {
		return nil
	}
	out := new(Task)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Task) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskList) DeepCopyInto(out *TaskList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Task, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskList.
func (in *TaskList) DeepCopy() *TaskList {
	if in == nil {
		return nil
	}
	out := new(TaskList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TaskList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskParam) DeepCopyInto(out *TaskParam) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskParam.
func (in *TaskParam) DeepCopy() *TaskParam {
	if in == nil {
		return nil
	}
	out := new(TaskParam)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskRef) DeepCopyInto(out *TaskRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskRef.
func (in *TaskRef) DeepCopy() *TaskRef {
	if in == nil {
		return nil
	}
	out := new(TaskRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskResource) DeepCopyInto(out *TaskResource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskResource.
func (in *TaskResource) DeepCopy() *TaskResource {
	if in == nil {
		return nil
	}
	out := new(TaskResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskResourceBinding) DeepCopyInto(out *TaskResourceBinding) {
	*out = *in
	out.ResourceRef = in.ResourceRef
	if in.ResourceSpec != nil {
		in, out := &in.ResourceSpec, &out.ResourceSpec
		if *in == nil {
			*out = nil
		} else {
			*out = new(PipelineResourceSpec)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskResourceBinding.
func (in *TaskResourceBinding) DeepCopy() *TaskResourceBinding {
	if in == nil {
		return nil
	}
	out := new(TaskResourceBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskRun) DeepCopyInto(out *TaskRun) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskRun.
func (in *TaskRun) DeepCopy() *TaskRun {
	if in == nil {
		return nil
	}
	out := new(TaskRun)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TaskRun) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskRunInputs) DeepCopyInto(out *TaskRunInputs) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]TaskResourceBinding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make([]Param, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskRunInputs.
func (in *TaskRunInputs) DeepCopy() *TaskRunInputs {
	if in == nil {
		return nil
	}
	out := new(TaskRunInputs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskRunList) DeepCopyInto(out *TaskRunList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TaskRun, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskRunList.
func (in *TaskRunList) DeepCopy() *TaskRunList {
	if in == nil {
		return nil
	}
	out := new(TaskRunList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TaskRunList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskRunOutputs) DeepCopyInto(out *TaskRunOutputs) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]TaskResourceBinding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make([]Param, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskRunOutputs.
func (in *TaskRunOutputs) DeepCopy() *TaskRunOutputs {
	if in == nil {
		return nil
	}
	out := new(TaskRunOutputs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskRunSpec) DeepCopyInto(out *TaskRunSpec) {
	*out = *in
	out.Trigger = in.Trigger
	in.Inputs.DeepCopyInto(&out.Inputs)
	in.Outputs.DeepCopyInto(&out.Outputs)
	if in.Results != nil {
		in, out := &in.Results, &out.Results
		if *in == nil {
			*out = nil
		} else {
			*out = new(Results)
			**out = **in
		}
	}
	if in.TaskRef != nil {
		in, out := &in.TaskRef, &out.TaskRef
		if *in == nil {
			*out = nil
		} else {
			*out = new(TaskRef)
			**out = **in
		}
	}
	if in.TaskSpec != nil {
		in, out := &in.TaskSpec, &out.TaskSpec
		if *in == nil {
			*out = nil
		} else {
			*out = new(TaskSpec)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		if *in == nil {
			*out = nil
		} else {
			*out = new(v1.Duration)
			**out = **in
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		if *in == nil {
			*out = nil
		} else {
			*out = new(core_v1.Affinity)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskRunSpec.
func (in *TaskRunSpec) DeepCopy() *TaskRunSpec {
	if in == nil {
		return nil
	}
	out := new(TaskRunSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskRunStatus) DeepCopyInto(out *TaskRunStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(duck_v1alpha1.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Results != nil {
		in, out := &in.Results, &out.Results
		if *in == nil {
			*out = nil
		} else {
			*out = new(Results)
			**out = **in
		}
	}
	if in.StartTime != nil {
		in, out := &in.StartTime, &out.StartTime
		if *in == nil {
			*out = nil
		} else {
			*out = new(v1.Time)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.CompletionTime != nil {
		in, out := &in.CompletionTime, &out.CompletionTime
		if *in == nil {
			*out = nil
		} else {
			*out = new(v1.Time)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]StepState, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskRunStatus.
func (in *TaskRunStatus) DeepCopy() *TaskRunStatus {
	if in == nil {
		return nil
	}
	out := new(TaskRunStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskSpec) DeepCopyInto(out *TaskSpec) {
	*out = *in
	if in.Inputs != nil {
		in, out := &in.Inputs, &out.Inputs
		if *in == nil {
			*out = nil
		} else {
			*out = new(Inputs)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Outputs != nil {
		in, out := &in.Outputs, &out.Outputs
		if *in == nil {
			*out = nil
		} else {
			*out = new(Outputs)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]core_v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]core_v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskSpec.
func (in *TaskSpec) DeepCopy() *TaskSpec {
	if in == nil {
		return nil
	}
	out := new(TaskSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskTrigger) DeepCopyInto(out *TaskTrigger) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskTrigger.
func (in *TaskTrigger) DeepCopy() *TaskTrigger {
	if in == nil {
		return nil
	}
	out := new(TaskTrigger)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestResult) DeepCopyInto(out *TestResult) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestResult.
func (in *TestResult) DeepCopy() *TestResult {
	if in == nil {
		return nil
	}
	out := new(TestResult)
	in.DeepCopyInto(out)
	return out
}

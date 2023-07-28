/*
Copyright 2023.

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

package v1

import (
	"errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var deploymentlog = logf.Log.WithName("deployment-resource")

func (r *Deployment) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-webapp-my-domain-v1-deployment,mutating=false,failurePolicy=fail,sideEffects=None,groups=webapp.my.domain,resources=deployments,verbs=create;update,versions=v1,name=vdeployment.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Deployment{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Deployment) ValidateCreate() (admission.Warnings, error) {
	deploymentlog.Info("validate create", "name", r.Name)

	return nil, r.validateKeptnMetric()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Deployment) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	deploymentlog.Info("validate update", "name", r.Name)

	return nil, r.validateKeptnMetric()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Deployment) ValidateDelete() (admission.Warnings, error) {
	deploymentlog.Info("validate delete", "name", r.Name)

	return nil, nil
}

func (r *Deployment) validateKeptnMetric() error {
	var allErrs field.ErrorList // defined as a list to allow returning multiple validation errors
	var err *field.Error
	if err = r.validateDeployment(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(
		schema.GroupKind{Group: "webapp", Kind: "Deployment"},
		r.Name,
		allErrs)
}

func (r *Deployment) validateDeployment() *field.Error {
	if r.spec.replicas < 3 {
		return field.Invalid(
			field.NewPath("spec").Child("replicas"),
			r.Spec.Range.Interval,
			errors.New("Forbidden! The time interval cannot be parsed. Please check for suitable conventions").Error(),
		)
	}
	return nil
}

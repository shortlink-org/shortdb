/*
Copyright 2022 Viktor Login.

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

package v1alpha1

import (
	"context"
	"strconv"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var shortdblog = logf.Log.WithName("shortdb-resource")

func (r *ShortDB) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, r).
		WithDefaulter(r).
		WithValidator(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-shortdb-shortdb-shortlink-v1alpha1-shortdb,mutating=true,failurePolicy=fail,sideEffects=None,groups=shortdb.shortdb.shortlink,resources=shortdbs,verbs=create;update,versions=v1alpha1,name=mshortdb.kb.io,admissionReviewVersions=v1

var _ admission.Defaulter[*ShortDB] = &ShortDB{}

// Default implements admission.Defaulter so a webhook will be registered for the type.
func (r *ShortDB) Default(_ context.Context, obj *ShortDB) error {
	shortdblog.Info("default", "name", obj.Name)

	// TODO(user): fill in your defaulting logic.
	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-shortdb-shortdb-shortlink-v1alpha1-shortdb,mutating=false,failurePolicy=fail,sideEffects=None,groups=shortdb.shortdb.shortlink,resources=shortdbs,verbs=create;update,versions=v1alpha1,name=vshortdb.kb.io,admissionReviewVersions=v1

var _ admission.Validator[*ShortDB] = &ShortDB{}

func (r *ShortDB) validateReplicas() error {
	var allErrs field.ErrorList

	if r.Spec.Deployments < 1 {
		fldPath := field.NewPath("spec").Child("deployments")
		allErrs = append(allErrs, field.Invalid(fldPath, strconv.Itoa(r.Spec.Deployments), "Deployments counter must be greater than zero"))
	}

	if len(allErrs) != 0 {
		return apierrors.NewInvalid(schema.GroupKind{Group: "shortdb.shortdb.shortlink", Kind: "ShortDB"}, r.Name, allErrs)
	}

	return nil
}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type.
func (r *ShortDB) ValidateCreate(_ context.Context, obj *ShortDB) (admission.Warnings, error) {
	shortdblog.Info("validate create", "name", obj.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, obj.validateReplicas()
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type.
func (r *ShortDB) ValidateUpdate(_ context.Context, _, newObj *ShortDB) (admission.Warnings, error) {
	shortdblog.Info("validate update", "name", newObj.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, newObj.validateReplicas()
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type.
func (r *ShortDB) ValidateDelete(_ context.Context, obj *ShortDB) (admission.Warnings, error) {
	shortdblog.Info("validate delete", "name", obj.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

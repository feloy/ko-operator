package controllers

import (
	"context"
	"fmt"

	kov1alpha1 "github.com/feloy/ko-operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

func (r *KoBuilderReconciler) setState(ctx context.Context, log logr.Logger, kobuilder *kov1alpha1.KoBuilder, state kov1alpha1.KoBuilderState) (err error) {
	log.Info(fmt.Sprintf("Set State * %s *", state))
	kobuilder.Status.State = state
	err = r.Status().Update(ctx, kobuilder)
	return
}

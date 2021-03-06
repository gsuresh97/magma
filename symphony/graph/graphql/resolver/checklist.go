// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

type checkListItemResolver struct{}

func (checkListItemResolver) Type(ctx context.Context, obj *ent.CheckListItem) (models.CheckListItemType, error) {
	return models.CheckListItemType(obj.Type), nil
}

type checkListItemDefinitionResolver struct{}

func (checkListItemDefinitionResolver) Type(ctx context.Context, obj *ent.CheckListItemDefinition) (models.CheckListItemType, error) {
	return models.CheckListItemType(obj.Type), nil
}

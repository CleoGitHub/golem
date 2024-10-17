package domainbuilder

import (
	"context"

	"github.com/cleogithub/golem/coredomaindefinition"
)

type Builder interface {
	WithModel(ctx context.Context, definition *coredomaindefinition.Model)
	WithRepository(ctx context.Context, definition *coredomaindefinition.Repository)
	WithRelation(ctx context.Context, definition *coredomaindefinition.Relation)
	WithCRUD(ctx context.Context, definition *coredomaindefinition.CRUD)
	WithUsecase(ctx context.Context, definition *coredomaindefinition.Usecase)
	Build(ctx context.Context) error
}

// Empty Builder will do nothing except returning a PanicBuilder
// to avoid using chained call without noticing that a EmptyBuilder war returned instead of the current Builder
type EmptyBuilder struct {
	err error
}

func (builder *EmptyBuilder) WithModel(ctx context.Context, definition *coredomaindefinition.Model) {
}

func (builder *EmptyBuilder) WithRepository(ctx context.Context, definition *coredomaindefinition.Repository) {
}

func (builder *EmptyBuilder) WithRelation(ctx context.Context, definition *coredomaindefinition.Relation) {
}

func (builder *EmptyBuilder) WithCRUD(ctx context.Context, definition *coredomaindefinition.CRUD) {
}

func (builder *EmptyBuilder) WithUsecase(ctx context.Context, definition *coredomaindefinition.Usecase) {
}

func (builder *EmptyBuilder) Build(ctx context.Context) error {
	return nil
}

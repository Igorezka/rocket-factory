package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	repoModel "github.com/Igorezka/rocket-factory/inventory/internal/repository/model"
)

type RepositorySuite struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	repository *repository

	generatedParts []repoModel.Part
}

func (s *RepositorySuite) SetupTest() {
	s.ctx = context.Background()

	s.repository = NewRepository()

	for _, v := range s.repository.data {
		s.generatedParts = append(s.generatedParts, v)
	}
}

func (s *RepositorySuite) TearDownTest() {}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

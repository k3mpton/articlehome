package service

import (
	"context"
	"fmt"
	publishhome "legi/newspapers/proto/gen/go"
	"log/slog"
)

type Saver interface {
	SaveArticle(
		ctx context.Context,
		NameArticle string,
		Description string,
		Rating float32,
	) (*publishhome.ArticleData, error)
}

type ArticleProvider interface {
	GetById(
		ctx context.Context,
		publishId int64,
	) (*publishhome.ArticleData, error)

	GetRandom(
		ctx context.Context,
	) (*publishhome.ArticleData, error)
}

type service struct {
	log      *slog.Logger
	saver    Saver
	provider ArticleProvider
}

func NewService(
	log *slog.Logger,
	saver Saver,
	provider ArticleProvider,
) service {
	return service{
		log:      log,
		saver:    saver,
		provider: provider,
	}
}

func (s *service) ArtCreate(
	ctx context.Context,
	Name string,
	Description string,
	Rating float32,
) (*publishhome.ArticleData, error) {
	const op = "service.ArtCreate"

	log := s.log.With(
		"op", op,
	)

	log.Info("process create Article...")

	art, err := s.saver.SaveArticle(ctx, Name, Description, Rating)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	log.Info("process Access created article!")

	return art, nil
}

func (s *service) GetArticle(
	ctx context.Context,
	publishId int64,
) (*publishhome.ArticleData, error) {
	const op = "service.GetArticle"

	log := s.log.With(
		"op", op,
	)

	log.Info("process get article...")

	art, err := s.provider.GetById(ctx, publishId)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	log.Info("process access getting article!")

	return art, nil

}

func (s *service) GetRandArticle(
	ctx context.Context,
) (*publishhome.ArticleData, error) {
	const op = "service.GetRandArticle"

	log := s.log.With(
		"op", op,
	)

	log.Info("proccess get rand article")

	art, err := s.provider.GetRandom(ctx)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	log.Info("process access get rand article! ")

	return art, nil
}

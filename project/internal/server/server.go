package server

import (
	"context"
	"fmt"
	publishhome "legi/newspapers/proto/gen/go"
	"log/slog"

	"google.golang.org/grpc"
)

type PublisherArticles interface {
	ArtCreate(
		ctx context.Context,
		Name string,
		Desritption string,
		Rating float32,
	) (*publishhome.ArticleData, error)

	GetArticle(
		ctx context.Context,
		Publish_id int64,
	) (*publishhome.ArticleData, error)

	GetRandArticle(
		ctx context.Context,
	) (*publishhome.ArticleData, error)
}

type Server struct {
	publishhome.UnimplementedPublishHomeServer
	Artcs PublisherArticles
}

func NewServer(grpcServer *grpc.Server, artcs PublisherArticles) {
	publishhome.RegisterPublishHomeServer(
		grpcServer,
		&Server{
			Artcs: artcs,
		},
	)
}

func (s *Server) CreatePublish(
	ctx context.Context,
	req *publishhome.CreateArticleRequest,
) (*publishhome.CreateArticleResponse, error) {
	const op = "server.Create"

	log := slog.With(
		"op", op,
	)

	log.Info("create article...")

	articleData, err := s.Artcs.ArtCreate(ctx, req.ArticleReq.NamePublish, req.ArticleReq.Description, req.ArticleReq.Rating)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	log.Info("Created Article!")

	return &publishhome.CreateArticleResponse{
		ArticleResp: articleData,
	}, nil
}

func (s *Server) GetById(
	ctx context.Context,
	req *publishhome.GetArticleByIDRequest,
) (*publishhome.GetArticleByIDResponse, error) {
	const op = "server.GetById"

	log := slog.With(
		"op", op,
	)

	log.Info("Geting article...")

	article, err := s.Artcs.GetArticle(ctx, req.ArticleId)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	return &publishhome.GetArticleByIDResponse{
		Article: article,
	}, nil
}

func (s *Server) GetIsTrendPublish(
	ctx context.Context,
	req *publishhome.GetRandomArticleIsTrendRequest,
) (*publishhome.GetRandomArticleIsTrendResponse, error) {
	const op = "server.GetRandomArticle"

	log := slog.With(
		"op", op,
	)

	log.Info("Gettign Random Article")

	article, err := s.Artcs.GetRandArticle(ctx)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)

	}

	return &publishhome.GetRandomArticleIsTrendResponse{
		Article: article,
	}, nil
}

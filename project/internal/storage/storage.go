package storage

import (
	"context"
	"fmt"
	"legi/newspapers/project/cmd/rediss"
	dbconnection "legi/newspapers/project/utils/DbConnection"
	publishhome "legi/newspapers/proto/gen/go"

	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"
)

type storage struct {
	postgresDBpool *pgxpool.Pool
	redis          *redis.Client
}

func NewStorage() *storage {
	return &storage{
		postgresDBpool: dbconnection.NewConnection(),
		redis:          rediss.InitRedis(),
	}
}

func (s *storage) SaveArticle(
	ctx context.Context,
	NameArticle string,
	Description string,
	Rating float32,
) (*publishhome.ArticleData, error) {
	const op = "storage.SaveArticle"

	query := `insert into Articles(name, description, Rating)
	values($1, $2, $3 )
	`
	s.postgresDBpool.QueryRow(context.Background(), query, NameArticle, Description, Rating)

	return &publishhome.ArticleData{
		NamePublish: NameArticle,
		Description: Description,
	}, nil
}

func (s *storage) GetById(
	ctx context.Context,
	publishId int64,
) (*publishhome.ArticleData, error) {
	const op = "storage.GetById"

	query := `select * from Articles where Id = $1`

	row := s.postgresDBpool.QueryRow(context.Background(), query, publishId)

	var Article publishhome.ArticleData

	err := row.Scan(
		&Article.PublishId,
		&Article.NamePublish,
		&Article.Description,
		&Article.Rating,
	)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	return &Article, nil
}

// func (s *storage) GetRandom(
// 	ctx context.Context,
// ) (*publishhome.ArticleData, error) {
// 	const op = "storage.GetByIdRandom"

// 	query := `select * from Articles order by RANDOM() LIMIT 1`

// 	row := s.postgresDBpool.QueryRow(context.Background(), query)

// 	var Article publishhome.ArticleData

// 	err := row.Scan(
// 		&Article.PublishId,
// 		&Article.NamePublish,
// 		&Article.Description,
// 		&Article.Rating,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("%v: %v", op, err)
// 	}

// 	return &Article, nil
// }

const (
	redisPoolKey = "articles:random_pool"
	poolSize     = 100
)

func (s *storage) GetRandom(ctx context.Context) (*publishhome.ArticleData, error) {
	const op = "storage.GetRandom"

	idStr, err := s.redis.LPop(ctx, redisPoolKey).Result()

	if err == redis.Nil {
		if err := s.fillRandomPool(ctx); err != nil {
			return nil, fmt.Errorf("%s: fill pool: %w", op, err)
		}
		idStr, _ = s.redis.LPop(ctx, redisPoolKey).Result()
	}

	if idStr == "" {
		return nil, fmt.Errorf("%s: no articles found", op)
	}

	query := `SELECT id, name, description, rating FROM Articles WHERE id = $1`
	var article publishhome.ArticleData
	err = s.postgresDBpool.QueryRow(ctx, query, idStr).Scan(
		&article.PublishId,
		&article.NamePublish,
		&article.Description,
		&article.Rating,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: db scan: %w", op, err)
	}

	return &article, nil
}

func (s *storage) fillRandomPool(ctx context.Context) error {
	query := fmt.Sprintf("SELECT id FROM Articles ORDER BY RANDOM() LIMIT %d", poolSize)

	rows, err := s.postgresDBpool.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []interface{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}

	if len(ids) > 0 {
		return s.redis.LPush(ctx, redisPoolKey, ids...).Err()
	}
	return nil
}

package main

import (
	"context"
	"fmt"
	publishhome "legi/newspapers/proto/gen/go"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {

	conn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := publishhome.NewPublishHomeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("=== Создание публикации ===")
	createResp, err := client.CreatePublish(ctx, &publishhome.CreateArticleRequest{
		ArticleReq: &publishhome.ArticleData{
			NamePublish: "fdsgdf fer",
			Description: "ffff",
			Rating:      4.8,
			CreatedAt:   timestamppb.Now(),
			UpdatedAt:   timestamppb.Now(),
		},
	})
	if err != nil {
		log.Printf("CreatePublish failed: %v", err)
	} else {
		fmt.Printf("Создана публикация")
		fmt.Printf("Название: %s\n", createResp.ArticleResp.NamePublish)
	}

	fmt.Println("\n=== Получение публикации по ID ===")
	getResp, err := client.GetById(ctx, &publishhome.GetArticleByIDRequest{
		ArticleId: 1,
	})
	if err != nil {
		log.Printf("GetById failed: %v", err)
	} else {
		fmt.Printf("Найдена публикация: %s\n", getResp.Article.NamePublish)
		fmt.Printf("Рейтинг: %.1f\n", getResp.Article.Rating)
	}

	fmt.Println("\n=== Получение трендовой публикации ===")
	trendResp, err := client.GetIsTrendPublish(ctx, &publishhome.GetRandomArticleIsTrendRequest{
		RandomArticleId: 1,
	})
	if err != nil {
		log.Printf("GetIsTrendPublish failed: %v", err)
	} else {
		fmt.Printf("Трендовая публикация: %s\n", trendResp.Article.NamePublish)
	}
}

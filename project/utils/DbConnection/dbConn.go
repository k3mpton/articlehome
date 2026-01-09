package dbconnection

import (
	"context"
	getenvfield "legi/newspapers/project/utils/GetEnvField"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewConnection() *pgxpool.Pool {
	dsn := getenvfield.Getenv("DATABASE_URL")

	// Парсим DSN в конфиг
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	// НАСТРАИВАЕМ ПУЛ под нагрузочное тестирование
	config.MaxConns = 90 // Увеличиваем лимит соединений
	config.MinConns = 5  // Минимальное количество соединений
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	// Таймауты на уровне соединения
	config.ConnConfig.ConnectTimeout = 5 * time.Second
	config.ConnConfig.RuntimeParams["statement_timeout"] = "3000"                    // 3 секунды
	config.ConnConfig.RuntimeParams["idle_in_transaction_session_timeout"] = "10000" // 10 секунд

	// Создаем пул с конфигом
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}

	// Пингуем с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping pool: %v", err)
	}

	// Логируем успешное создание пула
	log.Printf("PostgreSQL pool created. Max connections: %d", config.MaxConns)

	return pool
}

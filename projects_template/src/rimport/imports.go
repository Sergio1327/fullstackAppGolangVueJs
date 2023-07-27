package rimport

import (
	"log"
	"os"
	"projects_template/config"
	"projects_template/internal/repository/oracle"
	"projects_template/internal/transaction"
)

type RepositoryImports struct {
	Config         config.Config
	SessionManager transaction.SessionManager
	Repository     Repository
}

func NewRepositoryImports(
	sessionManager transaction.SessionManager,
) RepositoryImports {
	config, err := config.NewConfig(os.Getenv("CONF_PATH"))
	if err != nil {
		log.Fatalln(err)
	}

	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr:     config.RedisConnect(),
	// 	Password: config.Redis.Password,
	// 	DB:       0,
	// })

	// if err := redisClient.Ping(context.Background()); err.Err() != nil {
	// 	log.Fatalln("redis:", err)
	// }

	// ctrl := &gomock.Controller{}

	return RepositoryImports{
		Config:         config,
		SessionManager: sessionManager,
		Repository: Repository{
			Logger:   oracle.NewLogger(),
			Template: oracle.NewTemplate(),
		},
	}
}

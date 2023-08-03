package cmd

import (
	"context"
	"fmt"
	"github.com/mjedari/vgang-project/app/auth"
	"github.com/mjedari/vgang-project/app/collector"
	"github.com/mjedari/vgang-project/app/configs"
	"github.com/mjedari/vgang-project/app/router"
	"github.com/mjedari/vgang-project/app/wiring"
	"github.com/mjedari/vgang-project/domain/contracts"
	"github.com/mjedari/vgang-project/infra/healer"
	"github.com/mjedari/vgang-project/infra/rate_limiter"
	"github.com/mjedari/vgang-project/infra/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serving service.",
	Long:  `Serving service.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())

		serve(ctx)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
		fmt.Println()
		//for i := 10; i > 0; i-- {
		//	time.Sleep(time.Second * 1)
		//	fmt.Printf("\033[2K\rShutting down ... : %d", i)
		//}

		// Perform any necessary cleanup before exiting
		fmt.Println("\nService exited successfully.")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(ctx context.Context) {
	initWiring(ctx)

	// initiate auth and client
	authService := auth.NewAuthService(wiring.Wiring.GetStorage(), wiring.Wiring.Configs.OriginRemote)

	request := auth.NewLoginRequest()
	err := authService.Login(ctx, request)
	if err != nil {
		log.Fatalf("can not log into vgang: %v", err)
	}

	// initiate collector
	collectingService := collector.NewCollector(authService.Client,
		wiring.Wiring.GetStorage(), wiring.Wiring.Configs.Collector, wiring.Wiring.Configs.OriginRemote)
	go collectingService.Start(ctx)

	runHttpServer(ctx)
}

func runHttpServer(ctx context.Context) {
	r := router.NewRouter()

	address := net.JoinHostPort(configs.Config.Server.Host, configs.Config.Server.Port)
	logrus.WithField("HTTP_Host", configs.Config.Server.Host).
		WithField("HTTP_Port", configs.Config.Server.Port).
		Info("starting HTTP/REST vgang-project...")

	server := &http.Server{Addr: address, Handler: r}

	err := server.ListenAndServe()
	if err != nil {
		logrus.Fatal(err)
	}
}

func initWiring(ctx context.Context) {
	redisProvider, err := storage.NewRedis(configs.Config.Redis)
	if err != nil {
		logrus.Fatalf("Fatal error on create redis("+configs.Config.Redis.Host+":"+configs.Config.Redis.Port+")connection: %s \n", err)
	}

	rateLimiter := rate_limiter.NewRateLimiter(configs.Config.RateLimiter)

	wiring.Wiring = wiring.NewWire(redisProvider, rateLimiter, configs.Config)

	// init healer for services
	infraHealer := healer.NewHealerService([]contracts.IProvider{redisProvider}, 1)
	infraHealer.Start(ctx)

	logrus.Info("wiring initialized")
}

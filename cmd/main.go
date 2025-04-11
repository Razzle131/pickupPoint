package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Razzle131/pickupPoint/internal/handler"
	"github.com/Razzle131/pickupPoint/internal/repository/productRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/pvzRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/receptionRepo"
	"github.com/Razzle131/pickupPoint/internal/repository/userRepo"
	"github.com/Razzle131/pickupPoint/pkg/logger"
	"github.com/joho/godotenv"
)

//go:generate go tool oapi-codegen -config ../api/config.yaml ../api/swagger.yaml

func main() {
	godotenv.Load()

	cfg := handler.Config{
		Port: os.Getenv("SERVER_ADDRESS"),
		DSN:  os.Getenv("POSTGRES_CONN"),
	}

	// init logger
	logger.SetupLogging(slog.LevelDebug)

	slog.Debug("Debugging info enabled")

	slog.Info("Starting db", slog.String("DSN", os.Getenv("POSTGRES_CONN")))

	// dbConn, err := db.New(context.Background(), cfg.DSN)
	// if err != nil {
	// 	log.Error("Failed to connect to database", slog.String("error", err.Error()))
	// 	os.Exit(1)
	// }
	// defer dbConn.Close()

	// emailTest := types.Email("abobaboba@aa.ss")
	// res, err := emailTest.MarshalJSON()
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	return
	// }
	// slog.Info(string(res))

	// err = emailTest.UnmarshalJSON([]byte(`"abobaboba2@aa.ss"`))
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	return
	// }
	// slog.Info(string(emailTest))

	// t1 := time.Now()
	// t2 := time.Now().Add(time.Second * 10)

	// slog.Debug(fmt.Sprint(int(t1.Sub(t2))))
	// slog.Debug(fmt.Sprint(-2147483648 % 2147483647))

	myServer := handler.NewServer(userRepo.NewCache(), pvzRepo.NewCache(), productRepo.NewCache(), receptionRepo.NewCache())

	r := http.NewServeMux()

	r.HandleFunc("POST /dummyLogin", myServer.PostDummyLogin)
	r.HandleFunc("POST /login", myServer.PostLogin)
	r.HandleFunc("POST /register", myServer.PostRegister)
	r.HandleFunc("POST /pvz", myServer.PostPvz)
	r.HandleFunc("POST /products", myServer.PostProducts)
	r.HandleFunc("POST /receptions", myServer.PostReceptions)
	// r.HandleFunc("POST /api/auth", myServer.PostApiAuth)
	// r.HandleFunc("GET /api/buy/{item}", func(w http.ResponseWriter, r *http.Request) {
	// 	item := strings.Split(r.URL.Path, "/")[3]
	// 	myServer.GetApiBuyItem(w, r, item)
	// })
	// r.HandleFunc("GET /api/info", myServer.GetApiInfo)
	// r.HandleFunc("POST /api/sendCoin", myServer.PostApiSendCoin)

	srv := &http.Server{
		Addr:    "localhost:" + cfg.Port,
		Handler: r,
	}
	slog.Info("Starting server on address " + srv.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		//log.Error("Failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

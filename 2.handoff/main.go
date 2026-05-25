package main

// func NewIncidentStore(conf Config) (*mongo.Client, IncidentStore) {
// 	var client *mongo.Client = nil
// 	var store InstrumentedIncidentStore

// 	if conf.ConnectionString != "" {
// 		slog.Info("using mongo store", "db", conf.DatabaseName)
// 		client, err := mongo.Connect(options.Client().ApplyURI(conf.ConnectionString))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		mongoIncidentStore := NewMongoIncidentStore(client, conf.DatabaseName)
// 		store = InstrumentedIncidentStore{s: mongoIncidentStore}
// 	} else {
// 		slog.Info("no connection string, using in-memory store")
// 		store = InstrumentedIncidentStore{NewMemoryIncidentStore()}
// 	}
// 	return client, &store
// }

// func mainInit() (string, *http.Handler, *IncidentHandler) {
// 	config := loadConfig()
// 	client, IncidentStore := NewIncidentStore(config)
// 	registry := NewRegistry()
// 	promRegistry := prometheus.NewRegistry()
// 	NewMetrics(promRegistry)
// 	flagHandler := FlagHandler{store: CreateFlagStore()}
// 	incHandler := IncidentHandler{IncidentStore: IncidentStore, Registry: registry, FlagEvaluator: &flagHandler.store}
// 	router := getRouter(&incHandler, &flagHandler, client, promRegistry)
// 	return config.Port, &router, &incHandler
// }

// func main() {
// 	port, router, incHandler := mainInit()
// 	srv := http.Server{
// 		Addr:    ":" + port,
// 		Handler: *router,
// 	}
// 	go incHandler.Registry.run()
// 	defer close(incHandler.Registry.done)

// 	go func() {
// 		slog.Info(fmt.Sprintf("server starting port=%s", srv.Addr))
// 		err := srv.ListenAndServe()
// 		if err != nil && err != http.ErrServerClosed {
// 			log.Fatal(err)
// 		}
// 	}()
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
// 	<-quit

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	slog.Info("server shut down in <= 10 sec")
// 	srv.Shutdown(ctx)
// 	slog.Info("server shut down gracefully")
// }

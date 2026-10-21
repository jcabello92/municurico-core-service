package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"core-service/internal/config"
	dbPostgres "core-service/internal/database/postgresql"
	dbRedis "core-service/internal/database/redis"
	"core-service/internal/handlers"
	"core-service/internal/interceptors"
	reposPostgres "core-service/internal/repositories/postgresql"
	reposRedis "core-service/internal/repositories/redis"
	"core-service/internal/security"
	"core-service/internal/services"
	pb "core-service/proto"
)

func main() {
	cfg := config.LoadEnv()
	log.Printf("Iniciando servicio en modo: %s", cfg.AppEnv)

	// 1. Conexión, migraciones y seeders en PostgreSQL
	psqlDSN := dbPostgres.BuildPostgreSQLDSN(cfg.PSQLHost, cfg.PSQLPort, cfg.PSQLUser, cfg.PSQLPassword, cfg.PSQLDB)
	db := dbPostgres.ConnectPostgreSQL(cfg.PSQLDB, psqlDSN)
	dbPostgres.RunMigrations(db)
	dbPostgres.RunSeeders(db)

	// 2. Conexión a Redis
	redisDSN := dbRedis.BuildRedisDSN(cfg.RedisHost, cfg.RedisPort)
	redisClient := dbRedis.ConnectRedis(redisDSN, cfg.RedisPassword, cfg.RedisDB)

	// 3. Inicialización de componentes de seguridad
	hasher := security.NewPasswordHasher(cfg.PasswordBcryptCost)

	encryptor, err := security.NewEncryptor([]byte(cfg.EncryptionSecret))
	if err != nil {
		log.Fatalf("ERROR CRÍTICO: Fallo al inicializar Encryptor: %v", err)
	}

	bindexer, err := security.NewBindexer([]byte(cfg.BindexSecret))
	if err != nil {
		log.Fatalf("ERROR CRÍTICO: Fallo al inicializar Bindexer: %v", err)
	}

	jwtManager := security.NewJWT(cfg.JWTSecret)

	// 4. Inicialización de repositorios
	reposContainer := reposPostgres.NewContainer(db)
	redisSesionRepo := reposRedis.NewSesionRepository(redisClient)
	redisRecuperacionRepo := reposRedis.NewRecuperacionRepository(redisClient)

	// 5. Inicialización de servicios
	accesoService := services.NewAccesoService(
		reposContainer,
		redisSesionRepo,
		redisRecuperacionRepo,
		hasher,
		encryptor,
		bindexer,
		jwtManager,
		cfg.SessionTTL,
	)

	// 6. Inicialización de handlers
	accesoHandler := handlers.NewAccesoHandler(accesoService)

	// 7. Configuración e inicio del servidor gRPC
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.AuditoriaInterceptor()),
	)
	pb.RegisterAccesoServiceServer(grpcServer, accesoHandler)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("acceso", grpc_health_v1.HealthCheckResponse_SERVING)

	addr := fmt.Sprintf("%s:%s", cfg.GRPCHost, cfg.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("ERROR CRÍTICO: No se pudo abrir el puerto gRPC en %s: %v", addr, err)
	}

	go func() {
		log.Printf("Servidor gRPC escuchando en %s...", addr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("ERROR CRÍTICO: Fallo en ejecución del servidor gRPC: %v", err)
		}
	}()

	// 8. Espera de señales de apagado limpio (graceful shutdown)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Deteniendo servidor gRPC limpiamente...")
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	healthServer.SetServingStatus("acceso", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	grpcServer.GracefulStop()
	log.Println("Servicio detenido con éxito.")
}
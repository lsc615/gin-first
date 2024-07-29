package main

import (
	"context"
	"fmt"
	"github.com/shicli/gin-first/common"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shicli/gin-first/route"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// @title			gin first
// @version		1.0
// @description	This is gin
// @contact.name	shicli
func main() {
	InitConfig()
	common.InitDB()

	cleanup, err := InitTrace()
	log.Println(cleanup)
	if err != nil {
		log.Printf("Failed to initialize tracer: %v", err)
		return
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// 清理 OTLP exporter
		if err := cleanup(ctx); err != nil {
			log.Printf("Failed to shut down OTLP exporter: %v", err)
		}
	}()

	r := gin.Default()
	r.Use(otelgin.Middleware(serviceName))
	r = route.CollectRoute(r)
	port := viper.GetString("server.port")
	if err := r.Run("127.0.0.1:" + port); err != nil {
		panic(err)
	}
}

func InitConfig() {
	var ConfigPath string
	pflag.StringVarP(&ConfigPath, "", "c", "", "配置文件路径")
	pflag.Parse()

	if ConfigPath == "" {
		fmt.Println("请提供配置文件路径")
		return
	}

	//设置文件
	viper.SetConfigFile(ConfigPath)

	//读取文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("无法读取配置文件：%s", err))
	}
}

var (
	serviceName  = getEnv("SERVICE_NAME", "test")
	collectorURL = getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317") // Assuming localhost as default
	insecure     = getEnv("INSECURE", "true")
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func InitTrace() (func(context.Context) error, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 定义 grpc 连接是否采用安全模式
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if insecure == "true" || insecure == "1" {
		secureOption = otlptracegrpc.WithInsecure()
	}
	// 创建 otlp exporter
	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating otlp exporter: %w", err)
	}

	// 设置 resource
	resources, err := resource.New(
		ctx,
		resource.WithAttributes(
			// 我们可以使用该方法,来进行添加自定义属性
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		return nil, err
	}

	//创建 stdout exporter
	//exporter, err := stdouttrace.New(
	//	stdouttrace.WithPrettyPrint(),
	//)
	//if err != nil {
	//	return nil, fmt.Errorf("creating stdout exporter: %w", err)
	//}

	// 创建TracerProvider
	otlpProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		// 默认为 5s,为便于演示,设置为 1s
		trace.WithBatcher(exporter, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(resources),
	)

	// 创建TracerProvider
	otel.SetTracerProvider(otlpProvider)

	// 设置传播上下文的处理器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return exporter.Shutdown, nil
}

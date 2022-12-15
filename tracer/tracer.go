package tracer

import (
	"github.com/gin-gonic/gin"
	"github.com/joaootav/system_supermarket/infra/opentel"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

// Add zipkin log service
func SetupTracer(router *gin.Engine) {

	ot := opentel.NewOpenTel()
	ot.ServiceName = "GoApp"
	ot.ServiceVersion = "0.1"
	ot.ExporterEndpoint = "http://localhost:9411/api/v2/spans"
	Tracer = ot.GetTracer()
	router.Use(otelgin.Middleware(ot.ServiceName))
}

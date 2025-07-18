package middleware

import (
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
)

func TracingMiddleware(next http.Handler) http.Handler {
	return xray.Handler(xray.NewFixedSegmentNamer("KaigoInsuranceSystem"), next)
}

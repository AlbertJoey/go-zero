package trace

import (
	"context"
	"fmt"
	"github.com/SkyAPM/go2sky/propagation"
	"hash/crc32"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Exporter struct {
	Tracer *go2sky.Tracer
}

var _ sdktrace.SpanExporter = (*Exporter)(nil)

func NewSkywalking(endpoint, serviceName string) (*Exporter, error) {
	rp, err := reporter.NewGRPCReporter(endpoint, reporter.WithCheckInterval(time.Second))
	if err != nil {
		return nil, err
	}
	tracer, err := go2sky.NewTracer(serviceName, go2sky.WithReporter(rp))
	e := &Exporter{
		Tracer: tracer,
	}
	return e, nil
}

func (e *Exporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	for _, s := range spans {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		fmt.Println("-----------------start------------------")
		fmt.Println("name:", s.Name())
		fmt.Println("SpanContext:", s.SpanContext())
		fmt.Println("Parent:", s.Parent())
		fmt.Println("SpanKind:", s.SpanKind())
		fmt.Println("StartTime:", s.StartTime())
		fmt.Println("EndTime:", s.EndTime())
		fmt.Println("Attributes:", s.Attributes())
		fmt.Println("SpanContext:", s.SpanContext())
		fmt.Println("Links:", s.Links())
		fmt.Println("Events:", s.Events())
		fmt.Println("Status:", s.Status())
		fmt.Println("InstrumentationScope:", s.InstrumentationScope())
		fmt.Println("InstrumentationLibrary:", s.InstrumentationLibrary())
		fmt.Println("Resource:", s.Resource())
		fmt.Println("DroppedAttributes:", s.DroppedAttributes())
		fmt.Println("DroppedLinks:", s.DroppedLinks())
		fmt.Println("DroppedEvents:", s.DroppedEvents())
		fmt.Println("ChildSpanCount:", s.DroppedEvents())

		span, _, err := e.Tracer.CreateEntrySpan(ctx, s.Name(), func(key string) (string, error) {
			correlationContext := make(map[string]string)
			for _, kv := range s.Attributes() {
				correlationContext[string(kv.Key)] = kv.Value.AsString()
			}
			spanid := ([8]byte)(s.Parent().SpanID())
			sid := crc32.ChecksumIEEE(spanid[:]) / 2
			scx := propagation.SpanContext{
				Sample:                0,
				TraceID:               s.SpanContext().TraceID().String(),
				ParentSegmentID:       s.Parent().TraceID().String(),
				ParentSpanID:          int32(sid),
				ParentService:         s.Parent().TraceID().String(),
				ParentServiceInstance: s.Parent().TraceID().String(),
				ParentEndpoint:        s.Parent().TraceID().String(),
				AddressUsedAtClient:   s.SpanKind().String(),
				CorrelationContext:    correlationContext,
			}
			return scx.EncodeSW8(), nil
		})
		if err != nil {
			fmt.Println("err:", err)
		}
		span.SetOperationName(s.Name())
		span.SetPeer(s.Name())
		span.End()
	}
	return nil
}

// Shutdown stops the exporter flushing any pending exports.
func (e *Exporter) Shutdown(ctx context.Context) error {
	return nil
}

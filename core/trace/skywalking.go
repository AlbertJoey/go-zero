package trace

import (
	"context"
	"fmt"
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
		//span, _, err := e.Tracer.CreateEntrySpan(ctx, s.Name(), func(key string) (string, error) {
		//	scx := propagation.SpanContext{}
		//	if !s.Parent().TraceID().IsValid() { //parent
		//		spanid := ([8]byte)(s.SpanContext().SpanID())
		//		sid := crc32.ChecksumIEEE(spanid[:]) / 2
		//		scx = propagation.SpanContext{
		//			Sample:                1,
		//			TraceID:               s.SpanContext().TraceID().String(),
		//			ParentSegmentID:       s.Parent().SpanID().String(),
		//			ParentSpanID:          int32(sid),
		//			ParentService:         s.Name(),
		//			ParentServiceInstance: s.Name(),
		//			ParentEndpoint:        s.Name(),
		//			AddressUsedAtClient:   s.Name(),
		//		}
		//	} else { //child
		//		spanid := ([8]byte)(s.Parent().SpanID())
		//		sid := crc32.ChecksumIEEE(spanid[:]) / 2
		//		scx = propagation.SpanContext{
		//			Sample:                1,
		//			TraceID:               s.SpanContext().TraceID().String(),
		//			ParentSegmentID:       s.Parent().SpanID().String(),
		//			ParentSpanID:          int32(sid),
		//			ParentService:         s.Name(),
		//			ParentServiceInstance: s.Name(),
		//			ParentEndpoint:        s.Name(),
		//			AddressUsedAtClient:   s.Name(),
		//		}
		//	}
		//
		//	return scx.EncodeSW8(), nil
		//})
		//if err != nil {
		//	fmt.Println("err:", err)
		//}
		//span.SetComponent(8888)
		//span.Tag(go2sky.TagURL, s.Name())
		//span.SetSpanLayer(0)
		//span.Tag(go2sky.TagStatusCode, "200")
		//span.Tag(go2sky.TagURL, s.Name())
		//span.End()
	}
	return nil
}

// Shutdown stops the exporter flushing any pending exports.
func (e *Exporter) Shutdown(ctx context.Context) error {
	return nil
}

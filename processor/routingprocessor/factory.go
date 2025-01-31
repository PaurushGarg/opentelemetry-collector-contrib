// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package routingprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

const (
	// The value of "type" key in configuration.
	typeStr = "routing"
	// The stability level of the processor.
	stability = component.StabilityLevelBeta

	scopeName = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor"
	nameSep   = "/"

	processorKey             = "processor"
	metricSep                = "_"
	nonRoutedSpansKey        = "non_routed_spans"
	nonRoutedMetricPointsKey = "non_routed_metric_points"
	nonRoutedLogRecordsKey   = "non_routed_log_records"
)

// NewFactory creates a factory for the routing processor.
func NewFactory() processor.Factory {
	return processor.NewFactory(
		typeStr,
		createDefaultConfig,
		processor.WithTraces(createTracesProcessor, stability),
		processor.WithMetrics(createMetricsProcessor, stability),
		processor.WithLogs(createLogsProcessor, stability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		AttributeSource: defaultAttributeSource,
		ErrorMode:       ottl.PropagateError,
	}
}

func createTracesProcessor(_ context.Context, params processor.CreateSettings, cfg component.Config, nextConsumer consumer.Traces) (processor.Traces, error) {
	warnIfNotLastInPipeline(nextConsumer, params.Logger)
	return newTracesProcessor(params.TelemetrySettings, cfg)
}

func createMetricsProcessor(_ context.Context, params processor.CreateSettings, cfg component.Config, nextConsumer consumer.Metrics) (processor.Metrics, error) {
	warnIfNotLastInPipeline(nextConsumer, params.Logger)
	return newMetricProcessor(params.TelemetrySettings, cfg)
}

func createLogsProcessor(_ context.Context, params processor.CreateSettings, cfg component.Config, nextConsumer consumer.Logs) (processor.Logs, error) {
	warnIfNotLastInPipeline(nextConsumer, params.Logger)
	return newLogProcessor(params.TelemetrySettings, cfg)
}

func warnIfNotLastInPipeline(nextConsumer interface{}, logger *zap.Logger) {
	_, ok := nextConsumer.(component.Component)
	if ok {
		logger.Warn("another processor has been defined after the routing processor: it will NOT receive any data!")
	}
}

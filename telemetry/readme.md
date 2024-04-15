```
<dependency>
    <groupId>io.opentelemetry</groupId>
    <artifactId>opentelemetry-api</artifactId>
    <version>1.0.0</version> <!-- or the latest version -->
</dependency>
<dependency>
    <groupId>io.opentelemetry</groupId>
    <artifactId>opentelemetry-exporter-datadog</artifactId>
    <version>1.0.0</version> <!-- or the latest version -->
</dependency>
```
```
import io.opentelemetry.api.metrics.*;
import io.opentelemetry.api.metrics.common.Labels;
import io.opentelemetry.exporter.datadog.DatadogMetricExporter;
import io.opentelemetry.exporter.datadog.DatadogMetricExporterBuilder;

public class CustomMetricExample {

    public static void main(String[] args) {
        // Create a Datadog Metric Exporter
        MetricExporter exporter = configureDatadogExporter();

        // Initialize the MeterProvider
        MeterProvider meterProvider = configureMeterProvider(exporter);

        // Create a Meter
        Meter meter = meterProvider.get("custom_metrics");

        // Create a Counter metric
        Counter counterMetric = meter.counterBuilder("custom_counter")
                .setDescription("This is a custom counter metric")
                .setUnit("1")
                .build();

        // Increment the Counter metric
        counterMetric.add(1, Labels.empty());

        // Flush the exporter
        exporter.shutdown();
    }

    private static MetricExporter configureDatadogExporter() {
        // Configure Datadog exporter
        return DatadogMetricExporter.builder()
                .setApiKey("<YOUR_API_KEY>")
                .buildAndRegisterGlobal();
    }

    private static MeterProvider configureMeterProvider(MetricExporter exporter) {
        // Create a MeterProvider
        return MetricProvider.get()
                .builder()
                .setMetricExporter(exporter)
                .buildAndRegisterGlobal();
    }
}
```

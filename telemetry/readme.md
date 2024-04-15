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

```
import io.opentelemetry.api.common.Attributes;
import io.opentelemetry.api.metrics.LongCounter;
import io.opentelemetry.api.metrics.Meter;
import io.opentelemetry.api.metrics.MeterProvider;
import io.opentelemetry.sdk.common.CompletableResultCode;
import io.opentelemetry.sdk.metrics.SdkMeterProvider;
import io.opentelemetry.sdk.metrics.data.AggregationTemporality;
import io.opentelemetry.sdk.metrics.data.MetricData;
import io.opentelemetry.sdk.metrics.export.MetricExporter;
import io.opentelemetry.sdk.resources.Resource;

import java.io.IOException;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.Collection;

public class CustomMetricExample {
    public static void main(String[] args) {
        // Initialize the OpenTelemetry SDK
        MeterProvider meterProvider = SdkMeterProvider.builder()
                .setResource(Resource.create(Attributes.of("service.name", "my-service")))
                .build();

        // Register a custom exporter
        meterProvider.addMetricReader(new DatadogMetricExporter());

        // Get a meter instance
        Meter meter = meterProvider.get("my-service");

        // Create a custom metric
        LongCounter counter = meter.counterBuilder("my_custom_metric")
                .setDescription("A custom metric")
                .setUnit("requests")
                .build();

        // Record a value for the custom metric
        counter.add(10, Attributes.of("endpoint", "/myEndpoint"));
    }

    private static class DatadogMetricExporter implements MetricExporter {
        @Override
        public CompletableResultCode export(Collection<MetricData> metrics) {
            for (MetricData metric : metrics) {
                try {
                    sendMetricToDatadog(metric);
                } catch (IOException e) {
                    return CompletableResultCode.ofFailure();
                }
            }
            return CompletableResultCode.ofSuccess();
        }

        private void sendMetricToDatadog(MetricData metric) throws IOException {
            // Construct the Datadog API request
            String metricName = metric.getInstruments().get(0).getName();
            double metricValue = metric.getInstruments().get(0).getLastValue();
            Attributes attributes = metric.getInstruments().get(0).getAttributes();

            URL url = new URL("https://api.datadoghq.com/api/v1/series");
            HttpURLConnection connection = (HttpURLConnection) url.openConnection();
            connection.setRequestMethod("POST");
            connection.setRequestProperty("Content-Type", "application/json");
            connection.setRequestProperty("DD-API-KEY", "your_datadog_api_key");

            String payload = "{\"series\":[{\"metric\":\"" + metricName + "\",\"points\":[[" + System.currentTimeMillis() / 1000 + "," + metricValue + "]],\"tags\":[";
            for (String key : attributes.getKeys()) {
                payload += "\"" + key + ":" + attributes.get(key) + "\",";
            }
            payload = payload.substring(0, payload.length() - 1) + "]}]}";

            connection.setDoOutput(true);
            connection.getOutputStream().write(payload.getBytes());

            int responseCode = connection.getResponseCode();
            if (responseCode != HttpURLConnection.HTTP_OK) {
                throw new IOException("Error sending metric to Datadog: " + responseCode);
            }
        }

        @Override
        public CompletableResultCode flush() {
            return CompletableResultCode.ofSuccess();
        }

        @Override
        public CompletableResultCode shutdown() {
            return CompletableResultCode.ofSuccess();
        }
    }
}
```
opentelemetry-exporter-datadog

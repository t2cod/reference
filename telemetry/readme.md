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




If you're unable to download the opentelemetry-exporter-datadog artifact from the Maven repositories, you can manually copy the artifact to your local Maven repository.

Here are the steps to do that:

Download the opentelemetry-exporter-datadog artifact from the Sonatype Nexus repository:
Go to the Sonatype Nexus repository: https://oss.sonatype.org/content/repositories/snapshots/io/opentelemetry/exporter/opentelemetry-exporter-datadog/
Find the latest version of the artifact (e.g., 1.19.0-SNAPSHOT) and download the JAR file.
Locate your local Maven repository directory. The default location is usually:
Windows: %USERPROFILE%\.m2\repository
macOS/Linux: ~/.m2/repository
Navigate to the correct directory structure within your local Maven repository:
io/opentelemetry/exporter/opentelemetry-exporter-datadog/1.19.0-SNAPSHOT/
Create the directories if they don't already exist.
Copy the downloaded opentelemetry-exporter-datadog-1.19.0-SNAPSHOT.jar file to the directory you created in the previous step.
Update your project's pom.xml file to reference the local artifact:



```
import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import io.opentelemetry.api.common.Attributes;
import io.opentelemetry.api.metrics.LongCounter;
import io.opentelemetry.api.metrics.Meter;
import io.opentelemetry.api.metrics.MeterProvider;
import io.opentelemetry.sdk.common.CompletableResultCode;
import io.opentelemetry.sdk.metrics.SdkMeterProvider;
import io.opentelemetry.sdk.metrics.data.MetricData;
import io.opentelemetry.sdk.metrics.export.MetricExporter;
import io.opentelemetry.sdk.resources.Resource;

import java.io.IOException;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.util.Collection;

public class CustomMetricExample implements MetricExporter {
    public static void main(String[] args) {
        // Initialize the OpenTelemetry SDK
        MeterProvider meterProvider = SdkMeterProvider.builder()
                .setResource(Resource.create(Attributes.of("service.name", "my-service")))
                .build();

        // Register the custom exporter
        meterProvider.addMetricReader(new CustomMetricExample());

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

        JsonObject payload = new JsonObject();
        JsonArray series = new JsonArray();
        JsonObject dataPoint = new JsonObject();
        dataPoint.addProperty("metric", metricName);
        dataPoint.addProperty("points", System.currentTimeMillis() / 1000 + "," + metricValue);
        dataPoint.add("tags", buildTagsJson(attributes));
        series.add(dataPoint);
        payload.add("series", series);

        URL url = new URL("https://api.datadoghq.com/api/v1/series");
        HttpURLConnection connection = (HttpURLConnection) url.openConnection();
        connection.setRequestMethod("POST");
        connection.setRequestProperty("Content-Type", "application/json");
        connection.setRequestProperty("DD-API-KEY", "your_datadog_api_key");

        connection.setDoOutput(true);
        try (OutputStream os = connection.getOutputStream()) {
            os.write(payload.toString().getBytes(StandardCharsets.UTF_8));
        }

        int responseCode = connection.getResponseCode();
        if (responseCode != HttpURLConnection.HTTP_OK) {
            throw new IOException("Error sending metric to Datadog: " + responseCode);
        }
    }

    private JsonArray buildTagsJson(Attributes attributes) {
        JsonArray tagsArray = new JsonArray();
        for (String key : attributes.getKeys()) {
            JsonObject tagObject = new JsonObject();
            tagObject.addProperty(key, attributes.get(key));
            tagsArray.add(tagObject);
        }
        return tagsArray;
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
```


```
import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.api.metrics.LongCounter;
import io.opentelemetry.api.metrics.Meter;
import io.opentelemetry.api.metrics.MeterProvider;
import io.opentelemetry.api.metrics.common.Labels;
import io.opentelemetry.exporter.logging.LoggingMetricExporter;
import io.opentelemetry.exporter.logging.LoggingMetricExporterBuilder;
import io.opentelemetry.sdk.metrics.SdkMeterProvider;

import java.io.OutputStreamWriter;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.TimeUnit;

public class DatadogMetricsSender {

    private static final String DATADOG_API_URL = "https://api.datadoghq.com/api/v1/series?api_key=";

    public static void main(String[] args) {
        // Create a meter provider
        MeterProvider meterProvider = SdkMeterProvider.builder().buildAndRegisterGlobal();

        // Create a meter
        Meter meter = meterProvider.get("custom_metrics");

        // Create a counter metric
        LongCounter counterMetric = meter.counterBuilder("custom_counter")
                .setDescription("This is a custom counter metric")
                .setUnit("1")
                .build();

        // Increment the counter metric
        counterMetric.add(1, Labels.empty());

        // Create an OpenTelemetry object to get the default OpenTelemetry instance
        OpenTelemetry openTelemetry = OpenTelemetrySdk.get();

        // Add a logging metric exporter (for demonstration)
        LoggingMetricExporter exporter = new LoggingMetricExporterBuilder().build();
        meterProvider.getMetricProducer().registerMetricExporter(exporter);

        // Optionally, you can add other exporters such as Datadog exporter here
        MeterExporter datadogExporter = new MeterExporter();

        // Register the exporter
        meterProvider.getMetricProducer().registerMetricExporter(datadogExporter);

        // Flush the exporter before shutdown
        Runtime.getRuntime().addShutdownHook(new Thread(exporter::shutdown));
        Runtime.getRuntime().addShutdownHook(new Thread(datadogExporter::shutdown));

        // Sleep to demonstrate metric collection
        try {
            TimeUnit.SECONDS.sleep(5);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

    static class MeterExporter implements MetricExporter {
        @Override
        public void export(Collection<MetricRecord> records) {
            for (MetricRecord record : records) {
                // Iterate through metric records and send each one to Datadog
                try {
                    sendMetric(record.getName(), record.getValue(), System.currentTimeMillis() / 1000, "YOUR_API_KEY");
                } catch (Exception e) {
                    System.err.println("Failed to send metric to Datadog: " + e.getMessage());
                }
            }
        }

        @Override
        public CompletableResultCode flush() {
            // No need to implement flushing for HTTP exporter
            return CompletableResultCode.ofSuccess();
        }

        @Override
        public CompletableResultCode shutdown() {
            // Clean up resources if needed
            return CompletableResultCode.ofSuccess();
        }

        private void sendMetric(String metricName, long value, long timestamp, String apiKey) throws Exception {
            // Create JSON payload
            String payload = String.format("[{\"metric\":\"%s\",\"points\":[[%d,%d]],\"type\":\"gauge\"}]", metricName, timestamp, value);

            // Create HTTP connection
            URL url = new URL(DATADOG_API_URL + apiKey);
            HttpURLConnection conn = (HttpURLConnection) url.openConnection();
            conn.setRequestMethod("POST");
            conn.setRequestProperty("Content-Type", "application/json");
            conn.setDoOutput(true);

            // Write payload to the connection
            try (OutputStreamWriter writer = new OutputStreamWriter(conn.getOutputStream(), StandardCharsets.UTF_8)) {
                writer.write(payload);
                writer.flush();
            }

            // Check response code
            int responseCode = conn.getResponseCode();
            if (responseCode != HttpURLConnection.HTTP_OK) {
                throw new Exception("Failed to send metric to Datadog. Response code: " + responseCode);
            }

            // Close connection
            conn.disconnect();
        }
    }
}
```

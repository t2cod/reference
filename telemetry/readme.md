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

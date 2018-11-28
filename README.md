# HomeIOT Prometheus Exporter

A simple go application that exposes an API to create Prometheus gauge metrics.

I personally use it in my Home IOT network with Raspberry Pi's and Arduino's to expose sensor values to prometheus so that I can view them on Grafana, but this is so generic that it can be used to expose anything.


## Usage

Run the server with `go run main.go`

Metrics are exposed at `http://localhost:8080/metrics`

To add metrics make a post request to `/metrics` with the following JSON format in the body:
```json
{
  "metrics": [
    {
      "name": "Dummy Name",
      "type": "Dummy Type",
      "value": 4
    },{
      "name": "Another Dummy Name",
      "type": "Another Dummy Type",
      "value": 12.8
    }
  ]
}
```

This will create 2 gauge metrics called `iot_metric` with `name` and `type` as labels and `value` as the gauge value.
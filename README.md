# Metrics writer

[![Build Status](https://travis-ci.org/mrgleam/metrics_writer.svg?branch=master)](https://travis-ci.org/mrgleam/metrics_writer)

 Prometheus push gateway is great tools for push metrics from your application but you must count the metrics in your application before send to Prometheus push gateway.

 Metrics writer, you can send name and labels only via HTTP and Metrics writer will count your metrics automaticly.
## Example
```
go run main.go
```

```
curl -XPOST -d 'name=test&labels={"labels1":"test1","labels2":"test2"}' http://localhost:8080/counter/metrics
```

## SS

![ss1](https://raw.githubusercontent.com/mrgleam/metrics_writer/master/example/screenshot/Screen%20Shot%202561-04-23%20at%2021.30.53.png)

![ss2](https://raw.githubusercontent.com/mrgleam/metrics_writer/master/example/screenshot/Screen%20Shot%202561-04-23%20at%2021.30.39.png)

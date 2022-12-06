# loki-logger
A go based logger for Loki with [ZeroLog](https://github.com/rs/zerolog) adapter.


## Usage with Grafana backend plugins.

Step 1.  Setup a [Free Grafana Cloud Account with Loki/Logging](https://grafana.com/logs/)

Step 2.  Add settings to a plugin in your grafana.ini
```
[plugin.my-plugin-id]
# the logger to use:  console, file, loki
logger = loki
# the log level to capture
logger_level = 0
logger_url = https://logs-prod3.grafana.net/loki/api/v1/push
# logger_key is base64 encoded
logger_key = user:apikey
# size of buffer lines to keep in memory
logger_buffer = 5 
# comma delimited list of labels to apply to the log message
logger_labels = plugin:dynatrace
# log all http requests from the plugin to external rest apis
logger_http = true
```

Step 3.  Go get "github.com/scottlepp/loki-logger/pkg/log"

Step 4.  Add logging to your code. These will be sent to loki if your logger_level (above) is low enough.
```
log.Debug("foo")
```

Step 5. (optional) Log all incoming requests to the plugin

Typically when setting up your plugin you supply handlers for query, health, resources like so:
```
	im := datasource.NewInstanceManager(getInstance)
	handler := &Plugin{
		IM: im,
	}
	router := host.getCallResourceHandler()
	return datasource.ServeOpts{
		QueryDataHandler:    handler,
		CheckHealthHandler:  handler,
		CallResourceHandler: httpadapter.New(router),
	}
```
Just wrap the plugin that handles the calls:
```
	im := datasource.NewInstanceManager(getInstance)
	handler := &Plugin{
		IM: im,
	}
	logHandler := log.Handler{
		Plugin: handler,
	}
	router := host.getCallResourceHandler()
	return datasource.ServeOpts{
		QueryDataHandler:    logHandler,
		CheckHealthHandler:  logHandler,
		CallResourceHandler: httpadapter.New(router),
	}
```

Step 6. (optional) Log all http calls to external rest apis.
When creating your http client, wrap the transport 
```
c.Transport = log.NewHTTPLogger("my-plugin-id", c.Transport)
```

## Usage with [ZeroLog](https://github.com/rs/zerolog)

```
logger := &lokiLogger{
	URL:        "https://logs-prod3.grafana.net/loki/api/v1/push",
	Key:        "user:key", // base64 encoded
	BufferSize: 50,
	Level:      0,
	Labels:     "label1:1","label2:2"
}
  
zl := zerolog.New(lokiLogger).Level(zerolog.Level(logger.Level))
zl.Debug("foo")
```

# k8sprobe

A lightweight and extensible Kubernetes health probe library for managing **liveness** and **readiness probes** in Go.

This library simplifies the process of exposing application health status to Kubernetes through HTTP endpoints, ensuring seamless integration with Kubernetes' lifecycle management.

---

## Features

- **Simple API**: Create liveness and readiness probes with minimal setup.
- **Custom Health Probes**: Define and manage custom health checks with ease.
- **Thread-Safe**: Built-in synchronization for probe state updates.
- **Kubernetes-Ready**: Designed to integrate with Kubernetes health checks (`livenessProbe` and `readinessProbe`).
- **HTTP Handler**: Ready-to-use handlers for serving probe status over HTTP.

---

## Installation

Install the library with:

```shell script
go get github.com/Autodoc-Technology/k8sprobe
```

Then import it into your Golang application:

```textmate
import "github.com/Autodoc-Technology/k8sprobe"
```

---

## Usage

### Simple Example: Liveness Probe

Here’s a basic example to set up a **liveness probe** that is automatically invalidated after 10 seconds:

```textmate
package main

import (
	"github.com/Autodoc-Technology/k8sprobe"
	"net/http"
	"time"
)

func main() {
	// Create a liveness probe
	livenessProbe := k8sprobe.NewProbe(true)

	// Automatically invalidate the probe after 10 seconds
	go func() {
		time.Sleep(10 * time.Second)
		livenessProbe.SetValid(false, "Liveness probe invalid after timeout")
	}()

	// Create a manager and register the liveness probe
	manager := k8sprobe.NewManager()
	manager.RegisterProbe(k8sprobe.LivenessProbe, livenessProbe)

	// Serve the liveness probe over HTTP
	http.Handle("/healthz/"+k8sprobe.LivenessProbe.String(), k8sprobe.NewHttpHandler(manager))
	http.ListenAndServe(":8089", nil)
}
```

Run the application and access the liveness endpoint at:  
`http://localhost:8089/healthz/LivenessProbe`

### Custom Probes

You can also define custom health checks by implementing the `ValidityChecker` interface:

```textmate
type CustomProbe struct {
	isServiceUp bool
}

func (p CustomProbe) IsValid() (bool, string) {
	if p.isServiceUp {
		return true, "OK"
	}
	return false, "Service is down"
}
```

Register the custom probe to the manager:

```textmate
customProbe := CustomProbe{isServiceUp: true}
manager.RegisterProbe(k8sprobe.ReadinessProbe, customProbe)
```

---

### Endpoints

The library exposes the following HTTP endpoints for Kubernetes health checks:

- **Liveness Probe**: `/healthz/LivenessProbe`
- **Readiness Probe**: `/healthz/ReadinessProbe`

You can add more endpoints as needed by defining and registering probes to the manager.

---

## How It Works

1. **Probes**: Probes, such as `LivenessProbe` or `ReadinessProbe`, represent the application health state.
2. **Manager**: The `k8sprobe.Manager` manages these probes and aggregates their health statuses.
3. **HTTP Handler**: The `k8sprobe.NewHttpHandler()` exposes the probes via an HTTP endpoint for Kubernetes to query.

The library automatically handles these requests and responds with the correct status code.

- `200 OK`: The application is healthy.
- `503 Service Unavailable`: The application is unhealthy.

---

## Example Applications

The `example/` directory contains sample applications demonstrating various probe usage:

1. `simple_http_server`: A basic server with a liveness probe that automatically becomes invalid after 10 seconds.
2. `custom_probe`: Shows how to use a custom probe to monitor external application health.

To run an example, navigate to its directory and execute:

```shell script
go run main.go
```

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Make your feature or fix in a new branch.
3. Open a pull request for review.

Feel free to propose enhancements, report bugs, or request features by creating GitHub issues.

---

## License

This project is licensed under the [MIT License](LICENSE).

---

## About

Developed and maintained by **Autodoc Technology**. This library simplifies Kubernetes lifecycle integration for Go-based applications.

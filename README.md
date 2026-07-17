# Cowpoke

<!-- FLEET-BADGES:BEGIN -->
[![CI](https://github.com/tzervas/cowpoke/actions/workflows/fleet-ci.yml/badge.svg?branch=main)](https://github.com/tzervas/cowpoke/actions/workflows/fleet-ci.yml?query=branch%3Amain)
[![Security](https://github.com/tzervas/cowpoke/actions/workflows/fleet-security.yml/badge.svg?branch=main)](https://github.com/tzervas/cowpoke/actions/workflows/fleet-security.yml?query=branch%3Amain)
<!-- FLEET-BADGES:END -->

**CPU-only Kubernetes Agent for Intelligent Resource Scaling**

Cowpoke is a lightweight Kubernetes agent that uses statistical feature extraction and cosine similarity to provide unified vertical and horizontal pod scaling recommendations, addressing the common problem of HPA/VPA thrashing in dynamic workloads.

## Overview

Traditional Kubernetes autoscaling (HPA/VPA) often leads to thrashing behavior - rapid oscillation between scaling up and down - because they react to instantaneous metrics without understanding workload patterns. Cowpoke solves this by:

1. **Extracting temporal features** from CPU metrics: mean, coefficient of variation (CV), percentiles, skewness, kurtosis, trend slope, and autocorrelation
2. **Computing cosine similarity** between current and historical workload patterns
3. **Providing unified scaling recommendations** that consider both vertical (resource limits) and horizontal (replica count) dimensions

## Goals

- **Fix HPA/VPA Thrashing**: Use pattern similarity to avoid reactive scaling decisions
- **Minimal Resource Footprint**: <5MB RAM, CPU-only, no external dependencies
- **Security First**: Proper RBAC, runs as sidecar with minimal permissions
- **stdlib Preference**: Pure Go implementation using standard library where possible

## Features

### Statistical Feature Extraction

Cowpoke computes the following features from CPU metric time series:
- **Mean**: Average resource utilization
- **Coefficient of Variation (CV)**: Normalized variability measure
- **Percentiles**: Distribution characteristics (p50, p95, p99)
- **Skewness**: Distribution asymmetry
- **Kurtosis**: Distribution tail heaviness
- **Trend Slope**: Linear trend direction
- **Autocorrelation**: Temporal pattern persistence

### Similarity-Based Scaling

- Computes cosine similarity between feature vectors
- Identifies similar historical workload patterns
- Recommends scaling actions based on successful past responses

## Architecture

```
┌─────────────────────────────────────┐
│         Kubernetes Pod              │
│                                     │
│  ┌──────────┐      ┌────────────┐  │
│  │  Main    │      │  Cowpoke   │  │
│  │Container │      │  (Sidecar) │  │
│  │          │      │            │  │
│  │          │      │ • Features │  │
│  │          │      │ • Similarity│ │
│  │          │      │ • Scaling  │  │
│  └──────────┘      └────────────┘  │
│                                     │
└─────────────────────────────────────┘
         │                    │
         ▼                    ▼
    [App Traffic]     [Metrics API]
                            │
                            ▼
                    [Scaling Decision]
```

### Components

- **cmd/cowpoke**: Main agent entry point
- **internal/features**: Statistical feature extraction from metrics
- **internal/similarity**: Vector operations and cosine similarity computation
- **deploy/**: Kubernetes deployment manifests (sidecar, RBAC, ConfigMap)
- **config/samples/**: Example CRD configurations

## Quick Start

### Build

```bash
make build
```

### Test

```bash
make test
```

### Deploy as Sidecar

```bash
kubectl apply -f deploy/rbac.yaml
kubectl apply -f deploy/configmap.yaml
kubectl apply -f deploy/sidecar.yaml
```

### Configuration

Configure Cowpoke via ConfigMap:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cowpoke-config
data:
  window-size: "300"      # Time window in seconds
  similarity-threshold: "0.85"  # Cosine similarity threshold
  metrics-interval: "15"  # Metric collection interval
```

## Development

### Prerequisites

- Go 1.24+
- Docker (for container builds)
- kubectl (for deployment)

### Building

```bash
# Build binary
make build

# Run tests
make test

# Run linter
make lint

# Build Docker image
make docker
```

## Security

- Runs with minimal RBAC permissions (read-only access to metrics)
- No external network dependencies
- Processes metrics locally
- No persistent storage of sensitive data

## License

MIT License - See [LICENSE](LICENSE) for details

## Contributing

Contributions welcome! Please ensure:
- Tests pass (`make test`)
- Code is formatted (`make lint`)
- Resource constraints maintained (<5MB RAM)
- No new external dependencies without discussion

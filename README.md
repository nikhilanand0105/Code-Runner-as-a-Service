# Code Runner as a Service

A lightweight, concurrent code execution engine and sandbox API written in Go.

![Go Version](https://img.shields.io/badge/go-1.22-blue.svg)
![Build](https://img.shields.io/badge/build-passing-brightgreen)
![License](https://img.shields.io/badge/license-MIT-green.svg)

---

## Overview

`code-runner` is an isolated code-execution microservice that compiles and evaluates untrusted user-submitted code across multiple languages with strict timeouts and memory isolation limits.

### Features
- **Multi-Language Runtime Support:** Evaluates Python 3, Node.js / JavaScript, C++, Bash, and Go.
- **Process Isolation & Resource Guards:** Executes within sandboxed temporary workspaces with configurable deadlines (`context.WithTimeout`).
- **RESTful API:** Clean JSON endpoints for code execution and health telemetry.
- **Container-Ready:** Multi-stage Docker container packaged with multi-language execution runtimes.

---

## Architecture

```
Client (HTTP POST) ──► REST API Server (Go HTTP Mux)
                             │
                             ▼
                    Sandbox Controller
                             │
            ┌────────────────┼────────────────┐
            ▼                ▼                ▼
      Python Sandbox    NodeJS Sandbox    C++ Sandbox
       (subprocess)      (subprocess)     (clang++ / bin)
```

---

## API Specification

### `POST /api/v1/execute`

**Request Payload:**
```json
{
  "language": "python",
  "code": "print('Hello from Code Runner!')\nfor i in range(3):\n    print(i)",
  "timeout_ms": 2000
}
```

**Response Payload:**
```json
{
  "stdout": "Hello from Code Runner!\n0\n1\n2\n",
  "stderr": "",
  "exit_code": 0,
  "duration_ms": 38,
  "timed_out": false
}
```

---

## Quickstart

### Running with Docker

```bash
docker build -t code-runner .
docker run -p 8080:8080 code-runner
```

### Running Locally

```bash
go run ./cmd/server
```

---

## License

This project is licensed under the [MIT License](LICENSE).

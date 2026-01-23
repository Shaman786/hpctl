
# Host-Palace Infrastructure CLI (`hpctl`)

> **Architectural POC:** A production-grade Infrastructure-as-Code (IaC) tool demonstrating a Client–Server Microservices pattern using Go, Node.js, and Docker.

---

## 🚀 Project Overview

**hpctl** simulates a distributed cloud platform control plane. It decouples the user interface (CLI) from the infrastructure execution layer (API), closely mirroring real-world systems such as the AWS CLI or `kubectl`.

This project demonstrates:

- **Systems Programming:** Building a robust CLI in Go
- **Microservices Architecture:** Decoupling control logic (Go) from execution logic (Node.js)
- **Security Engineering:** Input validation and injection protection
- **Container Orchestration:** Programmatic Docker lifecycle management

---

## 🏗 Architecture

The system follows a strict **Client–Server Model**:

```mermaid
graph LR
    A[Client: hpctl CLI] -- REST (HTTP/1.1) --> B[Server: Cloud API]
    B -- Docker Socket --> C[Infrastructure: Containers]
````

### 1. Control Plane (Client)

* **Language:** Go (Golang) 1.22
* **Frameworks:** Cobra (command structure), Viper (configuration)
* **Responsibility:** Authenticates user intent, validates arguments, and renders structured output.

### 2. Execution Plane (Server)

* **Language:** Node.js (Express)
* **Responsibility:** Interfaces directly with the Docker Engine.
* **Security:** Regex-based middleware to prevent shell injection attacks.
* **Resource Management:** Enforces multi-tenant quotas (CPU and RAM limits).

---

## ⚡ Key Features

| Feature                | Description                                       | Tech Stack            |
| ---------------------- | ------------------------------------------------- | --------------------- |
| **Provisioning**       | Instantly spin up isolated micro-VMs (containers) | Go → `POST /servers`  |
| **Telemetry**          | Stream real-time container logs                   | Docker Streams → REST |
| **Resource Isolation** | Enforce 512 MB RAM / 0.5 CPU per instance         | Docker cgroups        |
| **Security**           | Input sanitization against RCE attacks            | Express middleware    |

---

## 🛠 Installation & Usage

### Prerequisites

* Go 1.22+
* Node.js & npm
* Docker Desktop (running)

---

### 1. Setup the Cloud API (Backend)

```bash
cd hp-cloud-api
npm install
node index.js
```

```text
🚀 HOST-PALACE MOCK CLOUD API Online at port 5000
```

---

### 2. Build the CLI (Frontend)

```bash
git clone https://github.com/Shaman786/hpctl.git
cd hpctl
go build -o hpctl main.go
```

```bash
./hpctl --help
```

---

### 3. Operational Workflow

#### Step 1: Create a Web Server

```bash
./hpctl vm create web-production-01 --image nginx
```

```text
✔ Success: Instance provisioned successfully (ID: hp-x9s8d7)
```

---

#### Step 2: Inspect Infrastructure

```bash
./hpctl vm list
```

```text
NAME                IMAGE   STATUS
web-production-01   nginx   Up 12 seconds
```

---

#### Step 3: Fetch Logs (Telemetry)

```bash
./hpctl vm logs web-production-01
```

```text
172.17.0.1 - - [23/Jan/2026...] "GET / HTTP/1.1" 200
```

---

#### Step 4: Decommission Resources

```bash
./hpctl vm destroy web-production-01
```

---

## 🔒 Security & Performance

* **Injection Protection:** All API inputs are validated using strict regex patterns (`^[a-zA-Z0-9-_]+$`) prior to execution.
* **Error Handling:** The Go client decodes structured API errors to provide actionable feedback.
* **Concurrency:** Supports parallel provisioning without port conflicts using Docker’s internal networking.

---

*Author: Fayaj Hossain*
*Built to demonstrate Infrastructure Engineering competencies.*


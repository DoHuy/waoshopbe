# Dropship System - AI-Powered E-Commerce Platform

The **Dropship System** is a next-generation, full-stack e-commerce platform engineered for high performance, global scalability, and intelligent user engagement. It combines a robust Go-based microservices backend with a modern Next.js storefront, integrated with cutting-edge AI capabilities.

Whether you are scaling a boutique store or managing a high-traffic global dropshipping operation, this system provides the infrastructure to handle complex logistics, real-time customer support, and semantic product discovery.

---

## ✨ Core Features

### 🤖 Intelligent AI Ecosystem
- **AI-Powered Chatbot**: A real-time WebSocket-based assistant that uses OpenAI's GPT models to help customers find products, answer FAQs, and check order statuses.
- **Semantic Product Search**: Powered by `pgvector`, the system understands the *intent* behind searches, allowing users to find products using natural language descriptions.

### ⚡ High-Performance Backend
- **Microservices Architecture**: Built with **go-zero**, ensuring clean separation of concerns and the ability to scale individual components (Gateway, RPC, MQ) independently.
- **Event-Driven Processing**: Uses **Kafka** to handle high-volume background tasks such as order processing and multi-channel notifications.
- **Distributed Caching**: Leverages **Redis** for ultra-low latency product catalog delivery and intelligent rate limiting.

### 🛍 Modern Storefront
- **Next.js 15+**: Utilizing the App Router for superior SEO, fast page loads, and a seamless "Single Page Application" feel.
- **Global Payments**: Deep integration with **PayPal**, supporting secure global transactions with automated backend webhook verification.
- **Responsive Design**: A mobile-first, conversion-optimized shopping experience.

### 🛠 Enterprise-Grade Infrastructure
- **Database as Code**: Managed schema migrations via **Atlas** and GORM.
- **Cloud-Native Storage**: High-performance asset delivery using **Cloudflare R2**.
- **Observability**: Built-in support for structured logging and infrastructure monitoring (Kafka UI, Redis Insight).

---

## 🏗 System Architecture

The Dropship System is designed with a high-performance, scalable microservices-oriented architecture.

### 1. Backend Ecosystem (Go-Zero)
Built using the **go-zero** framework, the backend emphasizes high concurrency and clean separation of concerns:
- **API Gateway**: Acts as the single entry point. It handles HTTP/REST requests, manages WebSocket connections for the AI chatbot, and implements JWT authentication and rate limiting.
- **RPC Services**: Core business logic (Orders, Products, Inventory) is implemented as gRPC services, ensuring high-speed internal communication and strong data typing.
- **Message Queue (MQ)**: An asynchronous layer using **Kafka** (3-broker cluster) to handle background tasks like Telegram/Email notifications and event-driven data updates.


### 2. 🤖 Advanced AI & RAG Integration

The Dropship System features a sophisticated AI layer that goes beyond simple chat, implementing a full **Retrieval-Augmented Generation (RAG)** pipeline.

#### **Retrieval-Augmented Generation (RAG)**
The system doesn't just rely on the LLM's internal knowledge. It dynamically retrieves relevant data from your specific product catalog to ground the AI's responses:
1.  **Vector Embeddings**: Every product in the database is processed through OpenAI's `text-embedding-3-small` (or similar) model to generate a 1536-dimensional vector representing its semantic meaning.
2.  **pgvector Storage**: These embeddings are stored in PostgreSQL using the `pgvector` extension, enabling high-performance similarity searches.
3.  **Semantic Retrieval**: When a user asks a question, the system converts the query into a vector and performs a "Nearest Neighbor" search to find the most relevant products, even if the exact keywords don't match.

#### **Agentic Function Calling**
The chatbot operates as an **AI Agent** using OpenAI's **Function Calling (Tools)** capability:
- **Dynamic Tool Use**: The LLM decides in real-time whether it needs to search the catalog (`search_products`), check an order status, or answer general questions.
- **Structured Data Bridge**: Function calling acts as a type-safe bridge between the unstructured natural language of the user and the structured gRPC/SQL environment of the backend.
- **Context-Aware Assistance**: The bot maintains conversation history via WebSockets, allowing for follow-up questions like "Does the first one come in blue?" after a search result.

#### **Technical Stack**
- **LLM**: OpenAI GPT-4o / GPT-3.5-Turbo.
- **Vector DB**: PostgreSQL + `pgvector`.
- **Real-time Comm**: WebSockets (via Go-Zero `rest.Route`).
- **Embeddings**: OpenAI API for real-time query vectorization.

### 4. Storage & Infrastructure
- **Cloudflare R2**: S3-compatible object storage for high-performance delivery of product images and assets.
- **Docker Orchestration**: The entire infrastructure (Postgres, Redis, Kafka, Apps) is containerized for consistent deployment across environments.

### 5. Frontend (Next.js)
- A modern **Next.js** application using the App Router for SEO-optimized product pages.
- Integrated with the **PayPal SDK** for secure, global payments.

---

## Prerequisites

Ensure you have the following installed on your system:

### General
- **Docker & Docker Compose**: For infrastructure (PostgreSQL, Redis, Kafka).
- **Make**: To run build and development commands.

### Backend (`dropshipbe`)
- **Go**: 1.24.0 or higher.
- **Protobuf Compiler (`protoc`)**: For generating code from `.proto` files.
- **goctl**: go-zero development tool.
  ```bash
  go install github.com/zeromicro/go-zero/tools/goctl@latest
  ```
- **Atlas**: For database migrations.
  ```bash
  curl -sSf https://atlasgo.sh | sh
  ```

### Frontend (`storefront`)
- **Node.js**: 18.18.0 or higher.
- **Yarn**: 1.22.x.
  ```bash
  npm install --global yarn
  ```

---

## 1. Backend Setup (`dropshipbe`)

Navigate to the backend directory:
```bash
cd dropshipbe
```

### Environment Variables
Copy the sample environment file and update the values:
```bash
cp env.sample .env
```
*Note: Ensure the database credentials match those in `dropship-deployment/docker-compose.yml`.*

### Start Infrastructure
Use Docker Compose to start PostgreSQL, Redis, and Kafka:
```bash
docker-compose -f dropship-deployment/docker-compose.yml up -d
```
Verify the services are running:
```bash
docker ps
```

### Database Migrations
Apply the database schema using Atlas via the Makefile:
```bash
make apply
```

### Run Services
You need to run the following components (preferably in separate terminals):

1.  **RPC Backend**:
    ```bash
    make rpc
    ```
2.  **API Gateway**:
    ```bash
    make gw
    ```
3.  **Message Queue (Consumer)**:
    ```bash
    make mq
    ```

---

## 2. Frontend Setup (`storefront`)

Navigate to the storefront directory:
```bash
cd ../storefront
```

### Environment Variables
Create a `.env.local` file:
```bash
cp .env.example .env.local
```
Then update `.env.local` with the following:
```env
NEXT_PUBLIC_API_URL=http://localhost:8888
NEXT_PUBLIC_PAYPAL_CLIENT_ID=your_paypal_client_id
```
*Note: Replace `your_paypal_client_id` with your actual PayPal Sandbox Client ID.*

### Install Dependencies
```bash
yarn install
```

### Run Development Server
```bash
yarn dev
```
The storefront will be available at `http://localhost:3000`.

---

## 3. Useful Commands

### Backend (`Makefile`)
- `make help`: Show all available commands.
- `make status`: Check migration status.
- `make gen`: Regenerate protobuf and gRPC code.
- `make diff name=migration_name`: Create a new migration based on GORM models.
- `make down`: Revert the latest migration.

---

## 🚀 Production Deployment & Scaling

To handle high traffic and ensure global availability, the system is designed to be deployed using modern cloud-native practices.

### 1. Frontend: Cloudflare Pages
For the `storefront` (Next.js), Cloudflare Pages is recommended:
- **Global Edge Network**: Your frontend is cached and served from 300+ cities globally, ensuring sub-100ms latency for users.
- **Automatic Scaling**: Cloudflare handles millions of concurrent requests without manual intervention.
- **Security**: Built-in WAF, DDoS protection, and SSL management.
- **Deployment**:
  1. Connect your GitHub/GitLab repository to Cloudflare Pages.
  2. Set the build command: `yarn build`.
  3. Set the output directory: `.next`.
  4. Configure `NEXT_PUBLIC_API_URL` to point to your K8s Ingress.


### 2. Backend: Kubernetes (K8s) Implementation Plan

For the Go-Zero backend, we use Kubernetes to orchestrate high-availability and extreme-scale workloads.

#### **High-Traffic Architecture**
1.  **Ingress Controller (NGINX/Traefik)**:
    -   Terminates SSL/TLS.
    -   Handles traffic shaping and routing to the `gateway` service.
    -   Configured with `proxy-body-size` and `keep-alive` optimizations for heavy WebSocket/API traffic.

2.  **Service Scaling (HPA & VPA)**:
    -   **Horizontal Pod Autoscaler (HPA)**: Automatically scales `gateway` and `rpc` pods based on custom metrics (Requests Per Second) or CPU/Memory utilization.
    -   **Vertical Pod Autoscaler (VPA)**: Continuously analyzes resource usage to optimize CPU/Memory requests/limits for individual services.

3.  **Resilience & Self-Healing**:
    -   **Liveness & Readiness Probes**: Built-in health checks for each service ensuring that traffic only hits healthy pods.
    -   **Pod Disruption Budgets (PDB)**: Ensures a minimum number of replicas are always available during cluster upgrades or maintenance.
    -   **Topology Spread Constraints**: Spreads pods across different availability zones to ensure 99.99% uptime.

4.  **Database & Infrastructure Layer**:
    -   **Cloud SQL / RDS**: Managed PostgreSQL with read-replicas for high-volume product lookups.
    -   **ElastiCache / MemoryDB**: Redis cluster for distributed caching and rate-limiting at scale.
    -   **Managed Kafka (MSK/Confluent)**: Multi-AZ Kafka cluster for reliable event-driven order processing and notifications.

#### **Step-by-Step K8s Roadmap**

| Phase | Goal | Key Components |
| :--- | :--- | :--- |
| **Phase 1: Foundation** | Containerization & Basic Deployment | Dockerfiles, Helm Charts, K8s Namespaces, ConfigMaps/Secrets. |
| **Phase 2: Traffic Control** | High Availability & Routing | Ingress Controller, Service-type: ClusterIP, External-DNS, Cert-Manager (Let's Encrypt). |
| **Phase 3: Autoscaling** | Handling Traffic Spikes | Metrics Server, HPA (CPU/Mem), VPA (Recommendation mode). |
| **Phase 4: Observability** | Monitoring & Debugging | Prometheus, Grafana, Loki (Logging), Tempo (Distributed Tracing). |
| **Phase 5: Performance** | Latency Optimization | Redis Sidecar/Cluster, Database Read-Replicas, CDN Cache Invalidation. |

#### **Example HPA Configuration**
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: dropship-gateway-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: dropship-gateway
  minReplicas: 3
  maxReplicas: 50
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### 3. Observability & Performance Monitoring

-   **Tracing (OpenTelemetry)**: Full request tracing from Gateway -> RPC -> DB to identify latency bottlenecks in high-traffic scenarios.
-   **Structured Logging**: All microservices emit JSON-formatted logs for easy ingestion into ELK/Loki stacks.
-   **Real-time Metrics**: Dashboards for monitoring Kafka consumer lag, Redis hit ratios, and Database query performance.

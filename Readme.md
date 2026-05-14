# Distributed gRPC Workflow Engine

A high-performance, distributed workflow orchestration engine designed to execute complex, dependency-driven tasks across a cluster of worker nodes. 

Workflows are defined as Directed Acyclic Graphs (DAGs) using simple YAML/JSON files, allowing developers to define resilient pipelines without writing deployment-specific code. The engine separates the control plane (Coordinator) from the data plane (Workers) using highly efficient, bidirectional gRPC streams.

##  Architecture Overview

This system operates on a strictly decoupled architecture, ensuring fault tolerance, exact-once execution semantics (where possible), and independent scalability.

* **The Coordinator (Control Plane):** The central brain of the system. It is responsible for parsing YAML workflow definitions, validating the DAG logic (cycle detection), persisting state to the database, and dispatching ready tasks to available workers. It *does not* execute user workloads.
* **Worker Nodes (Data Plane):** A pool of distributed, stateless workers that connect to the Coordinator via gRPC. They listen for dispatched tasks, execute the isolated business logic (e.g., HTTP requests, database queries), and stream the execution status and logs back to the Coordinator.
* **State Store:** A centralized relational database (e.g., PostgreSQL) owned exclusively by the Coordinator. It stores workflow templates, the real-time state of active workflow instances, and historical logs.

## Project Structure

The repository follows standard Go project layout conventions to maintain a clean separation of concerns.

```text
├── cmd/
│   ├── coordinator/      # Entry point for the central scheduling server
│   │   └── main.go
│   └── worker/           # Entry point for the distributed worker nodes
│       └── main.go
├── internal/
│   ├── api/              # HTTP REST handlers for clients to submit/trigger workflows
│   ├── dag/              # Graph algorithms: topological sorting, cycle detection, validation
│   ├── db/               # Database interactions, schema management, and state updates
│   ├── engine/           # The core orchestration loop (task dispatcher, dependency resolution)
│   └── grpc/             # gRPC server/client implementations and stream handlers
├── proto/                
│   └── workflow.proto    # Protobuf definitions defining the Coordinator-Worker contract
├── examples/             # Sample YAML/JSON workflow definitions for testing
├── deployments/          # Dockerfiles and Kubernetes manifests (HPA, Deployments)
├── go.mod
└── README.md
```

## Development Roadmap & To-Do List

Building this engine requires tackling specific distributed systems challenges in phases. 

### Phase 1: Foundation & Contracts
- [ ] Define the `workflow.proto` file (Bidirectional streams for task dispatch and status updates).
- [ ] Generate the Go gRPC code using `protoc`.
- [ ] Define the core Go structs representing a Workflow, Task, and Instance.
- [ ] Implement the YAML/JSON parser to unmarshal static definitions into Go structs.

### Phase 2: Graph Logic & Validation
*Note: This phase heavily utilizes graph algorithms.*
- [ ] Implement an adjacency list representation of the parsed tasks.
- [ ] Write a Topological Sort algorithm to map the execution order.
- [ ] Implement Cycle Detection (e.g., via Kahn's algorithm or DFS) to reject invalid workflows.

### Phase 3: The gRPC Communication Layer
- [ ] Build the gRPC Server in the Coordinator to accept worker registrations.
- [ ] Build the gRPC Client in the Worker to connect to the Coordinator and maintain a heartbeat.
- [ ] Implement bidirectional streaming: Coordinator pushes `TaskAssignment`, Worker returns `TaskStatus`.

### Phase 4: State Management & The Orchestration Loop
- [ ] Design the PostgreSQL schema (Tables: `workflows`, `instances`, `tasks`, `logs`).
- [ ] Write the Coordinator's "State Updater": accurately transition task states (Pending -> Running -> Success/Failed).
- [ ] Implement the "Dependency Resolver": when a task succeeds, check the DAG to find newly unblocked tasks and push them to the Dispatch Queue.
- [ ] Implement retry logic for transient worker failures based on YAML definitions.

### Phase 5: Task Execution 
- [ ] Implement a `TaskRunner` interface on the Worker node.
- [ ] Build specific task executors (e.g., `PostgresExecutor`, `HttpExecutor`, `SleepExecutor` for testing).
- [ ] Capture standard output/errors from task execution and stream logs back to the Coordinator.

### Phase 6: Infrastructure & Auto-scaling
- [ ] Containerize both the Coordinator and the Worker using Docker.
- [ ] Expose Prometheus metrics from the Coordinator (e.g., `pending_tasks_in_queue`).
- [ ] Create Kubernetes manifests to deploy the system, utilizing a Horizontal Pod Autoscaler (HPA) to scale the Worker deployment based on the exposed queue metrics.

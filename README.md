# 🚀 Uptime Monitor

A lightweight, robust URL uptime monitoring system built with **Go** and **PostgreSQL**. This service tracks user-submitted URLs, performs periodic health checks, and persists monitoring results for analysis.

## 🏗️ Architecture Overview

The system is currently designed as a **synchronous, interval-based monitoring engine**. It runs a background scheduling loop that ensures every URL is checked according to its defined interval.

### Core Components
- **API Server:** Handles user registration, login, and URL management via `net/http` and `chi`/`mux`.
- **Database:** Uses `PostgreSQL` for persistence and `sqlc` for type-safe database interactions.
- **Checker Service:** A background scheduler that fetches due tasks and synchronously validates URL health (status codes, latency).

## 💡 Current Implementation Details

- **Linear Processing:** Currently utilizes a single-goroutine scheduler for simplicity and reliability.
- **Resilient Scheduling:** Implements an `UpdateURLNextCheck` mechanism to ensure no duplicate checks occur and tasks are reliably rescheduled after each execution.
- **Robust Error Handling:** Handles network issues (e.g., DNS failures, connection timeouts) and database hiccups gracefully, ensuring the scheduler survives temporary infrastructure outages.

## 🔮 Future Scope & Roadmap

As the project scales, the following optimizations are planned:

### 1. High-Performance Concurrency (Worker Pool)
Transition from a linear scheduler to a **Concurrent Worker Pool**. This will decouple task dispatching from task execution, allowing thousands of URLs to be monitored in parallel without blocking the main scheduler.

### 2. Database Scaling (Batching)
To handle 1,000+ simultaneous checks, we will implement:
- **Batch Inserts:** Replace individual `INSERT` queries for check results with bulk inserts to reduce database round-trips.
- **Producer-Consumer-Writer Pattern:** Separate HTTP checkers (async) from database writers (batch) to minimize database connection pressure.

### 3. Monitoring Enhancements
- **Alerting System:** Send notifications (e.g., Slack, Email, Webhooks) when a URL transition state changes from `UP` to `DOWN`.
- **Advanced Checks:** Add support for regex response body matching, custom HTTP headers, and SSL certificate expiration monitoring.

---

## 🛠️ Tech Stack
- **Language:** Go
- **Database:** PostgreSQL
- **Database Tools:** `sqlc`, `goose` (migrations)
- **Networking:** Standard `net/http`

# 🚀 go-jobqueue

A background job queue system built in **Go** using **Clean Architecture**.

It allows developers to enqueue jobs and process them asynchronously with worker pools.  
Supports multiple notification mechanisms — **polling**, **webhook**, and **websocket**.

---

## 🧩 Architecture Overview

- **Domain Layer** — Defines core entities like `Job` and interfaces (`JobRepository`)
- **Use Case Layer** — Handles enqueueing, processing, and updating job states
- **Interface Layer** — Exposes HTTP endpoints and worker logic
- **Infrastructure Layer** — Implements Redis queue and notifiers
- **Frameworks / External** — Fiber, Redis, etc.

---

## ⚙️ Features

✅ Enqueue jobs via REST API  
✅ Background worker pool for processing jobs  
✅ Job status tracking (pending, processing, success, failed)  
✅ Redis-based queue management  
✅ Notification options:

- Polling via `GET /jobs/:id`
- Webhook callback
- WebSocket push events

---

## 🧰 Tech Stack

- **Language:** Go 1.22+
- **Framework:** Fiber
- **Queue:** Redis
- **Architecture:** Clean Architecture
- **Libraries:**
  - `github.com/gofiber/fiber/v2`
  - `github.com/redis/go-redis/v9`
  - `stretchr/testify` (for testing)

---

## 🚀 Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/supanut/go-jobqueue.git
cd go-jobqueue
```

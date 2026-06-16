# chirpy backend

Chirp is a minimal social network backend platform written in Go. This service provides a RESTful API for user management, authentication sessions, and content creation ("chirps"), powered by a PostgreSQL database layer.

## Motivation
This project was guided through boot.dev with the purpouse of building the backend logic of a webapp similar to X/Twitter. 

### Goal

The goal with `chirpy` is to be able to understand how http servers work under the hood. I think understanding how to actually build RESTful API's, build Web Apps without using a framework and write my own documentation really helped me see the bigger picture of how these processes work.
---

## 🛠️ Tech Stack

* **Language:** Go (Golang)
* **Database:** PostgreSQL
* **Migration Tool:** Goose
* **SQL Generator:** SQLC

---

## 🚀 Getting Started

### Prerequisites
Ensure you have the following installed on your machine:
* Go (1.20+ recommended)
* PostgreSQL
* Goose CLI binary

### 1. Environment Setup
Create a `.env` file in the root directory of the project and add your database connection string:

```env
DB_URL="postgres://YOUR_USERNAME:@localhost:5432/chirpy?sslmode=disable"
# Chirpy Backend

Chirp is a minimal social network backend platform written in Go. This service provides a RESTful API for user management, authentication sessions, and content creation ("chirps"), powered by a PostgreSQL database layer.

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
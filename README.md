# 🚀 Social Media Backend API (Golang + Gin)

![Project Preview](https://res.cloudinary.com/dgagbheuj/image/upload/v1776327544/dpoz8gelxr66vwfqh9ek.jpg)

A fully-featured **Social Media Backend API** built with **Golang (Gin)**, designed with clean architecture principles and production-level patterns.
This project demonstrates real-world backend engineering skills including authentication, real-time messaging, notifications, and scalable service structure.

---

## 📌 Features

### 🔐 Authentication & Users

* User registration & login (JWT-based)
* Secure password hashing
* Get current user profile
* Update profile
* Admin actions (delete users)

---

### 📝 Posts

* Create, update, delete posts
* Upload images (Cloudinary integration)
* Get all posts (with pagination)
* Get posts by user
* Notifications for followers when posting

---

### 🎬 Reels

* Create short media content
* Update & delete reels
* Fetch reels (global / user-specific)

---

### 💬 Comments

* Add comments to:

  * Posts
  * Reels
  * Comments (replies)
* Nested comments (replies system)
* Get comments with pagination
* Get replies for a specific comment
* Notifications when someone comments

---

### ❤️ Likes

* Toggle like/unlike
* Count likes
* Check if user liked content
* List users who liked

---

### 👥 Follow System

* Follow / Unfollow users
* Get followers / following lists
* Count followers & following
* Notifications on follow

---

### 💌 Messages (Real-Time)

* Send messages
* Get chat history
* Delete messages
* Mark messages as read
* Real-time chat using **WebSocket**

---

### 🔔 Notifications (Real-Time)

* Get all notifications
* Unread count
* Mark as read (single / all)
* Delete notifications
* Redis integration for fast unread tracking

---

## 🧠 Tech Stack

* **Language:** Go (Golang)
* **Framework:** Gin
* **Database:** MongoDB (Atlas)
* **Caching:** Redis
* **Realtime:** WebSocket
* **Cloud Storage:** Cloudinary
* **Authentication:** JWT

---

## 📂 Project Structure

```
backend/
│
├── app/              # Dependency Injection & Container setup
├── handlers/         # HTTP handlers (controllers)
├── services/         # Business logic
├── models/           # Data models
├── routes/           # API routes
├── middlewares/      # Auth & role middlewares
├── utils/            # Helpers (JWT, Cloudinary, Redis)
├── websocket/        # Real-time chat logic
└── main.go           # Entry point
```

---

## ⚙️ Environment Variables

Create a `.env` file in the root:

```env
PORT=8080

DB_NAME=go_learning
MONGO_URI=your_mongodb_uri

JWT_SECRET=your_secret
JWT_EXPIRES_IN=20d

CLOUDINARY_CLOUD_NAME=your_name
CLOUDINARY_API_KEY=your_key
CLOUDINARY_API_SECRET=your_secret

REDIS_URL=your_redis_url
```

---

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/ZenZN99/Virelo-Social-Media-API
cd Virelo-Social-Media-API
cd backend
```

---

### 2. Install dependencies

```bash
go mod tidy
```

---

### 3. Run the server

```bash
go run main.go
```

Server will run on:

```
http://localhost:8080
```

---

## 🔗 API Endpoints Overview

### Auth

```
POST   /api/auth/signup
POST   /api/auth/login
GET    /api/auth/me
PUT    /api/auth/update/profile
```

---

### Posts

```
POST   /api/post/create
GET    /api/post/posts
GET    /api/post/post/:postId
PUT    /api/post/update/:postId
DELETE /api/post/delete/:postId
```

---

### Comments

```
POST   /api/comment/create
GET    /api/comment/comments
GET    /api/comment/replies/:commentId
```

---

### Likes

```
POST   /api/likes/toggle
GET    /api/likes/count
GET    /api/likes/is-liked
```

---

### Follow

```
POST   /api/follow/
DELETE /api/follow/:userId
GET    /api/follow/followers/:userId
GET    /api/follow/following/:userId
```

---

### Messages

```
POST   /api/message/send
GET    /api/message/:receiverId
```

---

### Notifications

```
GET    /api/notification/
GET    /api/notification/unread-count
PUT    /api/notification/:id/read
PATCH  /api/notification/read-all
```

---

## 🧪 Testing

Use **Postman** or any API client:

* Set `Authorization` header:

```
Bearer <your_token>
```

---

## ⚡ Key Highlights

* Clean Architecture (Handlers → Services → Models)
* Dependency Injection (manual, Go-style)
* Real-time messaging (WebSocket)
* Scalable structure
* Production-ready patterns

---

## 🧑‍💻 Author

**Zen Al-Laham**
Full-Stack Engineer

---

## 📜 License

This project is for educational and portfolio purposes.

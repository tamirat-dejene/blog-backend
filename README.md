# G6 Blog Starter Project

A2SV Intern Team G6 – Blog Backend API  
Built with **Go** using **Clean Architecture**.

## 🚀 Features
- CRUD for blog posts
- User authentication and role-based access
- Blog filtering and search
- AI-powered content suggestions (planned)

## 🗂 Tech Stack
- Go (Gin, MongoDB driver)
- MongoDB
- Clean Architecture principles

## 📁 Project Structure
- `Delivery/` – Entry point, controllers, routers
- `Domain/` – Core domain models
- `Usecases/` – Business logic interfaces
- `Repositories/` – Data persistence layer
- `Infrastructure/` – DB, middleware, etc.
- `Configs/` – Environment & app config
- `Docs/` – Project documentation

## 👥 Team Contribution Guide

1. **Create a branch** from `main`:
```

git checkout -b feature/<your-name>-<task>

```

2. **Make your changes** and commit with a clear message:
```

git commit -m "Add: user login handler"

```

3. **Push to origin**:
```

git push origin feature/<your-name>-<task>

```

4. **Open a Pull Request** and request review from a teammate.
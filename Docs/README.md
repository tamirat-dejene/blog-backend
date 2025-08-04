# Blog Backend Authentication & User Management

This backend provides secure authentication, user management, and profile features for a blog platform. It uses JWT for access and refresh tokens, supports password reset, and integrates with ImageKit for avatar uploads.

---

## **Architecture Overview**

- **Controllers**: Handle HTTP requests (Gin framework).
- **Usecases**: Business logic for users, tokens, and password reset.
- **Repositories**: Database operations (MongoDB).
- **Domain**: Core models and interfaces.
- **Utils**: Cookie helpers.
- **Security**: JWT token generation/validation.
- **Storage**: ImageKit integration for file uploads.
- **Middleware**: Auth and role-based access control.

---

## **Authentication Flow**

### 1. **User Registration**

- **Endpoint**: `POST /api/auth/register`
- **Flow**:
  - Controller validates and maps request to domain user.
  - Usecase checks for existing username/email, hashes password, and saves user.
  - Repository inserts user into MongoDB.

### 2. **Login**

- **Endpoint**: `POST /api/auth/login`
- **Flow**:
  - Controller validates credentials.
  - Usecase fetches user and checks password.
  - JWT tokens (access & refresh) are generated.
  - Refresh token is saved/replaced in DB (old one revoked).
  - Tokens are set as HTTP-only cookies.

### 3. **Token Refresh**

- **Endpoint**: `POST /api/auth/refresh`
- **Flow**:
  - Controller validates refresh token from request/cookie.
  - Usecase checks token validity and expiry.
  - If near expiry, rotates (revokes old, saves new) refresh token.
  - New access token (and possibly refresh token) is issued and set in cookies.

### 4. **Logout**

- **Endpoint**: `POST /api/auth/logout`
- **Flow**:
  - Controller gets refresh token from cookie.
  - Usecase revokes and deletes the token in DB.
  - Cookies are cleared.

### 5. **Change Role**

- **Endpoint**: `PATCH /api/auth/change-role`
- **Flow**:
  - Controller checks initiator's role.
  - Usecase enforces role-change rules (only admin/superadmin can change roles).
  - Repository updates user role in DB.

### 6. **Forgot/Reset Password**

- **Endpoints**:
  - `POST /api/auth/forgot-password`
  - `POST /api/auth/reset-password`
- **Flow**:
  - Controller validates request.
  - Usecase generates and emails a reset token (forgot).
  - Usecase verifies token and updates password (reset).
  - Repository manages reset tokens in DB.

### 7. **Profile Update (with Avatar Upload)**

- **Endpoint**: `PATCH /api/user/profile`
- **Flow**:
  - Controller parses multipart form, reads avatar file.
  - Usecase uploads avatar to ImageKit, updates user fields.
  - Repository updates user in DB.

---

## **Key Files and Their Roles**

### **Controllers**

- `auth_controller.go`: Handles registration, login, logout, token refresh, role change, password reset.
- `user_controller.go`: Handles user profile update, including avatar upload.

### **Usecases**

- `user_usecase.go`: User registration, login, logout, profile update, role change.
- `refresh_token_usecase.go`: Refresh token CRUD and revocation logic.

### **Repositories**

- `user_repo.go`: MongoDB operations for users.
- `refresh_token_repo.go`: MongoDB operations for refresh tokens.
- `reset_password_repo.go`: MongoDB operations for password reset tokens.

### **Domain**

- `refresh_token.go`: Models and interfaces for tokens.
- `storage.go`: Storage service interface.

### **Utils**

- `cookies_services.go`: Helpers for setting, getting, deleting cookies.

### **Security**

- `jwt_service.go`: JWT token generation and validation.

### **Storage**

- `imagekit.go`: Uploads files to ImageKit, validates image data.

### **Middleware**

- `auth.go`: JWT authentication and role-based access control for routes.

---

## **Entity Relationships**

- **User** (1) --- (N) **RefreshToken**
- **User** (1) --- (N) **PasswordResetToken**

---

## **Token Handling**

- **Access Token**: Short-lived, for API authentication, stored in HTTP-only cookie.
- **Refresh Token**: Long-lived, for session renewal, stored in HTTP-only cookie and DB, rotated/revoked as needed.

---

## **Security Practices**

- Passwords are hashed before storage.
- All tokens are stored in HTTP-only cookies.
- Refresh tokens are revoked on logout and rotated on refresh.
- Role-based access enforced in middleware and usecases.

---

## **How to Run**

1. **Set up MongoDB and ImageKit credentials in your `.env` file.**
2. **Start the server:**
   ```
   go run Delivery/main.go
   ```
3. **Use Postman or similar tool to test endpoints.**

---

## **Example API Usage**

### Register

```json
POST /api/auth/register
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "password123",
  "first_name": "Alice",
  "last_name": "Smith"
}
```

### Login

```json
POST /api/auth/login
{
  "identifier": "alice",
  "password": "password123"
}
```

### Refresh Token

```json
POST /api/auth/refresh
{
  "refresh_token": "<token>"
}
```

---

## **Contributing**

- Follow Go best practices.
- Write unit tests for new features.
- Use clear commit messages.

---

## **License**

MIT

---

\*\*This backend is designed for secure, scalable user authentication and management in

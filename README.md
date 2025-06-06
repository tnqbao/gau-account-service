# Gau Account Service

**Tiếng Việt:**  
Gau Account Service là một microservice được phát triển bằng ngôn ngữ Go, chịu trách nhiệm chính trong việc quản lý tài khoản người dùng trong hệ thống. Dịch vụ cung cấp các API RESTful phục vụ cho việc đăng ký, đăng nhập, xác thực và cập nhật thông tin tài khoản. Kiến trúc của ứng dụng được thiết kế theo hướng module hóa và dễ mở rộng, hỗ trợ triển khai linh hoạt bằng Docker hoặc Kubernetes.

**English:**  
Gau Account Service is a microservice developed in Go, responsible for managing user accounts in the system. The service provides RESTful APIs for registration, login, authentication, and updating account information. The application architecture is modular and extensible, supporting flexible deployment with Docker or Kubernetes.

## Mục lục | Table of Contents

- [Tính năng | Features](#tính-năng--features)
- [Kiến trúc thư mục | Directory Structure](#kiến-trúc-thư-mục--directory-structure)
- [Yêu cầu hệ thống | System Requirements](#yêu-cầu-hệ-thống--system-requirements)
- [Hướng dẫn cài đặt | Installation](#hướng-dẫn-cài-đặt--installation)
- [Chạy bằng Docker | Run with Docker](#chạy-bằng-docker--run-with-docker)
- [Triển khai Kubernetes | Kubernetes Deployment](#triển-khai-kubernetes--kubernetes-deployment)
- [Thông tin liên hệ | Contact](#thông-tin-liên-hệ--contact)

## Tính năng | Features

- Đăng ký và đăng nhập người dùng  
  User registration and login
- Quản lý và xác thực tài khoản  
  Account management and authentication
- Hỗ trợ xác thực bằng JWT  
  JWT authentication support
- Quản lý thông tin hồ sơ người dùng  
  User profile management
- Middleware hỗ trợ logging, phân quyền  
  Middleware for logging and authorization
- Hệ thống repository tách biệt cho việc truy cập dữ liệu  
  Repository pattern for data access
- CI/CD tích hợp GitHub Actions  
  CI/CD with GitHub Actions

## Kiến trúc thư mục | Directory Structure

| Đường dẫn / Path         | Mô tả (VN)                                                      | Description (EN)                                   |
|-------------------------|------------------------------------------------------------------|----------------------------------------------------|
| `config/`               | Định nghĩa các cấu hình hệ thống và khởi tạo                     | System configuration and initialization            |
| `controller/`           | Định nghĩa logic xử lý request và gọi dịch vụ tương ứng          | Request handling logic and service calls           |
| `deploy/k8s-test/`      | Cấu hình manifest cho Kubernetes test                            | Kubernetes test manifests                          |
| `middlewares/`          | Các lớp middleware: xác thực, logging, phân quyền, v.v.          | Middleware: authentication, logging, authorization |
| `migrations/`           | Các script migration cơ sở dữ liệu                               | Database migration scripts                         |
| `models/`               | Định nghĩa các entity và mô hình dữ liệu                         | Entity and data model definitions                  |
| `providers/`            | Khởi tạo kết nối đến các tài nguyên bên ngoài (DB, Redis, etc.)  | External resource providers (DB, Redis, etc.)      |
| `repositories/`         | Truy cập và thao tác với dữ liệu từ CSDL                         | Data access and manipulation                       |
| `routes/`               | Định nghĩa các tuyến API và ánh xạ controller                    | API routes and controller mapping                  |
| `.github/workflows/`    | Cấu hình pipeline CI/CD sử dụng GitHub Actions                   | CI/CD pipeline configuration with GitHub Actions   |
| `main.go`               | Điểm vào chính của ứng dụng                                      | Application entry point                            |
| `Dockerfile`            | Dockerfile để build image                                        | Dockerfile for building image                      |
| `entrypoint.sh`         | Script khởi động ứng dụng trong container                        | Application startup script in container            |

## Yêu cầu hệ thống | System Requirements

- Go 1.20 hoặc mới hơn  
  Go 1.20 or newer
- PostgreSQL hoặc MySQL  
  PostgreSQL or MySQL
- Redis (tuỳ chọn)  
  Redis (optional)
- Docker (nếu dùng container)  
  Docker (if using containers)
- Kubernetes (nếu triển khai trên cụm)  
  Kubernetes (for cluster deployment)

## Hướng dẫn cài đặt | Installation

1. Clone repository:

   ```bash
   git clone https://github.com/tnqbao/gau-account-service.git
   cd gau-account-service
   ```

2. Cài đặt các phụ thuộc | Install dependencies:

   ```bash
   go mod download
   ```

3. Cấu hình biến môi trường | Configure environment variables:

   - Sao chép file `.env.example` thành `.env` và chỉnh sửa thông tin kết nối database, Redis, JWT, ...
   - Copy `.env.example` to `.env` and edit database, Redis, JWT, ... settings.

4. Chạy migration (nếu cần) | Run migrations (if needed):

   ```bash
   # Ví dụ với migrate tool
   migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up
   ```

5. Khởi động ứng dụng | Start the application:

   ```bash
   go run main.go
   ```

## Chạy bằng Docker | Run with Docker

1. Build Docker image:

   ```bash
   docker build -t gau-account-service .
   ```

2. Run container:

   ```bash
   docker run --env-file .env -p 8080:8080 gau-account-service
   ```

## Triển khai Kubernetes | Kubernetes Deployment

- Các file manifest mẫu nằm trong thư mục `deploy/k8s-test/`.
- Sample manifests are in the `deploy/k8s-test/` directory.

```bash
kubectl apply -f deploy/k8s-test/
```

## Thông tin liên hệ | Contact

- Tác giả / Author: Trần Nguyễn Quốc Bảo
- Email: quocbao.job106204@gmail.com
- GitHub: [https://github.com/tnqbao/gau-account-service](https://github.com/tnqbao/gau-account-service)

---
```

````
# Hệ Thống Backend Toán Cô Trà (TCT) - Kiến Trúc Sạch (Clean Architecture)

Hệ thống quản lý lớp học và tư vấn tuyển sinh Toán tư duy liên cấp **Toán Cô Trà** được tái cấu trúc hoàn chỉnh theo mô hình **Clean Architecture (Kiến Trúc Sạch)** chuẩn hóa bằng ngôn ngữ Go (Golang). Dự án tích hợp trang Landing Page giới thiệu kết hợp ứng dụng Quản trị Sales thời gian thực (SPA) và cổng tương tác với chatbot.

---

## 🌟 Tính Năng Nổi Bật

1. **Kiến Trúc Sạch Chuẩn Hóa (Clean Architecture)**:
   - Cách ly hoàn toàn mã nguồn nghiệp vụ lõi (`internal/domain`) khỏi các chi tiết kỹ thuật như cơ sở dữ liệu (`pkg/database`) hay phương thức phân phối HTTP (`internal/lead/delivery`).
   - Đảm bảo tính mở rộng cao, dễ dàng bảo trì và viết kiểm thử tự động (Unit Test).

2. **Cơ Chế Tự Động Sao Lưu An Toàn (Automatic Local Fallback)**:
   - Tự động phát hiện trạng thái kết nối cơ sở dữ liệu Postgres.
   - Nếu không có cơ sở dữ liệu Postgres hoặc cấu hình lỗi, hệ thống tự kích hoạt **Mock Repository luồng an toàn (Thread-safe In-Memory Mock Database)** được nạp sẵn dữ liệu mẫu sinh động. Giúp hệ thống chạy thử nghiệm ngay lập tức mà không cần cài đặt Postgres phức tạp.

3. **Thời Gian Thực Với Server-Sent Events (SSE)**:
   - Khi có phụ huynh đăng ký tư vấn mới tại Landing Page, thông tin đăng ký sẽ lập tức được đẩy tự động tới giao diện của đội ngũ Sales mà không cần F5/tải lại trang (Real-time stream).

4. **Trình Quản Trị Sales Tiện Ích & Đẹp Mắt (iOS-Style)**:
   - Giao diện quản lý đăng ký của đội ngũ Sales được thiết kế theo ngôn ngữ tinh gọn kiểu iOS, có nút chuyển trạng thái tư vấn dạng công tắc trượt (Toggle switch). Khi bật "Đã tư vấn", thẻ của học sinh lập tức đổi màu sáng ngọc lục bảo đầy tinh tế.

5. **Giới Hạn Tần Suất Truy Cập (Token Bucket Rate Limiter)**:
   - Bảo vệ hệ thống khỏi tấn công DDoS và spam form đăng ký bằng cách giới hạn số lượng request từ mỗi địa chỉ IP của Client một cách thông minh.

6. **Xuất & Nhập Danh Sách Excel (.xlsx) Chuyên Nghiệp**:
   - Tích hợp thư viện Excelize v2 xử lý bảng tính Excel tốc độ cao.
   - Hỗ trợ tải file mẫu chuẩn, xuất danh sách lớp học và học sinh.
   - Tự động kiểm tra tính hợp lệ dữ liệu (SĐT Việt Nam, khối lớp, loại lớp) và thực hiện Batch Upsert vào hệ thống.

---

## 📂 Sơ Đồ Cấu Trúc Thư Mục

```text
toan-co-tra-backend/
├── assets/                     # Tài nguyên ảnh tĩnh chính hãng của thương hiệu
│   ├── banner.png              # Banner quảng cáo chất lượng cao của hệ thống
│   └── cotra.jpg               # Chân dung nghệ thuật của Cô Trà giáo viên sáng lập
│
├── cmd/
│   └── api/
│       └── main.go             # Điểm khởi chạy hệ thống (Khởi tạo DB, Middleware, Router & Server)
│
├── config/
│   └── config.go               # Đọc tham số môi trường cấu hình (Port, Postgres, JWT, Chatbot)
│
├── internal/                   # Code private của hệ thống, không cho phép import ngoài phạm vi dự án
│   ├── domain/                 # Thực thể nghiệp vụ cốt lõi (Domain Models)
│   │   ├── class.go            # Struct phân loại lớp học (Cơ bản, Nâng cao, Chất lượng cao)
│   │   ├── lead.go             # Struct đăng ký của Phụ huynh (Họ tên, SĐT, Lớp, Hình thức học...)
│   │   └── user.go             # Struct tài khoản & phiên đăng nhập hệ thống (Teacher, Parent, Sales)
│   │
│   ├── lead/                   # Module quản lý Đăng ký tuyển sinh & Chatbot
│   │   ├── delivery/
│   │   │   ├── chatbot/
│   │   │   │   └── webhook.go  # Webhook chatbot (Zalo, Messenger...)
│   │   │   └── http/
│   │   │       └── handler.go  # API endpoints, Server-Sent Events stream, và Quản trị Sales
│   │   ├── repository/
│   │   │   └── postgres/
│   │   │       └── pg_lead.go  # Tương tác với PostgreSQL + Thread-safe Mock Fallback
│   │   └── usecase/
│   │       └── lead_uc.go      # Điều phối nghiệp vụ đăng ký tuyển sinh
│   │
│   └── middleware/             # Các bộ lọc trung gian bảo vệ hệ thống
│       ├── auth.go             # Bearer Token Auth kiểm tra phân quyền tài khoản
│       └── ratelimit.go        # IP-Based Token Bucket Rate Limiter chống spam
│
├── pkg/                        # Các thư viện dùng chung cho toàn bộ dự án
│   ├── database/
│   │   └── postgres.go         # Trình quản lý kết nối và Pooling PostgreSQL
│   ├── event/
│   │   └── dispatcher.go       # Bộ định tuyến sự kiện (Pub/Sub Event Broker) thời gian thực
│   └── logger/
│       └── zap.go              # Trình ghi log hệ thống chuẩn hóa
│
├── go.mod                      # Quản lý Golang Module
└── README.md                   # Tài liệu hướng dẫn sử dụng chi tiết này
```

---

## ⚙️ Cấu Hình Môi Trường (Environment Variables)

Hệ thống có thể được cấu hình linh hoạt thông qua các biến môi trường dưới đây (mặc định sẽ sử dụng các giá trị an toàn được thiết lập sẵn nếu không khai báo):

| Biến Môi Trường | Mô Tả | Giá Trị Mặc Định |
| :--- | :--- | :--- |
| `PORT` | Cổng HTTP của Web Server | `8080` |
| `DATABASE_URL` | Chuỗi kết nối tới cơ sở dữ liệu PostgreSQL | Mặc định trống (Kích hoạt chế độ In-memory Fallback) |
| `JWT_SECRET` | Khóa bí mật dùng cho phân quyền Token | `toan_co_tra_secret_key_2026` |
| `CHATBOT_SECRET` | Khóa bảo mật xác thực webhook từ chatbot | `super_secure_chatbot_webhook_token` |

---

## 🚀 Hướng Dẫn Chạy Thử Hệ Thống

### 1. Yêu Cầu Cài Đặt
- **Go**: Phiên bản 1.21 trở lên.

### 2. Tải Thư Viện Phụ Thuộc (Dependencies)
Truy cập vào thư mục chứa dự án `toan-co-tra-backend` và chạy lệnh:
```bash
go mod tidy
```

### 3. Biên Dịch (Compilation)
Để kiểm tra tính nhất quán và biên dịch dự án, chạy lệnh:
```bash
go build -o toan-co-tra-backend.exe ./cmd/api/main.go
```

### 4. Khởi Chạy Server
Khởi chạy hệ thống để trải nghiệm ngay lập tức trên máy tính cá nhân:
```bash
go run ./cmd/api/main.go
```

Khi khởi động thành công, màn hình Console sẽ hiển thị thông tin dạng:
```text
[INFO] Khởi tạo cấu hình môi trường hoàn tất.
[INFO] DATABASE_URL trống. Tự động chuyển sang chế độ Mock Database (Local Fallback) an toàn.
[INFO] Khởi tạo Event Dispatcher hoàn tất.
[INFO] Đăng ký API Route /api/v1/registrations thành công.
[INFO] Web Server Toán Cô Trà đang chạy tại địa chỉ http://localhost:8080
```

---

## 💻 Trải Nghiệm Giao Diện Premium (Single Page Application - SPA)

Mở trình duyệt web của bạn và truy cập:
👉 **[http://localhost:8080](http://localhost:8080)**

Tại đây bạn sẽ được trải nghiệm một giao diện SPA tích hợp vô cùng tinh tế, bao gồm:
1. **Phần Billboard**: Trình chiếu banner chính thống và giới thiệu về các lớp học chuyên sâu của Cô Trà.
2. **Phần Đăng Ký Học Thử**: Phụ huynh điền thông tin đăng ký (Họ tên, SĐT, Khối lớp, Hình thức học). Nút bấm "Đăng Ký Tư Vấn" được sửa lỗi phông chữ tiếng Việt chuẩn xác và có bo góc mượt mà.
3. **Phần Quản Trị Sales (Thời Gian Thực)**:
   - Bấm **"Đăng Nhập Đội Ngũ Sales"** bằng tài khoản mẫu:
     - **Tên đăng nhập**: `teacher1`
     - **Mật khẩu**: `password123`
   - Tại đây, bạn sẽ thấy danh sách các học sinh đăng ký.
   - Thử nghiệm đăng ký mới ở tab Phụ huynh: Giao diện Sales sẽ lập tức cập nhật Lead mới (có hiệu ứng nhấp nháy nổi bật và thông báo thành công thông qua Server-Sent Events).
   - Bạn có thể bấm vào công tắc **"Đã tư vấn" / "Chưa tư vấn"** để thay đổi trạng thái chăm sóc khách hàng.
4. **Phần Footer của Hệ Thống**:
   - Trưng bày đầy đủ thông tin thương hiệu Toán Cô Trà, số điện thoại liên lạc (`098 345 93 93`), và địa chỉ cơ sở học tập chính tại `Giang Biên, Long Biên, Hà Nội`.

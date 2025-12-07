[![Go Version](https://img.shields.io/badge/Go-1.23-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

---

## 📋 Mục lục

- [Mô tả Bài toán](#mô-tả-bài-toán)
- [Tính năng](#tính-năng)
- [Bắt đầu Nhanh](#bắt-đầu-nhanh)
- [Tài liệu API](#tài-liệu-api)
- [Kiến trúc & Quyết định Thiết kế](#kiến-trúc--quyết-định-thiết-kế)
- [Đánh đổi Kỹ thuật](#đánh-đổi-kỹ-thuật)
- [Thách thức & Giải pháp](#thách-thức--giải-pháp)
- [Kiểm thử](#kiểm-thử)
- [Hạn chế & Cải tiến Tương lai](#hạn-chế--cải-tiến-tương-lai)
- [Sẵn sàng cho Production](#sẵn-sàng-cho-production)

---

## 🎯 Mô tả Bài toán

Bài toán yêu cầu xây dựng một URL Shortener Service giống như Bit.ly:

User có một URL dài:
https://example.com/very/long/path/to/resource?param1=value1&param2=value2

Muốn rút gọn thành URL ngắn hơn:
http://short.url/abc123

Khi người dùng truy cập URL rút gọn → server tự động redirect về URL gốc

Hệ thống theo dõi được số lượt click

API hỗ trợ:

Tạo URL rút gọn

Redirect

Xem thông tin URL

Liệt kê các URL đã tạo
---

## ✨ Tính năng

### Chức năng Cốt lõi
- ✅ **Rút gọn URL**: Tạo các alias ngẫu nhiên gồm 6 ký tự
- ✅ **Xác thực Người dùng**: Đăng ký và đăng nhập bảo mật dựa trên JWT
- ✅ **Alias Tùy chỉnh**: Hỗ trợ mã rút gọn do người dùng định nghĩa
- ✅ **Chuyển hướng Nhanh**: Chuyển hướng 302 với việc theo dõi click bất đồng bộ (async)
- ✅ **Phân tích Click**: Tăng bộ đếm thời gian thực
- ✅ **Phân trang**: Liệt kê hiệu quả tất cả các URL (Chỉ dành cho Admin)

### Bảo mật & Xác thực
- ✅ **Xác thực URL**: Kiểm tra định dạng bằng regex
- ✅ **Chặn URL Riêng tư**: Ngăn chặn localhost và các địa chỉ IP riêng
- ✅ **Làm sạch Đầu vào**: Chỉ chấp nhận alias là chữ và số
- ✅ **Xử lý Va chạm**: Tự động thử lại với mã mới

### Hiệu năng
- ✅ **Chỉ mục Cơ sở dữ liệu**: Chỉ mục duy nhất (unique index) trên alias để tra cứu O(1)
- ✅ **Connection Pooling**: Cấu hình tối đa 25 kết nối
- ✅ **Thao tác Nguyên tử (Atomic)**: Đếm click không bị race condition

---

## Cách chạy project:

### 1. Clone Repository

```bash
git clone https://github.com/Faleeeee/URL_Shortener.git
cd URL_Shortener
```

### 2. Cấu hình Biến Môi trường

```bash
# Sao chép file môi trường mẫu
cp .env.example .env

# Chỉnh sửa .env với thông tin database của bạn
# Đảm bảo DATABASE_URL trỏ đến database PostgreSQL cục bộ của bạn
```

### 3. Tạo Database

Đảm bảo PostgreSQL đang chạy và tạo database:

```bash
createdb -U postgres url_shortener
```

### 4. Chạy Database Migrations

```bash
```bash
psql -U postgres -d url_shortener -f migrations/000001_create_urls_table.up.sql
```

### 5. Cài đặt Dependencies

```bash
go mod download
```

### 6. Chạy Service

```bash
go run cmd/api/main.go
```

Service sẽ bắt đầu tại `http://localhost:8080`

### 7. Truy cập Tài liệu Swagger

Mở trình duyệt của bạn tại:
```
http://localhost:8080/swagger/index.html
```

### 8. Luồng Xác thực

1. **Đăng ký** người dùng mới:
   ```bash
   curl -X POST http://localhost:8080/auth/register \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "password": "password123"}'
   ```

2. **Đăng nhập** để lấy token:
   ```bash
   curl -X POST http://localhost:8080/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "password": "password123"}'
   ```
   Sao chép `token` từ phản hồi.

3. **Sử dụng Token** cho các endpoint được bảo vệ:
   ```bash
   curl -X POST http://localhost:8080/url/shorten \
     -H "Authorization: Bearer <YOUR_TOKEN>" ...
   ```

---

## ⚙️ Cấu hình

Dịch vụ sử dụng biến môi trường để cấu hình. Tất cả các cài đặt được định nghĩa trong file `.env`.

### Biến Môi trường

| Biến | Mô tả | Mặc định | Bắt buộc |
|----------|-------------|---------|----------|
| `SERVER_PORT` | Cổng server lắng nghe | `8080` | Không |
| `DATABASE_URL` | Chuỗi kết nối PostgreSQL | - | Có |
| `JWT_SECRET` | Khóa bí mật để ký JWT token | - | Có |
| `JWT_EXPIRATION` | Thời gian hết hạn JWT token | `24h` | Không |

### Ví dụ file `.env`

```bash
# Cấu hình Server
SERVER_PORT=8080

# Cấu hình Database
DATABASE_URL=postgres://postgres:123456@localhost:5432/url_shortener?sslmode=disable

# Cấu hình JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=24h
```


---

## 🏗️ Kiến trúc & Quyết định Thiết kế

### 1. Lựa chọn Database: **PostgreSQL**

#### Tại sao là PostgreSQL?

✅ **Được chọn vì:**
- **Tuân thủ ACID**: Đảm bảo tính toàn vẹn dữ liệu cho các yêu cầu đồng thời
- **Ràng buộc Duy nhất (Unique Constraints)**: Ngăn chặn trùng lặp alias ở cấp database
- **Thao tác Nguyên tử**: `UPDATE ... SET count = count + 1` ngăn chặn race conditions
- **Đánh chỉ mục**: Tra cứu nhanh O(1) trên cột alias
- **Giao dịch (Transactions)**: Hỗ trợ các thao tác nhiều bước
- **Độ tin cậy**: Đã được kiểm chứng trong môi trường production

---

### 2. Tạo Mã Rút gọn: **Base62 + Cryptographic Random**

#### Thuật toán

```go
Ký tự: [0-9A-Za-z] = 62 khả năng
Độ dài: 6 ký tự
Tổng số tổ hợp: 62^6 = 56,800,235,584 (56.8 tỷ)
```

#### Tại sao chọn cách tiếp cận này?

✅ **Ưu điểm:**
- **Kháng va chạm cao**: 56.8 tỷ tổ hợp đảm bảo hầu như không có va chạm
- **Ngắn gọn**: Chỉ 6 ký tự (thân thiện với người dùng)
- **Không thể đoán trước**: Tính ngẫu nhiên mật mã ngăn chặn việc đoán URL
- **Stateless**: Không cần đồng bộ hóa bộ đếm phân tán

❌ **Các lựa chọn thay thế đã xem xét:**

| Cách tiếp cận | Tại sao không chọn |
|----------|----------------|
| **Auto-increment ID + base62** | Dễ đoán (rủi ro bảo mật), lộ số lượng URL |
| **MD5/SHA hash + cắt ngắn** | Có thể va chạm, mã dài hơn (8-10 ký tự) |
| **Snowflake ID** | Yêu cầu phối hợp phân tán, quá mức cần thiết |
| **UUID** | Quá dài (36 ký tự) cho URL "rút gọn" |

#### Xử lý Va chạm

```go
1. Tạo mã base62 ngẫu nhiên 6 ký tự
2. Thử INSERT vào database
3. Nếu vi phạm ràng buộc duy nhất → thử lại (tối đa 3 lần)
4. Nếu vẫn thất bại → trả về lỗi
```

**Xác suất Va chạm**: Với 1 triệu URL, xác suất ≈ 0.001% (không đáng kể)

---

### 3. Thiết kế API: **REST**

#### Tại sao REST thay vì GraphQL/gRPC?

✅ **REST được chọn vì:**
- **Đơn giản**: Dễ hiểu với mọi lập trình viên
- **Phù hợp hoàn hảo**: Các thao tác CRUD ánh xạ tự nhiên với các phương thức HTTP
- **Caching**: Caching của trình duyệt và CDN hoạt động ngay lập tức
- **Chuyển hướng**: Hỗ trợ chuyển hướng HTTP 302 gốc
- **Công cụ**: Swagger/OpenAPI cho tài liệu

---

### 4. Chiến lược Đồng thời (Concurrency)

#### Vấn đề: Race Conditions

**Kịch bản 1**: Hai người dùng tạo URL cùng lúc
**Giải pháp**: Ràng buộc duy nhất (unique constraint) của database trên cột `alias`

```sql
CREATE UNIQUE INDEX idx_alias ON urls(alias);
```

**Kịch bản 2**: Nhiều sự kiện click cho cùng một URL
**Giải pháp**: Cập nhật SQL nguyên tử (Atomic SQL update)

```sql
UPDATE urls SET click_count = click_count + 1 WHERE alias = ?
```

**Kịch bản 3**: Va chạm đọc-sửa-ghi (Read-modify-write)
**Giải pháp**: Sử dụng `QueryRow` + `Exec` với transactions

---

## Thách thức & Giải pháp

### Thách thức 1: Tạo URL Đồng thời

**Vấn đề**: Hai yêu cầu với cùng một URL dài đến cùng lúc

**Giải pháp**: Logic thử lại với exponential backoff
```go
for i := 0; i < MaxRetries; i++ {
    alias := GenerateShortCode()
    if err := repo.Create(alias); err == nil {
        return alias, nil
    }
    // Nếu trùng lặp, thử lại
}
```

**Thay thế đã xem xét**: Kiểm tra xem URL có tồn tại trước không → Vẫn có thể xảy ra Race condition

**Bài học**: Trong hệ thống phân tán, **Optimistic Locking** (thử lại khi lỗi) thường hiệu quả hơn Pessimistic Locking (khóa trước) khi tỷ lệ va chạm thấp.

---

### Thách thức 2: Ngăn chặn URL Riêng tư

**Vấn đề**: Người dùng có thể rút gọn `http://localhost:9090/admin` và chia sẻ nó

**Giải pháp**: Danh sách đen các mẫu riêng tư phổ biến
```go
if strings.Contains(host, "localhost") ||
   strings.HasPrefix(host, "127.") ||
   strings.HasPrefix(host, "192.168.") { ... }
```

**Hạn chế**: Không bắt được tất cả các dải riêng tư (ví dụ: `172.16-31.x.x`)

**Tương lai**: Sử dụng thư viện khớp CIDR để kiểm tra toàn diện

**Bài học**: Đừng bao giờ tin tưởng đầu vào từ người dùng (Zero Trust). Validation cần được thực hiện ở nhiều lớp (Application layer + Network layer).

---

### Thách thức 3: Race Conditions Bộ đếm Click

**Vấn đề**: Nhiều click → mất cập nhật

**Cách tiếp cận Tồi** (race condition):
```go
url := repo.FindByAlias(alias)
url.ClickCount++
repo.Update(url)  // Mất cập nhật!
```

**Cách tiếp cận Tốt** (nguyên tử):
```sql
UPDATE urls SET click_count = click_count + 1 WHERE alias = ?
```

**Bài học**: Luôn sử dụng các thao tác nguyên tử cho bộ đếm

**Bài học**: Hiểu rõ cơ chế khóa và tính nguyên tử (Atomicity) của database là cực kỳ quan trọng để đảm bảo tính đúng đắn của dữ liệu trong môi trường đa luồng.

---

### Thách thức 4: Tạo Code Swagger

**Vấn đề**: Tài liệu Swagger không đồng bộ với code

**Giải pháp**: Sử dụng chú thích `swag` trong code
```go
// @Summary Create a shortened URL
// @Param request body domain.ShortenRequest true "URL to shorten"
func (h *URLHandler) ShortenURL(c *gin.Context) { ... }
```

Sau đó tự động tạo:
```bash
swag init -g cmd/api/main.go
```

**Lợi ích**: Nguồn sự thật duy nhất (code)

**Bài học**: **Documentation-as-Code** giúp tài liệu luôn sống và chính xác, tránh việc tài liệu bị "thiu" (outdated) so với thực tế triển khai.

---

## Hạn chế & Cải tiến Tương lai

### Hạn chế Hiện tại

| Hạn chế | Tác động | Ưu tiên |
|------------|--------|----------|
| **Không Giới hạn Tốc độ** | Dễ bị lạm dụng | CAO |
| **Không Hết hạn URL** | Database tăng trưởng vô hạn | TRUNG BÌNH |
| **Không Dashboard Phân tích** | Thông tin chi tiết hạn chế | THẤP |
| **Không Tên miền Tùy chỉnh** | Chỉ localhost:8080 | THẤP |
| **Không Xác thực URL Độc hại** | Rủi ro lừa đảo (phishing) | TRUNG BÌNH |

### Cải tiến Tương lai

#### Giai đoạn 1: Bảo mật & Độ tin cậy
- [ ] **Giới hạn Tốc độ**: 100 yêu cầu/giờ mỗi IP
- [ ] **API Keys**: Xác thực cho các gói trả phí
- [ ] **Danh sách đen URL**: Chặn các tên miền độc hại đã biết
- [ ] **Hỗ trợ HTTPS**: Chứng chỉ TLS qua Let's Encrypt

#### Giai đoạn 2: Tính năng
- [ ] **Tạo Mã QR**: Tự động tạo mã QR cho các URL rút gọn
- [ ] **Hết hạn**: Tự động xóa sau N ngày/click
- [ ] **Bảo vệ Mật khẩu**: Bảo mật URL rút gọn bằng mật khẩu
- [ ] **Tên miền Tùy chỉnh**: Hỗ trợ `go.yourcompany.com`

#### Giai đoạn 3: Phân tích
- [ ] **Phân tích Click**: Theo dõi user agent, người giới thiệu (referrer), vị trí địa lý
- [ ] **Admin Dashboard**: Giao diện Web để quản lý URL
- [ ] **Thống kê Thời gian thực**: WebSocket cho cập nhật click trực tiếp

#### Giai đoạn 4: Quy mô
- [ ] **Redis Caching**: Cache các URL hot (quy tắc 80/20)
- [ ] **Read Replicas**: Mở rộng đọc PostgreSQL
- [ ] **Tích hợp CDN**: Cloudflare cho chuyển hướng toàn cầu
- [ ] **Database Sharding**: Phân vùng theo hash(alias)

---

## Sẵn sàng cho Production

### Còn thiếu gì cho Production?

| Yêu cầu | Trạng thái | Giải pháp |
|-------------|--------|----------|
| **SSL/TLS** | ❌ Chưa triển khai | Sử dụng Nginx reverse proxy + Let's Encrypt |
| **Giám sát** | ❌ Chưa triển khai | Thêm Prometheus + Grafana |
| **Theo dõi Lỗi** | ❌ Chưa triển khai | Tích hợp Sentry hoặc Rollbar |
| **CI/CD** | ❌ Chưa triển khai | GitHub Actions để test + deploy |
| **Load Balancer** | ❌ Chưa triển khai | Nginx hoặc AWS ALB |

### Kiến trúc Triển khai (Đề xuất)

```
Internet
   ↓
Cloudflare CDN (Bảo vệ DDoS, caching)
   ↓
Nginx Load Balancer (SSL termination, giới hạn tốc độ)
   ↓
Go Service (3 bản sao)
   ↓
PostgreSQL Primary + 2 Read Replicas
   ↓
Redis Cache (cache URL hot)
```

## Lược đồ Cơ sở dữ liệu

```sql
CREATE TABLE urls (
    id BIGSERIAL PRIMARY KEY,
    alias VARCHAR(16) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    click_count BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_alias ON urls(alias);
CREATE INDEX idx_created_at ON urls(created_at);
```

**Chiến lược Chỉ mục:**
- `idx_alias`: Chỉ mục duy nhất cho tra cứu alias O(1)
- `idx_created_at`: Cho các truy vấn phân tích (URL mới nhất trước)

---


**Công nghệ Sử dụng:**
- Go 1.23
- Gin Web Framework
- PostgreSQL 16
- Swagger
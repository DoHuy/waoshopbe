## What do you like about Golang?

### 1. High Performance
Go là một ngôn ngữ biên dịch, được chuyển đổi trực tiếp thành mã máy. Điều này giúp tốc độ thực thi của nó nhanh hơn đáng kể so với các ngôn ngữ thông dịch. Hiệu năng của Go thường được so sánh với C++.

### 2. Native Concurrency (Goroutines)
Một trong những đặc điểm nổi bật của Go chính là mô hình lập trình đồng thời. 
**Goroutines** là các "luồng" (threads) siêu nhẹ được quản lý bởi Go runtime chứ không phải hệ điều hành (OS). Bạn có thể khởi tạo hàng triệu luồng như vậy với chi phí bộ nhớ cực thấp. Kết hợp với  
**Channels** phục vụ cho việc giao tiếp, Go giúp việc triển khai các logic bất đồng bộ phức tạp trở nên an toàn và dễ dàng hơn.

### 3. Static Typing and Safety
Go giúp phát hiện lỗi ngay tại thời điểm biên dịch (compile time) thay vì lúc thực thi (runtime). Điều này giúp mã nguồn trở nên ổn định hơn và tối ưu hóa các công cụ hỗ trợ lập trình (như IDE hay tái cấu trúc mã - refactoring). Khác với một số ngôn ngữ tĩnh khác, khả năng suy luận kiểu (type inference) của Go giữ cho cú pháp luôn gọn gàng và ít rườm rà hơn.

### 4. Simplicity by Design
Go was designed at Google to be productive and readable. It has a small language specification, which means developers can become proficient quickly. There is usually one "idiomatic" way to solve a problem, which reduces friction in large teams.

### 5. Single Binary Deployment
Go được Google thiết kế với mục tiêu tối ưu hiệu suất công việc và khả năng đọc hiểu. Ngôn ngữ này có bộ quy tắc (specification) tinh gọn, giúp các lập trình viên nhanh chóng nắm vững và sử dụng thành thạo. Đặc biệt, Go thường chỉ có một cách giải quyết vấn đề theo kiểu "chuẩn mực" (idiomatic) duy nhất, điều này giúp giảm thiểu sự xung đột hoặc bất đồng ý kiến khi làm việc trong các đội ngũ lớn.

---

## Golang vs. PHP: Key Differences

| Feature | Golang | PHP |
| :--- | :--- | :--- |
| **Execution** | Compiled (Machine Code) | Interpreted (Zipped via OpCache/JIT) |
| **Concurrency** | Native (Goroutines & Channels) | Traditionally Synchronous (Requires extensions like Swoole for async) |
| **Typing** | Statically Typed | Dynamically Typed (with improved type hinting) |
| **Deployment** | Single Binary | Requires Web Server (Nginx/Apache) + PHP-FPM |
| **Performance** | Very High | Moderate (Fast for its class, but slower than Go) |
| **Ecosystem** | Strong for Microservices, Cloud, Infra | King of Web/CMS (Laravel, WordPress) |

### When to Choose Go
- **Microservices:** Với đặc tính nhẹ (ít tốn tài nguyên) và hiệu năng cao, Go là sự lựa chọn hoàn hảo cho các hệ thống phân tán.
- **High Concurrency:** Các ứng dụng như máy chủ trò chuyện (chat servers), truyền phát trực tuyến thời gian thực (real-time streaming), hoặc các cổng API (API gateways) chịu tải cao.
- **Cloud Native:** Go là ngôn ngữ của điện toán đám mây (Docker, Kubernetes và Terraform đều được xây dựng bằng Go).

### When to Choose PHP
- **Rapid Prototyping:** Đối với các ứng dụng CRUD thông thường, PHP (đặc biệt là khi kết hợp với Laravel) mang lại tốc độ phát triển cực kỳ nhanh chóng.
- **Content Management:** Nếu bạn đang xây dựng một trang blog hoặc một website tiêu chuẩn, nơi mà SEO và việc nhập liệu nội dung là những ưu tiên hàng đầu.
- **Shared Hosting:** PHP vẫn là ngôn ngữ được hỗ trợ rộng rãi nhất trên hầu hết tất cả các nhà cung cấp dịch vụ lưu trữ (hosting).

---

## What backend problems do you enjoy solving the most?
Thiết kế Hệ thống và Tối ưu hóa Luồng Dữ liệu.

- **Thiết kế Hệ thống**: Mình yêu thích thử thách việc lắp ghép các thành phần khác nhau lại với nhau—như Cơ sở dữ liệu, Bộ nhớ đệm (Redis), và Hàng đợi tin nhắn (Kafka/RabbitMQ)—để đảm bảo chúng vận hành nhịp nhàng mà không gây ra các điểm nghẽn (bottlenecks).

- **Tối ưu hiệu năng**: Cảm giác cực kỳ thỏa mãn khi biến một hệ thống đang trì trệ thành một "cỗ máy" hiệu suất cao, luôn duy trì thời gian phản hồi dưới 100ms.

- **Khả năng mở rộng** (Scalability): Việc xây dựng những hệ thống có thể mở rộng mượt mà từ chỗ phục vụ hàng trăm lên đến hàng triệu người dùng.

---
## What do you find most difficult in backend development?
Đảm bảo sự ổn định của hệ thống, xử lý lưu lượng truy cập cao, duy trì khả năng mở rộng và hỗ trợ vận hành trên quy mô toàn cầu.

---

## A specific project situation you solved that you are most proud of ?

### 1. Xử lý High-Concurrency trong tạo đơn hàng (Order Creation)
Thử thách: Tránh tình trạng Race Condition (tranh chấp dữ liệu), đảm bảo không bán quá số lượng tồn kho (Over-selling) và bảo vệ hệ thống khỏi các đợt "spam" request.

Giải pháp:
- Atomic Inventory Control với Redis Lua Script: Thay vì đọc tồn kho từ DB rồi mới trừ (dễ gây sai số khi chạy đa luồng), mình đẩy toàn bộ logic kiểm tra và trừ tồn kho vào một Lua Script chạy trực tiếp trên Redis. Vì Redis đơn luồng, script này đảm bảo tính nguyên tử (Atomicity) tuyệt đối—hoặc là trừ thành công và trả về token tạo đơn, hoặc là thông báo hết hàng ngay lập tức mà không có độ trễ mạng giữa các bước.
- Multi-layer Rate Limiting:  Chặn các bot hoặc cá nhân cố tình spam request từ một địa chỉ IP và giới hạn số lượng request tạo đơn đang được xử lý đồng thời cho cùng một sản phẩm để tránh làm nghẽn Worker xử lý phía sau.

### 2. Bảo mật cho Anonymous Chat (Chat ẩn danh giữa khách và Shop)
- Token Bucket/Leaky Bucket (Redis): Sử dụng Redis để giới hạn số lượng tin nhắn trong một khoảng thời gian (ví dụ: tối đa 5 tin nhắn/10 giây) từ cùng IP Address.
- Session-based Limiting: Khi một phiên chat ẩn danh được khởi tạo, cấp một JWT (Short-lived).
- Phòng tránh gửi các script mã độc thông qua chat
- Chỉ chia sẻ thông tin cần thiết liên quan tới sản phẩm không chia sẻ email, hay sdt cá nhân



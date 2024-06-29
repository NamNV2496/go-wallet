# Setup

## postgres

```txt
1. download image

	docker pull postgres:16.3-alpine

2. run container

	docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=wallet -p 5432:5432 -d postgres:16.3-alpine
```

## sqlc

```txt
1. Download sqlc binary

	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

2. confile file `sql.yaml`
3. run command

	sqlc generate

```

# Technologies

	1. sqlc - "https://github.com/sqlc-dev/sqlc"
	2. golang-migrate - "github.com/golang-migrate/migrate/v4"
	3. migration - "github.com/rubenv/sql-migrate" or "github.com/golang-migrate/migrate/v4/database/postgres"

# Ý tưởng
- 1 App ví điện tử (go-wallet)
- monolithic: gồm api gateway, wallet (manage transaction), account(manage balance and account)
- API 
	- tạo tk có gửi OTP qua email
	- login
	- lấy ra thông tin tài khoản theo token
	- Lấy ra danh sách các tài khoản gợi ý người nhận theo dữ liệu người dùng nhập (user_name hoặc số điện thoại)
	- api chuyển tiền xử lý chuyển tiền async qua kafka. (add case high workload and case fail) and idempotence
	- search lịch sử chuyển khoản (ngày từ ... đến ..., theo số stk người nhận)
	- nạp tiền: hold ở BE timeoout chờ call API fake to confirm through kafka

![alt text](image.png)

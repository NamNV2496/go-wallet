<!-- # Ý tưởng
- 1 App ví điện tử (go-wallet)
- monolithic: gồm api gateway, wallet (manage transaction), account(manage balance and account)
- API 
	- tạo tk có gửi OTP qua email
	- login
	- lấy ra thông tin tài khoản theo token
	- Lấy ra danh sách các tài khoản gợi ý người nhận theo dữ liệu người dùng nhập (user_name hoặc số điện thoại)
	- api chuyển tiền xử lý chuyển tiền async qua kafka. (add case high workload and case fail) and idempotence
	- search lịch sử chuyển khoản (ngày từ ... đến ..., theo số stk người nhận)
	- nạp tiền: hold ở BE timeoout chờ call API fake to confirm through kafka -->

<details>

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
</details>

# This is a basic wallet was implemented by Golang.

# Technologies

	1. sqlc - "https://github.com/sqlc-dev/sqlc"
	2. golang-migrate - "github.com/golang-migrate/migrate/v4"
	3. migration - "github.com/rubenv/sql-migrate" or "github.com/golang-migrate/migrate/v4/database/postgres"
	4. kafka - "github.com/IBM/sarama"
	5. jwt - "github.com/golang-jwt/jwt"
	6. paseto - "github.com/o1egl/paseto"


![alt text](docs/flow.png)

# Feature

	1. Login with jwt or paseto with session to renew token
	2. Get all accounts of user
	3. Create new user (email OTP confirm)
	4. Create new account
	5. Find account by username of phone
	6. Transfer money between 2 accounts (with confirm result through kafka)
	7. Get history of transfer money

# How to run
## Run with docker

	1. docker compose up

## Run with normal

	1. `make server`

Database

![alt text](docs/db.png)

Demostration

![alt text](docs/image.png)
![alt text](docs/image-1.png)
![alt text](docs/image-2.png)
![alt text](docs/image-3.png)
![alt text](docs/image-4.png)


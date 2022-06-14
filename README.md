# I. Hướng dẫn:
## 1. Get thư viện:
```bash
go get
```

hoặc
```bash
go mod tidy
```

## 3. Run SQL script:
Để tạo bảng mẫu, location:
```
scripts/sql/schema.sql
```

## 4. Download swag:
```bash
go get github.com/swaggo/swag/cmd/swag@v1.7.0
```

## 5. Download wire:
```bash
go get github.com/google/wire/cmd/wire@v0.5.0
```

### Windows:
```cmd
set ENV=dev&&swag init&&wire&&go run .
```

## 6. Go to:
http://localhost:8099/swagger/index.html


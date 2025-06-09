# Nepeta-cataria

### 1. clone project

```bash
git clone <repository-url>
cd Nepeta-cataria
```

### 2. install

```bash
go mod download

go install github.com/cosmtrek/air@latest
```

### 3. env

```bash
cp .env.example .env
```

### 4. 啟動服務

```bash
air
```

服務將在 `http://localhost:3100` 啟動

### 本地開發設定

```bash
# 安裝開發工具
go install github.com/cosmtrek/air@latest
go install github.com/swaggo/swag/cmd/swag@latest

# 產生 Swagger 文件
swag init --parseDependency --parseInternal
```

## API 文件

### 訪問 Swagger UI

swagger gin 不支援openapi v3
無提供swagger介面，[swagger-editor](https://editor.swagger.io/)也無法正確解析

### 更新 openapi 文件

修改了 API 後，需要重新生成 openapi 文檔：

```bash
swag init --parseDependency --parseInternal
```

並且將[swagger.yaml](./docs/swagger.yaml)丟進postman裡

## 專案結構說明

```
Nepeta-cataria/
├── .devcontainer/          # Dev Container 配置
├── boot/                   # 應用啟動和依賴注入
├── config/                 # 配置管理
├── controllers/            # HTTP 控制器
├── docs/                   # Swagger 生成的文件
├── global/                 # 全域變數和設定
├── middleware/             # 中間件
├── models/                 # 資料模型
├── pkg/                    # 外部套件和工具
├── repositories/           # 資料訪問層
├── routes/                 # 路由定義
├── services/               # 業務邏輯層
├── tools/                  # 開發工具
├── utils/                  # 工具函數
├── main.go                 # 應用入口
└── go.mod                  # Go 模組文件
```

## 其他

### 編譯makejudge工具
```bash
go build -o makejudge_tool ./tools/makejudge_cli/main.go
./makejudge_tool generate -problem_dir <dir>
./makejudge_tool testgen -problem_dir judges/testload -code_file judges/testload/test.py -lang python -output test_runner.py
```

### swagger path
```bash
/api/docs/index.html
```
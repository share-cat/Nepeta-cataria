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

## API 文件與

### 訪問 Swagger UI

swagger gin 不支援openapi v3
無提供swagger介面，[swagger-editor](https://editor.swagger.io/)也無法正確解析

### 更新 Swagger 文件

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
```


# TODO list



### 需要完成的功能

- [ ] **判題系統**
  - [ ] 代碼執行環境設置
  - [ ] 測試案例執行
  - [ ] 結果比對邏輯
  - [ ] 執行時間和記憶體使用統計

- [ ] **實時反饋**
  - [ ] 提交狀態即時更新
  - [ ] 執行結果顯示
  - [ ] 錯誤訊息處理

- [ ] **題目管理優化**
  - [x] 從 ZIP 檔案載入題目 ([`LoadJudgeController`](controllers/loadjudge_controller.go))
  - [ ] 題目測試資料管理
  - [ ] 題目標籤和分類系統完善



### 管理員功能
- [x] 基礎用戶 CRUD ([`UserController`](controllers/user_controller.go))
- [ ] **管理員介面**
  - [ ] 管理員登入頁面
  - [ ] 用戶管理面板
  - [ ] 題目管理面板
  - [ ] 系統統計面板

- [ ] **題目管理**
  - [ ] 題目新增介面 (圖形化)
  - [ ] 題目編輯介面
  - [ ] 測試資料上傳
  - [ ] 題目發布/下架

- [ ] **權限控制**
  - [ ] 管理員權限驗證
  - [ ] 角色分配系統
  - [ ] 操作日誌記錄

### 用戶管理
- [x] 用戶註冊/登入 API ([`AuthController`](controllers/auth_controller.go))
- [ ] 用戶個人資料頁面
- [ ] 密碼修改介面
- [ ] 頭像上傳功能



### 答題功能
- [x] 答案提交 API ([`SolutionController.SubmitSolution`](controllers/solution_controller.go))
- [x] 用戶答題記錄 ([`SolutionController.GetUserSolutions`](controllers/solution_controller.go))


### 成績系統
- [x] 答案評分邏輯 ([`SolutionService.UpdateSolutionResult`](services/solution_service.go))


### 社交功能
- [x] 討論區 API ([`DiscussionController`](controllers/discussion_controller.go))
- [x] 聊天系統 API ([`ChatSessionController`](controllers/chatsession_controller.go))



### 系統優化
- [ ] 代碼執行沙盒環境
- [ ] 檔案上傳安全性檢查
- [ ] API 速率限制
- [ ] 快取機制實作
- [ ] 錯誤處理優化

### 測試
- [ ] 單元測試撰寫
- [ ] API 整合測試
- [ ] 前端 E2E 測試


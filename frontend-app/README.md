## 專案需求
1. 下載並啟動docker desktop
2. 下載並安裝go
3. 下載或擁有pnpm並執行pnpm install

## 專案啟動指南
本專案包含前端和後端部分，以下為啟動步驟：

## 後端啟動
切換到後端目錄：

```jsx
cd go-task-api
echo "DB_URL=postgres://user:password@localhost:5432/mydb" > .env
echo "JWT_SECRET=mysecretkey" >> .env
echo "PORT=8080" >> .env
```

使用 Docker Compose 啟動後端服務：

```jsx
docker-compose up -d --build
go run main.go
```

注意： 啟動後，服務將在指定的容器端口運行，請確認 Docker 正常運作。

## 前端啟動
切換到前端應用程式目錄：

```jsx
cd frontend-app
```

使用 pnpm 啟動開發伺服器：

```jsx
pnpm dev
```

注意： 前端應用將在本地伺服器上運行，預設訪問 http://localhost:5173。

## 注意事項
請確保已安裝 Docker 和 pnpm。
如果出現連線錯誤，請確認後端服務是否正常啟動，並檢查 Docker 是否處於運行狀態。
若遇到端口衝突，請檢查是否有其他應用佔用相同端口，或修改 docker-compose.yml 和 Vite 設定檔。

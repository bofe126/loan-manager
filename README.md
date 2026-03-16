# 贷款管理系统 (Loan Manager)

一个现代化的贷款管理桌面应用，采用 Wails 框架构建，支持多种还款方式计算和可视化展示。

## ✨ 特性

- 📊 **多种还款方式**：支持等额本息、等额本金、先息后本三种计算方式
- 💰 **贷款管理**：完整的增删改查功能，支持批量录入
- 📈 **数据可视化**：趋势图表展示月还款和剩余金额变化
- 🎨 **现代化界面**：金融科技暗黑主题，玻璃态效果设计
- 💾 **数据持久化**：SQLite 数据库，数据存储在用户目录
- 🚀 **跨平台支持**：Windows、macOS、Linux 一键打包
- 📦 **单文件部署**：16MB 可执行文件，无需安装依赖

## 🛠️ 技术栈

### 后端
- **Go 1.21+**：高性能后端服务
- **GORM**：优雅的 ORM 框架
- **SQLite**：轻量级嵌入式数据库

### 前端
- **Vue 3**：渐进式 JavaScript 框架
- **TypeScript**：类型安全的开发体验
- **Vite**：极速的前端构建工具
- **Chart.js**：强大的图表库
- **Bootstrap 5**：响应式 UI 组件

### 桌面框架
- **Wails v2**：使用 Go + Web 技术构建原生桌面应用

## 📦 安装与使用

### 直接使用（推荐）

1. 下载最新版本的 `loan-manager-wails.exe`
2. 双击运行即可
3. 数据库会自动创建在 `C:\Users\<用户名>\.loan-manager\loans.db`

### 从源码构建

#### 前置要求
- Go 1.21+
- Node.js 16+
- Wails CLI v2

#### 安装 Wails CLI
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

#### 克隆项目
```bash
git clone <仓库地址>
cd loan-manager-wails
```

#### 开发模式
```bash
wails dev
```

#### 生产构建
```bash
wails build
```

构建产物位于 `build/bin/` 目录。

## 🎯 功能说明

### 贷款管理
- **新增贷款**：填写借款人、银行、金额、利率等信息
- **编辑贷款**：修改贷款详情
- **删除贷款**：删除不需要的贷款记录
- **还款记录**：记录每次还款金额，自动更新剩余贷款

### 还款方式

#### 等额本息
每月还款金额固定，包含本金和利息。适合收入稳定的借款人。

#### 等额本金
每月偿还相同本金，利息逐月递减。前期还款压力大，总利息较少。

#### 先息后本
每月只还利息，到期一次性偿还本金。适合短期周转。

### 数据可视化
- **统计卡片**：实时显示月还款总额、剩余贷款、贷款数量
- **趋势图表**：展示未来 24 个月的还款趋势和剩余金额变化
- **交互联动**：点击统计卡片切换图表视图

## 🎨 界面预览

应用采用金融科技暗黑主题设计：
- 深蓝色背景 (#0A0E27)
- 蓝色主色调 (#0080FF)
- 绿色强调色 (#39FF14)
- 玻璃态效果（Glassmorphism）
- 黄金分割比例布局

## 📁 项目结构

```
loan-manager-wails/
├── backend/                 # Go 后端
│   ├── database/           # 数据库初始化
│   ├── models/             # 数据模型
│   └── services/           # 业务逻辑
├── frontend/               # Vue 3 前端
│   ├── src/
│   │   ├── components/     # UI 组件
│   │   ├── views/          # 页面视图
│   │   ├── composables/    # 组合式函数
│   │   ├── router/         # 路由配置
│   │   └── types/          # TypeScript 类型
│   └── package.json
├── build/                  # 构建资源
│   ├── appicon.png         # 应用图标
│   └── windows/
│       └── icon.ico        # Windows 图标
├── app.go                  # Wails 应用主文件
├── main.go                 # Go 入口文件
└── wails.json              # Wails 配置
```

## 🔧 配置说明

### 数据库位置
- **Windows**: `C:\Users\<用户名>\.loan-manager\loans.db`
- **macOS**: `~/.loan-manager/loans.db`
- **Linux**: `~/.loan-manager/loans.db`

### 窗口尺寸
- 默认宽度：800px
- 默认高度：600px
- 可在 `main.go` 中修改

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 🙏 致谢

- [Wails](https://wails.io/) - 优秀的桌面应用框架
- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架
- [GORM](https://gorm.io/) - Go ORM 库
- [Chart.js](https://www.chartjs.org/) - 图表库

---

**Co-Authored-By: Claude Sonnet 4.6**

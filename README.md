# 炫酷工具箱服务

一个基于Magic UI设计灵感的炫酷在线工具集合平台，具有现代化的动画效果和科技感。

## 功能特点

- 🎨 **炫酷设计**: 受Magic UI启发的现代化界面，包含多种动画效果
- 🚀 **9个实用工具**: HTML转图片、JSON格式化、Base64编码/解码、正则表达式测试、URL编码/解码、哈希计算、时间戳转换、UUID生成器、颜色选择器
- ✨ **丰富动画**: 光束效果、粒子动画、渐变背景、霓虹光效、浮动图标等
- 📱 **响应式设计**: 完美适配桌面和移动设备
- 🔍 **实时搜索**: 支持工具名称和描述的实时搜索
- 🎯 **交互体验**: 流畅的悬停效果和过渡动画

## 技术栈

- **后端**: Go + Gin框架
- **前端**: HTML5 + Tailwind CSS + JavaScript
- **图标**: Font Awesome
- **动画**: CSS3动画和JavaScript交互
- **字体**: Google Fonts (Inter)

## 运行说明

### 前置要求

- Go 1.14+
- Node.js (可选，用于前端开发)

### 运行步骤

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd toolkit-service
   ```

2. **安装Go依赖**
   ```bash
   go mod tidy
   ```

3. **运行服务**
   ```bash
   go run main.go
   ```

4. **访问应用**
   打开浏览器访问 `http://localhost:8080`

### 自定义端口

设置环境变量 `PORT`:
```bash
export PORT=3000
go run main.go
```

## 项目结构

```
toolkit-service/
├── main.go              # 主程序入口
├── go.mod               # Go模块文件
├── go.sum               # 依赖版本锁定
├── build.sh             # 构建脚本
├── resources/static/    # 静态资源
│   ├── index.html       # 主页
│   ├── html2img/        # HTML转图片工具
│   ├── json-formatter/  # JSON格式化工具
│   ├── base64-encoder/  # Base64编码/解码工具
│   ├── regex-tester/    # 正则表达式测试工具
│   ├── url-encoder/     # URL编码/解码工具
│   ├── hash-calculator/ # 哈希计算工具
│   ├── timestamp-converter/ # 时间戳转换工具
│   ├── uuid-generator/  # UUID生成器工具
│   └── color-picker/    # 颜色选择器工具
├── .gitignore           # Git忽略文件
├── LICENSE              # 许可证
└── README.md           # 说明文档
```

## 设计特色

### 动画效果
- **背景渐变**: 15秒循环的多色渐变动画
- **光束效果**: 垂直移动的光束，营造科技感
- **粒子动画**: 浮动的粒子效果
- **霓虹光效**: 卡片和按钮的发光效果
- **悬停动画**: 卡片的3D变换和光效扫过

### 视觉设计
- **深色主题**: 现代化的深色背景
- **渐变色彩**: 紫色、粉色、蓝色的渐变搭配
- **毛玻璃效果**: 半透明的背景模糊效果
- **圆角设计**: 现代化的圆角元素
- **阴影效果**: 多层次的光影效果

## API接口

### 获取工具列表
```
GET /api/tools
```

响应格式:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "HTML转图片",
      "description": "将HTML代码转换为图片格式，方便分享和保存",
      "icon": "fas fa-image",
      "url": "/html2img"
    }
    // ... 更多工具
  ]
}
```

## 登录功能

- 默认密码: `187187187`
- 登录状态保存在localStorage中
- 支持登出功能

## 浏览器兼容性

- Chrome 80+
- Firefox 75+
- Safari 13+
- Edge 80+

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request来改进这个项目！

## 预览

主页包含:
- 全屏英雄区域with动画标题
- 实时搜索功能
- 9个工具卡片的网格布局
- 悬停动画和光效
- 响应式导航栏
- 关于部分with特色介绍
- 现代化页脚

每个工具页面都有独立的功能界面和一致的视觉风格。
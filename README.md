<div align="center">

# CloudFrog


### Secure File Sharing Made Simple

[![Go](https://img.shields.io/badge/go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![React](https://img.shields.io/badge/react-18.2.0-61DAFB?style=for-the-badge&logo=react)](https://reactjs.org)
[![License](https://img.shields.io/badge/license-MIT-green?style=for-the-badge)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-supported-2496ED?style=for-the-badge&logo=docker)](https://docker.com)

Cloud-based file sharing platform with automatic link expiration and secure SFTP storage.


</div>

---

## ✨ Features

<div align="center">
<img src="docs/screenshots/upload.png" alt="Upload Demo" width="600"/>
</div>

-  **Lightning Fast** uploads with progress tracking
-  **Secure Storage** using SFTP with encryption
-  **Auto-Expiring Links** for enhanced security
-  **Mobile Responsive** design
-  **Drag & Drop** file uploading
-  **Quick Copy** sharing links
-  **Rate Limiting** protection

## 🛠️ Tech Stack

### Backend
- **Framework:** [Gin](https://gin-gonic.com/) (Go)
- **Database:** MongoDB
- **Storage:** SFTP Server
- **Container:** Docker

### Frontend
- **Framework:** React 18 with Vite
- **Styling:** Tailwind CSS
- **Animation:** Framer Motion
- **Routing:** React Router v6
- **HTTP:** Axios

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
# Clone the repository
git clone gitrepo

# Enter the project directory
cd cloudfrog

# Start with Docker Compose
docker-compose up -d
```

Visit http://localhost:3000 to see the application.

### Manual Setup

#### Prerequisites
- Go 1.21+
- Node.js 18+
- SFTP Server
- MongoDB Database

#### Backend Setup
```bash
cd backend
cp .env.example .env    # Configure your environment variables
go mod download
go run main.go
```

#### Frontend Setup
```bash
cd frontend
cp .env.example .env    # Configure your environment variables
npm install
npm run dev
```

## ⚙️ Configuration

### Backend Environment Variables
```env
SFTP_HOST=your-sftp-host
SFTP_PORT=9333
SFTP_USER=your-username
SFTP_PASSWORD=your-password
SFTP_PATH=/path/to/storage
MONGODB_CONNECTION_STRING=your-mongodb-uri
BASE_URL=http://localhost:8080
```

### Frontend Environment Variables
```env
VITE_API_URL=http://localhost:8080
```

## 📝 API Documentation

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/upload` | POST | Upload a new file |
| `/download/:shortcode` | GET | Download a file |
| `/cleanup` | DELETE | Clean expired files (Admin) |

## 🔒 Security Features

-  Rate limiting protection
-  Automatic file expiration
-  Secure SFTP storage
-  Random shortcode generation
-  No public file listing


## 🐛 Troubleshooting

### Common Issues

1. **SFTP Connection Failed**
   ```bash
   # Check SFTP credentials
   sftp your-username@your-sftp-host -P 9333
   ```

2. **MongoDB Connection Issues**
   - Verify connection string
   - Check network access

3. **File Upload Failed**
   - Check file size limits
   - Verify SFTP permissions


---

<div align="center">
Made with ❤️ by Duypresso
</div>

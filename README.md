# 🔒 envSync

> **Platform Engineering & DevSecOps Config Sync Engine**  
> เครื่องมือ CLI และ API ส่วนกลางระดับองค์กร สำหรับจัดการ ซิงก์ และรักษาความปลอดภัยไฟล์ `.env` ป้องกัน Secret รั่วไหล และขจัดปัญหา Configuration Drift ระหว่างทีมพัฒนา

[![CI Pipeline](https://github.com/jakkayy/envSync/actions/workflows/ci.yml/badge.svg)](https://github.com/jakkayy/envSync/actions)
[![Security Audit](https://github.com/jakkayy/envSync/actions/workflows/security.yml/badge.svg)](https://github.com/jakkayy/envSync/actions)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/jakkayy/envSync)](https://github.com/jakkayy/envSync/releases)

---

## 📌 ปัญหาที่ `envSync` เข้ามาแก้ไข (Problem Statement)

1. **"ขอไฟล์ `.env` หน่อย":** นักพัฒนาใหม่เข้าทีมต้องคอยทักแชตขอไฟล์ `.env` จากเพื่อนร่วมทีม เสียเวลา Onboard
2. **Secret รั่วไหลผ่านแชต:** ทีมแอบส่ง API Keys, รหัสผ่านฐานข้อมูล หรือ JWT Secrets ผ่าน Slack, LINE, หรือ Email โดยไม่เข้ารหัส
3. **Configuration Drift:** เมื่อมีคนอัปเดต Config ในเครื่องตัวเองแล้วไม่ได้แจ้งคนอื่น ทำให้โปรเจกต์ของเพื่อนรันไม่ผ่านและเสียเวลา Debug เป็นชั่วโมง

---

## 🏗️ สถาปัตยกรรมระบบ (System Architecture)

```mermaid
graph TD
    subgraph เครื่องนักพัฒนา A (Developer Laptop A)
        DevA[ นักพัฒนา A ] -->|1. envsync push| CLIClientA[ envsync CLI Tool ]
        CLIClientA -->|Encrypts AES-256-GCM| LocalENVA[ Local .env File ]
    end

    subgraph ระบบส่วนกลาง (Central Platform Backend)
        CLIClientA <-->|2. HTTPS Encrypted Payload| BackendAPI[ envsync Central Server\n(Go / Gin REST API) ]
        BackendAPI <-->|3. Save Encrypted Payload & Audit| Database[( Database: PostgreSQL / SQLite )]
        BackendAPI -->|4. Webhook Notification| Webhook[ Slack / Discord Channel ]
    end

    subgraph เครื่องนักพัฒนา B (Developer Laptop B)
        Webhook -->|แจ้งเตือนเมื่อมีอัปเดต| DevB[ นักพัฒนา B ]
        DevB -->|5. envsync pull| CLIClientB[ envsync CLI Tool ]
        CLIClientB <-->|Fetch & Decrypt| BackendAPI
        CLIClientB -->|Updates & Show Diff| LocalENVB[ Local .env File ]
    end
```

---

## 🌟 ฟีเจอร์หลัก (Key Features)

- 🔒 **Zero-Knowledge Encryption (AES-256-GCM):** เข้ารหัสไฟล์ `.env` ในเครื่องผู้ใช้ก่อนส่งขึ้น Server หลังบ้านจะเห็นเฉพาะข้อมูลที่เข้ารหัสแล้วเท่านั้น
- 🚀 **Single Binary CLI:** เขียนด้วยภาษา Go ทำงานรวดเร็ว ไม่มี dependency ภายนอก รองรับ Linux, macOS, และ Windows
- 🎨 **Colorized Terminal Diff:** เปรียบเทียบความต่างระหว่างเครื่องตนเองกับ Server ส่วนกลาง พร้อมแสดงสีสันชัดเจน (เขียว/เหลือง/แดง)
- 📜 **Audit Trail & Rollback Engine:** บันทึกประวัติการเปลี่ยนแปลงย้อนหลัง ใครทำอะไร เมื่อไหร่ และสามารถย้อนกลับเวอร์ชันเดิมได้ทันที (`envsync rollback`)
- 🛡️ **Secret Masking & Detection:** ตรวจจับคีย์สำคัญ (AWS, JWT, DB URLs) ป้องกันการหลุด และซ่อนรหัสผ่านแสดงเป็น `******` เวลา print log
- ⚡ **In-memory Secret Injector (`envsync run`):** รันแอปพลิเคชันโดยฉีด Secret เข้าแรมโดยตรง โดยไม่ต้องสร้างไฟล์ `.env` ลงบน Disk
- ☸️ **Kubernetes Secret Exporter:** แปลง `.env` เป็น `secret.yaml` ของ Kubernetes ได้ด้วยคำสั่งเดียว
- 📢 **Webhooks Integrations:** ส่งการ์ดแจ้งเตือนเข้า Slack และ Discord อัตโนมัติเมื่อมีคนอัปเดต Config

---

## 📥 การติดตั้ง (Installation)

### 1. ติดตั้งผ่านสคริปต์อัตโนมัติ (Linux / macOS)
```bash
curl -fsSL https://raw.githubusercontent.com/jakkayy/envSync/main/scripts/install.sh | sh
```

### 2. Build จาก Source Code (ใช้ Go 1.22+)
```bash
git clone https://github.com/jakkayy/envSync.git
cd envSync
make build

# ไฟล์ binary จะอยู่ที่ bin/envsync และ bin/server
```

---

## 💻 คู่มือการใช้งาน CLI (Usage Guide)

### 1. เริ่มต้นใช้งานในโปรเจกต์ (`envsync init`)
เชื่อมต่อโปรเจกต์ปัจจุบันเข้ากับพื้นที่ส่วนกลาง:
```bash
$ envsync init --project payment-service

✔ Project 'payment-service' (ID: proj_4f14611c) initialized successfully.
✔ Config file .envsync.json created.
```

### 2. อัปเดตค่า Config ขึ้นส่วนกลาง (`envsync push`)
อ่านไฟล์ `.env`, เข้ารหัสด้วย AES-256-GCM แล้วส่งขึ้น Server พร้อมแจ้งเตือนทีม:
```bash
$ envsync push --env dev --message "Add Redis connection timeout"

🔒 Encrypting environment variables (AES-256-GCM)...
⬆️  Pushing 14 variables to project 'payment-service' (environment: dev)...
✅ Successfully updated! (Version: v4)
```

### 3. ดึงค่า Config ล่าสุดมาลงเครื่อง (`envsync pull`)
ดึงค่าล่าสุด ถอดรหัส และอัปเดตลงไฟล์ `.env` ในเครื่องพร้อมแสดงผล Diff:
```bash
$ envsync pull --env dev

⬇️  Fetching latest config for 'payment-service' (environment: dev)...
🔓 Decrypting environment variables...

Comparing local .env with remote (dev):
  + ADDED:    REDIS_TIMEOUT = 5000
  ~ MODIFIED: DB_HOST (localhost -> dev-db.internal.company.com)

Summary: 1 added, 1 modified, 0 removed, 10 unchanged.

✔ Local .env file updated successfully to version v4!
```

### 4. เปรียบเทียบความต่างโดยยังไม่อัปเดต (`envsync diff`)
```bash
$ envsync diff --env dev
```

### 5. ตรวจสอบประวัติการแก้ไข (`envsync history`)
```bash
$ envsync history

📜 Revision History for 'payment-service' (http://localhost:8080):

REV    DATE                 USER             MESSAGE
v4     2026-08-16 23:30     @naeiger         Add Redis connection timeout
v3     2026-08-16 22:15     @alex_dev        Updated DB_HOST
```

### 6. ย้อนกลับไปใช้เวอร์ชันก่อนหน้า (`envsync rollback`)
```bash
$ envsync rollback --env dev --to 3

⏪ Rolling back 'payment-service' (environment: dev) to version v3...
✔ Local .env file successfully rolled back to version v3!
```

### 7. ฉีด Secret เข้าแรมโดยไม่สร้างไฟล์ลงดิสก์ (`envsync run`)
```bash
$ envsync run --env dev -- node server.js
```

### 8. ส่งออกเป็น Kubernetes Secret Manifest (`envsync export k8s`)
```bash
$ envsync export k8s --name my-app-secret -n production -o secret.yaml

✔ Kubernetes Secret manifest successfully exported to 'secret.yaml'
```

---

## 🛠️ การรันระบบหลังบ้านในเครื่อง (Local Dev Stack)

คุณสามารถรัน API Server และ PostgreSQL ฐานข้อมูลส่วนกลางในเครื่องได้ด้วย Docker Compose:

```bash
# รัน API Server + PostgreSQL
docker-compose up -d

# ตรวจสอบสถานะการทำงาน
curl http://localhost:8080/health
```

---

## 🧪 การทดสอบระบบ (Testing & Quality Assurance)

```bash
# รัน Unit Tests และ Integration Tests
make test

# ตรวจสอบคุณภาพโค้ดด้วย Linter
make lint
```

---

## 🌟 นำไปเสนอในสัมภาษณ์งาน (Portfolio Showcase Tips)

หากคุณกำลังสัมภาษณ์งานตำแหน่ง **Platform Engineer / DevOps / DevEx Engineer** โปรเจกต์นี้พิสูจน์ทักษะ:
1. **Internal Tooling Craftsmanship:** สร้างวิศวกรรมเครื่องมือแก้ปัญหาจริงให้ทีมพัฒนา ไม่ใช่แค่ผู้ใช้เครื่องมือคนอื่น
2. **Zero-Trust & Client-side Security:** การเข้ารหัสลับข้อมูลความเสี่ยงสูงแบบ Zero-Knowledge
3. **CI/CD & Automation Mastery:** โครงสร้าง GitHub Actions, Security Scanning (Gosec/Trivy), Docker GHCR และ GoReleaser

---

## 📄 ใบอนุญาต (License)
โปรเจกต์นี้อยู่ภายใต้ใบอนุญาต **MIT License**

# 🗺️ แผนการทำงานแบบทีละขั้นตอน (Master Step-by-Step Execution Plan): `envSync`

**โปรเจกต์:** `envSync` (Platform Engineering & DevSecOps Config Management Tool)  
**เป้าหมาย:** สร้าง CLI Tool และ Backend API สำหรับจัดการไฟล์ `.env` ที่เข้ารหัส ป้องกัน Secret รั่วไหล และซิงก์ Config ระหว่างทีมได้อย่างปลอดภัย  
**รูปแบบการนำเสนอ:** กำหนด **ชื่อขั้นตอน**, **รายละเอียดการทำงาน**, **ไฟล์ที่แตะต้อง** และ **ข้อความ Commit ที่แนะนำ**

---

## 📌 สรุปภาพรวมของขั้นตอนการพัฒนา (Overview)

```
[Phase 0: เตรียมฐานระบบ] ➔ [Phase 1: CLI & เครื่องมือเข้ารหัสในเครื่อง] ➔ [Phase 2: Backend API & DB ส่วนกลาง] ➔ [Phase 3: DevSecOps & แจ้งเตือน] ➔ [Phase 4: Docker & CI/CD Release] ➔ [Phase 5: ฟีเจอร์ขั้นสูง]
```

---

## 🛠️ รายละเอียดขั้นตอนการทำงานทั้งหมด (Step-by-Step Breakdown)

---

### 📁 Phase 0: การจัดเตรียมรากฐานโปรเจกต์ (Project Foundations & Infrastructure Setup)

#### 🔹 ขั้นตอนที่ 1: วางสถาปัตยกรรมระบบและสร้างเอกสารแผนงาน
* **สิ่งที่ต้องทำ:** รวบรวมภาพรวมสถาปัตยกรรมระบบ (Mermaid Diagram), กำหนด Data Flow ระหว่าง CLI และ Central Server
* **ไฟล์ที่สร้าง/แก้ไข:** `docs/env_sync_cli_spec.md`, `docs/master_execution_plan.md`
* **Suggested Commit:** `docs: add comprehensive system architecture & implementation plan`

#### 🔹 ขั้นตอนที่ 2: เริ่มต้น Go Module และกำหนดเวอร์ชัน
* **สิ่งที่ต้องทำ:** รัน `go mod init github.com/jakkayy/envSync` ตั้งค่า Go 1.22+ เป็นเบสสำหรับฟีเจอร์ใหม่
* **ไฟล์ที่สร้าง/แก้ไข:** `go.mod`, `go.sum`
* **Suggested Commit:** `chore: initialize Go module (go.mod) targeting Go 1.22+`

#### 🔹 ขั้นตอนที่ 3: จัดโครงสร้างโฟลเดอร์ตามมาตรฐาน Go Standard Project Layout
* **สิ่งที่ต้องทำ:** สร้างไดเรกทอรี `cmd/envsync`, `cmd/server`, `internal/`, `pkg/`, `deployments/`
* **ไฟล์ที่สร้าง/แก้ไข:** โครงสร้างไดเรกทอรีโปรเจกต์
* **Suggested Commit:** `chore: setup standard Go project layout (cmd, internal, pkg, deployments)`

#### 🔹 ขั้นตอนที่ 4: ตั้งค่า Linters, Git Ignore และ Makefile สำหรับช่วยรันคำสั่ง
* **สิ่งที่ต้องทำ:** สร้าง `.gitignore`, ตั้งค่า `.golangci.yml` สำหรับตรวจคุณภาพโค้ด และเขียน `Makefile` เพื่อช่วยรัน `make build`, `make test`, `make lint`
* **ไฟล์ที่สร้าง/แก้ไข:** `.gitignore`, `.golangci.yml`, `Makefile`
* **Suggested Commit:** `chore: add .gitignore, .golangci.yml, and local development Makefile`

---

### 📁 Phase 1: การพัฒนา CLI Engine และตัวแปลงไฟล์ในเครื่อง (Core CLI & Local Parser)

#### 🔹 ขั้นตอนที่ 5: สร้างจุดเริ่มต้นของแอปพลิเคชัน CLI (Cobra Framework Setup)
* **สิ่งที่ต้องทำ:** ติดตั้ง `cobra-cli` และสร้าง Entrypoint ใน `cmd/envsync/main.go`
* **ไฟล์ที่สร้าง/แก้ไข:** `cmd/envsync/main.go`, `internal/cli/root.go`
* **Suggested Commit:** `feat(cli): initialize Cobra CLI application entrypoint in cmd/envsync`

#### 🔹 ขั้นตอนที่ 6: พัฒนาระบบ Flags และเวอร์ชันแอปพลิเคชัน (Global Flags Handler)
* **สิ่งที่ต้องทำ:** เพิ่ม flag `--version`, `--verbose`, และ `--config` ใน Root Command
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/cli/version.go`, `internal/cli/root.go`
* **Suggested Commit:** `feat(cli): add version and global flags handlers`

#### 🔹 ขั้นตอนที่ 7: ออกแบบโครงสร้างไฟล์ตั้งค่าโปรเจกต์ (`.envsync.json`)
* **สิ่งที่ต้องทำ:** สร้าง Data Model และ Struct ในการอ่าน/เขียนไฟล์ตั้งค่าประจำโฟลเดอร์โปรเจกต์ (เช่น project_id, default_env, server_url)
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/config/config.go`
* **Suggested Commit:** `feat(config): define local workspace config schema (.envsync.json)`

#### 🔹 ขั้นตอนที่ 8: พัฒนาคำสั่งเริ่มต้นโปรเจกต์ (`envsync init`)
* **สิ่งที่ต้องทำ:** พัฒนาคำสั่ง `envsync init --project <name>` เพื่อสร้างไฟล์ `.envsync.json` อัตโนมัติพร้อมตรวจสอบว่าเคย init หรือยัง
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/cli/init.go`
* **Suggested Commit:** `feat(config): implement envsync init command to create project config`

#### 🔹 ขั้นตอนที่ 9: พัฒนาตัวอ่านและแยกโครงสร้างไฟล์ `.env` (Dotenv Lexer & AST Parser)
* **สิ่งที่ต้องทำ:** เขียน Parser อ่านไฟล์ `.env` รองรับ Key-Value, การครอบด้วย Quotes (`""`, `''`), ค่าหลายบรรทัด (Multiline) และข้อความอธิบาย (Comments)
* **ไฟล์ที่สร้าง/แก้ไข:** `pkg/env/parser.go`, `pkg/env/model.go`
* **Suggested Commit:** `feat(env): implement lexer & AST parser for .env key-value file format`

#### 🔹 ขั้นตอนที่ 10: พัฒนาระบบจัดรูปแบบและบันทึกไฟล์ `.env` (Dotenv Formatter)
* **สิ่งที่ต้องทำ:** เขียนตัวเขียนไฟล์ลง Disk โดยคงบรรทัดว่างและ Comment เดิมไว้ไม่ให้หายไปเมื่อเกิดการอัปเดต
* **ไฟล์ที่สร้าง/แก้ไข:** `pkg/env/formatter.go`
* **Suggested Commit:** `feat(env): implement dotenv formatter preserving comments and empty lines`

#### 🔹 ขั้นตอนที่ 11: พัฒนาระบบเปรียบเทียบความแตกต่างของไฟล์ (Env Diff Engine)
* **สิ่งที่ต้องทำ:** สร้างอัลกอริทึมเปรียบเทียบ Key-Value ระหว่างไฟล์ในเครื่องกับไฟล์ปลายทาง จัดกลุ่มเป็น: Added (+), Modified (~), Deleted (-), Unchanged (=)
* **ไฟล์ที่สร้าง/แก้ไข:** `pkg/env/diff.go`
* **Suggested Commit:** `feat(diff): implement AST diff engine to compare local vs target env`

#### 🔹 ขั้นตอนที่ 12: พัฒนาหน้าจอแสดงผล Diff บน Terminal (Terminal UI Colors & Table)
* **สิ่งที่ต้องทำ:** แสดงผล Diff บน Terminal ให้สวยงามด้วยสีสัน (เขียว = เพิ่ม, เหลือง = แก้ไข, แดง = ลบ) ด้วย `lipgloss` หรือ `pterm`
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/ui/diff_printer.go`, `internal/cli/diff.go`
* **Suggested Commit:** `feat(diff): add colorized side-by-side Terminal UI table for env diffs`

#### 🔹 ขั้นตอนที่ 13: พัฒนาโมดูลเข้ารหัสข้อมูล (AES-256-GCM Cryptography Engine)
* **สิ่งที่ต้องทำ:** เขียนโมดูลเข้ารหัสและถอดรหัสข้อความด้วยอัลกอริทึม AES-256-GCM พร้อมสร้าง Random Nonce ป้องกัน Replay Attack
* **ไฟล์ที่สร้าง/แก้ไข:** `pkg/crypto/aes.go`
* **Suggested Commit:** `feat(crypto): implement AES-256-GCM symmetric cipher module`

#### 🔹 ขั้นตอนที่ 14: พัฒนาระบบแปลง Passphrase เป็น Encryption Key (Key Derivation Function - KDF)
* **สิ่งที่ต้องทำ:** ใช้ PBKDF2 / Argon2id ในการแปลงรหัสผ่านของผู้ใช้ให้กลายเป็น Cryptographic Key ขนาด 32 bytes ร่วมกับ Salt
* **ไฟล์ที่สร้าง/แก้ไข:** `pkg/crypto/kdf.go`
* **Suggested Commit:** `feat(crypto): implement key derivation function (PBKDF2/Argon2id)`

#### 🔹 ขั้นตอนที่ 15: เขียน Unit Test สำหรับระบบอ่านและแปลงไฟล์ `.env`
* **สิ่งที่ต้องทำ:** เขียนกรณีทดสอบ (Test Cases) สำหรับ `.env` รูปแบบแปลกๆ เช่น multiline value, special characters, empty lines
* **ไฟล์ที่สร้าง/แก้ไข:** `pkg/env/parser_test.go`
* **Suggested Commit:** `test(env): add unit tests for .env parser and formatter edge cases`

#### 🔹 ขั้นตอนที่ 16: เขียน Unit Test สำหรับระบบเข้ารหัสและถอดรหัส
* **สิ่งที่ต้องทำ:** ทดสอบว่าเมื่อเข้ารหัสแล้ว ถอดรหัสกลับมาจะได้ข้อความเดิม และทดสอบกรณีใส่ Key ผิดว่าจะต้องติด Error เสมอ
* **ไฟล์ที่สร้าง/แก้ไข:** `pkg/crypto/aes_test.go`
* **Suggested Commit:** `test(crypto): add unit tests for encryption, decryption, and KDF`

---

### 📁 Phase 2: การพัฒนา Backend API Server และฐานข้อมูลส่วนกลาง (Central API Backend & Storage)

#### 🔹 ขั้นตอนที่ 17: เริ่มต้นสร้าง Web API Server (Gin Framework Setup)
* **สิ่งที่ต้องทำ:** สร้างโปรแกรม API Server ใน `cmd/server/main.go` รองรับ HTTP REST API
* **ไฟล์ที่สร้าง/แก้ไข:** `cmd/server/main.go`, `internal/server/server.go`
* **Suggested Commit:** `feat(server): initialize Gin REST API server structure in cmd/server`

#### 🔹 ขั้นตอนที่ 18: พัฒนาระบบ Log และ Request Tracking (Structured Logger & Middleware)
* **สิ่งที่ต้องทำ:** เขียน Middleware สำหรับสร้าง Correlation ID (X-Request-ID) และบันทึก Log แบบ JSON ในทุก HTTP Request
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/server/middleware/logging.go`
* **Suggested Commit:** `feat(server): add structured logging middleware (zap/slog) and correlation IDs`

#### 🔹 ขั้นตอนที่ 19: พัฒนาระบบเชื่อมต่อฐานข้อมูล (Database Connection Pool)
* **สิ่งที่ต้องทำ:** เขียนตัวเชื่อมต่อ DB ที่รองรับทั้ง SQLite (สำหรับทดสอบในเครื่อง) และ PostgreSQL (สำหรับใช้งานจริง)
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/database/db.go`
* **Suggested Commit:** `feat(db): set up PostgreSQL/SQLite database connection pool`

#### 🔹 ขั้นตอนที่ 20: ออกแบบโครงสร้างตารางฐานข้อมูล (GORM Data Models)
* **สิ่งที่ต้องทำ:** นิยาม Data Structure ของ `Project`, `Environment`, `EnvVersion` (เก็บ Encrypted Payload) และ `AuditLog`
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/database/models.go`
* **Suggested Commit:** `feat(db): define GORM models (User, Project, Environment, EnvVersion, AuditLog)`

#### 🔹 ขั้นตอนที่ 21: พัฒนาระบบ Auto Migration สำหรับฐานข้อมูล
* **สิ่งที่ต้องทำ:** เพิ่มฟังก์ชันรัน Migration อัตโนมัติเมื่อ Server เริ่มทำงาน
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/database/migrate.go`
* **Suggested Commit:** `feat(db): add database migration scripts and schema auto-migrate`

#### 🔹 ขั้นตอนที่ 22: พัฒนาระบบยืนยันตัวตนผ่าน API Token (Authentication Middleware)
* **สิ่งที่ต้องทำ:** เขียน Middleware ตรวจสอบ `Authorization: Bearer <token>` ก่อนอนุญาตให้เข้าถึง API
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/server/middleware/auth.go`
* **Suggested Commit:** `feat(auth): implement JWT / API Key authentication middleware`

#### 🔹 ขั้นตอนที่ 23: พัฒนา API สร้างโปรเจกต์ใหม่ (`POST /api/v1/projects`)
* **สิ่งที่ต้องทำ:** สร้าง Endpoint สำหรับลงทะเบียนโปรเจกต์และสร้าง Secret Key ประจำโปรเจกต์
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/server/handlers/project.go`
* **Suggested Commit:** `feat(api): implement POST /api/v1/projects endpoint`

#### 🔹 ขั้นตอนที่ 24: พัฒนา API สำหรับการ Push ข้อมูลขึ้น Server (`POST /api/v1/sync/push`)
* **สิ่งที่ต้องทำ:** รับ Encrypted Payload จาก CLI ตรวจสอบเวอร์ชันล่าสุด บันทึกสร้างเป็น Version ใหม่ในฐานข้อมูล
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/server/handlers/push.go`
* **Suggested Commit:** `feat(api): implement POST /api/v1/sync/push payload upload endpoint`

#### 🔹 ขั้นตอนที่ 25: พัฒนา API สำหรับการ Pull ข้อมูลจาก Server (`GET /api/v1/sync/pull`)
* **สิ่งที่ต้องทำ:** ส่งคืน Encrypted Payload เวอร์ชันล่าสุดของ environment นั้นๆ กลับไปให้ CLI
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/server/handlers/pull.go`
* **Suggested Commit:** `feat(api): implement GET /api/v1/sync/pull payload retrieval endpoint`

#### 🔹 ขั้นตอนที่ 26: เชื่อมต่อคำสั่ง `envsync push` ใน CLI เข้ากับ API Server
* **สิ่งที่ต้องทำ:** ให้คำสั่ง `envsync push` อ่านไฟล์ `.env`, เข้ารหัสด้วย AES-256 แล้วยิง HTTP POST ไปที่ API Server
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/cli/push.go`, `pkg/client/api_client.go`
* **Suggested Commit:** `feat(cli): integrate envsync push command with REST API server`

#### 🔹 ขั้นตอนที่ 27: เชื่อมต่อคำสั่ง `envsync pull` ใน CLI เข้ากับ API Server
* **สิ่งที่ต้องทำ:** ให้คำสั่ง `envsync pull` ดึง Encrypted Payload จาก Server, ถอดรหัส และอัปเดตลงไฟล์ `.env` ในเครื่อง
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/cli/pull.go`
* **Suggested Commit:** `feat(cli): integrate envsync pull command with local file updater`

---

### 📁 Phase 3: ระบบความปลอดภัย ประวัติการใช้งาน และการแจ้งเตือน (DevSecOps & Integrations)

#### 🔹 ขั้นตอนที่ 28: พัฒนาระบบบันทึกประวัติการใช้งานแบบ Async (Audit Trail Logger)
* **สิ่งที่ต้องทำ:** บันทึกทุกคำสั่ง push/pull ลงตาราง `audit_logs` (ใคร, ทำอะไร, เมื่อไหร่, สภาพแวดล้อมไหน) แบบ Background Worker
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/audit/logger.go`
* **Suggested Commit:** `feat(audit): implement async audit trail logger for server actions`

#### 🔹 ขั้นตอนที่ 29: พัฒนา API ดึงประวัติการแก้ไข (`GET /api/v1/projects/:id/history`)
* **สิ่งที่ต้องทำ:** สร้าง Endpoint คืนค่ารายการเวอร์ชันย้อนหลังและ Audit Logs ของโปรเจกต์
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/server/handlers/history.go`
* **Suggested Commit:** `feat(audit): implement GET /api/v1/projects/:id/history audit API endpoint`

#### 🔹 ขั้นตอนที่ 30: พัฒนาคำสั่งดูประวัติการแก้ไข (`envsync history`)
* **สิ่งที่ต้องทำ:** สร้างคำสั่ง CLI ดึงประวัติจาก Server มาแสดงเป็นตารางเรียงตามเวลาบน Terminal
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/cli/history.go`
* **Suggested Commit:** `feat(cli): implement envsync history command with timeline output`

#### 🔹 ขั้นตอนที่ 31: พัฒนาคำสั่งย้อนคืนเวอร์ชัน (`envsync rollback`)
* **สิ่งที่ต้องทำ:** คำสั่งย้อน Config กลับไปใช้เวอร์ชันก่อนหน้า (เช่น `envsync rollback --to v3`)
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/cli/rollback.go`
* **Suggested Commit:** `feat(cli): implement envsync rollback command to restore previous version`

#### 🔹 ขั้นตอนที่ 32: พัฒนาระบบส่งแจ้งเตือนเข้า Slack (Slack Incoming Webhook Integration)
* **สิ่งที่ต้องทำ:** ส่ง Slack Rich Card Notification เข้าช่อง `#dev-team` เมื่อมีการ push config ใหม่
* **ไฟล์ที่สร้าง/แก้ไข:** `pkg/notifier/slack.go`
* **Suggested Commit:** `feat(integrations): implement Slack incoming webhook notifier for config changes`

#### 🔹 ขั้นตอนที่ 33: พัฒนาระบบส่งแจ้งเตือนเข้า Discord (Discord Webhook Integration)
* **สิ่งที่ต้องทำ:** ส่ง Discord Embed Notification เมื่อมีการอัปเดต Config
* **ไฟล์ที่สร้าง/แก้ไข:** `pkg/notifier/discord.go`
* **Suggested Commit:** `feat(integrations): implement Discord webhook notification payload`

#### 🔹 ขั้นตอนที่ 34: พัฒนาระบบตรวจจับ Secret รั่วไหล (Secret Detector Regex Engine)
* **สิ่งที่ต้องทำ:** สแกนค่าใน `.env` เพื่อหา Secret สำคัญ (เช่น AWS Access Key, Private Key, JWT) แล้วแจ้งเตือนคำเตือน
* **ไฟล์ที่สร้าง/แก้ไข:** `pkg/security/detector.go`
* **Suggested Commit:** `feat(security): implement secret detection regex engine (AWS keys, JWTs, DB passwords)`

#### 🔹 ขั้นตอนที่ 35: พัฒนาระบบซ่อนรหัสผ่านใน Terminal Log (Secret Masking Middleware)
* **สิ่งที่ต้องทำ:** แปลงรหัสผ่านหรือ API Key ให้กลายเป็น `******` เมื่อพิมพ์ออกหน้าจอ Terminal
* **ไฟล์ที่สร้าง/แก้ไข:** `pkg/security/masker.go`
* **Suggested Commit:** `feat(security): implement terminal output secret masking middleware`

#### 🔹 ขั้นตอนที่ 36: เขียน Integration Test สำหรับทดสอบการทำงานร่วมกันระหว่าง CLI และ Server
* **สิ่งที่ต้องทำ:** ทดสอบกระบวนการตั้งแต่ init -> push -> pull -> history -> rollback แบบอัตโนมัติ
* **ไฟล์ที่สร้าง/แก้ไข:** `tests/integration/sync_test.go`
* **Suggested Commit:** `test(server): add API integration test suite for push/pull/audit flows`

---

### 📁 Phase 4: โครงสร้างพื้นฐาน, Docker และ CI/CD Pipeline (Infrastructure & Release Automation)

#### 🔹 ขั้นตอนที่ 37: จัดทำ Dockerfile สำหรับ API Server (Multi-stage Build)
* **สิ่งที่ต้องทำ:** เขียน Dockerfile สองขั้นตอน (Build stage + Minimal Run stage บน Alpine/Distroless) เพื่อให้ได้ Container Image ขนาดเล็กปลอดภัย
* **ไฟล์ที่สร้าง/แก้ไข:** `deployments/docker/Dockerfile`
* **Suggested Commit:** `chore(docker): create multi-stage Dockerfile for server (Alpine/Distroless)`

#### 🔹 ขั้นตอนที่ 38: จัดทำ Docker Compose สำหรับสภาพแวดล้อมการพัฒนาในเครื่อง
* **สิ่งที่ต้องทำ:** เขียน `docker-compose.yml` เพื่อรัน API Server + PostgreSQL DB ด้วยคำสั่งเดียว
* **ไฟล์ที่สร้าง/แก้ไข:** `docker-compose.yml`
* **Suggested Commit:** `chore(docker): create Docker Compose setup for local dev stack (API + Postgres)`

#### 🔹 ขั้นตอนที่ 39: ตั้งค่า GitHub Actions สำหรับตรวจโค้ดและรัน Unit Test อัตโนมัติ
* **สิ่งที่ต้องทำ:** สร้าง Workflow รัน `golangci-lint` และ `go test -v ./...` ทุกครั้งที่มีการเปิด Pull Request หรือ Push
* **ไฟล์ที่สร้าง/แก้ไข:** `.github/workflows/ci.yml`
* **Suggested Commit:** `ci: add GitHub Actions workflow for linting (golangci-lint) and unit tests`

#### 🔹 ขั้นตอนที่ 40: ตั้งค่า GitHub Actions สำหรับสแกนความปลอดภัย (Gosec & Trivy Scanner)
* **สิ่งที่ต้องทำ:** สแกนหาช่องโหว่ความปลอดภัยใน Go Source Code และ Docker Image อัตโนมัติใน CI
* **ไฟล์ที่สร้าง/แก้ไข:** `.github/workflows/security.yml`
* **Suggested Commit:** `ci: add GitHub Actions step for security scanning (Gosec & Trivy vulnerability scanner)`

#### 🔹 ขั้นตอนที่ 41: ตั้งค่า GitHub Actions สำหรับสร้างและส่ง Docker Image เข้า GHCR
* **สิ่งที่ต้องทำ:** Build Docker Image และ Push ขึ้น GitHub Container Registry อัตโนมัติเมื่อ Merge เข้า branch `main`
* **ไฟล์ที่สร้าง/แก้ไข:** `.github/workflows/docker-publish.yml`
* **Suggested Commit:** `ci: add GitHub Actions workflow to build and push Docker images to GHCR`

#### 🔹 ขั้นตอนที่ 42: ตั้งค่า GoReleaser สำหรับสร้างไฟล์ Executable (Cross-Compilation Config)
* **สิ่งที่ต้องทำ:** ตั้งค่าไฟล์ `.goreleaser.yaml` สำหรับ Build Binary สำหรับ Linux, macOS, และ Windows
* **ไฟล์ที่สร้าง/แก้ไข:** `.goreleaser.yaml`
* **Suggested Commit:** `ci: add GoReleaser configuration (.goreleaser.yaml) for cross-compilation`

#### 🔹 ขั้นตอนที่ 43: ตั้งค่า GitHub Actions สำหรับออก Release อัตโนมัติ (Automated Release Workflow)
* **สิ่งที่ต้องทำ:** ตั้งค่า Trigger ให้ GoReleaser รันเมื่อมีการสร้าง Git Tag (เช่น `git tag v1.0.0`)
* **ไฟล์ที่สร้าง/แก้ไข:** `.github/workflows/release.yml`
* **Suggested Commit:** `ci: add GitHub Actions workflow to trigger GoReleaser on git tags`

#### 🔹 ขั้นตอนที่ 44: เขียนเอกสารการใช้งานและสคริปต์ติดตั้งสำหรับผู้ใช้ (README & Install Script)
* **สิ่งที่ต้องทำ:** จัดทำเอกสารคู่มือการใช้งาน README.md และสคริปต์ `install.sh` สำหรับดาวน์โหลด Binary มาลงเครื่องผู้ใช้
* **ไฟล์ที่สร้าง/แก้ไข:** `README.md`, `scripts/install.sh`
* **Suggested Commit:** `docs: add installation script and comprehensive README documentation`

---

### 📁 Phase 5: ฟีเจอร์ระดับสูงสำหรับสาย Platform Engineer (Advanced Boosters)

#### 🔹 ขั้นตอนที่ 45: พัฒนาระบบส่งออกเป็น Kubernetes Secret Manifest (`envsync export k8s`)
* **สิ่งที่ต้องทำ:** เพิ่มคำสั่งสร้างไฟล์ `secret.yaml` สำหรับเอาไปปรับใช้ใน Kubernetes Cluster
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/cli/export.go`, `pkg/exporter/k8s.go`
* **Suggested Commit:** `feat(k8s): implement Kubernetes Secret manifest generator (K8s Exporter)`

#### 🔹 ขั้นตอนที่ 46: พัฒนาคำสั่งฉีด Secret เข้า Memory โดยไม่สร้างไฟล์ลง Disk (`envsync run`)
* **สิ่งที่ต้องทำ:** คำสั่ง `envsync run --env prod -- node server.js` เพื่ออ่าน Secret ดึงเข้า Process Environment โดยไม่สร้างไฟล์ `.env` บนดิสก์
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/cli/run.go`
* **Suggested Commit:** `feat(cli): implement envsync run subcommand for in-memory secret injection`

#### 🔹 ขั้นตอนที่ 47: พัฒนาระบบคำแนะนำคำสั่งอัตโนมัติใน Terminal (Shell Autocompletion Generator)
* **สิ่งที่ต้องทำ:** เพิ่มคำสั่ง `envsync completion bash/zsh/fish` เพื่อสร้างสคริปต์ Auto-complete เวลาผู้ใช้กด Tab
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/cli/completion.go`
* **Suggested Commit:** `feat(cli): add subcommand for shell auto-completion generation (bash, zsh, fish)`

#### 🔹 ขั้นตอนที่ 48: พัฒนาระบบ Metrics สำหรับการติดตามผลการทำงาน (Prometheus Endpoint)
* **สิ่งที่ต้องทำ:** เพิ่ม Endpoint `/metrics` ใน API Server สำหรับส่งข้อมูลทราฟฟิกและประสิทธิภาพให้ Prometheus
* **ไฟล์ที่สร้าง/แก้ไข:** `internal/server/middleware/metrics.go`
* **Suggested Commit:** `feat(metrics): add Prometheus /metrics endpoint to backend server`

#### 🔹 ขั้นตอนที่ 49: จัดทำเอกสาร ADR และ OpenAPI / Swagger Specification
* **สิ่งที่ต้องทำ:** เขียนเอกสารการตัดสินใจเลือกสถาปัตยกรรม (Architecture Decision Records) และ Swagger API Docs
* **ไฟล์ที่สร้าง/แก้ไข:** `docs/adr/0001-use-aes-256-gcm.md`, `docs/swagger.yaml`
* **Suggested Commit:** `docs: add architectural decision records (ADRs) and OpenAPI/Swagger spec`

---

## 🎯 เงื่อนไขความเสร็จสิ้นของงาน (Definition of Done - DoD)
1. ทุกขั้นตอนมี **ชื่อขั้นตอนที่ชัดเจน** พร้อมรายละเอียดงานและผลลัพธ์
2. **การทำ Git Commit** ในแต่ละขั้นตอนจะใช้ข้อความที่เป็นมาตรฐาน (Suggested Commit)
3. โค้ดทั้งหมดผ่านการตรวจ Linting และสแกนความปลอดภัย 100%

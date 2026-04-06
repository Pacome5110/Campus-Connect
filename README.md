# CampusConnect

**Öğrenci:** Pacome berınyuy fondzenyuy
**Okul No:** 24080410151

Üniversite etkinlik platformu – NestJS + Go polyglot backend. Node REST/GraphQL ile veri yönetimini yaparken, performans odaklı loglama, analitik ve metrik süreçleri eşzamanlı olarak Go tarafından webhooklar aracılığıyla dağıtık bir yapıyla çözülmektedir.

## Teknolojiler
- NestJS (TypeScript) — REST + GraphQL (Main Backend - :3000)
- Go (Gin) — Notification + Analytics (Auxiliary Service - :8080)
- PostgreSQL — Veritabanı
- Docker — Containerization

## Kurulum

### Gereksinimler
- Node.js 18+, Go 1.21+, PostgreSQL 15+

### NestJS Service
```bash
cd nestjs-service
npm install
cp .env.example .env
npx prisma migrate dev --name init
npm run start:dev
```

### Go Service
```bash
cd go-service
go mod download
cp .env.example .env
go run main.go
```

## Environment Variables
- `DATABASE_URL`: Ortak Postgres bağlantı URL'si (Prisma).
- `JWT_SECRET`: Bearer token imzalama şifresi (AuthService).
- `API_KEY`: NestJS'nin Go servis uç noktalarına erişmek için kullandığı sır.
- `GO_WEBHOOK_URL`: Go hook API'si, asenkron veriler için NestJS kullanır.

## API Endpoints

| Servis  | Metot | Path               | Açıklama                     |
|---------|-------|--------------------|------------------------------|
| NestJS  | POST  | /api/v1/auth/req   | Kullanıcı Kayıt işlemi       |
| NestJS  | POST  | /api/v1/auth/login | Token alımı (JWT)            |
| NestJS  | GET   | /api/v1/users      | Tüm kullanıcıları listele    |
| NestJS  | GET   | /api/v1/events     | Tüm etkinlikler              |
| NestJS  | POST  | /api/v1/events     | Etkinlik Oluştur (ADMIN)     |
| NestJS  | GraphQL| /graphql          | GraphQL endpoint             |
| Go      | POST  | /webhook           | NestJS event dinleyicisi     |
| Go      | GET   | /api/analytics     | Metrikleri getir (API_KEY)   |
| Go      | POST  | /api/notifications | Bildirim yolla (API_KEY)     |

## Örnek Request / Response

1. **Register**
```bash
curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"name":"A","email":"a@m.com","password":"123"}'
```
*Yanıt:* `{"id": "uuid", "name": "A", "email": "a@m.com", "role": "USER"}`

2. **Login**
```bash
curl -X POST http://localhost:3000/api/v1/auth/login -H "Content-Type: application/json" -d '{"email":"a@m.com","password":"123"}'
```
*Yanıt:* `{"access_token": "eyJ..."}`

3. **Get Events**
```bash
curl -X GET http://localhost:3000/api/v1/events -H "Authorization: Bearer <TOKEN>"
```
*Yanıt:* `[{"id": "...", "title": "Buluşma"}]`

4. **Go Analytics**
```bash
curl -X GET http://localhost:8080/api/analytics -H "X-API-Key: my-secret-key"
```
*Yanıt:* `{"active_users": 530, "event_count": 142}`

5. **Go Rate Limit Test (Hata Senaryosu)**
```bash
for i in {1..10}; do curl -X GET http://localhost:8080/health; done
```
*Yanıt 6. istekten itibaren:* `{"error": "Too Many Requests. Rate limit exceeded."}`

## Mimari Kararlar
- **Neden NestJS + Go?** NestJS esnek, iyi organize edilmiş modüler yapısıyla iş kuralları (CRUD) için harikadır; Go is hafif thread (Goroutine) yönetimleri sayesinde arka plandaki işlemler (Notification / Analytics write) için çok güçlü ve az bellek yakar.
- **Webhook Tercihi:** İki servis arasındaki güçlü bağımsızlığı sağlamak ve mesajların HTTP protokolü üzerinden asenkron aktarımını kolaylaştırmak için Event-Driven Architecture kurgulanmıştır.
- **Rate Limiting (Token Bucket / Mutex Counter):** Sistemin abuse edilmesini veya DDoS vektörlerini azaltmak ve API Gateway benzeri davranışta bulunmak için bellekte basit concurrency-safe kilitli (sync.Mutex) sayaç yapısı kullanıldı.

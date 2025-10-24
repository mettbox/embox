# EmBox

> A media management app using luckycloud.de

Concept:

- Upload media files via the app. Thumbnails are stored in the backend, originals in LuckyCloud
- Media can be organised into albums
- Users can be invited and managed by admins
- Admins can see everything, other users only the media and albums that are public
- All media is automatically organised into folders by date in LuckyCloud

## Getting started

### 1. Database

```sh
docker compose up -d
```

### 2. API

```sh
cd api
```

Create file for the environment variables:

```sh
cp .env.example .env
```

Development with hot reload using air:

```sh
go install github.com/air-verse/air@latest
```

```sh
air
```

w/o hot reload use:

```sh
go run cmd/api/main.go
```

### 3. APP

In a second terminal:

```sh
cd app
nvm use
npm run dev
```

## Resources

- [LuckyCloud REST API](https://docs.luckycloud.de/en/cloud-storage/rest-api)
- [LuckyCloud API Documentation](https://storage.luckycloud.de/published/api-dokumentation/home.md)
- [uberspace Web Backend](https://manual.uberspace.de/web-backends/)

## Deployment: uberspace

- CPU: x86_64 (Intel/AMD)
- OS: CentOS 7 (glibc 2.17)
- go1.25.1

Setup backend:

- App: `/home/mettbox/public/em.mettbox.de` => https://em.mettbox.de
- Api: `/home/mettbox/em` => https://em.mettbox.de/api


```sh
uberspace web backend set em.mettbox.de/api --http --port 2705 --remove-prefix
```

Run Server:

```sh
cd /home/mettbox/em
./main
```

Status:

```sh
uberspace web backend list
```

Setup Database

@TODO

### Build API

In uberspace terminal:

```sh
go build -o /home/mettbox/em/main ./cmd/api/main.go
chmod +x /home/mettbox/em/main
```

Run binary deteched:

```sh
nohup ./main > /dev/null 2>&1 &
```

Stop server:

```sh
ps -aux
kill <PID>
```

### Build App

```sh
npm run build
````

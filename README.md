# go-keycloack
Repository to test Authorization and Authentication inside Keycloack with Golang 👩‍💻

# Steps:

- Start Keycloack 
```
docker run -p 8080:8080 -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin quay.io/keycloak/keycloak:14.0.0
```

- Install libs from golang 
```
go get lib_name
```

- Run golang file 
```
go run main.go
```

# Don't forget to change credentials from .env file

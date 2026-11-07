# Integración Android ForgeFit

## URL del backend

El emulador Android usa por defecto `http://10.0.2.2:8080/api/v1/`. Para usar otro host, por ejemplo un teléfono conectado a la misma red, compila con:

```powershell
.\gradlew.bat :app:assembleDebug -PforgefitApiBaseUrl=http://192.168.1.20:8080/api/v1/
```

Sustituye la IP por la de la computadora donde corre el servidor. La API escucha en el puerto configurado por `PORT` (8080 por defecto). Para iniciar Go desde `Servidor/`:

```powershell
go run ./cmd/api
```

Configura MySQL en `Servidor/.env`; importa `database.sql` antes de iniciar.

## Contrato GPS para el mapa

Todas las rutas de carrera requieren `Authorization: Bearer <token>`:

| Método y ruta | Uso |
| --- | --- |
| `POST /runs/start` | Inicia una sesión. Recibe `initial_lat` y `initial_lng` opcionales. |
| `GET /runs/active` | Devuelve la sesión activa y sus `route_points`, ordenados por `sequence_order`. |
| `POST /runs/{id}/gps` | Añade `{ "latitude", "longitude", "speed_kmh", "altitude_m" }`. |
| `POST /runs/{id}/pause` | Pausa la sesión. |
| `POST /runs/{id}/resume` | Reanuda la sesión. |
| `POST /runs/{id}/finish` | Finaliza y guarda las métricas. |
| `GET /runs/history` | Devuelve las carreras del usuario autenticado. |

Cada punto incluye `latitude`, `longitude`, `sequence_order` y `timestamp`; Android puede convertir esa lista en la polilínea de Google Maps. Las carreras y sus puntos quedan asociados al usuario autenticado. La llave de Google Maps es de cliente Android y no se envía al backend.

## Capas Android

Las pantallas envían eventos a ViewModels; los casos de uso llaman interfaces del dominio; las implementaciones `data` usan el cliente OkHttp compartido en `core/network`. La sesión agrega el token Bearer y elimina credenciales locales cuando la API responde 401.

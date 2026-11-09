# ForgeFit Backend API

Servidor Backend desarrollado en **Go (Golang)** aplicando **Arquitectura Hexagonal (Ports & Adapters)** para la aplicación móvil **ForgeFit - Tu entrenador personal digital**.

---

## Arquitectura Hexagonal

El proyecto está diseñado bajo los principios de la Arquitectura Hexagonal para garantizar el desacoplamiento estricto entre el núcleo del negocio (dominio y casos de uso) y las tecnologías externas (HTTP Gin, MySQL y JWT).

```
Servidor/
├── cmd/
│   └── api/
│       └── main.go                         # Punto de entrada y cableado de dependencias
├── internal/
│   ├── core/                               # NÚCLEO DEL DOMINIO
│   │   ├── domain/                         # Entidades de Dominio, DTOs y Errores
│   │   │   ├── user.go                     # Entidad User y UserProfile
│   │   │   ├── activity.go                 # Entidades DailyActivity y WeeklyActivityDay
│   │   │   ├── run.go                      # Entidades Run y RunRoutePoint (GPS)
│   │   │   ├── routine.go                  # Entidades Routine y RoutineExercise
│   │   │   ├── nutrition.go                # Entidades DailyNutrition, HydrationLog, MealLog
│   │   │   └── dashboard.go                # DTO DashboardSummary
│   │   ├── ports/                          # Puertos (Interfaces de Entrada y Salida)
│   │   │   ├── user_repository.go / user_service.go
│   │   │   ├── user_profile_service.go
│   │   │   ├── activity_service.go
│   │   │   ├── run_service.go
│   │   │   ├── routine_service.go
│   │   │   ├── nutrition_service.go
│   │   │   └── dashboard_service.go
│   │   └── services/                       # Casos de Uso (Lógica de Negocio)
│   │       ├── user_service.go
│   │       ├── user_profile_service.go
│   │       ├── activity_service.go
│   │       ├── run_service.go
│   │       ├── routine_service.go
│   │       ├── nutrition_service.go
│   │       └── dashboard_service.go
│   └── adapters/                           # ADAPTADORES
│       ├── handlers/                       # Adaptadores HTTP (Gin Controllers)
│       │   ├── user_handler.go
│       │   ├── profile_handler.go
│       │   ├── activity_handler.go
│       │   ├── run_handler.go
│       │   ├── routine_handler.go
│       │   ├── nutrition_handler.go
│       │   ├── dashboard_handler.go
│       │   └── middlewares/                # Middlewares HTTP (JWT Auth, CORS)
│       │       ├── auth_middleware.go
│       │       └── cors_middleware.go
│       ├── repository/                     # Adaptadores de Salida (MySQL Repositories)
│       │   ├── mysql_user_repository.go
│       │   ├── mysql_user_profile_repository.go
│       │   ├── mysql_activity_repository.go
│       │   ├── mysql_run_repository.go
│       │   ├── mysql_routine_repository.go
│       │   └── mysql_nutrition_repository.go
│       └── security/                       # Adaptadores de Seguridad (JWT, Bcrypt, Google)
│           ├── bcrypt_adapter.go
│           ├── jwt_adapter.go
│           └── google_auth_adapter.go
├── pkg/                                    # Paquetes Compartidos
│   ├── database/                           # Pool de conexiones MySQL
│   │   └── mysql.go
│   └── config/                             # Carga de variables de entorno
│       └── config.go
├── database.sql                            # Script DDL completo de la base de datos MySQL
├── .env.example                            # Plantilla de variables de entorno
├── go.mod / go.sum
└── README.md
```

---

## Módulos y Funcionalidades

1. **Autenticación y Usuarios**: Registro local con validación de contraseña y Login con JWT.
2. **Perfil de Usuario**: Peso, altura, edad, género, nivel de actividad, total de entrenamientos, racha de días activos y kilómetros acumulados.
3. **Dashboard Principal**: Endpoint unificado para la pantalla de inicio (saludo personalizado, resumen de pasos, calorías quemadas/consumidas, última carrera, próxima rutina y gráfica semanal).
4. **Pasos y Actividad Diaria**: Conteo de pasos, objetivo diario (ej. 10,000 pasos), distancia, calorías quemadas y minutos activos.
5. **Carrera GPS en Vivo (`Running`)**: Control de sesión (iniciar, pausar, reanudar, finalizar), registro de puntos de geolocalización (latitud, longitud, altitud), velocidad actual, ritmo (min/km), frecuencia cardíaca (ppm) y elevación (+m).
6. **Rutinas de Entrenamiento**: Rutinas programadas ("Fuerza completa", "Día de espalda") y marcado de ejercicios/series completadas con peso en kg o peso corporal.
7. **Nutrición e Hidratación**: Lectura del consumo y comidas cargadas en la base de datos, además de ingesta de agua (registros de +150ml, +250ml, +500ml) y macronutrientes.

---

## Endpoints de la API REST (`/api/v1`)

### Autenticación Pública
- `POST /api/v1/auth/register` - Registro de nuevo usuario.
- `POST /api/v1/auth/login` - Inicio de sesión con correo y contraseña.

### Endpoints Protegidos (requiere header `Authorization: Bearer <token>`)

#### Dashboard Principal
- `GET /api/v1/dashboard` - Obtiene el resumen consolidado de la pantalla de inicio.

#### Perfil y Usuario
- `GET /api/v1/users/profile` - Obtiene la información del usuario autenticado.
- `GET /api/v1/users/profile/details` - Obtiene los datos corporales y métricas del perfil.
- `PUT /api/v1/users/profile` - Actualiza peso, altura, edad, género y nivel de actividad.
- `GET /api/v1/users/settings` y `PUT /api/v1/users/settings` - Consulta y actualiza nombre y correo.
- `GET /api/v1/users/notification-preferences` y `PUT /api/v1/users/notification-preferences` - Consulta y actualiza preferencias de notificación.
- `PUT /api/v1/users/password` - Ruta reservada; responde `501` porque el cambio está deshabilitado.
- `DELETE /api/v1/users/account` - Anonimiza la cuenta, elimina sus datos asociados y bloquea sus tokens activos.

#### Actividad y Pasos
- `GET /api/v1/activity/today` - Obtiene el resumen de pasos y actividad del día.
- `POST /api/v1/activity/steps` - Sincroniza los pasos y distancia registrados por el dispositivo.
- `GET /api/v1/activity/weekly` - Obtiene los datos de actividad de los últimos 7 días para la gráfica.

#### Carreras y GPS en Vivo
- `POST /api/v1/runs/start` - Inicia una nueva sesión de carrera GPS.
- `GET /api/v1/runs/active` - Obtiene la carrera activa en curso (si existe).
- `POST /api/v1/runs/:id/pause` - Pausa la sesión de carrera.
- `POST /api/v1/runs/:id/resume` - Reanuda la sesión de carrera.
- `POST /api/v1/runs/:id/finish` - Finaliza la sesión de carrera y guarda las métricas globales.
- `POST /api/v1/runs/:id/gps` - Registra coordenadas GPS en tiempo real para el mapa.

#### Rutinas y Ejercicios
- `GET /api/v1/routines/today` - Obtiene la rutina de ejercicios de hoy con la lista de ejercicios.
- `POST /api/v1/routines/exercises/:id/complete` - Marca o desmarca un ejercicio como completado.

#### Nutrición e Hidratación
- `GET /api/v1/nutrition/today` - Obtiene el consumo de agua, calorías y macros de hoy.
- `POST /api/v1/nutrition/water` - Registra consumo de agua (+150ml, +250ml, +500ml).
- Las comidas mostradas por `GET /api/v1/nutrition/today` se leen de la base de datos; no existe endpoint para registrarlas.

El contrato OpenAPI con esquemas de request y response está en [openapi.yaml](openapi.yaml).

---

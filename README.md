# Order Status Service

Este microservicio se encarga de manejar los estados de las órdenes dentro del sistema, y un catálogo base de estados administrable por usuarios con rol de administrador.
Forma parte del ecosistema de microservicios, y se comunica con el microservicio de autenticación (auth-service) para validar los tokens JWT y los permisos.


## Roles y permisos
* Usuario
  * Ver el estado de sus propias órdenes.
  * No puede modificar ni crear estados.
* Administrador
  * Crear nuevos estados en el catálogo base.
  * Ver el catálogo completo (que contiene a los estados base).
  * Cambiar el estado de cualquier orden.

 
## Requisitos previos
* Docker y Docker Compose instalados.
* Microservicio de autenticación (prod-auth-go) en ejecución.
* Microservicio de ordenes (prod-orders-go) en ejecución.
* Base de datos MongoDB accesible (puede estar en otro contenedor).


## Build y ejecución con Docker
En la raíz del proyecto abrir una nueva terminal y ejecutar:
``` bash
docker-compose up --build
```


## Autenticación
Cada endpoint que modifica información requiere un token JWT válido.
El token se valida comunicándose con el microservicio de autenticación configurado en AUTH_SERVICE_URL.
Una vez obtenido el token correspondiente según el caso de uso (admin o user), agregar en Postman o el cliente HTTP:
``` bash
Authorization: Bearer <TOKEN>
```


## Endpoints

### 1. Estados base del catálogo (solo administradores)
#### Ver todos los estados del catálogo
``` pgsql
GET http://localhost:8080/admin/status/catalog
Authorization: Bearer <TOKEN_ADMIN>
```

Respuesta:
``` JSON
[
  { "name": "Pendiente" },
  { "name": "En preparación" },
  { "name": "Enviado" },
  { "name": "Entregado" },
  { "name": "Cancelado" }
]
```


#### Crear nuevo estado en el catálogo
``` pgsql
POST /admin/status/catalog
Authorization: Bearer <TOKEN_ADMIN>
```

Body:
``` JSON
{
  "name": "Rechazado"
}
```

Respuesta:
``` JSON
{
  "message": "status added to catalog"
}
```


### 2. Estados de órdenes reales
#### (Automático) Crear un nuevo estado al realizar una orden
Cuando el microservicio de órdenes registra una nueva orden, debe hacer un POST al siguiente endpoint para inicializar su estado en “Pendiente”:
```
POST /status/init
```

Body:
``` JSON
{
  "order_id": "30c08621-2eca-418b-a314-3112f3106230",
  "user_id": "69123bca4f816b60d535c741"
}
```

Respuesta esperada:
``` JSON
{
    "id": "6913e3833d346cc666bf096b",
    "order_id": "30c08621-2eca-418b-a314-3112f3106230",
    "user_id": "69123bca4f816b60d535c741",
    "status": "Pendiente",
    "updated_at": "2025-11-12T01:31:47.710277285Z"
}
```


#### Cambiar el estado de una orden (solo admin)
``` pgsql
PUT /status/:order_id
Authorization: Bearer <TOKEN_ADMIN>
```

Body:
``` JSON
{
  "status": "Enviado"
}
```

Respuesta:
``` JSON
{
  "order_id": "6913d0ed4db2631092337add",
  "status": "Enviado"
}
```


#### Ver los estados de las órdenes del usuario actual autenticado
``` sql
GET /status
Authorization: Bearer <TOKEN_USER>
```

Respuesta:
``` JSON
[
  {
    "order_id": "6913d0ed4db2631092337add",
    "status": "Pendiente"
  },
  {
    "order_id": "6913bf285e55b075693ca1d8",
    "status": "Entregado"
  }
]
```


#### Obtener el listado de todas las órdenes y sus estados (sólo admin)
``` sql
GET status/all
Authorization: Bearer <TOKEN_ADMIN>
```

Respuesta:
``` JSON
[
    {
        "id": "6913d0ed4db2631092337add",
        "order_id": "88242cab-53de-4f1f-9a82-7d99165f6014",
        "user_id": "6913bf285e55b075693ca1d8",
        "status": "Enviado",
        "updated_at": "2025-11-12T02:44:56.276Z"
    },
    {
        "id": "6913e3833d346cc666bf096b",
        "order_id": "30c08621-2eca-418b-a314-3112f3106230",
        "user_id": "69123bca4f816b60d535c741",
        "status": "Pendiente",
        "updated_at": "2025-11-12T01:31:47.71Z"
    }
{
        "id": "6913e3833d346cc666bf096c",
        "order_id": "30c08621-2eca-348b-a314-3112f323819",
        "user_id": "34975zca4f816b60d538m357",
        "status": "En Preparación",
        "updated_at": "2025-11-12T02:40:32.71Z"
    }
]
```


#### Obtener ordenes por estado

``` sql
GET /status/filter?status=Enviado
Authorization: Bearer <TOKEN_ADMIN>
```

Respuesta:
``` JSON
[
    {
        "id": "6913d0ed4db2631092337add",
        "order_id": "88242cab-53de-4f1f-9a82-7d99165f6014",
        "user_id": "6913bf285e55b075693ca1d8",
        "status": "Enviado",
        "updated_at": "2025-11-12T02:44:56.276Z"
    }
]
```

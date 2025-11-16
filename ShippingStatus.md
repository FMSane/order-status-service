# Shipping Status

## Modelo de Datos

### Catálogo de estados

``` JSON
{
  "id": string,
  "name": string,
  "created_at": string
}
```

Estado de orden

``` JSON
{
    "id": string,
    "order_id": string,
    "user_id": string,
    "status_id": string,
    "status": string,
    "shipping": {
        "address_line1": string,
        "city": string,
        "country": string
    },
    "created_at": string,
    "updated_at": string
}
```

## API

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

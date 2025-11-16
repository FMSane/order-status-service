# Order Status Service

Este microservicio se encarga de manejar los estados de las órdenes dentro del sistema, y un catálogo base de estados administrable por usuarios con rol de administrador.
Forma parte del ecosistema de microservicios, y se comunica con el microservicio de autenticación (auth-service) para validar los tokens JWT y los permisos.


## Roles y permisos
* Usuario
  * Ver el estado de sus propias órdenes.
  * Sólo puede cancelar una orden, si es que esta no está en estado "Enviado", "Entregado" ni "Rechazado". 
  * No puede modificar ni crear estados base.
* Administrador
  * Crear nuevos estados en el catálogo base.
  * Ver el catálogo completo (que contiene a los estados base).
  * Cambiar el estado de cualquier orden, excepto al estado "Cancelado".
  * Puede "Rechazar" una orden, si es que esta no está en estado "Cancelado", "Enviado" ni "Entregado".

* Otras consideraciones
  * Al establecer el estado de una orden, el sistema comprobará que ese estado no sea el actual de la orden, para así proceder a actualizarlo.
  * Si una orden posee estado "Cancelado", "Rechazado" o "Entregado", ya no se podrá cambiar el estado (estados finales).
 
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


## Modelo de Datos y API's

[Estados de Envío](ShippingStatus.md)

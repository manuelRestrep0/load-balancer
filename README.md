# Balanceador de Carga

Trabajo final de Sistemas Operativos

Autores: Jair Santiago Leal, Juan Manuel Restrepo

## Descripción

Este proyecto implementa un balanceador de carga HTTP en Go que distribuye peticiones entre múltiples servidores backend utilizando el algoritmo round-robin. El sistema está diseñado para demostrar conceptos de concurrencia, sincronización y comunicación entre procesos.

## Componentes

El proyecto está compuesto por cuatro archivos principales:

- main.go: Punto de entrada del programa. Inicializa los backends, crea el balanceador, implementa un worker pool con 4 workers y configura los endpoints HTTP.

- balancer.go: Implementa la lógica del balanceador de carga con el algoritmo round-robin usando operaciones atómicas para garantizar la distribución equitativa de peticiones.

- backend.go: Simula servidores backend que procesan peticiones con un tiempo de respuesta aleatorio entre 200 y 800 milisegundos.

- request.go: Define la estructura Request que encapsula la información necesaria para procesar cada petición HTTP.

## Funcionamiento

El sistema inicia tres servidores backend en los puertos 9001, 9002 y 9003. El balanceador principal escucha en el puerto 8080 y distribuye las peticiones entrantes entre los backends disponibles.

Cuando llega una petición al balanceador, esta se envía a través de un canal a uno de los cuatro workers del worker pool. El worker selecciona el siguiente backend usando round-robin, reenvía la petición y devuelve la respuesta al cliente.

El algoritmo round-robin garantiza que las peticiones se distribuyan de manera equitativa entre los backends, rotando secuencialmente entre ellos.

## Características

- Distribución de carga mediante algoritmo round-robin
- Worker pool con 4 workers para procesamiento concurrente
- Uso de goroutines y channels para manejo de concurrencia
- Operaciones atómicas para contadores thread-safe
- Endpoint de estadísticas en /stats que muestra cuántas peticiones procesó cada worker

## Uso

Para ejecutar el programa:

```
go run .
```

El balanceador estará disponible en http://localhost:8080 y los backends en los puertos 9001, 9002 y 9003.

Para ver las estadísticas de los workers, acceder a http://localhost:8080/stats

## Conceptos de Sistemas Operativos Aplicados

- Concurrencia: Implementada mediante goroutines que permiten procesar múltiples peticiones simultáneamente
- Sincronización: Uso de channels para comunicación segura entre goroutines y operaciones atómicas para contadores compartidos
- Comunicación entre procesos: Los workers se comunican con el handler principal a través de channels
- Planificación: El worker pool actúa como un planificador que distribuye el trabajo entre los workers disponibles


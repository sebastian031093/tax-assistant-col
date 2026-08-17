# 📚 Material de Estudio: Fundamentos de Go y Arquitectura Web Nativa

Este documento resume los conceptos clave, las herramientas nativas y las decisiones de diseño aplicadas durante la primera fase del proyecto **`tax-assistant-col`**.

---

## 1. Conceptos Básicos de Go Implementados

### 🔹 Handlers (Manejadores)

En Go, un **handler** es el encargado de recibir una petición HTTP (`http.Request`) y escribir una respuesta (`http.ResponseWriter`).

- **La Interfaz `http.Handler`**: Cualquier estructura que implemente el método `ServeHTTP(ResponseWriter, *Request)` se considera un handler válido.
- **`http.ResponseWriter`**: Es una **interfaz** nativa que se comunica directamente con la conexión de red del cliente. No se usa como puntero porque las interfaces en Go ya actúan como referencias internas.
- **`*http.Request`**: Es un **puntero a una estructura** que contiene toda la metadata de la petición (Headers, Body, Cookies, URL). Se pasa como puntero para evitar duplicar datos pesados en memoria (eficiencia) y permitir que los middlewares modifiquen su contexto de forma segura.

### 🔹 Concurrencia y Goroutines

- **`go func() { ... }()`**: Permite ejecutar funciones de forma asíncrona. Se utilizó para lanzar el método bloqueante `server.ListenAndServe()` en un hilo ligero de Go, liberando el flujo principal para escuchar las señales de apagado del sistema operativo.

### 🔹 Contextos (`context.Context`)

Herramienta nativa para gestionar los tiempos de vida, propagar señales de cancelación y pasar datos a través de los límites de la API y los hilos de ejecución.

---

## 2. Librerías Nativas Utilizadas (`Standard Library`)

Go promueve el uso de su robusta librería estándar antes de saltar a frameworks de terceros:

| Paquete                 | Componente Clave       | Propósito en el Proyecto                                                                                                                        |
| :---------------------- | :--------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------- |
| **`net/http`**          | `http.ServeMux`        | El enrutador nativo. Utiliza coincidencia por prefijos y permite restringir métodos HTTP mediante sintaxis explícita (`GET /ruta`).             |
| **`net/http`**          | `http.Server`          | Estructura de configuración avanzada para el servidor (timeouts de lectura, escritura e inactividad).                                           |
| **`net/http/httptest`** | `httptest.NewRecorder` | Captura las respuestas HTTP en memoria para realizar afirmaciones (_assertions_) en las pruebas unitarias sin levantar un servidor real de red. |
| **`os` & `strconv`**    | `os.LookupEnv`         | Permiten la lectura y el parseo de variables de entorno del sistema operativo con valores de respaldo (_fallbacks_).                            |
| **`os/signal`**         | `signal.NotifyContext` | Intercepta señales críticas enviadas por el kernel del sistema operativo (como interrupciones de teclado o comandos de apagado de la nube).     |

---

## 3. Arquitectura del Proyecto e Implementación

Hemos estructurado el proyecto siguiendo el estándar idiomático de la comunidad (**Standard Go Project Layout**):

```text
tax-assistant-col/
├── cmd/
│   └── api/
│       └── main.go       # Punto de entrada (Inicia la lectura y arranca el ciclo)
├── internal/             # Carpeta protegida por el compilador para código privado
│   ├── config/
│   │   └── config.go     # Estructura de configuración y carga de entorno
│   └── server/
│       ├── handlers.go   # Controladores lógicos (Respuestas estructuradas en JSON)
│       ├── routes.go     # Mapeo de endpoints y captura del Error 404 global
│       ├── server.go     # Orquestación del servidor y Graceful Shutdown
│       └── server_test.go# Pruebas de tabla dirigidas por datos (Table-Driven Tests)
└── go.mod                # Identificador profesional del módulo y control de dependencias
```

### ⚙️ Mecanismo de Apagado Controlado (Graceful Shutdown)

Cuando el servidor recibe una señal de detención (`Ctrl + C` o `SIGTERM`), el ciclo de vida se gestiona de la siguiente manera:

1. **`signal.NotifyContext`** captura la señal y cancela el contexto principal de la aplicación.
2. El servidor invoca **`server.Shutdown()`**, lo que detiene inmediatamente la aceptación de nuevas conexiones.
3. Se genera un **`context.WithTimeout`** de 10 segundos, otorgando una ventana de tolerancia para que las peticiones que ya estaban en curso terminen de procesarse limpiamente antes de cerrar el proceso por completo.

### 🧪 Pruebas Dirigidas por Tablas (Table-Driven Tests)

Estructura idiomática que define un _slice_ de structs anónimos para evaluar múltiples casos de prueba sobre el enrutador en un único bucle. Permite validar de forma centralizada códigos de respuesta (`200 OK`, `405 Method Not Allowed`, `404 Not Found`) y contenido JSON mediante sub-pruebas aisladas (`t.Run`).

---

## 4. Recursos para Consulta en Profundidad

- **[Go Packages - net/http](https://go.dev)**: Documentación oficial detallada sobre la interfaz `Handler`, el funcionamiento de `ServeMux` y la configuración de `http.Server`.
- **[Go User Manual - Effective Go](https://go.dev)**: Guía oficial con las mejores prácticas para escribir código claro, idiomático y eficiente en Go.
- **[Golang Project Layout (Standard)](https://github.com)**: Repositorio de referencia de la comunidad para entender la separación de conceptos (`cmd/`, `internal/`, `pkg/`).
- **[Go Modules Reference](https://go.dev)**: Manual técnico completo sobre cómo el motor de Go gestiona el archivo `go.mod`, la resolución de importaciones y la integridad con `go.sum`.

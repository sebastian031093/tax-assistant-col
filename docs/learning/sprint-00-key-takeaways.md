# Go fundamentals

0. QUE ES UN HANDLER EN GO:
   En Go, un handler (manejador) es un bloque de código o una estructura encargada de procesar las peticiones HTTP entrantes y generar una respuesta. Funciona como el puente entre el cliente (un navegador o app) y tu servidor.
   La interfaz http.HandlerEn el paquete nativo net/http, un handler es cualquier tipo que implemente la interfaz http.Handler, la cual requiere un único método llamado ServeHTTP

1)  ¿Por qué r es un puntero (\*http.Request)?
    Se utiliza un puntero al objeto de la petición (Request) por dos razones fundamentales:Evitar copias pesadas en memoria (Eficiencia): Una solicitud HTTP puede contener muchísima información (cabeceras, cookies, formularios, datos del navegador, etc.). Si pasaras el objeto por valor (sin el asterisco), Go tendría que duplicar toda esa estructura en memoria para cada petición. Pasar un puntero copia únicamente la dirección de memoria, lo cual es instantáneo y consume recursos mínimos.
    Permitir modificaciones y compartir estado: Durante el ciclo de vida de una petición, es común que la información se actualice. Por ejemplo, al leer el cuerpo (Body) de la solicitud, añadir datos de autenticación o guardar variables de contexto. Usar un puntero asegura que todos los componentes (como los middlewares) lean y escriban exactamente sobre el mismo objeto original.

2)  ¿Por qué w NO lleva un asterisco (http.ResponseWriter)?A primera vista parece confuso que r sea un puntero y w no. La respuesta es sencilla: http.ResponseWriter es una interfaz, no una estructura (struct).

En Go, las interfaces ya actúan internamente como referencias. Detrás de escena, el servidor web de Go le pasa a tu función un objeto oculto que implementa los métodos de ResponseWriter, y esa implementación ya maneja los punteros necesarios para escribir directamente en la conexión de red del cliente. Ponerle un asterisco extra sería un error de sintaxis.

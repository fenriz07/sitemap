# Exercise #5: Sitemap Builder

[![exercise status: released](https://img.shields.io/badge/exercise%20status-released-green.svg?style=for-the-badge)](https://gophercises.com/exercises/sitemap)

## Detalles

A sitemap is basically a map of all of the pages within a specific domain. They are used by search engines and other tools to inform them of all of the pages on your domain.

Una forma de construirlos es visitando primero la página raíz del sitio web y haciendo una lista de cada enlace en esa página que vaya a una página en el mismo dominio. Por ejemplo, en`calhoun.io` puede encontrar un enlace a `calhoun.io/hire-me/` junto con varios otros enlaces.

Una vez que haya creado la lista de enlaces, puede visitar cada uno y agregar nuevos enlaces a su lista. Al repetir este paso una y otra vez, eventualmente visitaría todas las páginas del dominio a las que se puede acceder siguiendo los enlaces desde la página raíz.

En este ejercicio, su objetivo es crear un creador de mapas de sitio como el descrito anteriormente. El usuario final ejecutará el programa y le proporcionará una URL. (*hint - use a flag or a command line arg for this!*) que usarás para comenzar el proceso.

Una vez que haya determinado todas las páginas de un sitio, el creador de su mapa del sitio debería generar los datos en el siguiente formato XML:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>http://www.example.com/</loc>
  </url>
  <url>
    <loc>http://www.example.com/dogs</loc>
  </url>
</urlset>
```

*Note: This should be the same as the [standard sitemap protocol](https://www.sitemaps.org/index.html)*

Donde cada página se enumera en su propio `<url>` etiqueta e incluye el `<loc>`etiqueta dentro de ella.

Para completar este ejercicio, recomiendo hacer primero [link parser exercise](https://github.com/gophercises/link) y usar el paquete creado en él para analizar sus páginas HTML para enlaces.

A partir de ahí, es probable que necesite encontrar una manera de determinar si un enlace va al mismo dominio o a uno diferente. Si va a un dominio diferente, no deberíamos incluirlo en nuestro generador de mapas de sitio, pero si es el mismo dominio, deberíamos hacerlo. Recuerde que los enlaces al mismo dominio pueden tener el formato de `/just-the-path` o `https://domain.com/with-domain`, pero ambos van al mismo dominio.

### Notes

**1. Tenga en cuenta que los enlaces pueden ser cíclicos.**

Es decir, página`abc.com` may link to page `abc.com/about`, y luego la página acerca puede volver a la página de inicio (`abc.com`). Estos ciclos también pueden ocurrir en muchas páginas, por ejemplo, puede tener:

```
/about -> /contact
/contact -> /pricing
/pricing -> /testimonials
/testimonials -> /about
```

Donde el ciclo toma 4 enlaces para finalmente alcanzarlo, pero de hecho hay un ciclo.

Es importante recordar esto porque no desea que su programa entre en un bucle infinito donde sigue visitando las mismas páginas una y otra vez. Si tiene problemas con esto, el ejercicio de bonificación puede ayudar a aliviar el problema temporalmente, pero cubriremos cómo evitarlo por completo en las transmisiones de pantalla para este ejercicio.

**2. Los siguientes paquetes serán útiles ...**

- [net/http](https://golang.org/pkg/net/http/) -para iniciar solicitudes GET a cada página en su mapa del sitio y obtener el HTML en esa página
- la rama `solución` de [github.com/gophercises/link](https://github.com/gophercises/link) - no podrá "go get" este paquete porque no está comprometido con master, pero si completa el ejercicio localmente, puede usar el código de él en este ejercicio. Si esto causa confusión o problemas, ¡te ayudaré a descubrir cómo hacer todo esto!<jon@calhoun.io>

- [encoding/xml](https://golang.org/pkg/encoding/xml/) -para imprimir la salida XML al final
- [flag](https://golang.org/pkg/flag/) - analizar indicadores proporcionados por el usuario como el dominio del sitio web

Probablemente me faltan algunos paquetes aquí, así que no se preocupe si está utilizando otros. Esta es solo una lista aproximada de paquetes que espero usar cuando codifique la solución 😁

## Bonus

Como ejercicios de bonificación también puede agregar un `depth` bandera que define el número máximo de enlaces a seguir al crear un mapa del sitio. Por ejemplo, si tenía una profundidad máxima de 3 y los siguientes enlaces:

```
a->b->c->d
```

Entonces su creador de sitemaps no visitaría ni incluiría `d` porque debes seguir más de 3 enlaces para acceder a la página.

Por otro lado, si los enlaces de la página fueran así:

```
a->b->c->d
b->d
```

Donde también hay un enlace a la página `d` De la página `b`,entonces su creador de sitemaps debe incluir `d` porque se puede acceder en 3 enlaces.

*Hint - I find using a BFS ([breadth-first search](https://en.wikipedia.org/wiki/Breadth-first_search)) is the best way to achieve this bonus exercise without doing extra work, but it isn't required and you could likely come up with a working solution without using a BFS.*

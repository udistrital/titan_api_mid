# titan_api_mid

Gestor de reglas para el aplicativo [titan_api_crud](https://github.com/udistrital/titan_api_crud).

## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones
* [Golang](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/golang.md)
* [BeeGo](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/beego.md)
* [Docker](https://docs.docker.com/engine/install/ubuntu/)
* [Docker Compose](https://docs.docker.com/compose/)

### Variables de Entorno
```shell
TITAN_API_MID__NSCRUD=[Descripción]
TITAN_API_MID__NSKRONOS=[Descripción]
TITAN_API_MID__NSARGOMID=[Descripción]
TITAN_API_MID__NSARGOCRUD=[Descripción]
TITAN_API_MID__NSWSO2ARGO=[Descripción]
TITAN_API_MID__NSWSO2COLAS=[Descripción]
TITAN_API_MID__NSARGOAMAZON=[Descripción]
```
Prueba
**NOTA:** Las variables se pueden ver en el fichero conf/app.conf y están identificadas con TITTITAN_API_MID__...

### Ejecución del Proyecto
```shell
#1. Obtener el repositorio con Go
go get github.com/udistrital/titan_api_mid

#2. Moverse a la carpeta del repositorio
cd $GOPATH/src/github.com/udistrital/titan_api_mid

# 3. Moverse a la rama **develop**
git pull origin develop && git checkout develop

# 4. alimentar todas las variables de entorno que utiliza el proyecto.
TITAN_API_MID__NSCRUD=8080 TITAN_API_MID__NSARGOMID=127.0.0.1:27017 TITAN_API_MID__SOME_VARIABLE=some_value bee run
```

### Ejecución docker-compose
```shell
#1. Clonar el repositorio
git clone -b develop https://github.com/udistrital/titan_api_mid

#2. Moverse a la carpeta del repositorio
cd titan_api_mid

#3. Crear un fichero con el nombre **custom.env**
# En windows ejecutar:* ` ni custom.env`
touch custom.env

#4. Crear la network **back_end** para los contenedores
docker network create back_end

#5. Ejecutar el compose del contenedor
docker-compose up --build

#6. Comprobar que los contenedores estén en ejecución
docker ps
```

### Ejecución Pruebas
Pruebas unitarias
```shell
# Not Data
```
## Estado CI

| Develop | Relese 0.0.1 | Master |
| -- | -- | -- |
| [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/titan_api_mid/status.svg?ref=refs/heads/develop)](https://hubci.portaloas.udistrital.edu.co/udistrital/titan_api_mid) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/titan_api_mid/status.svg?ref=refs/heads/release/0.0.1)](https://hubci.portaloas.udistrital.edu.co/udistrital/titan_api_mid) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/titan_api_mid/status.svg)](https://hubci.portaloas.udistrital.edu.co/udistrital/titan_api_mid) |

## Licencia

This file is part of titan_api_mid.

titan_api_mid is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

titan_api_mid is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with titan_api_mid. If not, see https://www.gnu.org/licenses/.

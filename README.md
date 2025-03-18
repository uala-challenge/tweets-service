# Tweets Service
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=uala-challenge_tweets-service&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=uala-challenge_tweets-service)



## Descripción
El servicio **Tweets Service** es responsable de almacenar los tweets creados por los usuarios y publicarlos en un tópico de mensajería para su procesamiento asincrónico por otros servicios de la aplicación.

## Características principales
- Creación y almacenamiento de tweets con toda su información.
- Publicación de eventos en un tópico de mensajería para procesamiento asíncrono.
- API REST para la gestión de tweets.
- Logs detallados para monitoreo y depuración.

## Tecnologías utilizadas
- **Golang**: Lenguaje de desarrollo principal.
- **DynamoDB**: Base de datos utilizada para almacenar los tweets.
- **SNS**: Broker de mensajería para publicación de eventos.
- **Docker**: Contenedorización para despliegue.
- **REST API**: Exposición de endpoints.
  dynamo:
  endpoint: 'http://localhost:4566'

sns:
endpoint: 'http://localhost:4566'
- **Logrus**: Manejo de logs estructurados.

## Instalación y configuración
### Prerrequisitos
- Tener instalado **Go** (versión 1.18 o superior).
- Contar con **Docker** y **Docker Compose** (para pruebas locales).
- Configurar las variables de entorno adecuadas.

### Clonar el repositorio
```bash
  git clone https://github.com/uala-challenge/tweets-service.git
  cd tweets-service
```

### Configuración
El servicio utiliza un archivo de configuración en formato YAML. Antes de ejecutar, asegúrate de configurar correctamente `config.yaml`:
```yaml
database:
  dynamo_table: "tweets"

broker:
  type: kafka  # Opciones: kafka, rabbitmq
  url: "localhost:9092"
  topic: "tweets-events"

api:
  port: 8083

logging:
  level: "info"
```

### Ejecutar en entorno local
#### Usando Go directamente
```bash
  go run main.go
```

#### Usando Docker
```bash
  docker-compose up --build
```

## API REST
El servicio expone los siguientes endpoints:

### **Crear un tweet**
```
POST /tweets
```
#### **Descripción:**
Permite a un usuario crear un nuevo tweet. El tweet se almacena en **DynamoDB** y luego se publica un evento en el tópico de mensajería para su procesamiento asíncrono.

#### **Payload esperado:**
```json
{
  "user_id": "12345",
  "content": "Este es un nuevo tweet",
  "timestamp": "2025-03-17T12:00:00Z"
}
```

#### **Respuesta:**
```json
{
  "tweet_id": "abc123",
  "message": "Tweet creado exitosamente"
}
```

### **Publicación en el tópico de mensajería**
Una vez almacenado el tweet, se publica un mensaje en el tópico **tweets-events**:
```json
{
  "tweet_id": "abc123",
  "user_id": "12345",
  "content": "Este es un nuevo tweet",
  "timestamp": "2025-03-17T12:00:00Z"
}
```
Este evento será consumido por el servicio **Event Processor** para su procesamiento asíncrono.


## Testing
Para ejecutar las pruebas unitarias:
```bash
  go test ./...
```

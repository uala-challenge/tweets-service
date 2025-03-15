## Descripción

Este módulo forma parte de la solución cash-manager-core, diseñada para modernizar la gestión de transacciones financieras en una arquitectura distribuida. Cash-manager-core es una aplicación desarrollada en Go, orientada a paquetes y basada en los principios de arquitectura limpia, lo que garantiza modularidad, escalabilidad y facilidad de mantenimiento.

El módulo de Procesamiento de Mensajes es responsable de consumir mensajes provenientes de tópicos en una cola de mensajes, aplicar validaciones y lógica de negocio, y finalmente almacenar los datos en sistemas de bases de datos relacionales y no relacionales. Este módulo asegura que todas las transacciones se registren correctamente y se clasifiquen de acuerdo con su tipología (provisión, retiro, dispensación), asegurando la consistencia y precisión del sistema.


## Tecnologías Utilizadas

- **Lenguaje:** Go (Golang).
- **Mensajería:** Apache Kafka o RabbitMQ.
- **Bases de Datos:** Oracle Database (relacional), MongoDB/Redis (no relacional).
- **Plataforma:** Oracle Cloud Infrastructure (OCI).
- **Logs y Monitorización:**
    - ElasticSearch para manejo de logs.
    - Splunk para el monitoreo de transacciones y análisis operativos.
    - New Relic para monitoreo de métricas de aplicación.


#!/bin/sh
echo "Esperando a que LocalStack inicie completamente..."
sleep 5

echo "Creación de ambiente local con LocalStack..."

echo "Creando tabla DynamoDB 'UalaChallenge' en LocalStack..."
aws --endpoint-url=http://localstack:4566 --region us-east-1 dynamodb create-table \
    --table-name UalaChallenge \
    --attribute-definitions \
        AttributeName=PK,AttributeType=S \
        AttributeName=SK,AttributeType=S \
        AttributeName=GSI1PK,AttributeType=S \
        AttributeName=GSI1SK,AttributeType=S \
        AttributeName=GSI2PK,AttributeType=S \
        AttributeName=GSI2SK,AttributeType=S \
    --key-schema \
        AttributeName=PK,KeyType=HASH \
        AttributeName=SK,KeyType=RANGE \
    --global-secondary-indexes \
        "[
            {
                \"IndexName\": \"GSI1\",
                \"KeySchema\": [
                    {\"AttributeName\": \"GSI1PK\", \"KeyType\": \"HASH\"},
                    {\"AttributeName\": \"GSI1SK\", \"KeyType\": \"RANGE\"}
                ],
                \"Projection\": {\"ProjectionType\": \"ALL\"}
            },
            {
                \"IndexName\": \"GSI2\",
                \"KeySchema\": [
                    {\"AttributeName\": \"GSI2PK\", \"KeyType\": \"HASH\"},
                    {\"AttributeName\": \"GSI2SK\", \"KeyType\": \"RANGE\"}
                ],
                \"Projection\": {\"ProjectionType\": \"ALL\"}
            }
        ]" \
    --billing-mode PAY_PER_REQUEST

echo "Tabla DynamoDB creada exitosamente."

SNS_TOPIC_ARN=$(aws --endpoint-url=http://localstack:4566 --region us-east-1 sns create-topic --name uala-challenge --query 'TopicArn' --output text)
echo "SNS Topic creado: $SNS_TOPIC_ARN"

SQS_QUEUE_URL=$(aws --endpoint-url=http://localstack:4566 --region us-east-1 sqs create-queue --queue-name tweets --query 'QueueUrl' --output text)
SQS_QUEUE_ARN=$(aws --endpoint-url=http://localstack:4566 --region us-east-1 sqs get-queue-attributes --queue-url "$SQS_QUEUE_URL" --attribute-names QueueArn --query 'Attributes.QueueArn' --output text)

echo "SQS Queue creada: $SQS_QUEUE_URL"
echo "ARN de la cola SQS: $SQS_QUEUE_ARN"

SUBSCRIPTION_ARN=$(aws --endpoint-url=http://localstack:4566 --region us-east-1 sns subscribe \
    --topic-arn "$SNS_TOPIC_ARN" \
    --protocol sqs \
    --notification-endpoint "$SQS_QUEUE_ARN" \
    --query 'SubscriptionArn' --output text)

echo "Suscripción de SQS a SNS creada: $SUBSCRIPTION_ARN"

aws --endpoint-url=http://localstack:4566 --region us-east-1 sns set-subscription-attributes \
    --subscription-arn "$SUBSCRIPTION_ARN" \
    --attribute-name RawMessageDelivery \
    --attribute-value true

echo "RawMessageDelivery activado para SNS-SQS."
echo "Configuración de LocalStack finalizada correctamente."
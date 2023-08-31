package database

import (
	"api-pendencias/model"
	"api-pendencias/utils"
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
)

type Connection struct {
	Logs            *utils.Logger
	DynamodbCliente *dynamodb.Client
	Ctx             *gin.Context
}

func (database *Connection) SelectAllPendencias() ([]model.Pendencia, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Pendencia"),
		ExpressionAttributeNames: map[string]string{
			"#P": "Pago",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberBOOL{Value: false},
		},
		FilterExpression: aws.String("#P = :v1"),
	}

	result, err := database.DynamodbCliente.Scan(context.TODO(), input)
	if err != nil {
		utils.HandleError("E", "Erro ao buscar pendencias", err)
		return nil, err
	}

	var pendencias []model.Pendencia
	for _, item := range result.Items {
		var p model.Pendencia
		err = attributevalue.UnmarshalMap(item, &p)
		if err != nil {
			utils.HandleError("E", "Error ao fazer unmarshal de Pendencia", err)
			return nil, err
		}
		pendencias = append(pendencias, p)
	}

	return pendencias, nil
}

func (database *Connection) SelectPendenciasByName() ([]model.Pendencia, error) {
	nome := database.Ctx.Param("nome")
	// Define the expression attributes
	exprAttrNames := map[string]string{
		"#P": "Pago",
		"#N": "Nome",
	}

	exprAttrValues := map[string]types.AttributeValue{
		":pagoVal": &types.AttributeValueMemberBOOL{Value: false},
		":nomeVal": &types.AttributeValueMemberS{Value: nome},
	}

	// Define the query input
	input := &dynamodb.QueryInput{
		TableName:                 aws.String("Pendencia"),
		KeyConditionExpression:    aws.String("#N = :nomeVal"),
		FilterExpression:          aws.String("#P = :pagoVal"),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
	}

	resp, err := database.DynamodbCliente.Query(context.TODO(), input)
	if err != nil {
		utils.HandleError("E", "Error na query ao pegar Pendencia por nome", err)
		return nil, err
	}

	var pendencias []model.Pendencia
	for _, item := range resp.Items {
		var p model.Pendencia
		if err := attributevalue.UnmarshalMap(item, &p); err != nil {
			utils.HandleError("E", "Error ao fazer unmarshal de Pendencia por nome", err)
			return nil, err
		}
		pendencias = append(pendencias, p)
	}

	return pendencias, nil
}

func (database *Connection) CreatePendencia(p model.Pendencia) error {

	item, err := attributevalue.MarshalMap(p)
	if err != nil {
		utils.HandleError("E", "Erro ao criar o json de criar Pendencia", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Pendencia"),
		Item:      item,
	}

	_, err = database.DynamodbCliente.PutItem(context.Background(), input)
	if err != nil {
		utils.HandleError("E", "Erro ao criar Pendencia no banco", err)
		return err
	}

	return nil
}

func (database *Connection) UpdatePendencia(p model.Pendencia) error {
	if p.Nome == "" || p.Dia == "" {
		return errors.New("nome and Dia must be provided")
	}

	updateExpression := "SET Description = :description, Valor = :valor, Pago = :pago"
	expressionValues, err := attributevalue.MarshalMap(map[string]interface{}{
		":description": p.Description,
		":valor":       p.Valor,
		":pago":        p.Pago,
	})

	if err != nil {
		utils.HandleError("E", "Erro no marshal de atualizar pendencia", err)
		return err
	}

	key, err := attributevalue.MarshalMap(map[string]string{
		"Nome": p.Nome,
		"Dia":  p.Dia,
	})
	if err != nil {
		utils.HandleError("E", "Erro no marshal de atualizar pendencia", err)
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String("Pendencia"),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionValues,
	}

	_, err = database.DynamodbCliente.UpdateItem(context.Background(), input)
	if err != nil {
		utils.HandleError("E", "Erro ao atualizar pendencia no banco de dados", err)
		return err
	}

	return nil
}

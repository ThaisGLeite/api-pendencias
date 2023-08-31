package pendencia

import (
	"api-pendencias/database"
	"api-pendencias/model"
	"api-pendencias/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Database struct {
	Database *database.Connection
}

func (db *Database) GetAllPendencias(c *gin.Context) {
	db.Database.Ctx = c
	pendencias, err := db.Database.SelectAllPendencias()
	if err != nil {
		utils.HandleError("E", "Falha ao pegar todas as pendencias", err)
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, pendencias)
}

func (db *Database) GetPendencia(c *gin.Context) {
	db.Database.Ctx = c
	pendencias, err := db.Database.SelectPendenciasByName()
	if err != nil {
		utils.HandleError("E", "Falha ao pegar todas as pendencias por nome", err)
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, pendencias)
}

func (db *Database) CreatePendencia(c *gin.Context) {
	var pendencia model.Pendencia
	err := c.BindJSON(&pendencia)
	if err != nil {
		utils.HandleError("E", "create Pendencia", err)
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	db.Database.Ctx = c
	err = db.Database.CreatePendencia(pendencia)
	if err == nil {
		c.IndentedJSON(http.StatusOK, "pendencia criada com sucesso")
		return
	}
	c.IndentedJSON(http.StatusInternalServerError, err)

}

func (db *Database) UpdatePendencia(c *gin.Context) {
	var pendencia model.Pendencia
	err := c.BindJSON(&pendencia)
	if err != nil {
		utils.HandleError("E", "Update Pendencia", err)
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	db.Database.Ctx = c
	err = db.Database.UpdatePendencia(pendencia)
	if err == nil {
		c.IndentedJSON(http.StatusOK, "pendencia atualizada com sucesso")
		return
	}
	c.IndentedJSON(http.StatusInternalServerError, err)

}

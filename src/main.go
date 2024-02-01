package main


import (
   "fmt"
   "gorm.io/driver/postgres"
   "gorm.io/gorm"
   "log"
   "net/http"
   "github.com/gin-gonic/gin"
  
)

type Todo struct {
	gorm.Model
	Title     string
	Completed bool
 }


var db *gorm.DB

func migrateDB() {
   migrator := db.Migrator()

   // Verificar se a tabela já existe
   if !migrator.HasTable(&Todo{}) {
      // Se a tabela não existe, realizar a migração
      if err := db.AutoMigrate(&Todo{}); err != nil {
         log.Fatal(err)
      }
   }
}

func initDB() {

   dsn := "host=localhost port=5432 user=postgres password=AIzaSyCwz48Sy5zsYhq5_4QsggwQV5ca3XrQpkw dbname=quickbotdb sslmode=disable"
   var err error
   db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
   if err != nil {
      log.Fatal(err)
   }
  // Executar migrações
  migrateDB()
}

func main() {

   initDB()
   defer func() {
      sqlDB, err := db.DB()
      if err != nil {
         panic(err)
      }
      sqlDB.Close()
   }()


   r := gin.Default()

   // Rotas
   r.GET("/todos/:id", getTodos)
   r.POST("/todos", createTodo)
   r.PUT("/todos/:id", updateTodo)
   r.DELETE("/todos/:id", deleteTodo)

   // Iniciar o servidor
   port := 8080
   fmt.Printf("API está rodando em http://localhost:%d\n", port)
   log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

func getTodos(c *gin.Context) {
   var todo Todo
   id := c.Params.ByName("id")
   result := db.First(&todo, id)
   
   if result.Error != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
      return
   }

   c.JSON(http.StatusOK, todo)
}

func createTodo(c *gin.Context) {
   var todo Todo
   if err := c.ShouldBindJSON(&todo); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
   }
   db.Create(&todo)
   c.JSON(http.StatusCreated, todo)
}

func updateTodo(c *gin.Context) {
   id := c.Params.ByName("id")
   var todo Todo
   if err := db.First(&todo, id).Error; err != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
      return
   }

   if err := c.ShouldBindJSON(&todo); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
   }

   db.Save(&todo)
   c.JSON(http.StatusOK, todo)
}

func deleteTodo(c *gin.Context) {
   id := c.Params.ByName("id")
   var todo Todo
   if err := db.First(&todo, id).Error; err != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
      return
   }

   db.Delete(&todo)
   c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}

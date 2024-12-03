package internal

import (
	"context"
	"encoding/json"
	"log"
	"time"

	//"github.com/ValeHenriquez/example-rabbit-go/tasks-server/controllers"
	//"github.com/ValeHenriquez/example-rabbit-go/tasks-server/models"
	"github.com/FelipeGeraldoblufus/product-microservice-go/controllers"
	"github.com/FelipeGeraldoblufus/product-microservice-go/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Handler(d amqp.Delivery, ch *amqp.Channel) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var response models.Response
	log.Println(" [.] Received a message")

	var Payload struct {
		Pattern string          `json:"pattern"`
		Data    json.RawMessage `json:"data"`
		ID      string          `json:"id"`
	}
	var err error
	err = json.Unmarshal(d.Body, &Payload)

	actionType := Payload.Pattern

	//dataJSON, err := json.Marshal(Payload.Data)
	failOnError(err, "Failed to marshal data")
	switch actionType {
	case "GET_PRODUCT":
		log.Println(" [.] Getting product by ID")
	
		var err error
		var productJson []byte
		var product models.Product
	
		// Aquí, Payload.Data debería ser solo un número (ID)
		log.Printf("Received Payload: %s", Payload.Data)
	
		// Convertir el Payload.Data a uint (ID del producto)
		var productID string
		if err := json.Unmarshal(Payload.Data, &productID); err != nil {
			log.Printf("Error unmarshalling data: %v", err)
			response = models.Response{
				Success: "error",
				Message: "Error parsing request data",
				Data:    []byte(err.Error()),
			}
			break
		}
	
		log.Printf("Searching for product with ID: %d", productID)
	
		// Llamar a la función para obtener el producto por ID
		product, err = controllers.GetByProductID(productID)
		if err != nil {
			// Si no se encuentra el producto o ocurre otro error
			log.Printf("Error getting product by ID: %v", err)
			response = models.Response{
				Success: "error",
				Message: "Error getting product",
				Data:    []byte(err.Error()),
			}
		} else {
			// Si todo está bien, devolver el producto en formato JSON
			productJson, err = json.Marshal(product) // Serializar el producto a JSON
			if err != nil {
				log.Printf("Error serializing product: %v", err)
				response = models.Response{
					Success: "error",
					Message: "Error serializing product",
					Data:    []byte(err.Error()),
				}
			} else {
				// Enviar la respuesta con el producto serializado como JSON
				response = models.Response{
					Success: "success",
					Message: "Product retrieved",
					Data:    productJson, // Enviar los datos como JSON
				}
			}
		}


	case "FIND_ALL":
		log.Println(" [.] Getting all products")
		
		var err error
		var productsJson []byte
		var products []models.Product
		
		// Llamar a la función para obtener todos los productos
		products, err = controllers.GetAllProducts()
		if err != nil {
			// Si ocurre un error al obtener los productos
			log.Printf("Error getting all products: %v", err)
			response = models.Response{
				Success: "error",
				Message: "Error getting products",
				Data:    []byte(err.Error()),
			}
		} else {
			// Si todo está bien, devolver los productos en formato JSON
			productsJson, err = json.Marshal(products) // Serializar los productos a JSON
			if err != nil {
				log.Printf("Error serializing products: %v", err)
				response = models.Response{
					Success: "error",
					Message: "Error serializing products",
					Data:    []byte(err.Error()),
				}
			} else {
				// Enviar la respuesta con los productos serializados como JSON
				response = models.Response{
					Success: "success",
					Message: "Products retrieved",
					Data:    productsJson, // Enviar los datos como JSON
				}
			}
		}
	
	

	case "GET_USERBYNAME":
		log.Println(" [.] Getting product by Name")
		var data struct {
			Name string `json:"username"`
		}
		var err error
		var userJson []byte
		var users models.User

		err = json.Unmarshal(Payload.Data, &data)
		users, err = controllers.GetByUser(data.Name)

		userJson, err = json.Marshal(users)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error getting product",
				Data:    []byte(err.Error()),
			}
		} else {
			response = models.Response{
				Success: "succes",
				Message: "Product retrieved",
				Data:    userJson,
			}
		}

	case "EDIT_PRODUCT":
		log.Println(" [.] Editing product by Name")
	
		// Cambiar la estructura para reflejar el JSON con 'updateDTO'
		var data struct {
			UpdateDTO struct {
				Product        string `json:"product"`
				NewNameProduct string `json:"newnameProduct"`
				NewPrice       int    `json:"newPrice"`
				NewStock       int    `json:"newStock"`
				NewDescription string `json:"newDescription"`
				NewCategory    string `json:"newCategory"`
			} `json:"updateDTO"`
		}
	
		var err error
		var userJson []byte
		var producto models.Product
	
		// Log para verificar los datos antes de deserializar
		log.Printf("Received data before unmarshalling: %s", string(Payload.Data))
	
		// Decodificar los datos recibidos
		err = json.Unmarshal(Payload.Data, &data)
		if err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			response = models.Response{
				Success: "error",
				Message: "Error decoding JSON",
				Data:    []byte(err.Error()),
			}
			break
		}
	
		// Log para verificar los datos después del unmarshalling
		log.Printf("Decoded data: %+v", data)
	
		// Verificar que el campo 'Product' no esté vacío
		if data.UpdateDTO.Product == "" {
			log.Println("Error: product name is empty")
			response = models.Response{
				Success: "error",
				Message: "Product name cannot be empty",
				Data:    []byte("Product name cannot be empty"),
			}
			break
		}
	
		// Llamada a la función para actualizar el producto
		producto, err = controllers.UpdateProduct(
			data.UpdateDTO.Product, 
			data.UpdateDTO.NewNameProduct, 
			data.UpdateDTO.NewPrice, 
			data.UpdateDTO.NewStock, 
			data.UpdateDTO.NewDescription, 
			data.UpdateDTO.NewCategory,
		)
		if err != nil {
			log.Printf("Error updating product: %v", err)
			response = models.Response{
				Success: "error",
				Message: "Error updating product",
				Data:    []byte(err.Error()),
			}
			break
		}
	
		// Convertir el resultado a JSON y preparar la respuesta
		userJson, err = json.Marshal(producto)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			response = models.Response{
				Success: "error",
				Message: "Error marshaling JSON",
				Data:    []byte(err.Error()),
			}
		} else {
			log.Printf("Product updated successfully: %+v", producto)
			response = models.Response{
				Success: "success",
				Message: "Product updated",
				Data:    userJson,
			}
		}
	
	
	case "CREATE_PRODUCT":
		log.Println(" [.] Creating product")
	
		// Estructura para deserializar los datos recibidos
		var data struct {
			Name        string  `json:"name"`
			Price       int `json:"price"`
			Stock       int     `json:"stock"`
			Description string  `json:"description"`
			Category string `json:"category"` 
		}
	
		var err error
		var dataJson []byte
		var product models.Product
	
		// Log de depuración para verificar los datos recibidos
		log.Printf("Received data: %+v\n", Payload.Data)
	
		// Deserializar el payload JSON
		err = json.Unmarshal(Payload.Data, &data)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error decoding JSON",
				Data:    []byte(err.Error()),
			}
			break
		}
	
		// Crear el producto utilizando los datos deserializados
		product, err = controllers.CreateProduct(data.Name, data.Price, data.Stock, data.Description, data.Category)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error creating product",
				Data:    []byte(err.Error()),
			}
			break
		}
	
		// Serializar el producto creado a JSON
		dataJson, err = json.Marshal(product)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error marshaling JSON",
				Data:    []byte(err.Error()),
			}
		} else {
			response = models.Response{
				Success: "success",
				Message: "Product created",
				Data:    dataJson,
			}
		}
	

	case "DELETE_PRODUCT":
		log.Println(" [.] Deleting product")
		var data struct {
			Name string `json:"name"`
		}
		var err error
		var dataJson []byte
		var product models.Product
		// Log de depuración para verificar los datos recibidos
		log.Printf("Received data: %+v\n", data.Name)

		err = json.Unmarshal(Payload.Data, &data)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error decoding JSON",
				Data:    []byte(err.Error()),
			}
			break
		}

		err = controllers.DeleteProductByName(data.Name)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error Deleting product",
				Data:    []byte(err.Error()),
			}
			break
		}
		dataJson, err = json.Marshal(product)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error marshaling JSON",
				Data:    []byte(err.Error()),
			}
		} else {
			response = models.Response{
				Success: "success",
				Message: "Product deleted",
				Data:    dataJson,
			}
		}

	case "EDIT_USER":
		log.Println(" [.] Editing user")
		var data struct {
			CurrentUsername string `json:"currentUsername"`
			NewUsername     string `json:"newUsername"`
		}
		var err error
		// Log de depuración para verificar los datos recibidos
		log.Printf("Received data: %+v\n", data.CurrentUsername, data.NewUsername)

		err = json.Unmarshal(Payload.Data, &data)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error decoding JSON",
				Data:    []byte(err.Error()),
			}
			break
		}

		// Llama a la función para editar el usuario
		_, err = controllers.EditUser(data.CurrentUsername, data.NewUsername)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error editing user",
				Data:    []byte(err.Error()),
			}
			break
		}

		response = models.Response{
			Success: "success",
			Message: "User edited successfully",
			Data:    nil, // No necesitas enviar datos específicos en la respuesta
		}

	case "CREATE_USER":
		log.Println(" [.] Creating user")
		var data struct {
			Username string `json:"username"`
		}
		var err error

		// Log de depuración para verificar los datos recibidos
		log.Printf("Received data: %+v\n", data)

		err = json.Unmarshal(Payload.Data, &data)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error decoding JSON",
				Data:    []byte(err.Error()),
			}
			break
		}

		// Verificar que el campo necesario (username) no esté vacío
		if data.Username == "" {
			response = models.Response{
				Success: "error",
				Message: "Username is required",
				Data:    nil,
			}
			break
		}

		// Llama a la función para crear el usuario
		createdUser, err := controllers.CreateUser(data.Username)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error creating user",
				Data:    []byte(err.Error()),
			}
			break
		}

		// Convertir createdUser a formato JSON y luego a []byte
		userData, err := json.Marshal(createdUser)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error encoding user data",
				Data:    []byte(err.Error()),
			}
			break
		}

		response = models.Response{
			Success: "success",
			Message: "User created successfully",
			Data:    userData,
		}

	case "DELETE_USER":
		log.Println(" [.] Deleting user")
		var data struct {
			Username string `json:"username"`
		}
		var err error
		// Log de depuración para verificar los datos recibidos
		log.Printf("Received data: %+v\n", data.Username)

		err = json.Unmarshal(Payload.Data, &data)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error decoding JSON",
				Data:    []byte(err.Error()),
			}
			break
		}

		// Llama a la función para eliminar el CartItem
		err = controllers.DeleteUser(data.Username)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error deleting cartitem",
				Data:    []byte(err.Error()),
			}
			break
		}

		response = models.Response{
			Success: "success",
			Message: "User deleted successfully",
			Data:    nil, // No necesitas enviar datos específicos en la respuesta
		}

	case "CREATE_CATEGORY":
		log.Println(" [.] Creating category")
		//log.Println("data ", Payload.Data.Data)
		//log.Println("data JSON", dataJSON)

		/*var category models.Category
		err := json.Unmarshal(Payload.Data.Data, &category)
		failOnError(err, "Failed to unmarshal category")

		log.Println("category ", category)

		categoryJson, err := json.Marshal(category)
		failOnError(err, "Failed to marshal category")

		//err = json.Unmarshal(categoryJson, &category)

		_, err = controllers.CreateCategory(category)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error creating category",
				Data:    []byte(err.Error()),
			}
		} else {
			response = models.Response{
				Success: "succes",
				Message: "Category created",
				Data:    categoryJson,
			}
		}*/

		/*case "GET_TOP3POPULARPRODUCTS":
		log.Println(" [.] Getting top 3 popular products")

		products, err := controllers.GetTop3PopularProducts()
		failOnError(err, "Failed to get products")
		productsJSON, err := json.Marshal(products)
		failOnError(err, "Failed to marshal products")

		response = models.Response{
			Success: "succes",
			Message: "Products retrieved",
			Data:    productsJSON,
		}*/
	}

	responseJSON, err := json.Marshal(response)
	failOnError(err, "Failed to marshal response")

	err = ch.PublishWithContext(ctx,
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: d.CorrelationId,
			Body:          responseJSON,
		})
	failOnError(err, "Failed to publish a message")

	d.Ack(false)
}

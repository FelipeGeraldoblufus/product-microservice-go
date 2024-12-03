package controllers

import (
	"errors"
	db "github.com/FelipeGeraldoblufus/product-microservice-go/config"
	"github.com/FelipeGeraldoblufus/product-microservice-go/models"
	"gorm.io/gorm"

	"fmt"
	"math/rand"
	"time"
)



func CreateUser(username string) (*models.User, error) {
	// Crear un nuevo usuario sin el carrito (carrito ha sido eliminado)
	newUser := models.User{
		Username: username,
	}

	// Verificar si el nombre de usuario ya existe en la base de datos
	var existingUser models.User
	if err := db.DB.Where("username = ?", newUser.Username).First(&existingUser).Error; err == nil {
		// Si el usuario ya existe, devolver un error
		return nil, errors.New("username already exists")
	}

	// Guardar el nuevo usuario en la base de datos
	if err := db.DB.Save(&newUser).Error; err != nil {
		// Si ocurre un error al guardar, devolverlo
		return nil, err
	}

	// Devolver el usuario creado
	return &newUser, nil
}

func GetUser(usuario string) ([]models.User, error) {
	var user []models.User
	err := db.DB.Find(&user).Error

	return user, err
}

func GetByUser(username string) (models.User, error) {
	var users models.User
	err := db.DB.Where("username = ?", username).Find(&users).Error

	return users, err
}

// Función para obtener un producto por su ID
func GetByProductID(productID string) (models.Product, error) {
    var product models.Product

    // Buscar el producto por su product_id en la base de datos
    if err := db.DB.Where("product_id = ?", productID).First(&product).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            // Si no se encuentra el producto
            return models.Product{}, errors.New("product not found")
        }
        return models.Product{}, err
    }

    // Devolver el producto encontrado
    return product, nil
}

func GetAllProducts() ([]models.Product, error) {
	var products []models.Product

	// Consulta para obtener todos los productos
	if err := db.DB.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func UpdateProduct(productoIngresado string, newName string, newPrice int, newStock int, newDescription string, newCategory string) (models.Product, error) {
	// Inicia una transacción
	tx := db.DB.Begin()
	defer func() {
		// Recupera la transacción en caso de error y finaliza la función
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Consulta la base de datos para obtener el producto existente por su nombre
	var producto models.Product
	if err := tx.Where("name = ?", productoIngresado).First(&producto).Error; err != nil {
		tx.Rollback()
		return producto, err
	}

	// Verifica si el nombre está siendo cambiado y si existe otro producto con el mismo nombre
	if producto.Name != newName {
		var duplicateProduct models.Product
		if err := tx.Where("name = ?", newName).First(&duplicateProduct).Error; err == nil {
			// Ya existe un producto con el nuevo nombre
			tx.Rollback()
			return producto, errors.New("product with the same name already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			// Otro error al buscar el producto duplicado
			tx.Rollback()
			return producto, err
		}
	}

	// Actualiza los campos del producto existente con los nuevos valores
	if newName != "" {
		producto.Name = newName
	}
	if newPrice > 0 {
		producto.Price = newPrice
	}
	if newStock >= 0 {
		producto.Stock = newStock
	}
	if newDescription != "" {
		producto.Description = newDescription
	}
	if newCategory != "" {
		producto.Category = newCategory
	}

	// Guarda los cambios en la base de datos
	if err := tx.Save(&producto).Error; err != nil {
		// Ocurrió un error al guardar en la base de datos, realiza un rollback
		tx.Rollback()
		return producto, err
	}

	// Confirma la transacción
	tx.Commit()

	// Devuelve el producto actualizado
	return producto, nil
}

func generateProductID() string {
	// Semilla para el generador aleatorio, utilizando la hora actual para mayor unicidad
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(1000000) // Generar un número aleatorio de 6 dígitos
	return fmt.Sprintf("product-%d-%d", time.Now().Unix(), randomNumber)
}

// CreateProduct crea un nuevo producto con el nombre proporcionado
// Si el producto ya existe, devuelve un error.
func CreateProduct(name string, price int, stock int, description string, category string) (models.Product, error) {
	// Verificar si el producto ya existe en la base de datos
	var existingProduct models.Product
	if err := db.DB.Where("name = ?", name).First(&existingProduct).Error; err == nil {
		// Si el producto ya existe, devolver un error
		return models.Product{}, errors.New("product with the same name already exists")
	}

	// Validar los datos (opcional, pero recomendado)
	if price <= 0 {
		return models.Product{}, errors.New("price must be greater than zero")
	}
	if stock < 0 {
		return models.Product{}, errors.New("stock cannot be negative")
	}
	if description == "" {
		return models.Product{}, errors.New("description cannot be empty")
	}

	if category == "" {
		return models.Product{}, errors.New("description cannot be empty")
	}

	// Crear un nuevo producto
	newProduct := models.Product{
		Name:        name,
		Price:       price,
		Stock:       stock,
		Description: description,
		Category: category,
	}

	// Generar un product_id único manualmente
	newProduct.ProductID = generateProductID()

	// Iniciar una transacción
	tx := db.DB.Begin()

	// Manejo de errores de la transacción
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Intentar almacenar el nuevo producto en la base de datos
	if err := tx.Create(&newProduct).Error; err != nil {
		tx.Rollback() // Deshacer la transacción si hay un error
		return models.Product{}, err
	}

	// Confirmar la transacción si no hay errores
	tx.Commit()

	// Devolver el producto creado
	return newProduct, nil
}


func DeleteProductByName(nameProduct string) error {
	// Abre una transacción
	tx := db.DB.Begin()

	// Maneja los errores de la transacción
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Busca el producto por nombre
	var product models.Product
	if err := tx.Where("name = ?", nameProduct).First(&product).Error; err != nil {
		tx.Rollback() // Deshace la transacción en caso de error
		return err
	}

	// Elimina el producto
	if err := tx.Delete(&product).Error; err != nil {
		tx.Rollback() // Deshace la transacción en caso de error
		return err
	}

	// Confirma la transacción si no hay errores
	tx.Commit()

	return nil
}

func EditUser(currentUsername string, newUsername string) (*models.User, error) {
	// Buscar el usuario actual en la base de datos
	var existingUser models.User
	if err := db.DB.Where("username = ?", currentUsername).First(&existingUser).Error; err != nil {
		return nil, err
	}

	// Modificar el nombre de usuario
	existingUser.Username = newUsername

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingUser).Error; err != nil {
		return nil, err
	}

	// Devolver el usuario actualizado
	return &existingUser, nil
}

func DeleteUser(usuario string) error {
	// Abre una transacción
	tx := db.DB.Begin()

	// Maneja los errores de la transacción
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Busca el producto por nombre
	var user models.User
	if err := tx.Where("username = ?", usuario).First(&user).Error; err != nil {
		tx.Rollback() // Deshace la transacción en caso de error
		return err
	}

	// Elimina el usuario
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback() // Deshace la transacción en caso de error
		return err
	}

	// Confirma la transacción si no hay errores
	tx.Commit()

	return nil
}



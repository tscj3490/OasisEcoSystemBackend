package domain

import (
	validator "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type Product struct {
	Entity         `bson:",inline"`
	Name           string `json:"name" bson:"name"`
	RegistryNumber int64  `json:"registryNumber" bson:"registryNumber"`
	Holder         string `json:"holder" bson:"holder"`
	Formula        string `json:"formula" bson:"formula"`
}

//Validate validates the Product struct
func (product *Product) Validate() (bool, error) {
	//Default validation
	_, err := validator.ValidateStruct(product)
	if err != nil {
		return false, err
	}

	//Custom validations

	return true, err
}

// CreateProductByProduct creates a new product with the product values brought
func CreateProductByProduct(product *Product) (*bson.ObjectId, error) {
	product.InitializeNewData()

	if _, err := product.Validate(); err != nil {
		log.Errorf("Error in domain.CreateProductByProduct.Validate -> error: %s", err.Error())
		return nil, err
	}

	if err := DB.C("products").Insert(product); err != nil {
		log.Errorf("Error in domain.CreateProductByProduct.Insert -> error: %v", err.Error())
		return nil, err
	}

	return &product.ID, nil
}

func GetProducts() ([]Product, error) {
	result := []Product{}

	if err := DB.C("products").Find(bson.M{"active": true}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

func GetProductByID(productID bson.ObjectId) (*Product, error) {
	var product Product

	if err := DB.C("products").Find(bson.M{"active": true, "_id": productID}).One(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

func DeleteProductByID(productID bson.ObjectId) error {

	if err := DB.C("products").UpdateId(productID, bson.M{"$set": bson.M{"active": false}}); err != nil {
		return err
	}

	return nil
}

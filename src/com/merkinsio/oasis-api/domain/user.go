package domain

import (
	"com/merkinsio/oasis-api/constants"
	"time"

	validator "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// User Model that holds the user data
type User struct {
	Entity                `bson:",inline"` // Inherit all fields into this struct
	Email                 string           `bson:"email" json:"email" valid:"email,required"`
	FirstName             string           `bson:"firstName" json:"firstName" valid:"required"`
	LastName              string           `bson:"lastName" json:"lastName"`
	FullName              string           `bson:"fullName" json:"fullName"`
	Password              string           `bson:"password" json:"password,omitempty" valid:"-"`
	Role                  string           `bson:"role" json:"role" valid:"-"`
	OrganizationID        bson.ObjectId    `bson:"organizationId" json:"organizationId"`
	SocialAccounts        SocialAccounts   `bson:"socialAccounts" json:"socialAccounts"`
	JoinedBy              string           `json:"joinedBy" bson:"joinedBy"`
	TermOfServices        bool             `json:"termOfServices" bson:"termOfServices"`
	PasswordRecoveryToken string           `json:"passwordRecoveryToken" bson:"passwordRecoveryToken"`
	Phone                 string           `json:"phone" bson:"phone"`
	CooperativeID         *bson.ObjectId   `json:"cooperativeId" bson:"cooperativeId"`
	Address               Address          `json:"address" bson:"address"`
}

// User Model that holds the user data
type UserAndOwner struct {
	Owner bool `json:"owner" bson:"owner"`
	User  User `json:"user" bson:"user"`
}

//GetUserByID Retrieves a User by its ID
func GetUserByID(ID bson.ObjectId) (*User, error) {
	user := User{}

	if err := DB.C("users").Find(bson.M{"_id": ID}).Select(bson.M{"password": 0}).One(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

//CreateUserByUser Inserts a new user in the database by
// the given user data
func CreateUserByUser(user *User) (*bson.ObjectId, error) {
	user.InitializeNewData()
	password, _ := GenerateStringHashedPassword(user.Password)
	user.Password = password
	normalizedEmail, err := validator.NormalizeEmail(user.Email)
	if user.JoinedBy == "" {
		user.JoinedBy = constants.JoinedEmail
	}
	if user.Role == "" {
		user.Role = constants.RoleFarmer
	}
	if err != nil {
		log.Errorf("Error in domain.CreateUserByUser.NormalizeEmail -> error: %v", err.Error())
		return nil, err
	}
	user.Email = normalizedEmail

	fullName := user.FirstName + " " + user.LastName

	user.FullName = fullName

	if user.OrganizationID == "" {
		user.OrganizationID = bson.NewObjectId()
	}
	user.TermOfServices = false

	// Validate the user struct using the `valid` tags in the model before inserting in the db
	if _, err := validator.ValidateStruct(user); err != nil {
		log.Errorf("Error in domain.CreateUserByUser.ValidateStruct -> error: %s", err.Error())
		return nil, err
	}

	if err := DB.C("users").Insert(user); err != nil {
		log.Errorf("Error in domain.CreateUserByUser.Insert -> error: %s", err.Error())
		return nil, err
	}

	return &user.ID, nil
}

//GetUserByEmail Retrieves the user from the db by the given email
func GetUserByEmail(email string) (*User, error) {
	user := User{}
	normalized, _ := validator.NormalizeEmail(email)

	if err := DB.C("users").Find(bson.M{"email": normalized}).Select(bson.M{"password": 0}).One(&user); err != nil {
		log.Errorf("Error in domain.GetUserByEmail -> error: %v", err.Error())
		return nil, err
	}

	return &user, nil
}

//GetUserByToken Retrieves the user from the db by the given token
func GetUserByToken(token string) (*User, error) {
	user := User{}

	if err := DB.C("users").Find(bson.M{"passwordRecoveryToken": token}).Select(bson.M{"password": 0}).One(&user); err != nil {
		log.Errorf("Error in domain.GetUserByToken -> error: %v", err.Error())
		return nil, err
	}

	return &user, nil
}

//GetUserByGoogleID Gets the user by its google user id
func GetUserByGoogleID(googleID string) (*User, error) {
	var user User
	if err := DB.C("users").Find(bson.M{
		"socialAccounts.google.id": googleID,
	}).One(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserGoogleSocialDataByEmail Updates the user social google data
func UpdateUserGoogleSocialDataByEmail(email string, google SocialAccount) error {

	if err := DB.C("users").Update(bson.M{"email": email}, bson.M{
		"$set": bson.M{
			"socialAccounts.google": google,
			"updatedAt":             time.Now().UTC(),
		},
	}); err != nil {
		return err
	}

	return nil
}

// UpdatePasswordRecoveryByToken Updates the user social google data
func UpdatePasswordRecoveryByToken(token string, password string) error {

	if err := DB.C("users").Update(bson.M{"passwordRecoveryToken": token}, bson.M{
		"$set": bson.M{
			"password":              password,
			"passwordRecoveryToken": "",
		},
	}); err != nil {
		return err
	}

	return nil
}

// UpdateUserFacebookSocialDataByEmail Updates the user social facebook data
func UpdateUserFacebookSocialDataByEmail(email string, facebook SocialAccount) error {

	if err := DB.C("users").Update(bson.M{"email": email}, bson.M{
		"$set": bson.M{
			"socialAccounts.facebook": facebook,
			"updatedAt":               time.Now().UTC(),
		},
	}); err != nil {
		return err
	}

	return nil
}

// TermOfServicesAcceptedByUserID
func TermOfServicesAcceptedByUserID(user User) (*User, error) {

	if err := DB.C("users").Update(bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"termOfServices": user.TermOfServices,
		},
	}); err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserByIdWithRecoveryToken
func UpdateUserByIdWithRecoveryToken(user User, token string) (*User, error) {

	if err := DB.C("users").Update(bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"passwordRecoveryToken": token,
		},
	}); err != nil {
		return nil, err
	}

	return &user, nil
}

//UpdateUserByUser updates passed user with its values
func UpdateUserByUser(newUser User) error {
	var user User

	// Update NewUser Data (If the new data has an empty field, it will be completed with the old stored data)
	if err := DB.C("users").Find(bson.M{"_id": newUser.ID}).One(&user); err != nil {
		return err
	}
	if newUser.FirstName != "" {
		user.FirstName = newUser.FirstName
	}
	if newUser.LastName != "" {
		user.LastName = newUser.LastName
	}
	user.FullName = user.FirstName + " " + user.LastName
	if newUser.Phone != "" {
		user.Phone = newUser.Phone
	}
	if newUser.Email != "" {
		user.Email = newUser.Email
	}
	if newUser.Password != "" {
		newPassword, _ := GenerateStringHashedPassword(newUser.Password)
		user.Password = newPassword
	}
	if newUser.CooperativeID != nil {
		user.CooperativeID = newUser.CooperativeID
	}
	user.Role = newUser.Role
	user.UpdatedAt = time.Now().UTC()

	if err := DB.C("users").UpdateId(user.ID, user); err != nil {
		if err.Error() != "not found" {
			log.Errorf("Error in domain.UpdateUserByID -> error: %v", err.Error())
		}
		return err
	}

	return nil
}

// GetEmployeesByCooperativeID Returns the users wich role is not farmer
func GetEmployeesByCooperativeID(cooperativeID *bson.ObjectId) ([]User, error) {
	result := []User{}

	if err := DB.C("users").Find(bson.M{
		"cooperativeId": cooperativeID,
		"active":        true,
		"$and": []bson.M{
			bson.M{"role": bson.M{"$ne": constants.RoleFarmer}},
			bson.M{"role": bson.M{"$ne": constants.RoleAdmin}},
		},
	}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

//DeleteUserByID Delete a User by its ID
func DeleteUserByID(userID bson.ObjectId) error {

	if err := DB.C("users").UpdateId(userID, bson.M{"$set": bson.M{"active": false}}); err != nil {
		return err
	}

	return nil
}

func GetClientsByCooperativeID(cooperativeID *bson.ObjectId) ([]User, error) {
	result := []User{}
	if err := DB.C("users").Find(bson.M{
		"cooperativeId": cooperativeID,
		"active":        true,
		"$and": []bson.M{
			bson.M{"role": bson.M{"$ne": constants.RoleFarmer}},
		},
	}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

// GetAllEmployees Returns the users wich role is not farmer
func GetAllEmployees() ([]User, error) {
	result := []User{}

	if err := DB.C("users").Find(bson.M{
		"active": true,
		"$and": []bson.M{
			bson.M{"role": bson.M{"$ne": constants.RoleFarmer}},
			//bson.M{"role": bson.M{"$ne": constants.RoleAdmin}},
		},
	}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

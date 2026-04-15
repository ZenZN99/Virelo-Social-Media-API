package services

import (
	"backend/models"
	"backend/utils"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userModel    *mongo.Collection
	tokenService *utils.TokenService
	cloudinaryService   *utils.CloudinaryService
}

func NewUserService(userModel *mongo.Collection, tokenService *utils.TokenService, cloudinaryService *utils.CloudinaryService) *UserService {
	return &UserService{
		userModel:    userModel,
		tokenService: tokenService,
		cloudinaryService:   cloudinaryService,
	}
}

func (s *UserService) SignUp(ctx context.Context, data models.UserModel) (map[string]interface{}, error) {
	switch {
	case data.FullName == "" || data.Email == "" || data.Password == "":
		return nil, errors.New("all fields are required")
	case !utils.IsValidEmail(data.Email):
		return nil, errors.New("Invalid Email address")
	case len(data.Password) < 8:
		return nil, errors.New("password must be at least 8 characters")

	case len(data.Password) > 40:
		return nil, errors.New("the password must be at charset 40")
	}

	var existing models.UserModel

	err := s.userModel.FindOne(ctx, bson.M{"email": data.Email}).Decode(&existing)

	if err == nil {
		return nil, errors.New("Email already registred")
	}

	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(data.Password), 12)

	if err != nil {
		return nil, err
	}

	newUser := models.UserModel{
		ID:       primitive.NewObjectID(),
		FullName: data.FullName,
		Email:    data.Email,
		Password: string(hashed),
		Role:     string(utils.GetRole()),
		Avatar: models.Avatar{
			URL:      "https://res.cloudinary.com/dgagbheuj/image/upload/v1763194734/avatar-default-image_yc4xy4.jpg",
			PublicID: "",
		},
		Bio:       "No bio yet",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	token, err := s.tokenService.GenerateToken(utils.TokenPayload{
		UserID: newUser.ID.Hex(),
		Role:   newUser.Role,
	})

	if err != nil {
		return nil, err
	}

	_, err = s.userModel.InsertOne(ctx, newUser)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": "Account created successfully",
		"user":    newUser,
		"token":   token,
	}, nil

}

func (s *UserService) Login(ctx context.Context, data models.UserModel) (map[string]interface{}, error) {
	if data.Email == "" || data.Password == "" {
		return nil, errors.New("all fields are required")
	}

	var user models.UserModel

	err := s.userModel.FindOne(ctx, bson.M{"email": data.Email}).Decode(&user)

	if err != nil {
		return nil, errors.New("Incorrect password or email address")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))

	if err != nil {
		return nil, errors.New("Incorrect password or email address")
	}

	token, err := s.tokenService.GenerateToken(utils.TokenPayload{
		UserID: user.ID.Hex(),
		Role:   string(user.Role),
	})

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": "Logged  successfully",
		"user":    user,
		"token":   token,
	}, nil
}

func Logout(c *gin.Context) {

	c.SetCookie(
		"token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(200, gin.H{
		"success": "Logged out successfully",
	})
}

func (s *UserService) Me(ctx context.Context, userId string) (models.UserModel, error) {

	var user models.UserModel

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return models.UserModel{}, errors.New("invalid user id")
	}

	err = s.userModel.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)

	if err != nil {
		return models.UserModel{}, errors.New("User not found!")
	}

	return user, nil
}

func (s *UserService) UpdateProfile(
	ctx context.Context,
	userID string,
	avatarPath string,
	bio string,
) (map[string]interface{}, error) {

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	var user models.UserModel

	err = s.userModel.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)

	if err != nil {
		return nil, errors.New("user not found")
	}

	updatedData := bson.M{}

	if avatarPath != "" {

		if strings.Contains(user.Avatar.URL, "res.cloudinary.com") {

			parts := strings.Split(user.Avatar.URL, "/")
			file := parts[len(parts)-1]
			publicID := strings.Split(file, ".")[0]

			if publicID != "" {
				_ = s.cloudinaryService.Delete("users/avatars/" + publicID)
			}
		}

		uploadURL, err := s.cloudinaryService.UploadFile(
			ctx,
			avatarPath,
			"users/avatars",
			"image/jpeg",
		)

		if err != nil {
			return nil, err
		}

		updatedData["avatar"] = uploadURL
	}

	if bio != "" {
		updatedData["bio"] = bio
	}

	_, err = s.userModel.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updatedData},
	)

	if err != nil {
		return nil, err
	}

	err = s.userModel.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)

	if err != nil {
		return nil, errors.New("user not found update")
	}

	return map[string]interface{}{
		"success": "Profile updated successfully",
		"user":    user,
	}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, userID string) ([]models.UserModel, error) {

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	filter := bson.M{
		"_id": bson.M{
			"$ne": objectID,
		},
	}

	cursor, err := s.userModel.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.UserModel

	for cursor.Next(ctx) {
		var user models.UserModel
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetUserById(ctx context.Context, id string) (models.UserModel, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.UserModel{}, errors.New("invalid user id")
	}

	var user models.UserModel

	err = s.userModel.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return models.UserModel{}, errors.New("user not found")
	}

	return user, nil
}

func (s *UserService) DeleteUserById(ctx context.Context, id string) (map[string]interface{}, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	var user models.UserModel

	err = s.userModel.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Avatar.URL != "" && strings.Contains(user.Avatar.URL, "res.cloudinary.com") {

		parts := strings.Split(user.Avatar.URL, "/")
		file := parts[len(parts)-1]
		publicID := strings.Split(file, ".")[0]

		if strings.Contains(user.Avatar.URL, "avatar-default-image") {
			_ = s.cloudinaryService.Delete("users/avatars/" + publicID)
		}
	}

	_, err = s.userModel.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": "user deleted successfully",
	}, nil
}

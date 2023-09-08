package controller

import (
	"database/sql"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/logger"
	"github.com/Enthreeka/go-bot-storage/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type userController struct {
	userRepo repository.User

	log *logger.Logger
}

func NewUserController(userRepo repository.User, log *logger.Logger) User {
	return &userController{
		userRepo: userRepo,
		log:      log,
	}
}

// CheckUser it`s a function for getting data about user from storage, if user not exist
// in error handler we create a new user
func (u *userController) CheckUser(update *tgbotapi.Update) (*model.User, error) {
	userID := update.Message.Chat.ID

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			userModel := &model.User{
				ID:        userID,
				Nickname:  update.Message.From.UserName,
				FirstName: update.Message.From.FirstName,
				LastName:  update.Message.From.LastName,
			}

			user, err = u.createUser(userModel)
			if err != nil {
				return nil, err
			}

			return user, nil
		}
		return nil, err
	}

	return user, nil
}

func (u *userController) createUser(user *model.User) (*model.User, error) {
	createdUser, err := u.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	u.log.Info("the user has been successfully created")
	return createdUser, nil
}

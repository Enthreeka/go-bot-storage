package controller

import (
	"context"
	"database/sql"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/logger"
	"github.com/Enthreeka/go-bot-storage/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

type userController struct {
	userRepoSqlite repository.User
	userRepoRedis  repository.User

	log *logger.Logger
}

func NewUserController(userRepo repository.User, userRepoRedis repository.User, log *logger.Logger) User {
	return &userController{
		userRepoSqlite: userRepo,
		userRepoRedis:  userRepoRedis,
		log:            log,
	}
}

// CheckUser it`s a function for getting data about user from storage, if user not exist
// in error handler we create a new user
func (u *userController) CheckUser(update *tgbotapi.Update) (*model.User, error) {
	userID := update.Message.Chat.ID

	// The first step. We are looking for user in redis
	user, err := u.userRepoRedis.GetByID(context.Background(), userID)
	// The second step. If redis has nil, we fall into the check another db
	if err == redis.Nil {

		// The third step. We are looking for users in sqlite/postgres
		user, err := u.userRepoSqlite.GetByID(context.Background(), userID)
		// If we have an error we must check this error on ErrNoRows
		if err != nil {
			// The fourth step. Checking ErrNoRows in general database
			if err == sql.ErrNoRows {
				userModel := &model.User{
					ID:        userID,
					Nickname:  update.Message.From.UserName,
					FirstName: update.Message.From.FirstName,
					LastName:  update.Message.From.LastName,
				}

				// The fifth step. Creating user in redis and general database
				user, err = u.createUser(userModel)
				if err != nil {
					u.log.Error("failed to create new user: %v", err)
					return nil, err
				}

				u.log.Info("Create new user - id:[%d],username:[%s]", userID, update.Message.From.UserName)
				return user, nil
			}
			u.log.Error("failed to get by id: %v", err)
			return nil, err
		}

		// The fourth step in case sqlite/postgres error == nil. Creating user only in redis
		// with getting user value from general database
		_, err = u.userRepoRedis.Create(context.Background(), user)
		if err != nil {
			return nil, NewErrController("redis", err)
		}
		// Return user from general database
		return user, err
	}

	// If user exist in redis we return his value
	return user, nil
}

func (u *userController) createUser(user *model.User) (*model.User, error) {
	createdUser, err := u.userRepoSqlite.Create(context.Background(), user)
	if err != nil {
		return nil, NewErrController("sqlite", err)
	}
	_, err = u.userRepoRedis.Create(context.Background(), user)
	if err != nil {
		return nil, NewErrController("redis", err)
	}

	return createdUser, nil
}

package database

import (
	"errors"

	"github.com/google/uuid"
	"github.com/marktrs/simple-todo/logger"
	"github.com/marktrs/simple-todo/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedUserData() error {
	log := logger.Log
	log.Info().Msg("seeding user data...")

	// Create users data
	users := []model.User{
		{
			Username: "test01",
			Password: "1111",
		},
		{
			Username: "test02",
			Password: "2222",
		},
		{
			Username: "test03",
			Password: "3333",
		},
		{
			Username: "test04",
			Password: "4444",
		},
	}

	for _, user := range users {
		user := user // prevent implicit memory aliasing in for loop
		bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if err != nil {
			log.Fatal().AnErr("error", err).Msg("failed to generate password for seed user")
		}

		user.ID = uuid.New().String()
		user.Password = string(bytes)

		if err := DB.FirstOrCreate(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				log.Warn().Msg("users already exists, skipping...")
				break
			}
			log.Fatal().AnErr("error", err).Msg("failed to seed user data")
			return err
		}
	}

	return nil
}

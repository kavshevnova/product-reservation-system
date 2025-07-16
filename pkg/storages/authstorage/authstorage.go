package authstorage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-redis/redis/v8"
	"github.com/kavshevnova/product-reservation-system/pkg/domain/models"
	"time"
)

type StorageUsers struct {
	client *redis.Client
}

func NewUsersStorage(addr, password string, db int) (*StorageUsers, error) {
	const op = "storages.NewUsersStorage"

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	//Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%s %s", op, err)
	}
	return &StorageUsers{client: client}, nil
}

func (s *StorageUsers) SaveUser(ctx context.Context, email string, passhash []byte) (uid int64, err error) {
	const op = "storages.authstorage.SaveUser"

	//Проверяем существование пользователя
	if exists, err := s.client.Exists(ctx, "user:email:"+email).Result(); err != nil {
		return 0, fmt.Errorf("%s %s", op, err)
	} else if exists == 1 {
		return 0, models.ErrUserExists
	}

	//Формируем id
	uid, err = s.client.Incr(ctx, "global:user:id").Result()
	if err != nil {
		return 0, fmt.Errorf("%s %s", op, err)
	}

	//Сохраняем данные пользователя
	userData := map[string]interface{}{
		"id":       uid,
		"email":    email,
		"passhash": passhash,
	}

	//Используем транзакцию для атомарности
	_, err = s.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(ctx, fmt.Sprintf("user:%d", uid), userData) //создаем хэш в редди с ключом id и сохраняем данные пользователя в виде полей хэша
		pipe.Set(ctx, "user:email:"+email, uid, 0)             //создаем для поиска по email
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("%s %s", op, err)
	}
	return uid, nil
}

func (s *StorageUsers) User(ctx context.Context, email string) (models.User, error) {
	const op = "storages.authstorage.User"

	//Получаем id пользователя по email
	uid, err := s.client.Get(ctx, "user:email:"+email).Int64()
	if err != nil {
		if err == redis.Nil {
			return models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s %s", op, err)
	}

	result, err := s.client.HGetAll(ctx, fmt.Sprintf("user:%d", uid)).Result()
	if err != nil {
		return models.User{}, fmt.Errorf("%s %s", op, err)
	}
	if len(result) == 0 {
		return models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
	}

	user := models.User{
		UserID:   uid,
		Email:    result["email"],
		Passhash: []byte(result["passhash"]),
	}
	return user, nil
}

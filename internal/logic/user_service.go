package logic

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/namnv2496/go-wallet/internal/databaseaccess"
	db "github.com/namnv2496/go-wallet/internal/databaseaccess/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type UserLogic interface {
	CreateUser(ctx context.Context, userName string, password string, fullname string, email string, phone string, role string) (db.User, error)
	GetUser(ctx context.Context, username string) (db.User, error)
	GetUsersByUsernameOrPhone(ctx context.Context, Username string, phone string, limit int32) ([]db.GetUsersByUsernameOrPhoneRow, error)
	UpdateUser(ctx context.Context, password string, fullname string, email string, phone string, username string) (db.User, error)
	VerifyEmail(ctx context.Context, username string) error
}

var _ UserLogic = (*userLogic)(nil)

type userLogic struct {
	database *databaseaccess.Database
}

func NewUserLogic(
	database *databaseaccess.Database,
) (UserLogic, error) {
	return &userLogic{
		database: database,
	}, nil
}

func (u userLogic) GetUser(ctx context.Context, username string) (db.User, error) {

	return u.database.GetUser(ctx, pgtype.Text{
		String: username,
		Valid:  true,
	})
}

func (u userLogic) GetUsersByUsernameOrPhone(
	ctx context.Context,
	username string,
	phone string,
	limit int32,
) ([]db.GetUsersByUsernameOrPhoneRow, error) {

	user := formatUsername(&username)
	phoneFormatted := formatPhone(&phone)

	var usernameField pgtype.Text
	if user != "" {
		usernameField = pgtype.Text{
			String: user,
			Valid:  true,
		}
	} else {
		usernameField = pgtype.Text{
			Valid: false,
		}
	}
	var arg db.GetUsersByUsernameOrPhoneParams
	if phoneFormatted != "" {
		arg = db.GetUsersByUsernameOrPhoneParams{
			Username: usernameField,
			Phone:    phoneFormatted,
			Limit:    limit,
		}
	} else {
		arg = db.GetUsersByUsernameOrPhoneParams{
			Username: usernameField,
			Limit:    limit,
		}
	}
	log.Println(arg)
	return u.database.GetUsersByUsernameOrPhone(ctx, arg)
}

func (u userLogic) CreateUser(
	ctx context.Context,
	userName string,
	password string,
	fullname string,
	email string,
	phone string,
	role string,
) (db.User, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, err
	}
	arg := db.CreateUserParams{
		Username: pgtype.Text{
			String: userName,
			Valid:  true,
		},
		HashedPassword:  string(hashPassword),
		FullName:        fullname,
		Email:           email,
		Phone:           phone,
		IsEmailVerified: false,
	}
	return u.database.CreateUser(ctx, arg)
}

func (u userLogic) UpdateUser(ctx context.Context, password string, fullname string, email string, phone string, username string) (db.User, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, err
	}
	arg := db.UpdateUserParams{
		HashedPassword: pgtype.Text{
			String: string(hashPassword),
			Valid:  true,
		},
		PasswordChangedAt: pgtype.Timestamptz{
			Time: time.Now(),
		},
		FullName: pgtype.Text{
			String: fullname,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: email,
			Valid:  true,
		},
		Phone: pgtype.Text{
			String: phone,
			Valid:  true,
		},
		Username: pgtype.Text{
			String: username,
			Valid:  true,
		},
	}
	return u.database.UpdateUser(ctx, arg)
}

func (u userLogic) VerifyEmail(ctx context.Context, username string) error {

	arg := db.VerifyEmailParams{
		Username: pgtype.Text{
			String: username,
			Valid:  true,
		},
		IsEmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
	}
	_, err := u.database.VerifyEmail(ctx, arg)
	if err != nil {
		return err
	}
	return nil
}

func formatPhone(phone *string) string {
	if phone != nil && *phone != "" {
		return "%" + *phone + "%"
	}
	return ""
}

func formatUsername(username *string) string {
	if username != nil && *username != "" {
		return "%" + *username + "%"
	}
	return ""
}

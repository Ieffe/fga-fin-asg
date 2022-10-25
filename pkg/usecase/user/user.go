package user

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"fin-asg/pkg/domain/claim"
	"fin-asg/pkg/domain/message"
	"fin-asg/pkg/domain/user"
	"fin-asg/pkg/usecase/crypto"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type UserUsecaseImpl struct {
	userRepo user.UserRepo
}

func NewUserUsecase(userRepo user.UserRepo) user.UserUseCase {
	return &UserUsecaseImpl{userRepo: userRepo}
}

func (u *UserUsecaseImpl) GetUserByIdSvc(ctx context.Context, id uint64) (result user.User, errMsg message.ErrorMessage) {
	log.Printf("%T - GetUserById is invoked]\n", u)
	defer log.Printf("%T - GetUserById executed\n", u)

	// get user from repository (database)
	log.Println("getting user from user repository")
	result, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		// ini berarti ada yang salah dengan connection di database
		log.Println("error when fetching data from database: " + err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type: "INTERNAL ERROR",
		}
		return result, errMsg
	}
	// check user id > 0 ?
	log.Println("checking user id")
	if result.ID <= 0 {
		// kalau tidak berarti user not found
		log.Println("user is not found: " + strconv.FormatUint(id,10))

		
		errMsg = message.ErrorMessage{
			Error: err,
			Type: "NOT_FOUND",
		}
		return result, errMsg
	}
	return result, errMsg
}

func (u *UserUsecaseImpl) InsertUserSvc(ctx context.Context, input user.User) (result user.User, errMsg message.ErrorMessage) {
	log.Printf("%T - InsertUserSvc is invoked\n", u)
	defer log.Printf("%T - InsertUserSvc executed\n", u)

	// user input validation
	if isValid, err := govalidator.ValidateStruct(input); !isValid {
		switch err.Error() {
		case "username is required":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "USERNAME_IS_EMPTY",
			}
			return result, errMsg
		case "email is required":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "EMAIL_IS_EMPTY",
			}
			return result, errMsg
		case "invalid email format":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "WRONG_EMAIL_FORMAT",
			}
			return result, errMsg
		case "password is required":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "PASSWORD_IS_EMPTY",
			}
			return result, errMsg
		case "password has to have a minimum length of 6 characters":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "INVALID_PASSWORD_FORMAT",
			}
			return result, errMsg

		default:
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "INVALID_PAYLOAD",
			}
			return result, errMsg
		}
	}

	log.Println("calling register user repo")
	err := u.userRepo.InsertUser(ctx, &input)
	if err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "idx_users_username"`) {
			err = errors.New("username has already been registered")
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "USERNAME_REGISTERED",
			}
			return result, errMsg
		}

		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "idx_users_email"`) {
			err = errors.New("email has already been registered")
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "EMAIL_REGISTERED",
			}

			return result, errMsg
		}

	}

	if err != nil {
		log.Printf("error when fetching data from database: %s\n", err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "INTERNAL_CONNECTION_PROBLEM",
		}
		return result, errMsg
	}

	return input, errMsg
}

func (u *UserUsecaseImpl) UpdateUserSvc(ctx context.Context, id uint64, email string, username string) (idToken string, errMsg message.ErrorMessage) {
	log.Printf("%T - UpdateUserSvc is invoked\n", u)
	defer log.Printf("%T - UpdateUserSvc executed\n", u)

	if email == "" {
		errMsg := message.ErrorMessage{
			Error: errors.New("email is required"),
			Type:  "EMAIL_IS_EMPTY",
		}
		return idToken, errMsg
	}

	// email validation
	if !govalidator.IsEmail(email) {
		errMsg := message.ErrorMessage{
			Error: errors.New("invalid email format"),
			Type:  "WRONG_EMAIL_FORMAT",
		}
		return idToken, errMsg
	}

	// username validation
	if username == "" {
		errMsg := message.ErrorMessage{
			Error: errors.New("username is required"),
			Type:  "USERNAME_IS_EMPTY",
		}
		return idToken, errMsg
	}

	result, err := u.userRepo.UpdateUser(ctx, id, email, username)

	if err != nil {
		log.Printf("error when fetching data from database: %s\n", err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "INTERNAL_CONNECTION_PROBLEM",
		}
		return idToken, errMsg
	}

	claimId := claim.ID{
		JWTID:    uuid.New(),
		Username: result.Username,
		Email:    result.Email,
		DOB:      time.Time(result.Dob),
	}
	idToken, _ = crypto.CreateJWT(ctx, claimId)

	return idToken, errMsg
	
}

func (u *UserUsecaseImpl) DeleteUserSvc(ctx context.Context, id uint64) ( errMsg message.ErrorMessage) {
	log.Printf("%T - UpdateUserSvc is invoked\n", u)
	defer log.Printf("%T - UpdateUserSvc executed\n", u)

	log.Println("calling delete user repo")
	err := u.userRepo.DeleteUser(ctx, id)

	if err != nil {
		log.Printf("error when fetching data from database: %s\n", err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "INTERNAL_CONNECTION_PROBLEM",
		}
		return errMsg
	}

	return errMsg
}
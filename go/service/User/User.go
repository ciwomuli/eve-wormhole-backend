package user

import (
	"eve-wormhole-backend/go/dao"
	"eve-wormhole-backend/go/entity"
	"eve-wormhole-backend/go/service/esi"
	"eve-wormhole-backend/go/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func Callback(c *gin.Context) (map[string]any, error) {
	token, err := esi.EveSSOCallback(c)

	if err != nil {
		return nil, utils.WrapError(err)
	}

	if err != nil {
		return nil, utils.WrapError(err)
	}
	v, err := esi.E.SSO.Verify(esi.E.SSO.TokenSource(token))
	if err != nil {
		return nil, utils.WrapError(err)
	}

	account, err := GetUserAccountByCharacterID(v.CharacterID)
	userId, exist := c.Get("userId")
	logrus.Debugf("User ID from context: %v", userId)
	if exist && userId.(uint) != 0 {
		if err == gorm.ErrRecordNotFound {
			account := &entity.UserAccount{
				UserID:        userId.(uint),
				CharacterID:   v.CharacterID,
				CharacterName: v.CharacterName,
				AccessToken:   token.AccessToken,
				RefreshToken:  token.RefreshToken,
				Expiry:        token.Expiry,
			}
			err := SaveUserAccount(account)
			if err != nil {
				return nil, utils.WrapError(err)
			}
			return gin.H{"token": utils.GenerateJWT(userId.(uint), "")}, nil
		} else if err != nil {
			return nil, utils.WrapError(err)
		} else {
			account.UserID = userId.(uint)
			account.CharacterID = v.CharacterID
			account.CharacterName = v.CharacterName
			account.AccessToken = token.AccessToken
			account.RefreshToken = token.RefreshToken
			account.Expiry = token.Expiry
			err := SaveUserAccount(account)
			if err != nil {
				return nil, utils.WrapError(err)
			}
			return gin.H{"token": utils.GenerateJWT(userId.(uint), "")}, nil
		}
	} else {
		if err == gorm.ErrRecordNotFound {
			user := &entity.User{
				Name: v.CharacterName,
			}
			err := SaveUser(user)
			if err != nil {
				return nil, utils.WrapError(err)
			}
			account := &entity.UserAccount{
				UserID:        user.ID,
				CharacterID:   v.CharacterID,
				CharacterName: v.CharacterName,
				AccessToken:   token.AccessToken,
				RefreshToken:  token.RefreshToken,
				Expiry:        token.Expiry,
			}
			err = SaveUserAccount(account)
			if err != nil {
				return nil, utils.WrapError(err)
			}
			return gin.H{"token": utils.GenerateJWT(user.ID, "")}, nil
		} else if err != nil {
			return nil, utils.WrapError(err)
		} else {
			return gin.H{"token": utils.GenerateJWT(account.UserID, "")}, nil
		}
	}
}

func SaveUser(user *entity.User) error {
	if err := dao.SqlSession.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func SaveUserAccount(account *entity.UserAccount) error {
	if err := dao.SqlSession.Save(account).Error; err != nil {
		return err
	}
	return nil
}

func GetUserAccountByCharacterID(characterID int32) (*entity.UserAccount, error) {
	var account entity.UserAccount
	if err := dao.SqlSession.Where("character_id = ?", characterID).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func GetUserByID(userID uint) (*entity.User, error) {
	var user entity.User
	if err := dao.SqlSession.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByIDWithAccounts(userID uint) (*entity.User, error) {
	var user entity.User
	if err := dao.SqlSession.Preload("Accounts").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllUserAccounts() ([]entity.UserAccount, error) {
	var accounts []entity.UserAccount
	if err := dao.SqlSession.Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func Token(u *entity.UserAccount) oauth2.TokenSource {
	token := &oauth2.Token{
		AccessToken:  u.AccessToken,
		RefreshToken: u.RefreshToken,
		Expiry:       u.Expiry,
	}
	tokSrc := esi.E.SSO.TokenSource(token)
	token, err := tokSrc.Token()
	if err != nil {
		return nil
	}
	u.AccessToken = token.AccessToken
	u.RefreshToken = token.RefreshToken
	u.Expiry = token.Expiry
	SaveUserAccount(u)
	return tokSrc
}

func UpdateUserActiveTimebyID(userId uint) error {
	if err := dao.SqlSession.Model(&entity.User{}).Where("id = ?", userId).Update("active_time", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

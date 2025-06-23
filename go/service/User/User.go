package User

import (
	"eve-wormhole-backend/go/dao"
	entity "eve-wormhole-backend/go/entity/User"
	"eve-wormhole-backend/go/service/ESI"
	"eve-wormhole-backend/go/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// wrapError 包装错误并添加行号信息

func Callback(c *gin.Context) error {
	s := sessions.Default(c)

	token, err := ESI.EveSSOCallback(c)
	if err != nil {
		if s.Get("login") == true {
			return utils.WrapError(err)
		} else {
			return utils.WrapError(err)
		}
	}
	v, err := ESI.E.SSO.Verify(ESI.E.SSO.TokenSource(token))
	if err != nil {
		if s.Get("login") == true {
			return utils.WrapError(err)
		} else {
			return utils.WrapError(err)
		}
	}

	account, err := GetUserAccountByCharacterID(v.CharacterID)
	if s.Get("login") == true && s.Get("userId") != nil {
		userId := s.Get("userId").(uint)
		if err == gorm.ErrRecordNotFound {
			account := &entity.UserAccount{
				UserID:        userId,
				CharacterID:   v.CharacterID,
				CharacterName: v.CharacterName,
				AccessToken:   token.AccessToken,
				RefreshToken:  token.RefreshToken,
				Expiry:        token.Expiry,
			}
			err := SaveUserAccount(account)
			if err != nil {
				return utils.WrapError(err)
			}
			return nil
		} else if err != nil {
			return utils.WrapError(err)
		} else {
			account.UserID = userId
			account.CharacterID = v.CharacterID
			account.CharacterName = v.CharacterName
			account.AccessToken = token.AccessToken
			account.RefreshToken = token.RefreshToken
			account.Expiry = token.Expiry
			err := SaveUserAccount(account)
			if err != nil {
				return utils.WrapError(err)
			}
			return nil
		}
	} else {
		if err == gorm.ErrRecordNotFound {
			user := &entity.User{
				Name: v.CharacterName,
			}
			err := CreateUser(user)
			if err != nil {
				return utils.WrapError(err)
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
				return utils.WrapError(err)
			}
			s.Set("userId", user.ID)
			s.Set("login", true)
			if err := s.Save(); err != nil {
				return utils.WrapError(err)
			}
			return nil
		} else if err != nil {
			return utils.WrapError(err)
		} else {
			userId := account.UserID
			s.Set("userId", userId)
			s.Set("login", true)
			if err := s.Save(); err != nil {
				return utils.WrapError(err)
			}
			return nil
		}
	}
}

func CreateUser(user *entity.User) error {
	if err := dao.SqlSession.Create(user).Error; err != nil {
		return utils.WrapError(err)
	}
	return nil
}

func SaveUserAccount(account *entity.UserAccount) error {
	if err := dao.SqlSession.Save(account).Error; err != nil {
		return utils.WrapError(err)
	}
	return nil
}

func GetUserAccountByCharacterID(characterID int32) (*entity.UserAccount, error) {
	var account entity.UserAccount
	if err := dao.SqlSession.Where("character_id = ?", characterID).First(&account).Error; err != nil {
		return nil, utils.WrapError(err)
	}
	return &account, nil
}

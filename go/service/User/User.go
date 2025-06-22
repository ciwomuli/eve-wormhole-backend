package User

import (
	"eve-wormhole-backend/go/dao"
	entity "eve-wormhole-backend/go/entity/User"
	"eve-wormhole-backend/go/service/ESI"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Callback(c *gin.Context) (string, error) {
	s := sessions.Default(c)

	token, err := ESI.EveSSOCallback(c)
	if err != nil {
		if s.Get("login") == true {
			return "/esi", err
		} else {
			return "/login", err
		}
	}
	v, err := ESI.E.SSO.Verify(ESI.E.SSO.TokenSource(token))
	if err != nil {
		if s.Get("login") == true {
			return "/esi", err
		} else {
			return "/login", err
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
				return "/esi", err
			}
			return "/esi", nil
		} else if err != nil {
			return "/esi", err
		} else {
			account.UserID = userId
			account.CharacterID = v.CharacterID
			account.CharacterName = v.CharacterName
			account.AccessToken = token.AccessToken
			account.RefreshToken = token.RefreshToken
			account.Expiry = token.Expiry
			err := SaveUserAccount(account)
			if err != nil {
				return "/esi", err
			}
			return "/esi", nil
		}
	} else {
		if err == gorm.ErrRecordNotFound {
			user := &entity.User{
				Name: v.CharacterName,
			}
			err := CreateUser(user)
			if err != nil {
				return "/login", err
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
				return "/login", err
			}
			s.Set("userId", user.ID)
			s.Set("login", true)
			if err := s.Save(); err != nil {
				return "/login", err
			}
			return "/home", nil
		} else if err != nil {
			return "/login", err
		} else {
			userId := account.UserID
			s.Set("userId", userId)
			s.Set("login", true)
			if err := s.Save(); err != nil {
				return "/login", err
			}
			return "/home", nil
		}
	}
}

func CreateUser(user *entity.User) error {
	if err := dao.SqlSession.Create(user).Error; err != nil {
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

package user

import (
	"time"

	"os"
	"github.com/dgrijalva/jwt-go"
	"github.com/iffigues/musicroom/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"fmt"
)


type TokenDetails struct {
  AccessToken  string
  RefreshToken string
  AccessUuid   string
  RefreshUuid  string
  AtExpires    int64
  RtExpires    int64
}

func ExtractToken(r *http.Request) string {
  bearToken := r.Header.Get("Authorization")
  strArr := strings.Split(bearToken, " ")
  if len(strArr) == 2 {
     return strArr[1]
  }
  return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
  tokenString := ExtractToken(r)
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
     if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
     }
     return []byte(os.Getenv("ACCESS_SECRET")), nil
  })
  if err != nil {
     return nil, err
  }
  return token, nil
}

func TokenValid(r *http.Request) error {
  token, err := VerifyToken(r)
  fmt.Println(token)
  if err != nil {
     return err
  }
  if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
     return err
  }
  return nil
}


func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
  token, err := VerifyToken(r)
  fmt.Println(token)
  if err != nil {
     return nil, err
  }
  claims, ok := token.Claims.(jwt.MapClaims)
  if ok && token.Valid {
     accessUuid, ok := claims["access_uuid"].(string)
     if !ok {
        return nil, err
     }
     userId := claims["user_id"].(string)
     return &AccessDetails{
        AccessUuid: accessUuid,
        UserId:   userId,
     }, nil
  }
  return nil, err
}

func (u *UserUtils)CreateToken(user User, ha bool) (*TokenDetails, error) {

	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = util.Uid().String()
	if !ha {
		td.AtExpires = 0
	}

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = util.Uid().String()
	if !ha {
		td.RtExpires = 0
	}
	var err error
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	if !ha {
		os.Setenv("ACCESS_SECRET", "ntm")
	}
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = user.Uid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = user.Uid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}


func (u *UserUtils)CreateAuth(userid string, td *TokenDetails) (err error) {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := u.Client.Set(td.AccessUuid, userid, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := u.Client.Set(td.RefreshUuid, userid, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (u *UserUtils)DeleteAuth(givenUuid string) (int64, error) {
  deleted, err := u.Client.Del(givenUuid).Result()
  if err != nil {
     return 0, err
  }
  return deleted, nil
}

func (u *UserUtils)FetchAuth(authD *AccessDetails) (string, error) {
  userid, err := u.Client.Get(authD.AccessUuid).Result()
  if err != nil {
     return "", err
  }
  return userid, nil
}

func (u *UserUtils)Refresh(c *gin.Context) {
  mapToken := map[string]string{}
  if err := c.ShouldBindJSON(&mapToken); err != nil {
     c.JSON(http.StatusUnprocessableEntity, err.Error())
     return
  }
  refreshToken := mapToken["refresh_token"]

  //verify the token
  os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
  token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
     //Make sure that the token method conform to "SigningMethodHMAC"
     if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
     }
     return []byte(os.Getenv("REFRESH_SECRET")), nil
  })
  //if there is an error, the token must have expired
  if err != nil {
     c.JSON(http.StatusUnauthorized, "Refresh token expired")
     return
  }
  //is token valid?
  if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
     c.JSON(http.StatusUnauthorized, err)
     return
  }
  //Since token is valid, get the uuid:
  claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
  if ok && token.Valid {
     refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
     if !ok {
        c.JSON(http.StatusUnprocessableEntity, err)
        return
     }
     userId, _ := claims["user_id"]
     uu := User{
     	Uid: userId.(string),
     }
     //Delete the previous Refresh Token
     deleted, delErr := u.DeleteAuth(refreshUuid)
     if delErr != nil || deleted == 0 { //if any goes wrong
        c.JSON(http.StatusUnauthorized, "unauthorized")
        return
     }
 //Create new pairs of refresh and access tokens
     ts, createErr := u.CreateToken(uu, true)
     if  createErr != nil {
        c.JSON(http.StatusForbidden, createErr.Error())
        return
     }
 //save the tokens metadata to redis
saveErr := u.CreateAuth(userId.(string), ts)
 if saveErr != nil {
        c.JSON(http.StatusForbidden, saveErr.Error())
    return
}
 tokens := map[string]string{
       "access_token":  ts.AccessToken,
  "refresh_token": ts.RefreshToken,
}
     c.JSON(http.StatusCreated, tokens)
  } else {
     c.JSON(http.StatusUnauthorized, "refresh expired")
  }
}

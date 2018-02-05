package main

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"

	"github.com/robbert229/jwt"
	"time"
)

type User struct {
	UserId   string
	Password string
}

const (
	ValidUser = "arslan"
	ValidPass = "topcu"
	SecretKey = "secretkey"
)

func main() {

	m := martini.Classic()

	m.Use(martini.Static("static"))
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.HTML(201, "index", nil)
	})

	m.Get("/auth/:user/:pass", func(params martini.Params, r render.Render)  {
		UserId := params["user"]
		Password := params["pass"]

		if UserId == ValidUser && Password == ValidPass {

			algorithm := jwt.HmacSha256(SecretKey)

			claims := jwt.NewClaim()
			claims.Set("Role", "Admin")
			claims.SetTime("exp", time.Now().Add(time.Minute))

			tokenString, err := algorithm.Encode(claims)
			if err != nil {
				panic(err)
			}

			data := map[string]string{
				"token": tokenString,
			}

			r.JSON(200, data)
		}else{


			r.JSON(200, "error")
		}

	})

	// Check Key is ok
	m.Get("/token/:token", func(params martini.Params, r render.Render)  {

		token := params["token"]

		algorithm := jwt.HmacSha256(SecretKey)
		if algorithm.Validate(token) != nil {
			r.JSON(401, "error")
		}

		loadedClaims, err := algorithm.Decode(token)
		if err != nil {
			panic(err)
		}

		role, err := loadedClaims.Get("Role")

		data := map[string]string{
			"role": role.(string),
		}

		r.JSON(200, data)
	})

	m.RunOnAddr(":5000")
}



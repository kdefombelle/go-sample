package main

import (
	"compress/flate"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/kdefombelle/go-sample/authn"
	"github.com/kdefombelle/go-sample/db"
	"github.com/kdefombelle/go-sample/http/rest"
	"github.com/kdefombelle/go-sample/jwt"
	"github.com/kdefombelle/go-sample/logger"
	"github.com/kdefombelle/go-sample/nursery/account"
	"github.com/kdefombelle/go-sample/nursery/plant"

	"gopkg.in/yaml.v2"

	_ "github.com/kdefombelle/go-sample/api" // This line is necessary for go-swagger to find the docs
)

type config struct {
	db.Parameters `yaml:"database"`
	Origin        string `yaml:"origin"`
}

func main() {
	logger.Logger.Infof("Application starting")
	//command line options
	help := flag.Bool("help", false, "option, print usage")
	file := flag.String("f", "nursery.yml", "the configuration file")
	certificate := flag.String("c", "nursery.rsa", "the certificate file") // openssl genrsa -out nursery.rsa
	flag.Parse()
	if *help {
		flag.Usage()
	}
	logger.Logger.Infof("Configuration file [%s]", *file)
	logger.Logger.Infof("Certificate file [%s]", *certificate)
	cfg := readConfiguration(file)

	// create db connection
	logger.Logger.Infof("Create DB connection")
	var connectionFactory db.ConnectionFactory
	db, close, err := connectionFactory.Create(cfg.Parameters)
	defer close()
	if err != nil {
		log.Fatalf("error while connection to database [%s]", err)
		panic(err)
	}

	// create application services and controllers
	logger.Logger.Infof("Create application services and controllers")
	filePrivateKeyReader := &jwt.FilePrivateKeyReader{KeyPath: *certificate}
	jwtService := jwt.Service{
		PrivateKeyReader: filePrivateKeyReader,
	}

	// initialize account
	accountRepository := &account.DbRepository{
		Db: db,
	}
	authnController := rest.AuthnController{
		AccountService: &account.Service{
			AccountRepository: accountRepository,
			Encrypter:         &authn.Md5Encrypter{},
		},
		JwtService: jwtService,
	}

	// initialize plant
	plantController := rest.PlantController{
		PlantService: &plant.Service{
			PlantRepository: &plant.DbRepository{
				Db: db,
			},
		},
	}

	logger.Logger.Info("Initialise routing")
	r := chi.NewRouter()

	// create middlewares
	logger.Logger.Infof("Create middlewares")
	jwtMiddleware := rest.JwtMiddleware{
		JwtService: jwtService,
	}
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Logger)
	compressor := middleware.NewCompressor(flate.DefaultCompression)
	r.Use(compressor.Handler)

	r.Route("/login", func(r chi.Router) {
		r.Post("/", authnController.Signin)
	})

	// plant endpoints
	r.Route("/plant", func(r chi.Router) {
		r.With(jwtMiddleware.Check).Post("/", plantController.AddPlant)
		r.With(jwtMiddleware.Check).Get("/{id}", plantController.GetPlant)
	})

	logger.Logger.Info("Application started")
	http.ListenAndServe(":3000", r)
}

func readConfiguration(c *string) config {
	f, err := os.Open(*c)
	if err != nil {
		log.Fatalf("Cannot read configuration: [%v]", err)
		panic("Cannot read configuration")
	}
	defer f.Close()

	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Cannot decode configuration: [%v]", err)
		panic("Cannot decode configuration")
	}
	return cfg
}

package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"

	"g0/internal/common"
	"g0/internal/config"
	"g0/internal/database"
	httpcustom "g0/internal/http"
	"g0/internal/middleware"
	"g0/internal/routes"

	"github.com/labstack/echo/v5"
)

//go:embed dist docs/swagger.json docs/swagger.yaml
var embeddedFiles embed.FS

func main() {

	fmt.Println("DatabaseUrl:", config.GetDatabasePath())

	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	defer func(DB *sql.DB) {
		if DB == nil {
			return
		}

		if err := DB.Close(); err != nil {
			common.Logger.LogError().
				Err(err).
				Msg("Erro ao fechar banco de dados")
		}
	}(database.DB)

	common.NewLogger()

	e := echo.New()

	e.HTTPErrorHandler = httpcustom.HTTPErrorHandler

	e.Use(middleware.LoggingMiddleware)

	//API
	routes.Route(e)

	//Swagger
	e.GET("/openapi/swagger.json", func(c *echo.Context) error {
		data, err := embeddedFiles.ReadFile("docs/swagger.json")
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				"swagger.json não encontrado",
			)
		}

		return c.Blob(
			http.StatusOK,
			"application/json; charset=utf-8",
			data,
		)
	})

	e.GET("/openapi/swagger.yaml", func(c *echo.Context) error {
		data, err := embeddedFiles.ReadFile("docs/swagger.yaml")
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				"swagger.yaml não encontrado",
			)
		}

		return c.Blob(
			http.StatusOK,
			"application/yaml; charset=utf-8",
			data,
		)
	})

	// FrontEnd
	distFS := echo.MustSubFS(
		embeddedFiles,
		"dist",
	)

	e.StaticFS(
		"/assets/",
		echo.MustSubFS(distFS, "assets"),
	)

	e.FileFS(
		"/darkIcon.png",
		"darkIcon.png",
		distFS,
	)

	e.FileFS(
		"/lightIcon.png",
		"lightIcon.png",
		distFS,
	)

	e.FileFS(
		"/favicon.ico",
		"favicon.ico",
		distFS,
	)

	e.FileFS(
		"/",
		"index.html",
		distFS,
	)

	e.GET("/*", func(c *echo.Context) error {
		return c.FileFS(
			"index.html",
			distFS,
		)
	})

	//Assinatura
	asciiArt := `
                                             .......                                                                                                     
                                             .-@--......                                                                                                 
                                         ...=%-%:.@@@%*=.....                                                                                            
                                        ..@#+*=@@@%+++++=-+#.....                                                                                        
                                       .%##+%#*@#-@#%#%%%#+=@@%#%-..                   .                 ............                                    
                                      .@% #+=+@+=-+##%####%#:...-.* ...... ........=#%*=......       ....%%..=#+#@=+@...                                 
                                     .-#@#*+##%=@@@+=*#+@=+@@@@@@@@+%%#%@@@@#@@@@##%@=#%+@#%%*-.....+%@-+.*:@@=-:@*##@@-..                               
                                    .:%*==#@%@=@@.:-##@=+##@@.....+=@@*#@+%%%@+**@#*#@#@*@-+=*#+###-.+@=#@%#+=@@-@=#%*+##-.                              
                                     +@#=@@==#@+:--@#@+@%@@=%#%#@=%@*=#@+@+%#=+###-:+-##%#@@=@*@@@@%@+%@@=%@%%-%@@::##@%#@%.                             
                                       #@@++%%=+*=*%@=*=%@.=#+%@%#@**#@#:##%%+-@+#@%@@@%++#@#%#@%.*@@@%=@@=*%@##-@@@#:-#**@#-.                           
                                     =@%#@*+#@##+=@@***@@-+##.=%+@%+%%@+@%%@*#=%*@#+#-#+#%%##@.==%%++@@#+%@:*%#*#+%=@+@=*#@#@..                          
                                     :%@%%=**%-=-%*##=+@+.-.-%@@#%##*@%@@@+%#+#%@#+=#%%%#@@%%%@@+:.-#%#%%+%@.#@##+@.@%@.@#==+=-.                         
                                    .%%%#**%@+=+@@+@-*@#.-@@@%@@@@@@@+=++@##%#@@%=%@@%@@@#=#@%%#+#=@:.%%@@.@#-@%#.*@*:@#-%+@@==-                         
                                    =@%+#%%@@#..@-@==#@=@=%@%%@@#==++#@%%##*#@*-++=@#++%=+@@@@@@%@@=+#.=@@#=%-+@%++-@#*@++%@#++%.                        
                                    %##%####@@-@#@+=*@#.##@#@#*:.@%#%*++**%%%#+-****@@@*-=..#+%%@@@*%=@.=*@+@--+@##-*@+%#-*@*@#*:.                       
                                   -#%+###*#@#@#%+=+@*=@@@@#-:=%..-...=---=##@*##%@@@+:%-@@@-=.#=*%@@.-@-#@#@@+=+%%:@+#%#=#@%+#=%.                       
                                   +###%%%#@%*+=+*%%=@@@*-:-++%@%+%@@@@@@@@@%%%@%%*+%@@=#..--=+==#*%#@*=%.%@.%%-#%%%#*@#@==@+#%#@..                      
                                  .%+##%@##%+%+#+@+-@@+=-=%@%%%@@@@%@@@%*+%@@@@@@@@@%@@-@@@#-#*+==-=@+@.%:#@=@@%=#@+#%=+%%-@+*+ +.                       
                                 .-@+%#%%+@=%==@-%-@@***#@%%+%@%+%%+..=*%@@+=-.--++@@*=#@@#@#%##*#%=+@=@.:+@=##+*=@##%@+@@-=*+#=:                        
                                 .%#-@%%#+@+=+@-@+@@##=#%=+@@##*#@++=@@%@....-@@-.=..@#+-=@@@*##*%+*++##@..@++@=+=*@%=%=#@#*   +.                        
                                ..%##@%%%*@#+%%:@.@@%=%#@@%#+%-@@-=@#@#@@%@=#@++*+=++@@#%%#*@@#@%%#+#+%=#@.%**@#++*#@@-@*#++.  .                         
                                .###%*###*@=+#*%-#@%=%%%+#@+@@+%=#%@*#%+%@%+==-+==*@@%**#%#+*#*#%%%#%%*@*@.:@=%=%++##@#*+@#@.                            
                                .*# @=*##*@++@+%.*@=#%#*#%%-@@-*#=+++*##@@@%%%%@%%@@###++########+####*%%*@.@#@####+@#@%+%%-=                            
                                 :%#@+-@*#@##@=@=%@+%%%+#%%+#@#-+**+===:-.:--=-=----+-=++++*++%+%@#%+##+@#-#%+%#@%=:@.@@%-@#%                            
                                 =#*%#*@#@%*+%*#-@@%#*+##%@=*@*@-:==+==+#+=###=%%#+=-+**++=*+=@++%#%+=%#@#*=-###-%@-@@.@@@-@+.                           
                                 =##+%++%@++%##@.@%+#@+%@@%+*%+#@@*####=@@#-+#+%*%@@#=+++++**+@#=###@*#+@%*=+#%##*@-=#@.=@@-@.                           
                                 ##%+-*@#@+#@*##+@##@**@@%=%#%*+#+##++#*..@@@@@@@@@.-+*##+==%%%+*+%%=#####@==*%#%@%@+:@@..%+%..                          
                                 -=%#@%#@:*@##%#-@%-@+#@*=**##%=+####*##@-..=......#*#%*+*##+@#**=#%*+#%%+%#=+%#%*#=@#.@@=@+@                            
                                 .#%%+=@=-%#+#%#*@%+@=@%=+#+%%*#-+=+**##@@@@@@@@@@#*+=++#+=%%%##++%##+%%#%*%++=%#%##++=:@:+..                            
                                 .%+#=@=+%@##@#=@+@+@%@#+*#+=#@%+=+%##**+++===+==+#*+#*+**%#%%#+*+%##+#%*%:%*#+%%+#@+%*-@=*.                             
                                  =#+@.%@@%+%#@=@=%=@+@#+****#=*@#=#+##**###%%#%#+%%##*+=*@=%**+++@%#+%%*@-%@=#=%-@@ @**#+..                             
                                  .+#@+##***%#%=%#@+#=@##%=*#=@=-%@#*#+###%+%*@@=#%*###*+*+%%*+++*#%#+%%#@.%+-%*##=+ ......                              
                                  .@=.=@#@#%*%%+%##+#=@%#%%####@=-*##*####%+#*%@+%%%##+##+#%%++#+@##*##%-@-%%+*@#+*.                                     
                                   @+#@+%*#*%#%#*%.#@-@#%####*=%@==+#%%#%@+##*%%#=%@+**#*@%%*=%%%%####%#+@+%=*+=%%*.                                     
 - Whisky                          .................................................................................
	`

	fmt.Println(asciiArt)

	//Servidor
	if err := e.Start("0.0.0.0:3001"); err != nil {
		common.Logger.LogError().
			Err(err).
			Msg("Erro ao iniciar servidor")
	}
}

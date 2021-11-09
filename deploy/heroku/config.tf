resource "heroku_config" "common" {
  # Environment variables
  vars = {
    APP_ADDR = ":8080"
    APP_RUN_MODE = "release"
    APP_ENV = "release"
    DB_CONNECTION = "postgres"
    DB_HOST = "database-service"
    DB_PORT = "5432"
    DB_DATABASE = "backend"
    OAUTH2_REDIRECT_URL = ""
    FRONTEND_URL = ""
    APP_SECRET_KEY = ""
    DB_USERNAME = ""
    DB_PASSWORD = ""
    OAUTH2_CLIENT_ID = ""
    OAUTH2_CLIENT_SECRET = ""
  }
}

resource "heroku_app_config_association" "default" {
  app_id = heroku_app.default.id

  vars = heroku_config.common.vars
}

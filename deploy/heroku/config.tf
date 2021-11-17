resource "heroku_config" "common" {
  # Environment variables
  vars = {
    APP_RUN_MODE = "release"
    APP_ENV = "release"
    DB_CONNECTION = "postgres"
    OAUTH2_REDIRECT_URL = ""
    FRONTEND_URL = ""
    APP_SECRET_KEY = ""
    OAUTH2_CLIENT_ID = ""
    OAUTH2_CLIENT_SECRET = ""
  }
}

resource "heroku_app_config_association" "default" {
  app_id = heroku_app.default.id

  vars = heroku_config.common.vars
}

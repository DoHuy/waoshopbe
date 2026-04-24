data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./tools/atlas_loader.go",
  ]
}

env "local" {
  src = data.external_schema.gorm.url

  dev = "postgres://huydv:SecretPassword123@localhost:5432/uk_dropship_db?sslmode=disable"

  migration {
    dir = "file://dropship-deployment/sql/migrations?format=golang-migrate"
  }
}
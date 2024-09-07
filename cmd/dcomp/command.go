package main

var commandList = map[string]string{
	"build": "docker compose build",
	"pull":  "docker compose pull",
	"start": "docker compose up -d --remove-orphans",
	"stop":  "docker compose down",
	"ps":    "docker compose ps",
	"logs":  "docker compose logs",
}
var commandBase = "docker compose"

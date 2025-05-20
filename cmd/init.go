package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	projectType string
	framework   string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Go project with a specific structure",
	Run: func(cmd *cobra.Command, args []string) {
		wd, _ := os.Getwd()
		fmt.Printf("Generating project in %s\n", wd)
		fmt.Printf("Framework: %s, Type: %s\n", framework, projectType)

		if projectType == "microservice" {
			if framework == "gin" {
				generateMicroserviceGin(wd)
			} else {
				fmt.Println("Only 'gin' framework is supported for now.")
			}
		} else {
			fmt.Println("Only 'microservice' type is supported for now.")
		}
	},
}

func init() {
	initCmd.Flags().StringVar(&projectType, "type", "", "Project type: microservice")
	initCmd.Flags().StringVar(&framework, "framework", "", "Framework to use: gin")
	rootCmd.AddCommand(initCmd)
}

func generateMicroserviceGin(basePath string) {
	// Directories to create
	dirs := []string{
		"cmd/server",
		"internal/app/handler",
		"internal/app/service",
		"internal/app/repository",
		"internal/app/router",
		"internal/app/config",
		"pkg/model",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(basePath, dir)
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", fullPath, err)
			return
		}
	}

	// File templates
	files := map[string]string{
		"cmd/server/main.go":                   mainGo,
		"internal/app/router/router.go":        routerGo,
		"internal/app/handler/user_handler.go": userHandlerGo,
		"internal/app/service/user_service.go": userServiceGo,
		"internal/app/repository/user_repo.go": userRepoGo,
		"internal/app/config/config.go":        configGo,
		"pkg/model/user.go":                    userModelGo,
		".env":                                 envFile,
		"Dockerfile":                           dockerfile,
		"Makefile":                             makefile,
	}

	for path, content := range files {
		fullPath := filepath.Join(basePath, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			fmt.Printf("Failed to write file %s: %v\n", fullPath, err)
			return
		}
	}

	fmt.Println("Project generated successfully!")
}

var mainGo = `package main

import (
	"github.com/gin-gonic/gin"
	"internal/app/router"
)

func main() {
	r := gin.Default()
	router.SetupRoutes(r)
	r.Run()
}
`

var routerGo = `package router

import (
	"github.com/gin-gonic/gin"
	"internal/app/handler"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/users", handler.GetUsers)
}
`

var userHandlerGo = `package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Users fetched successfully"})
}
`

var userServiceGo = `package service

func FetchUsers() []string {
	return []string{"Alice", "Bob", "Charlie"}
}
`

var userRepoGo = `package repository

func GetAllUsersFromDB() []string {
	// Imagine DB call here
	return []string{"Alice", "Bob", "Charlie"}
}
`

var configGo = `package config

import "os"

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
`

var userModelGo = `package model

type User struct {
	ID   int
	Name string
}
`

var envFile = `PORT=8080
`

var dockerfile = `FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main ./cmd/server/main.go

EXPOSE 8080
CMD ["./main"]
`

var makefile = `build:
	go build -o bin/app ./cmd/server/main.go

run:
	go run ./cmd/server/main.go
`

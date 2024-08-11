package prisma

import (
	"log"

	"github.com/jianshao/chrome-exts/CleanTracks/backend/prisma/db"
	"github.com/joho/godotenv"
)

var gPrisma *db.PrismaClient

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	gPrisma = db.NewClient()
	if err := gPrisma.Connect(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
}

func GetPrismaClient() *db.PrismaClient {
	return gPrisma
}

func Close() {
	if gPrisma != nil {
		gPrisma.Disconnect()
		gPrisma = nil
	}
}

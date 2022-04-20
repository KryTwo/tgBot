package service

import (
	"github.com/jmoiron/sqlx"
	"log"
	e "main/lib"
	"main/repository"
	"math/rand"
	"time"
)

// NewUser checks for the presence of the user in the database, if not, adds it
func NewUser(db *sqlx.DB, uName string) string {
	isUser, err := repository.FindUser(db, uName)
	if err != nil {
		log.Fatal(e.Wrap("Can't do FindUser", err))
	}

	if !isUser {
		err := repository.InsertUserToDB(db, uName)
		if err != nil {
			log.Fatal(e.Wrap("Can't insert a user to DB", err))
		}

		return msgPosAnswAddUser(uName) // "Скоро ты будешь пидором, обещаю, "
	}
	return msgNegAnswAddUser(uName) // "Ты уже в базе, будущий пидор"
}

// Stats returns statistics for all users
func Stats(db *sqlx.DB) string {

	stats, err := repository.GetStats(db)
	if err != nil {
		log.Fatal("Can't get stats", err)
	}

	return statsToText(stats)
}

func StatsMe(db *sqlx.DB, uName string) string {
	if b, _ := repository.FindUser(db, uName); !b {
		return msgUserNotFind(uName)
	}
	res, err := repository.GetStatsMe(db, uName)
	if err != nil {
		log.Fatal(e.Wrap("Can't get stats me", err))
	}

	return statsMeToText(res)
}

func IncRandUser(db *sqlx.DB, uName string) string {
	randNum := randNum()
	if err := repository.IncrCounter(db, randNum); err != nil {
		log.Fatal(e.Wrap("Can't increment counter to user", err))
	}

	return msgResultIncRand(uName)
}

func randNum() int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(3)
	return r + 1
}

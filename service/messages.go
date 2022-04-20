package service

import (
	"bufio"
	"log"
	"main/structs"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func msgPosAnswAddUser(u string) string {
	return "Добро пожаловать в клуб! \nСкоро ты станешь пидором, @" + u
}
func msgNegAnswAddUser(u string) string {
	return "А не дохуя тебя там будет, @" + u + "?"
}

func msgUserNotFind(u string) string {
	return "@" + u + " ты недостоин. Вступи в клуб: /make_me_pidor"
}

// statsToText convert []structs.Users to string
func statsToText(str []structs.Users) string {

	var res string = "Статистика пидоров:\n"
	for v := range str {
		res = res + "@" + str[v].Username + " " +
			strconv.Itoa(str[v].Count) + " раз(а)\n"
	}

	return strings.TrimSpace(res)
}

// statsMeToText
func statsMeToText(str structs.Users) string {
	var res string = "@" + str.Username + ", ты был пидором: " + strconv.Itoa(str.Count) + " раз(а)"
	return res
}

func msgResultIncRand(u string) string {
	file, err := os.Open("text")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	randsource := rand.NewSource(time.Now().UnixNano())
	randgenerator := rand.New(randsource)

	lineNum := 1
	var pick string
	for scanner.Scan() {
		line := scanner.Text()
		// Instead of 1 to N it's 0 to N-1
		roll := randgenerator.Intn(lineNum)

		if roll == 0 {
			pick = line
		}

		lineNum += 1
	}

	return strings.Replace(pick, "*", "@"+u, 1)
}

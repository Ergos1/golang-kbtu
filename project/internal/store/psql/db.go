package psql

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	_ "log"
	"os"
	"github.com/jmoiron/sqlx"
)

func getEnv() (map[string]string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("[error] Error during getting path")
	}
	godotenv.Load(pwd + "\\.env")
	env, err := godotenv.Read()
	if err != nil {
		return nil, fmt.Errorf("[error] Error during reading .env file")
	}
	return env, nil
}

func getPsqlInfo() (string, error ){
	env, err := getEnv()
	if err != nil {
		return "", err
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
    "password=%s dbname=%s sslmode=disable",
    env["HOST"], env["PORT"], env["USER"], env["PASSWORD"], env["DBNAME"])
	return psqlInfo, nil
}

func NewDB() *sqlx.DB{
	psqlInfo, err := getPsqlInfo()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
//
//func main(){
//	psqlInfo, err := getPsqlInfo()
//	if err != nil {
//		log.Fatal(err)
//	}
//	db, err := sqlx.Open("postgres", psqlInfo)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//	rows, err := db.Queryx("Select id, balance from Wallet")
//	for rows.Next(){
//		var wallet models.Wallet
//		err = rows.StructScan(&wallet)
//
//		fmt.Println(wallet)
//	}
//}
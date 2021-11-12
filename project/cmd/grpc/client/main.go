package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"example.com/api"
	"google.golang.org/grpc"
)

const (
	port = ":5000"
)
func getAllExample(ctx context.Context, conn *grpc.ClientConn){
	getAllCollections(ctx, conn)
	getAllNonFungibleTokens(ctx, conn)
	getAllUsers(ctx, conn)
	getAllWallets(ctx, conn)
}

func getAllCollections(ctx context.Context, conn *grpc.ClientConn ) {
	collectionClient := api.NewCollectionServiceClient(conn)
	if collections, err := collectionClient.All(ctx, &api.Empty{}); err != nil {
		log.Fatal(err)
	} else {
		log.Print("Collections: \n")
		for k, v := range collections.Collections {
			log.Printf("%v %v", k, v)
		}
	}
}
func getAllUsers(ctx context.Context, conn *grpc.ClientConn ) {
	userClient := api.NewUserServiceClient(conn)
	if users, err := userClient.All(ctx, &api.Empty{}); err != nil {
		log.Fatal(err)
	} else {
		log.Print("Users: \n")
		for k, v := range users.Users {
			log.Printf("%v %v", k, v)
		}
	}
}

func getAllNonFungibleTokens(ctx context.Context, conn *grpc.ClientConn ) {
	nftClient := api.NewNonFungibleTokenServiceClient(conn)
	if nfts, err := nftClient.All(ctx, &api.Empty{}); err != nil {
		log.Fatal(err)
	} else {
		log.Print("Nfts: \n")
		for k, v := range nfts.NonFungibleTokens {
			log.Printf("%v %v", k, v)
		}

	}
}

func getAllWallets(ctx context.Context, conn *grpc.ClientConn ) {
	walletClient := api.NewWalletServiceClient(conn)
	if wallets, err := walletClient.All(ctx, &api.Empty{}); err != nil {
		log.Fatal(err)
	} else {
		log.Print("Wallets: \n")
		for k, v := range wallets.Wallets {
			log.Printf("%v %v", k, v)
		}
	}
}

func createAll(ctx context.Context, conn *grpc.ClientConn) {

	collectionClient := api.NewCollectionServiceClient(conn)
	nftClient := api.NewNonFungibleTokenServiceClient(conn)
	userClient := api.NewUserServiceClient(conn)
	walletClient := api.NewWalletServiceClient(conn)

	for i := uint64(1); i <= 10; i++{
		_, err := walletClient.Create(ctx, &api.Wallet{
			Id:      i,
			Balance: 10000000,
		})
		if err != nil {
			log.Fatal(err)
		}
		_, err = userClient.Create(ctx, &api.User{
			Id:       i,
			WalletId: i,
			Username: "Yergeldi" + strconv.Itoa(int(i)),
			Email:    "Yergeldi@mail.ru", // must be unique
			Password: "hdjGHjdskaghljsdkghfdskjg",
		})
		if err != nil {
			log.Fatal(err)
		}
		_, err = collectionClient.Create(ctx, &api.Collection{
			Id:          i,
			OwnerId:     1,
			Name:        "Yergeldi",
			Symbol:      "YED",
			Description: "ASDASDAS",
		})
		if err != nil {
			log.Fatal(err)
		}
		_, err = nftClient.Create(ctx, &api.NonFungibleToken{
			Id:           i,
			Likes:        0,
			CollectionId: i,
			OwnerId:      i,
			Price:        1230123213,
			Royalties:    10,
			Title:        "Yergeldi",
			Description:  "YESDFSD",
			Properties:   map[string]string{"Mana" : strconv.Itoa(int(i * 100))},
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}


func main(){
	ctx := context.Background()

	connStartTime := time.Now()
	conn, err := grpc.Dial("localhost" + port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect to %s: %v", port, err)
	}
	log.Printf("connected in %d microsec", time.Now().Sub(connStartTime).Microseconds())

	userClient := api.NewUserServiceClient(conn)

	// CREATE EXAMPLE
	createAll(ctx, conn)

	// GET ALL EXAMPLE
	getAllExample(ctx, conn)

	// FIELDS WITH ID
	validId := uint64(1)
	//invalidId := uint64(101)

	// GET BY ID EXAMPLE

	if user, err := userClient.ByID(ctx, &api.Id{Id:validId}); err != nil {
		log.Fatal(err)
	} else {
		log.Print(user)
	}

	//if user, err := userClient.ByID(ctx, &api.Id{Id:invalidId}); err != nil {
	//	log.Fatal(err)
	//} else {
	//	log.Print(user)
	//}


	// UPDATE EXAMPLE

	_, err = userClient.Update(ctx, &api.User{
		Id:       validId,
		WalletId: 2,
		Username: "ASDasda",
		Email:    "ASdasd@gmail.com",
		Password: "asdasdadegw",
	})

	if err != nil {
		log.Fatal(err)
	}

	if user, err := userClient.ByID(ctx, &api.Id{Id:validId}); err != nil {
		log.Fatal(err)
	} else {
		log.Print(user)
	}

	// DELETE EXAMPLE

	if _, err := userClient.Delete(ctx, &api.Id{Id: validId}); err != nil {
		log.Fatal(err)
	} else {
		log.Print("[Success] User was deleted")
	}

	getAllUsers(ctx, conn)
}
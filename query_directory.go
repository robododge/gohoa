package gohoa

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DirQueryService struct {
	DBService
}

func NewDirQueryService() *DirQueryService {
	dbSvc := createDBService("directory")
	return &DirQueryService{dbSvc}
}

func (s *DirQueryService) FindAllMembers(members *[]Member) error {
	ctx := context.TODO()
	cursor, err := s.collection.Find(ctx, bson.D{})
	if err != nil {
		log.Println("Error finding all", err)
		return err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var member Member
		err := cursor.Decode(&member)
		if err != nil {
			log.Println("Error decoding member: ", err)
			return err
		}
		*members = append(*members, member)
	}
	return nil
}

func (s *DirQueryService) FindMemberById(memberId int) (Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{primitive.E{Key: "memberid", Value: memberId}}
	var member Member
	err := s.collection.FindOne(ctx, filter).Decode(&member)
	if err != nil {
		log.Println("Error finding member by id: ", err)
		return member, err
	}
	return member, nil
}

func (s *DirQueryService) FindCountByStreetName() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pipeline := bson.A{
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$paddress.streetname"},
					{"total", bson.D{{"$sum", 1}}},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"total", -1}}}},
	}
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Error aggregating by street name: ", err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Println("Error decoding aggregate result: ", err)
		}
		fmt.Printf("Street: %s, total: %d\n", result["_id"], result["total"])
	}
}

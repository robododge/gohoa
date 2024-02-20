package gohoa

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DirLoaderService struct {
	DbService
}

func NewDirLoader() *DirLoaderService {
	dbSvc := createDbService("directory")
	return &DirLoaderService{dbSvc}
}

func (s *DirLoaderService) PopulateMongoFromJson() {
	config := GetConfig()
	allMembers := NewAllMembers()
	allMembers.PopulateFromJsonFile(config.SlimMembersJson)

	err := s.BulkInsert(allMembers.Members)
	if err != nil {
		log.Println("Error upserting members: ", err)
	}

}

func (s *DirLoaderService) RevalidateMongoFromJson() {
	config := GetConfig()
	allMembers := NewAllMembers()
	allMembers.PopulateFromJsonFile(config.SlimMembersJson)

	err := s.Revalidate(allMembers.Members)
	if err != nil {
		log.Println("Error revalidating members: ", err)
	}

}

func (s *DirLoaderService) BulkInsert(members []Member) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mi := make([]interface{}, len(members))
	for i, m := range members {
		m.ID = fmt.Sprintf("%d-%d", m.MemberId, m.PAddress.AddressID)
		mi[i] = m
	}

	res, err := s.collection.InsertMany(ctx, mi)
	if err != nil {
		log.Println("Error bulk inserting members: ", err)
		return err
	}
	log.Printf("Inserted %d members\n", len(res.InsertedIDs))
	return nil

}

func (s *DirLoaderService) Revalidate(members []Member) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mi := make([]interface{}, len(members))
	for i, m := range members {
		m.ID = fmt.Sprintf("%d-%d", m.MemberId, m.PAddress.AddressID)
		mi[i] = m
	}
	log.Println("Phase 1 Revalidating members : ", len(members))
	update1 := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "abandoned", Value: time.Now()}}}}
	res, err := s.collection.UpdateMany(ctx, bson.D{}, update1)
	if err != nil {
		log.Println("Error phase1 revalidte : ", err)
		return err
	}
	log.Printf("Phase 1 revalidated %d members\n", res.ModifiedCount)

	log.Printf("Phase 2, start revalidated %d members\n", res.ModifiedCount)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel2()
	for i, m := range members {
		id := fmt.Sprintf("%d-%d", m.MemberId, m.PAddress.AddressID)
		// filter := bson.D{{"_id", id}}
		update2 := bson.D{
			primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "refresh", Value: time.Now()}}},
			primitive.E{Key: "$unset", Value: bson.D{primitive.E{Key: "abandoned", Value: ""}}},
		}
		res, err = s.collection.UpdateByID(ctx2, id, update2)
		if err != nil {
			log.Println("Phase 2 error revalidating members: ", err)
			return err
		}
		log.Printf("Done with member[%d] %s updated %d\n", i, id, res.ModifiedCount)
	}
	log.Println("Phase 2, finished")
	return nil
}

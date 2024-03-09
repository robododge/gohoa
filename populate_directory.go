package gohoa

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DirLoaderService struct {
	DBService
}

func NewDirLoader() *DirLoaderService {
	dbSvc := createDBService("directory")
	return &DirLoaderService{dbSvc}
}

func (s *DirLoaderService) PopulateMongoFromJSON() {
	config := GetConfig()
	allMembers := NewAllMembers()
	allMembers.PopulateFromJsonFile(config.SlimMembersJson)
	allMembers.DeDupeMembers()

	err := s.BulkInsert(allMembers.Members)
	if err != nil {
		log.Println("Error upserting members: ", err)
	}

}

func (s *DirLoaderService) RevalidateMongoFromJSON() {
	config := GetConfig()
	allMembers := NewAllMembers()
	allMembers.PopulateFromJsonFile(config.SlimMembersJsonReval)
	allMembers.DeDupeMembers()

	err := s.Revalidate(allMembers.Members)
	if err != nil {
		log.Println("Error revalidating members: ", err)
	}

}

func populateAllMembersID(members []Member) []interface{} {
	mi := make([]interface{}, 0, len(members))
	for _, m := range members {
		if newID, addrErr := CreateMongoIDForDiretory(&m); addrErr == nil {
			m.ID = newID
			mi = append(mi, m)
		} else {
			log.Println("Error creating mongo id: ", addrErr)
		}

	}
	return mi
}

func (s *DirLoaderService) BulkInsert(members []Member) error {
	mi := populateAllMembersID(members)

	dbSession, err := s.client.StartSession()
	if err != nil {
		log.Println("Error starting session: ", err)
		return err
	}
	defer dbSession.EndSession(context.TODO())

	_, err = dbSession.WithTransaction(context.TODO(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		res, err2 := s.collection.InsertMany(sessCtx, mi)
		if err2 == nil {
			log.Printf("Inserted %d members\n", len(res.InsertedIDs))
		}
		return res, err2
	})
	if err != nil {
		log.Println("Error bulk inserting members: ", err)
		return err
	}
	return nil

}

func (s *DirLoaderService) Revalidate(members []Member) error {

	log.Println("Revalidating total members : ", len(members))

	dbSession, err := s.client.StartSession()
	if err != nil {
		log.Println("Error starting session: ", err)
		return err
	}
	defer dbSession.EndSession(context.TODO())

	log.Println("Transcation Started..")
	result, err := dbSession.WithTransaction(context.TODO(), func(sessCtx mongo.SessionContext) (interface{}, error) {

		update1 := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "abandoned", Value: time.Now()}}}}
		res, err := s.collection.UpdateMany(sessCtx, bson.D{}, update1)
		if err != nil {
			log.Println("Error phase1 revalidte : ", err)
			return res, err
		}
		log.Printf("Phase 1 revalidated %d members\n", res.ModifiedCount)

		log.Printf("Phase 2, start revalidated %d members\n", res.ModifiedCount)
		for i, m := range members {
			m2 := m
			if id, addrErr := CreateMongoIDForDiretory(&m2); addrErr == nil {

				// filter := bson.D{{"_id", id}}
				update2 := bson.D{
					primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "refresh", Value: time.Now()}}},
					primitive.E{Key: "$unset", Value: bson.D{primitive.E{Key: "abandoned", Value: ""}}},
				}
				res, err = s.collection.UpdateByID(sessCtx, id, update2)
				if err != nil {
					log.Println("Phase 2 error revalidating members: ", err)
					return res, err
				}
				log.Printf("Done with member[%d] %s updated %d\n", i, id, res.ModifiedCount)

				if res.ModifiedCount == 0 {
					log.Printf("Member %s at address %d %s, did not exist, inserting fresh !\n ", m2.MemberName, m2.PAddress.Number, m2.PAddress.StreetName)
					m2.ID = id
					res2, err := s.collection.InsertOne(sessCtx, m2)
					if err != nil {
						log.Println("Phase 2 error inserting members: ", err)
						return res2, err
					}
					if res2.InsertedID != nil {
						log.Printf("Inserted new member %s at address %d %s, id: %s\n", m2.MemberName, m2.PAddress.Number, m2.PAddress.StreetName, res2.InsertedID)
					}

				}
			} else {
				log.Println("Address error creating mongo id: ", addrErr)
			}

		}
		log.Println("Phase 2, finished")
		return res, err
	})
	log.Println("Transcation finished, last result: ", result)
	return err
}

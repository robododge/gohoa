package gohoa

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

type ContactReqService struct {
	DBService
}

func NewContactReqService() *ContactReqService {
	dbSvc := createDBService("mb_events")
	return &ContactReqService{dbSvc}
}

func (s *ContactReqService) CreateContactRequest(cr ContactRequest) error {
	objID := primitive.NewObjectID()
	strObjID := fmt.Sprintf("CONTACT-%s", objID.Hex())
	cr.ID = strObjID
	cr.Type = "contact"
	cr.CreateDate = primitive.NewDateTimeFromTime(time.Now())
	_, err := s.collection.InsertOne(context.TODO(), cr)
	return err
}

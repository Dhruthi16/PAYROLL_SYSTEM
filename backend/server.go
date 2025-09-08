package main

import (
	"context"
	"log"
	"time"

	"PAYROLL_SYSTEM/backend/models"
	"PAYROLL_SYSTEM/backend/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/emptypb"
)

type payrollServer struct {
	proto.UnimplementedPayrollServiceServer
	collection *mongo.Collection
}

func newServer(uri string) *payrollServer {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Mongo connection failed: %v", err)
	}

	coll := client.Database("payroll_db").Collection("payrolls")
	return &payrollServer{collection: coll}
}

func (s *payrollServer) CreatePayroll(ctx context.Context, req *proto.CreatePayrollRequest) (*proto.Payroll, error) {
	doc := models.Payroll{
		EmpID:   req.EmpId,
		EmpName: req.EmpName,
		Salary:  req.Salary,
		Month:   req.Month,
	}
	res, err := s.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()
	return &proto.Payroll{
		Id:      id,
		EmpId:   req.EmpId,
		EmpName: req.EmpName,
		Salary:  req.Salary,
		Month:   req.Month,
	}, nil
}

func (s *payrollServer) GetPayroll(ctx context.Context, req *proto.GetPayrollRequest) (*proto.Payroll, error) {
	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	var doc models.Payroll
	if err := s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&doc); err != nil {
		return nil, err
	}

	return &proto.Payroll{
		Id:      doc.ID,
		EmpId:   doc.EmpID,
		EmpName: doc.EmpName,
		Salary:  doc.Salary,
		Month:   doc.Month,
	}, nil
}

func (s *payrollServer) UpdatePayroll(ctx context.Context, req *proto.UpdatePayrollRequest) (*proto.Payroll, error) {
	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"emp_name": req.EmpName,
			"salary":   req.Salary,
			"month":    req.Month,
		},
	}
	_, err = s.collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return nil, err
	}

	return s.GetPayroll(ctx, &proto.GetPayrollRequest{Id: req.Id})
}

func (s *payrollServer) DeletePayroll(ctx context.Context, req *proto.DeletePayrollRequest) (*emptypb.Empty, error) {
	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *payrollServer) ListPayrolls(ctx context.Context, req *proto.ListPayrollsRequest) (*proto.ListPayrollsResponse, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var res []*proto.Payroll
	for cur.Next(ctx) {
		var doc models.Payroll
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		res = append(res, &proto.Payroll{
			Id:      doc.ID,
			EmpId:   doc.EmpID,
			EmpName: doc.EmpName,
			Salary:  doc.Salary,
			Month:   doc.Month,
		})
	}
	return &proto.ListPayrollsResponse{Payrolls: res}, nil
}

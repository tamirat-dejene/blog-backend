package repositories

import (
	"context"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"

	"go.mongodb.org/mongo-driver/bson"
)

type OTPRepository struct {
	db         mongo.Database
	collection string
}

func NewOTPRepository(db mongo.Database, collection string) domain.IOTPRepository {
	return &OTPRepository{
		db:         db,
		collection: collection,
	}
}

// SaveOTP
func (r *OTPRepository) SaveOTP(ctx context.Context, otp *domain.OTP) error {
	otpModel := mapper.OtpFromDomain(otp)
	collection := r.db.Collection(r.collection)
	if collection == nil {
		return fmt.Errorf("database collection is not initialized")
	}
	_, err := collection.InsertOne(ctx, otpModel)
	if err != nil {
		return err
	}
	return nil
}

// FindOTPByEmail
func (r *OTPRepository) FindOTPByEmail(ctx context.Context, email string) (*domain.OTP, error) {
	var otpModel mapper.OtpDB
	err := r.db.Collection(r.collection).FindOne(ctx, bson.M{"email": email}).Decode(&otpModel)
	if err != nil {
		if err == mongo.ErrNoDocuments() {
			return nil, domain.ErrOTPNotFound
		}
		return nil, err
	}
	return mapper.OtpToDomain(&otpModel), nil
}

// delete OTP by email
func (r *OTPRepository) DeleteOTPByID(ctx context.Context, id string) error {
	_, err := r.db.Collection(r.collection).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete OTP with id %s: %w", id, err)
	}
	return nil
}

// update otp by id
func (r *OTPRepository) UpdateOTPByID(ctx context.Context, otp *domain.OTP) error {
	otpModel := mapper.OtpFromDomain(otp)
	_, err := r.db.Collection(r.collection).UpdateOne(ctx, bson.M{"_id": otpModel.ID}, bson.M{"$set": otpModel})
	if err != nil {
		return fmt.Errorf("failed to update OTP with id %s: %w", otp.ID, err)
	}
	return nil
}

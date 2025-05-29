package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"meaningfullname/Yurt_Mart/review/internal/model"
	"meaningfullname/Yurt_Mart/review/internal/repository"
	"meaningfullname/Yurt_Mart/review/proto"
)

type ReviewService struct {
	proto.UnimplementedReviewServiceServer
	repo repository.ReviewRepository
}

func NewReviewService(repo repository.ReviewRepository) *ReviewService {
	return &ReviewService{
		repo: repo,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, req *proto.CreateReviewRequest) (*proto.CreateReviewResponse, error) {
	review := &model.Review{
		ID:        uuid.New().String(),
		UserID:    req.UserId,
		ProductID: req.ProductId,
		Rating:    req.Rating,
		Comment:   req.Comment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, review); err != nil {
		return nil, status.Error(codes.Internal, "failed to create review")
	}

	return &proto.CreateReviewResponse{
		Review: &proto.Review{
			Id:        review.ID,
			UserId:    review.UserID,
			ProductId: review.ProductID,
			Rating:    review.Rating,
			Comment:   review.Comment,
			Timestamp: review.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *ReviewService) GetProductReviews(ctx context.Context, req *proto.GetProductReviewsRequest) (*proto.GetProductReviewsResponse, error) {
	reviews, err := s.repo.GetByProductID(ctx, req.ProductId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get product reviews")
	}

	protoReviews := make([]*proto.Review, len(reviews))
	for i, review := range reviews {
		protoReviews[i] = &proto.Review{
			Id:        review.ID,
			UserId:    review.UserID,
			ProductId: review.ProductID,
			Rating:    review.Rating,
			Comment:   review.Comment,
			Timestamp: review.CreatedAt.Format(time.RFC3339),
		}
	}

	return &proto.GetProductReviewsResponse{
		Reviews: protoReviews,
	}, nil
}

func (s *ReviewService) GetUserReviews(ctx context.Context, req *proto.GetUserReviewsRequest) (*proto.GetUserReviewsResponse, error) {
	reviews, err := s.repo.GetByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user reviews")
	}

	protoReviews := make([]*proto.Review, len(reviews))
	for i, review := range reviews {
		protoReviews[i] = &proto.Review{
			Id:        review.ID,
			UserId:    review.UserID,
			ProductId: review.ProductID,
			Rating:    review.Rating,
			Comment:   review.Comment,
			Timestamp: review.CreatedAt.Format(time.RFC3339),
		}
	}

	return &proto.GetUserReviewsResponse{
		Reviews: protoReviews,
	}, nil
}

func (s *ReviewService) UpdateReview(ctx context.Context, req *proto.UpdateReviewRequest) (*proto.UpdateReviewResponse, error) {
	review, err := s.repo.GetByID(ctx, req.ReviewId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "review not found")
	}

	review.Rating = req.Rating
	review.Comment = req.Comment
	review.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, review); err != nil {
		return nil, status.Error(codes.Internal, "failed to update review")
	}

	return &proto.UpdateReviewResponse{
		Review: &proto.Review{
			Id:        review.ID,
			UserId:    review.UserID,
			ProductId: review.ProductID,
			Rating:    review.Rating,
			Comment:   review.Comment,
			Timestamp: review.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, req *proto.DeleteReviewRequest) (*proto.DeleteReviewResponse, error) {
	if err := s.repo.Delete(ctx, req.ReviewId); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete review")
	}

	return &proto.DeleteReviewResponse{
		Success: true,
	}, nil
}

func (s *ReviewService) GetAverageProductRating(ctx context.Context, req *proto.GetAverageProductRatingRequest) (*proto.GetAverageProductRatingResponse, error) {
	avg, err := s.repo.GetAverageRating(ctx, req.ProductId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get average rating")
	}

	return &proto.GetAverageProductRatingResponse{
		AverageRating: avg,
	}, nil
}

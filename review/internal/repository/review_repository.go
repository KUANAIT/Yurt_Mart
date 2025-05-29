package repository

import (
	"context"
	"time"

	"meaningfullname/Yurt_Mart/common/database"
	"meaningfullname/Yurt_Mart/review/internal/model"

	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *model.Review) error
	GetByID(ctx context.Context, id string) (*model.Review, error)
	GetByProductID(ctx context.Context, productID string) ([]*model.Review, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.Review, error)
	Update(ctx context.Context, review *model.Review) error
	Delete(ctx context.Context, id string) error
	GetAverageRating(ctx context.Context, productID string) (float64, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository() ReviewRepository {
	return &reviewRepository{
		db: database.GetDB(),
	}
}

func (r *reviewRepository) Create(ctx context.Context, review *model.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *reviewRepository) GetByID(ctx context.Context, id string) (*model.Review, error) {
	var review model.Review
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetByProductID(ctx context.Context, productID string) ([]*model.Review, error) {
	var reviews []*model.Review
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) GetByUserID(ctx context.Context, userID string) ([]*model.Review, error) {
	var reviews []*model.Review
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) Update(ctx context.Context, review *model.Review) error {
	review.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *reviewRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Review{}, "id = ?", id).Error
}

func (r *reviewRepository) GetAverageRating(ctx context.Context, productID string) (float64, error) {
	var avg float64
	err := r.db.WithContext(ctx).Model(&model.Review{}).
		Where("product_id = ?", productID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avg).Error
	if err != nil {
		return 0, err
	}
	return avg, nil
}

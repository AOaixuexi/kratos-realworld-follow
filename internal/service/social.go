package service

import (
	"context"

	pb "helloworld/api/web/v1"
	"helloworld/internal/biz"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertArticle(do *biz.Article) *pb.Article {
	return &pb.Article{
		Slug:           do.Slug,
		Title:          do.Title,
		Description:    do.Description,
		Body:           do.Body,
		TagList:        do.TagList,
		CreatedAt:      timestamppb.New(do.CreatedAt),
		UpdatedAt:      timestamppb.New(do.UpdatedAt),
		Favorited:      do.Favorited,
		FavoritesCount: do.FavoritesCount,
		Author:         convertProfile(do.Author),
	}
}

func convertComment(do *biz.Comment) *pb.Comment {
	return &pb.Comment{
		Id:        uint32(do.ID),
		CreatedAt: timestamppb.New(do.CreatedAt),
		UpdatedAt: timestamppb.New(do.UpdatedAt),
		Body:      do.Body,
		Author:    convertProfile(do.Author),
	}
}

func convertProfile(do *biz.Profile) *pb.Profile {
	return &pb.Profile{
		Username:  do.Username,
		Bio:       do.Bio,
		Image:     do.Image,
		Following: do.Following,
	}
}

// 具体的业务逻辑实现
func (s *WebService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.ProfileReply, error) {
	rv, err := s.sc.GetProfile(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	return &pb.ProfileReply{
		Profile: convertProfile(rv),
	}, nil
}
func (s *WebService) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.ProfileReply, error) {
	rv, err := s.sc.FollowUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	return &pb.ProfileReply{
		Profile: convertProfile(rv),
	}, nil
}
func (s *WebService) UnfollowUser(ctx context.Context, req *pb.UnfollowUserRequest) (*pb.ProfileReply, error) {
	rv, err := s.sc.UnfollowUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	return &pb.ProfileReply{
		Profile: convertProfile(rv),
	}, nil
}
func (s *WebService) ListArticles(ctx context.Context, req *pb.ListArticlesRequest) (*pb.MultipleArticlesReply, error) {
	rv, err := s.sc.ListArticles(ctx)
	if err != nil {
		return nil, err
	}
	articles := make([]*pb.Article, 0)
	for _, x := range rv {
		articles = append(articles, convertArticle(x))
	}
	return &pb.MultipleArticlesReply{Articles: articles}, nil
}
func (s *WebService) FeedArticles(ctx context.Context, req *pb.FeedArticlesRequest) (*pb.MultipleArticlesReply, error) {
	rv, err := s.sc.ListArticles(ctx,
		biz.ListLimit(req.Limit),
		biz.ListOffset(req.Offset),
	)
	if err != nil {
		return nil, err
	}
	articles := make([]*pb.Article, 0)
	for _, x := range rv {
		articles = append(articles, convertArticle(x))
	}
	return &pb.MultipleArticlesReply{Articles: articles}, nil
}
func (s *WebService) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.SingleArticleReply, error) {
	rv, err := s.sc.GetArticle(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	return &pb.SingleArticleReply{
		Article: convertArticle(rv),
	}, nil
}
func (s *WebService) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.SingleArticleReply, error) {
	rv, err := s.sc.CreateArticle(ctx, &biz.Article{
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		TagList:     req.Article.TagList,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SingleArticleReply{
		Article: convertArticle(rv),
	}, nil
}
func (s *WebService) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.SingleArticleReply, error) {
	rv, err := s.sc.UpdateArticle(ctx, &biz.Article{
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		TagList:     req.Article.TagList,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SingleArticleReply{
		Article: convertArticle(rv),
	}, nil
}
func (s *WebService) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.SingleArticleReply, error) {
	err := s.sc.DeleteArticle(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	return &pb.SingleArticleReply{
		Article: &pb.Article{
			Slug: req.Slug,
		},
	}, nil
}
func (s *WebService) AddComment(ctx context.Context, req *pb.AddCommentRequest) (*pb.SingleCommentReply, error) {
	rv, err := s.sc.AddComment(ctx, req.Slug, &biz.Comment{
		Body: req.Comment.Body,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SingleCommentReply{
		Comment: convertComment(rv),
	}, nil
}
func (s *WebService) GetComments(ctx context.Context, req *pb.AddCommentRequest) (*pb.MultipleCommentsReply, error) {
	rv, err := s.sc.ListComments(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	comments := make([]*pb.Comment, 0)
	for _, x := range rv {
		comments = append(comments, convertComment(x))
	}
	return &pb.MultipleCommentsReply{Comments: comments}, nil
}
func (s *WebService) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (reply *pb.SingleCommentReply, err error) {
	reply = &pb.SingleCommentReply{}

	return &pb.SingleCommentReply{
		Comment: &pb.Comment{
			Id: uint32(req.Id),
		},
	}, nil
}
func (s *WebService) FavoriteArticle(ctx context.Context, req *pb.FavoriteArticleRequest) (*pb.SingleArticleReply, error) {
	rv, err := s.sc.FavoriteArticle(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	return &pb.SingleArticleReply{
		Article: convertArticle(rv),
	}, nil
}
func (s *WebService) UnfavoriteArticle(ctx context.Context, req *pb.UnfavoriteArticleRequest) (*pb.SingleArticleReply, error) {
	rv, err := s.sc.UnfavoriteArticle(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	return &pb.SingleArticleReply{
		Article: convertArticle(rv),
	}, nil
}
func (s *WebService) GetTags(ctx context.Context, req *pb.GetTagsRequest) (reply *pb.TagListReply, err error) {
	rv, err := s.sc.GetTags(ctx)
	if err != nil {
		return nil, err
	}
	tags := make([]string, len(rv))
	for i, x := range rv {
		tags[i] = string(x)
	}
	reply = &pb.TagListReply{Tags: tags}
	return reply, nil
}

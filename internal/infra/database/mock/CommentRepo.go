// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	domain "trainee/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// CommentRepo is an autogenerated mock type for the CommentRepo type
type CommentRepo struct {
	mock.Mock
}

// DeleteComment provides a mock function with given fields: id
func (_m *CommentRepo) DeleteComment(id int64) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetComment provides a mock function with given fields: id
func (_m *CommentRepo) GetComment(id int64) (domain.Comment, error) {
	ret := _m.Called(id)

	var r0 domain.Comment
	if rf, ok := ret.Get(0).(func(int64) domain.Comment); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(domain.Comment)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCommentsByPostID provides a mock function with given fields: postID, offset
func (_m *CommentRepo) GetCommentsByPostID(postID int64, offset int) ([]domain.Comment, error) {
	ret := _m.Called(postID, offset)

	var r0 []domain.Comment
	if rf, ok := ret.Get(0).(func(int64, int) []domain.Comment); ok {
		r0 = rf(postID, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Comment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, int) error); ok {
		r1 = rf(postID, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveComment provides a mock function with given fields: comment
func (_m *CommentRepo) SaveComment(comment domain.Comment) (domain.Comment, error) {
	ret := _m.Called(comment)

	var r0 domain.Comment
	if rf, ok := ret.Get(0).(func(domain.Comment) domain.Comment); ok {
		r0 = rf(comment)
	} else {
		r0 = ret.Get(0).(domain.Comment)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.Comment) error); ok {
		r1 = rf(comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateComment provides a mock function with given fields: comment
func (_m *CommentRepo) UpdateComment(comment domain.Comment) (domain.Comment, error) {
	ret := _m.Called(comment)

	var r0 domain.Comment
	if rf, ok := ret.Get(0).(func(domain.Comment) domain.Comment); ok {
		r0 = rf(comment)
	} else {
		r0 = ret.Get(0).(domain.Comment)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.Comment) error); ok {
		r1 = rf(comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCommentRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewCommentRepo creates a new instance of CommentRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCommentRepo(t mockConstructorTestingTNewCommentRepo) *CommentRepo {
	mock := &CommentRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

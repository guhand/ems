package service

import (
	apperror "ems/app/model/app_error"
	"ems/app/model/request"
	"ems/app/model/response"
	"ems/domain"
)

type noticeService struct {
	noticeRepository     domain.NoticeRepository
	departmentRepository domain.DepartmentRepository
}

func NewNoticeService(noticeRepository domain.NoticeRepository,
	departmentRepository domain.DepartmentRepository) domain.NoticeService {
	return &noticeService{noticeRepository, departmentRepository}
}

func (s *noticeService) CreateNotice(departmentMemberID uint, req *request.CreateNotice) error {
	isDepartmentMemberExists, err := s.departmentRepository.IsDepartmentMemberExists(departmentMemberID)

	if err != nil {
		return err
	}

	if !isDepartmentMemberExists {
		return apperror.DataNotFoundError("Employee")
	}

	if err := s.noticeRepository.CreateNotice(departmentMemberID, req); err != nil {
		return err
	}

	return nil
}

func (s *noticeService) FetchActiveUserNotices(filters *request.CommonRequest) ([]response.FetchActiveUserNotices, error) {
	data, err := s.noticeRepository.FetchActiveUserNotices(filters)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *noticeService) FetchNotice(departmentMemberID uint) (*response.FetchActiveUserNotices, error) {
	isDepartmentMemberExists, err := s.departmentRepository.IsDepartmentMemberExists(departmentMemberID)

	if err != nil {
		return nil, err
	}

	if !isDepartmentMemberExists {
		return nil, apperror.DataNotFoundError("Employee")
	}

	data, err := s.noticeRepository.FetchNotice(departmentMemberID)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *noticeService) ApproveNotice(departmentMemberID, approvedBy uint, req *request.ApproveNotice) error {
	isDepartmentMemberExists, err := s.departmentRepository.IsDepartmentMemberExists(departmentMemberID)

	if err != nil {
		return err
	}

	if !isDepartmentMemberExists {
		return apperror.DataNotFoundError("Employee")
	}

	if err := s.noticeRepository.ApproveNotice(departmentMemberID, approvedBy, req); err != nil {
		return err
	}

	return nil
}

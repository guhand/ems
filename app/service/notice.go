package service

import (
	apperror "ems/app/model/app_error"
	"ems/app/model/request"
	"ems/app/model/response"
	"ems/domain"
	"ems/utils"
	"fmt"
)

type noticeService struct {
	noticeRepository     domain.NoticeRepository
	departmentRepository domain.DepartmentRepository
}

func NewNoticeService(noticeRepository domain.NoticeRepository,
	departmentRepository domain.DepartmentRepository) domain.NoticeService {
	return &noticeService{noticeRepository, departmentRepository}
}

func (s *noticeService) ApplyNotice(departmentMemberID uint, req *request.ApplyNotice) error {
	isDepartmentMemberExists, err := s.departmentRepository.IsDepartmentMemberExists(departmentMemberID)

	if err != nil {
		return err
	}

	if !isDepartmentMemberExists {
		return apperror.DataNotFoundError("user")
	}

	if err := s.noticeRepository.ApplyNotice(departmentMemberID, req); err != nil {
		return err
	}

	return nil
}

func (s *noticeService) FetchActiveUserNotices(roleID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	data, err := s.noticeRepository.FetchActiveUserNotices(roleID, filters)

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
		return nil, apperror.DataNotFoundError("user")
	}

	data, err := s.noticeRepository.FetchNotice(departmentMemberID)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *noticeService) ApproveNotice(approvedBy uint, req *request.ApproveNotice) error {
	isDepartmentMemberExists, err := s.departmentRepository.IsDepartmentMemberExists(uint(req.DepartmentMemberID))

	if err != nil {
		return err
	}

	if !isDepartmentMemberExists {
		return apperror.DataNotFoundError("user")
	}

	isNoticeExistsByUser, err := s.noticeRepository.IsApproveExistsByUser(uint(req.DepartmentMemberID))

	if err != nil {
		return err
	}

	if !isNoticeExistsByUser {
		return fmt.Errorf("notice not found for the user")
	}

	if err := s.noticeRepository.ApproveNotice(uint(req.DepartmentMemberID), approvedBy, req); err != nil {
		return err
	}

	return nil
}

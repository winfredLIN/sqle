package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

var (
	ErrProjectNotExist = func(projectName string) error {
		return errors.New(errors.DataNotExist, fmt.Errorf("project [%v] is not exist", projectName))
	}
	ErrProjectArchived = errors.New(errors.ErrAccessDeniedError, fmt.Errorf("project is archived"))
)

type GetProjectReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetProjectResV1 struct {
	controller.BaseRes
	Data      []*ProjectListItem `json:"data"`
	TotalNums uint64             `json:"total_nums"`
}

type ProjectListItem struct {
	Name           string     `json:"name"`
	Desc           string     `json:"desc"`
	CreateUserName string     `json:"create_user_name"`
	CreateTime     *time.Time `json:"create_time"`
	Archived       bool       `json:"archived"`
}

// GetProjectListV1
// @Summary 获取项目列表
// @Description get project list
// @Tags project
// @Id getProjectListV1
// @Security ApiKeyAuth
// @Param page_index query uint32 true "page index"
// @Param page_size query uint32 true "size of per page" default(50)
// @Success 200 {object} v1.GetProjectResV1
// @router /v1/projects [get]
func GetProjectListV1(c echo.Context) error {
	req := new(GetProjectReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)

	user := controller.GetUserName(c)

	mp := map[string]interface{}{
		"limit":            limit,
		"offset":           offset,
		"filter_user_name": user,
	}

	s := model.GetStorage()
	projects, total, err := s.GetProjectsByReq(mp)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	resp := []*ProjectListItem{}
	for _, project := range projects {
		resp = append(resp, &ProjectListItem{
			Name:           project.Name,
			Desc:           project.Desc,
			CreateUserName: project.CreateUserName,
			CreateTime:     &project.CreateTime,
			Archived:       project.Status == model.ProjectStatusArchived,
		})
	}

	return c.JSON(http.StatusOK, &GetProjectResV1{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      resp,
		TotalNums: total,
	})
}

type GetProjectDetailResV1 struct {
	controller.BaseRes
	Data ProjectDetailItem `json:"data"`
}

type ProjectDetailItem struct {
	Name           string     `json:"name"`
	Desc           string     `json:"desc"`
	CreateUserName string     `json:"create_user_name"`
	CreateTime     *time.Time `json:"create_time"`
	Archived       bool       `json:"archived"`
}

// GetProjectDetailV1
// @Summary 获取项目详情
// @Description get project detail
// @Tags project
// @Id getProjectDetailV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} v1.GetProjectDetailResV1
// @router /v1/projects/{project_name}/ [get]
func GetProjectDetailV1(c echo.Context) error {
	projectName := c.Param("project_name")
	userName := controller.GetUserName(c)
	s := model.GetStorage()
	err := CheckIsProjectMember(userName, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrProjectNotExist(projectName))
	}

	return c.JSON(http.StatusOK, GetProjectDetailResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: ProjectDetailItem{
			Name:           project.Name,
			Desc:           project.Desc,
			CreateUserName: project.CreateUser.Name,
			CreateTime:     &project.CreatedAt,
			Archived:       project.Status == model.ProjectStatusArchived,
		},
	})
}

type CreateProjectReqV1 struct {
	Name string `json:"name" valid:"required"`
	Desc string `json:"desc"`
}

// CreateProjectV1
// @Summary 创建项目
// @Description create project
// @Accept json
// @Produce json
// @Tags project
// @Id createProjectV1
// @Security ApiKeyAuth
// @Param project body v1.CreateProjectReqV1 true "create project request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects [post]
func CreateProjectV1(c echo.Context) error {
	return createProjectV1(c)
}

type UpdateProjectReqV1 struct {
	Desc *string `json:"desc"`
}

// UpdateProjectV1
// @Summary 更新项目
// @Description update project
// @Accept json
// @Produce json
// @Tags project
// @Id updateProjectV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Param project body v1.UpdateProjectReqV1 true "create project request"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/ [patch]
func UpdateProjectV1(c echo.Context) error {
	req := new(UpdateProjectReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectName := c.Param("project_name")
	err = CheckIsProjectManager(user.Name, projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	archived, err := s.IsProjectArchived(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if archived {
		return controller.JSONBaseErrorReq(c, ErrProjectArchived)
	}

	sure, err := s.CheckUserCanUpdateProject(projectName, user.ID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !sure {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("you can not modify this project"))
	}

	attr := map[string]interface{}{}
	if req.Desc != nil {
		attr["desc"] = *req.Desc
	}

	return controller.JSONBaseErrorReq(c, s.UpdateProjectInfoByID(projectName, attr))
}

// DeleteProjectV1
// @Summary 删除项目
// @Description delete project
// @Id deleteProjectV1
// @Tags project
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/ [delete]
func DeleteProjectV1(c echo.Context) error {
	return deleteProjectV1(c)
}

type GetProjectTipsReqV1 struct {
	FunctionalModule string `json:"functional_module" query:"functional_module"`
}

type GetProjectTipsResV1 struct {
	controller.BaseRes
	Data []ProjectTipResV1 `json:"data"`
}

type ProjectTipResV1 struct {
	Name string `json:"project_name"`
}

// GetProjectTipsV1
// @Summary 获取项目提示列表
// @Description get project tip list
// @Tags project
// @Id getProjectTipsV1
// @Security ApiKeyAuth
// @Param functional_module query string false "functional module" Enums(operation_record)
// @Success 200 {object} v1.GetProjectTipsResV1
// @router /v1/project_tips [get]
func GetProjectTipsV1(c echo.Context) error {
	req := new(GetProjectTipsReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()

	data := []ProjectTipResV1{}

	switch req.FunctionalModule {
	case "operation_record":
		projectNameList, err := s.GetOperationRecordProjectNameList()
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		for _, projectName := range projectNameList {
			data = append(data, ProjectTipResV1{
				Name: projectName,
			})
		}
	default:
		projects, err := s.GetProjectTips(controller.GetUserName(c))
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		for _, project := range projects {
			data = append(data, ProjectTipResV1{
				Name: project.Name,
			})
		}
	}

	return c.JSON(http.StatusOK, GetProjectTipsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    data,
	})
}

// ArchiveProjectV1
// @Summary 归档项目
// @Description archive project
// @Accept json
// @Produce json
// @Tags project
// @Id archiveProjectV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/archive [post]
func ArchiveProjectV1(c echo.Context) error {
	return archiveProjectV1(c)
}

// UnarchiveProjectV1
// @Summary 取消归档项目
// @Description archive project
// @Accept json
// @Produce json
// @Tags project
// @Id unarchiveProjectV1
// @Security ApiKeyAuth
// @Param project_name path string true "project name"
// @Success 200 {object} controller.BaseRes
// @router /v1/projects/{project_name}/unarchive [post]
func UnarchiveProjectV1(c echo.Context) error {
	return unarchiveProjectV1(c)
}

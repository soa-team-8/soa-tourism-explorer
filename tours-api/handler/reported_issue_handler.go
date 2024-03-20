package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
	"tours/model"
	"tours/service"
	"tours/utils"
)

type ReportedIssueHandler struct {
	ReportedIssueService *service.ReportedIssueService
	HttpUtils            *utils.HttpUtils
}

func (e *ReportedIssueHandler) Create(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	categoryStr := vars["cat"]
	descriptionStr := vars["desc"]
	priorityStr := vars["prior"]
	tourIdStr := vars["tourID"]
	userIdStr := vars["userID"]

	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}
	tourID, err := strconv.Atoi(tourIdStr)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	reportedIssue, err := e.ReportedIssueService.Create(categoryStr, descriptionStr, uint64(priority), uint64(tourID), uint64(userID))
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusCreated, reportedIssue)
}

func (e *ReportedIssueHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	reportedIssues, err := e.ReportedIssueService.GetAll()
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, reportedIssues)
}

func (e *ReportedIssueHandler) GetAllByAuthor(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	reportedIssues, err := e.ReportedIssueService.GetByAuthor(id)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, reportedIssues)
}

func (e *ReportedIssueHandler) GetAllByTourist(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	reportedIssues, err := e.ReportedIssueService.GetByTourist(id)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, reportedIssues)
}

func (e *ReportedIssueHandler) AddComment(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	comment, err := e.HttpUtils.Decode(req.Body, &model.ReportedIssueComment{})
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	reportedIssue, err := e.ReportedIssueService.AddComment(id, *comment.(*model.ReportedIssueComment))
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, reportedIssue)
}

func (e *ReportedIssueHandler) AddDeadline(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	date, err := e.HttpUtils.Decode(req.Body, &time.Time{})
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	reportedIssue, err := e.ReportedIssueService.AddDeadline(id, *date.(*time.Time))
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, reportedIssue)
}

func (e *ReportedIssueHandler) PenalizeAuthor(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	reportedIssue, err := e.ReportedIssueService.PenalizeAuthor(id)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, reportedIssue)
}

func (e *ReportedIssueHandler) Close(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	reportedIssue, err := e.ReportedIssueService.Close(id)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, reportedIssue)
}

func (e *ReportedIssueHandler) Resolve(resp http.ResponseWriter, req *http.Request) {
	id, err := e.HttpUtils.GetIDFromRequest(req)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusBadRequest)
		return
	}

	reportedIssue, err := e.ReportedIssueService.Resolve(id)
	if err != nil {
		e.HttpUtils.HandleError(resp, err, http.StatusInternalServerError)
		return
	}

	e.HttpUtils.WriteJSONResponse(resp, http.StatusOK, reportedIssue)
}

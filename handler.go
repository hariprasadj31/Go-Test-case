package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"regexp"

	"github.com/google/uuid"
)

func ValidatePayeeRef(w http.ResponseWriter, r *http.Request) {

	var req AssociateWrapper

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	errors := ValidateRequest(req)

	if len(errors) > 0 {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"success"}`))
}

func ValidateRequest(req AssociateWrapper) []ErrorResponse {

	var errors []ErrorResponse

	if req.AssociateRequests == nil {
		errors = append(errors, ErrorResponse{
			Id:      "F00-101108",
			Name:    "AssociateRequests",
			Message: "AssociateRequests should not be null",
		})
		return errors
	}

	duplicateMap := make(map[string]bool)

	for i, r := range req.AssociateRequests {

		indexBusiness := fmt.Sprintf("AssociateRequests[%d].BusinessId", i)
		indexPayee := fmt.Sprintf("AssociateRequests[%d].PayeeRef", i)

		// BusinessId required
		if r.BusinessId == "" {
			errors = append(errors, ErrorResponse{
				Id:      "F00-101110",
				Name:    indexBusiness,
				Message: "BusinessId is required",
			})
		} else {
			_, err := uuid.Parse(r.BusinessId)
			if err != nil {
				errors = append(errors, ErrorResponse{
					Id:      "F00-101109",
					Name:    indexBusiness,
					Message: "BusinessId is invalid",
				})
			}
		}

		// PayeeRef required
		if r.PayeeRef == "" {
			errors = append(errors, ErrorResponse{
				Id:      "F00-101112",
				Name:    indexPayee,
				Message: "PayeeRef is required",
			})
		}

		// PayeeRef max length
		if len(r.PayeeRef) > 50 {
			errors = append(errors, ErrorResponse{
				Id:      "F00-101111",
				Name:    indexPayee,
				Message: "PayeeRef can only have maximum of 50 characters",
			})
		}

		// Duplicate check
		key := r.BusinessId + "_" + r.PayeeRef

		if duplicateMap[key] {
			errors = append(errors, ErrorResponse{
				Id:      "F00-101113",
				Name:    "AssociateRequests.PayeeRef",
				Message: fmt.Sprintf("Duplicate PayeeRef [%s] exists for the same BusinessId.", r.PayeeRef),
			})
		}

		duplicateMap[key] = true
	}

	return errors
}
func ValidateBusinessEndpoint(w http.ResponseWriter, r *http.Request) {

	var req BusinessRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	errors := ValidateBusiness(req)

	if len(errors) > 0 {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"business validated successfully"}`))
}

func ValidateBusiness(req BusinessRequest) []ErrorResponse {

	var errors []ErrorResponse

	// PayerRef validation
	if IsNullOrWhiteSpace(req.PayerRef) {
		errors = append(errors, ErrorResponse{
			Id:      "B00-20002",
			Name:    "PayerRef",
			Message: "PayerRef is required",
		})
	}

	// Email validation
	if IsNullOrWhiteSpace(req.Email) {
		errors = append(errors, ErrorResponse{
			Id:      "B00-20003",
			Name:    "Email",
			Message: "Email is required",
		})
	} else {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

		if !emailRegex.MatchString(req.Email) {
			errors = append(errors, ErrorResponse{
				Id:      "B00-20004",
				Name:    "Email",
				Message: "Email is invalid",
			})
		}
	}

	// EIN or SSN required
	if IsNullOrWhiteSpace(req.EINorSSN) {
		errors = append(errors, ErrorResponse{
			Id:      "B00-20005",
			Name:    "EINorSSN",
			Message: "EINorSSN is required",
		})
	}

	// EIN validation
	if req.IsEIN {

		if IsNullOrWhiteSpace(req.BusinessNm) {
			errors = append(errors, ErrorResponse{
				Id:      "B00-20001",
				Name:    "BusinessNm",
				Message: "Business name is required when EIN is used",
			})
		}

	} else { // SSN validation

		if IsNullOrWhiteSpace(req.FirstNm) {
			errors = append(errors, ErrorResponse{
				Id:      "B00-20006",
				Name:    "FirstNm",
				Message: "First name is required when SSN is used",
			})
		}

		if IsNullOrWhiteSpace(req.LastNm) {
			errors = append(errors, ErrorResponse{
				Id:      "B00-20007",
				Name:    "LastNm",
				Message: "Last name is required when SSN is used",
			})
		}

	}

	return errors
}

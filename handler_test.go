package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	name           string
	payload        string
	expectedStatus int
}

func executeRequest(handler http.HandlerFunc, payload string) *httptest.ResponseRecorder {

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBufferString(payload),
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	return w
}

func TestAssociateEndpoint(t *testing.T) {

	testCases := []TestCase{

		{
			name:           "F00-101108 AssociateRequests Null",
			payload:        `{}`,
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "F00-101110 BusinessId Required",
			payload: `{
			"AssociateRequests":[
			{
			"BusinessId":"",
			"PayeeRef":"PayeeRef001"
			}]}`,
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "F00-101109 BusinessId Invalid",
			payload: `{
			"AssociateRequests":[
			{
			"BusinessId":"123",
			"PayeeRef":"PayeeRef001"
			}]}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "F00-101112 PayeeRef Required",
			payload: `{
			"AssociateRequests":[
			{
			"BusinessId":"870c9cf0-eed3-4258-90f7-b07e11c8c47c",
			"PayeeRef":""
			}]}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "F00-101111 PayeeRef Max Length",
			payload: `{
			"AssociateRequests":[
			{
			"BusinessId":"870c9cf0-eed3-4258-90f7-b07e11c8c47c",
			"PayeeRef":"ABCDEFGHIJKLMNOPQRSTUfhgfhgfhgfjhfgjgfjgfujgfjghVWXYZABCDYZefgrgrehgrthrt"
			}]}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "F00-101113 Duplicate PayeeRef",
			payload: `{
			"AssociateRequests":[
			{
			"BusinessId":"870c9cf0-eed3-4258-90f7-b07e11c8c47c",
			"PayeeRef":"PayeeRef001"
			},
			{
			"BusinessId":"870c9cf0-eed3-4258-90f7-b07e11c8c47c",
			"PayeeRef":"PayeeRef001"
			}]}`,
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "Valid Request",
			payload: `{
			"AssociateRequests":[
			{
			"BusinessId":"870c9cf0-ee3e-4258-90f7-b07e11c8c47c",
			"PayeeRef":"PayeeRef001"
			}]
			}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			res := executeRequest(ValidatePayeeRef, tc.payload)

			if res.Code != tc.expectedStatus {
				t.Errorf("Expected %d but got %d", tc.expectedStatus, res.Code)
			}
		})
	}
}
func TestValidateBusinessEndpoint(t *testing.T) {

	testCases := []TestCase{

		{
			name: "PayerRef Required",
			payload: `{
			  "BusinessNm": "Snowdaze LLC",
			  "PayerRef": "y8y",
			  "IsEIN": true,
			  "EINorSSN": "23-3456789",
			  "Email": "james@sample.com"
			}`,
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "Email Required",
			payload: `{
			  "BusinessNm": "Snowdaze LLC",
			  "PayerRef": "Snow123",
			  "IsEIN": true,
			  "EINorSSN": "23-3456789",
			  "Email": ""
			}`,
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "Invalid Email",
			payload: `{
			  "BusinessNm": "Snowdaze LLC",
			  "PayerRef": "Snow123",
			  "IsEIN": true,
			  "EINorSSN": "23-3456789",
			  "Email": "invalidemail"
			}`,
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "EIN Missing Business Name",
			payload: `{
			  "BusinessNm": "",
			  "PayerRef": "Snow123",
			  "IsEIN": true,
			  "EINorSSN": "23-3456789",
			  "Email": "james@sample.com"
			}`,
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "SSN Missing First Name",
			payload: `{
			  "FirstNm": "",
			  "LastNm": "Smith",
			  "PayerRef": "Snow123",
			  "IsEIN": false,
			  "EINorSSN": "123456789",
			  "Email": "james@sample.com"
			}`,
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "SSN Missing Last Name",
			payload: `{
			  "FirstNm": "James",
			  "LastNm": "",
			  "PayerRef": "Snow123",
			  "IsEIN": false,
			  "EINorSSN": "123456789",
			  "Email": "james@sample.com"
			}`,
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "EINorSSN Required",
			payload: `{
			  "BusinessNm": "Snowdaze LLC",
			  "PayerRef": "Snow123",
			  "IsEIN": true,
			  "EINorSSN": "",
			  "Email": "james@sample.com"
			}`,
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "Valid EIN Request",
			payload: `{
			  "BusinessNm": "Snowdaze LLC",
			  "PayerRef": "Snow123",
			  "IsEIN": true,
			  "EINorSSN": "23-3456789",
			  "Email": "james@sample.com",
			  "ContactNm": "James Smith"
			}`,
			expectedStatus: http.StatusOK,
		},

		{
			name: "Valid SSN Request",
			payload: `{
			  "FirstNm": "James",
			  "LastNm": "Smith",
			  "PayerRef": "Snow123",
			  "IsEIN": false,
			  "EINorSSN": "123456789",
			  "Email": "james@sample.com"
			}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			res := executeRequest(ValidateBusinessEndpoint, tc.payload)

			if res.Code != tc.expectedStatus {
				t.Errorf("Expected %d but got %d", tc.expectedStatus, res.Code)
			}

		})
	}
}

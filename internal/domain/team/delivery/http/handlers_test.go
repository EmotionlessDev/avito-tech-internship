package http

import (
	"testing"

	"github.com/EmotionlessDev/avito-tech-internship/internal/validator"
)

func TestAddTeamRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request AddTeamRequest
		expErr  map[string]string
	}{
		{
			name: "valid request",
			request: AddTeamRequest{
				TeamName: "TeamA",
				Members: []AddTeamMemberRequest{
					{ID: "u1", Name: "Alice", IsActive: true},
				},
			},
			expErr: nil,
		},
		{
			name: "missing team name",
			request: AddTeamRequest{
				TeamName: "",
				Members: []AddTeamMemberRequest{
					{ID: "u1", Name: "Alice", IsActive: true},
				},
			},
			expErr: map[string]string{
				"team_name": "team_name is required",
			},
		},
		{
			name: "member missing id and name",
			request: AddTeamRequest{
				TeamName: "TeamB",
				Members: []AddTeamMemberRequest{
					{ID: "", Name: "", IsActive: false},
				},
			},
			expErr: map[string]string{
				"members[0].id":   "member id is required",
				"members[0].name": "member name is required",
			},
		},
		{
			name: "multiple members with some missing fields",
			request: AddTeamRequest{
				TeamName: "TeamC",
				Members: []AddTeamMemberRequest{
					{ID: "u2", Name: "", IsActive: true},
					{ID: "", Name: "Bob", IsActive: false},
				},
			},
			expErr: map[string]string{
				"members[0].name": "member name is required",
				"members[1].id":   "member id is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			tt.request.Validate(v)
			if len(v.Errors) != len(tt.expErr) {
				t.Errorf("expected %d errors, got %d", len(tt.expErr), len(v.Errors))
			}

			for key, expMsg := range tt.expErr {
				if msg, exists := v.Errors[key]; !exists {
					t.Errorf("expected error for key %q but none found", key)
				} else if msg != expMsg {
					t.Errorf("for key %q, expected message %q, got %q", key, expMsg, msg)
				}
			}
		})
	}
}

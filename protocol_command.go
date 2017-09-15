package kr

type CreateTeamRequest struct {
	Name string `json:"name"`
}

type CreateTeamResponse struct {
	PrivateKeySeed *[]byte `json:"seed,omitempty"`
	Error          *string `json:"error,omitempty"`
}
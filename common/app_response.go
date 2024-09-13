package common

type successRes struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
	Paging  interface{} `json:"paging,omitempty"`
	Filter  interface{} `json:"filter,omitempty"`
	Token   interface{} `json:"token,omitempty"`
}

func NewSuccessResponse(data, paging, filter, token interface{}) *successRes {
	return &successRes{
		Data:   data,
		Paging: paging,
		Filter: filter,
		Token:  token,
	}
}

func SimpleSuccessResponseWithToken(data, token interface{}) *successRes {
	return NewSuccessResponse(data, nil, nil, token)
}

func SimpleSuccessResponse(data interface{}) *successRes {
	return NewSuccessResponse(data, nil, nil, nil)
}

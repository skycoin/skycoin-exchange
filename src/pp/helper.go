package pp

func PtrBool(b bool) *bool {
	return &b
}

func PtrInt32(i int32) *int32 {
	return &i
}

func PtrInt64(i int64) *int64 {
	return &i
}

func PtrUint64(i uint64) *uint64 {
	return &i
}

func PtrString(s string) *string {
	return &s
}

func MakeErrRes(err error) *EmptyRes {
	res := &EmptyRes{}
	res.Result = MakeResult(ErrCode_WrongFormat, err.Error())
	return res
}

func MakeErrResWithCode(code ErrCode) *EmptyRes {
	res := &EmptyRes{}
	res.Result = MakeResultWithCode(code)
	return res
}

func MakeResult(code ErrCode, reason string) *Result {
	result := &Result{}
	result.SetCode(code)
	result.SetReason(reason)
	return result
}

func MakeResultWithCode(code ErrCode) *Result {
	return MakeResult(code, code.String())
}

func (r *Result) SetCode(code ErrCode) {
	r.Code = PtrInt32(int32(code))
	r.Success = PtrBool(code == ErrCode_Success)
	r.Reason = PtrString(code.String())
}

func (r *Result) SetReason(reason string) {
	r.Reason = PtrString(reason)
}

func (r *Result) SetCodeAndReason(code ErrCode, reason string) {
	r.SetCode(code)
	r.SetReason(reason)
}

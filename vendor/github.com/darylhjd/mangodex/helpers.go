package mangodex

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"
)

type Demographic string
type Status string
type ReadStatus string
type ContentRating string
type ListVisibility string

const (
	ShonenDemographic Demographic = "shonen"
	ShoujoDemographic Demographic = "shoujo"
	JoseiDemographic  Demographic = "josei"
	SeinenDemograpic  Demographic = "seinen"

	OngoingStatus   Status = "ongoing"
	CompletedStatus Status = "completed"
	HiatusStatus    Status = "hiatus"
	CancelledStatus Status = "cancelled"

	ReadingReadStatus    ReadStatus = "reading"
	OnHoldReadStatus     ReadStatus = "on_hold"
	PlanToReadReadStatus ReadStatus = "plan_to_read"
	DroppedReadStatus    ReadStatus = "dropped"
	ReReadingReadStatus  ReadStatus = "re_reading"
	CompletedReadStatus  ReadStatus = "completed"

	SafeRating       ContentRating = "safe"
	SuggestiveRating ContentRating = "suggestive"
	EroticaRating    ContentRating = "erotica"
	PornRating       ContentRating = "pornographic"

	PublicList  ListVisibility = "public"
	PrivateList ListVisibility = "private"
)

var (
	testClient = NewDexClient()
	user, pwd  = os.Getenv("USER"), os.Getenv("PASSWORD")
)

type ResponseType interface {
	GetResult() string
}

type Response struct {
	Result string `json:"result"`
}

func (r *Response) GetResult() string {
	return r.Result
}

// checkErrorAndResult : Helper function to check success of request by error and status code.
func checkErrorAndResult(err error, r ResponseType) error {
	switch {
	case err != nil:
		return err
	case r.GetResult() != "ok":
		return errors.New(r.GetResult())
	default:
		return nil
	}
}

// responseOp : Convenience function for simple operations that return a ResponseType.
func (dc *DexClient) responseOp(ctx context.Context, method, path string, body io.Reader, r ResponseType) error {
	u, _ := url.Parse(BaseAPI)
	u.Path = path

	// Default ResponseType will be a Response struct
	if r == nil {
		res := Response{}
		r = &res
	}

	_, err := dc.RequestAndDecode(ctx, method, u.String(), body, &r)
	return checkErrorAndResult(err, r)
}

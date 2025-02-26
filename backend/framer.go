package backend

import "github.com/kosimas/grafana-plugin-sdk-go/data"

// FrameResponse creates a DataResponse that contains the Framer's data.Frames.
func FrameResponse(f data.Framer) *DataResponse {
	frames, err := f.Frames()
	return &DataResponse{
		Frames: frames,
		Error:  err,
	}
}

// FrameResponseWithError creates a DataResponse with the error's contents (if not nil), and the Framer's data.Frames.
// This function is particularly useful if you have a function that returns `(StructX, error)`, where StructX implements Framer, which is a very common pattern.
func FrameResponseWithError(f data.Framer, err error) *DataResponse {
	if err != nil {
		return &DataResponse{
			Error: err,
		}
	}

	return FrameResponse(f)
}

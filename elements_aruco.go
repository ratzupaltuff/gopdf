package gopdf

import (
	"errors"
	"fmt"
	"image"
	"image/color"

	"github.com/raceresult/gopdf/pdf"
)

var errArucoMarkerIDOutOfRange = errors.New("aruco marker id out of range")

// ArUco4x4Element draws one marker from the OpenCV DICT_4X4_1000 set.
//
// The marker is rendered without a white quiet zone around it.
// The caller is responsible for leaving any required margin.
// For better marker distance, use the smallest ID range that fits your use case
// (for example, 0..49, 0..99, or 0..249 instead of the full 0..999 range).
type ArUco4x4Element struct {
	Left, Top, Size Length
	ID              int
	Color           Color
	Rotate          float64
	Transparency    float64
}

// NewArUco4x4Element creates a new marker element for the given marker id.
func NewArUco4x4Element(id int) *ArUco4x4Element {
	return &ArUco4x4Element{ID: id}
}

// Build adds the marker to the content stream.
func (q *ArUco4x4Element) Build(page *pdf.Page) (string, error) {
	marker, err := aruco4x4MarkerImage(q.ID)
	if err != nil {
		return "", err
	}

	page.GraphicsState_q()
	defer page.GraphicsState_Q()

	err = build2DCode(
		page,
		q.Left, q.Top, q.Size,
		q.Rotate, q.Transparency, q.Color,
		marker,
	)

	return "", err
}

func aruco4x4MarkerImage(id int) (image.Image, error) {
	if id < 0 || id >= len(aruco4x4Markers) {
		return nil, fmt.Errorf("%w: id %d (valid range 0..%d)", errArucoMarkerIDOutOfRange, id, len(aruco4x4Markers)-1)
	}

	// One module black border + inner 4x4 data pattern (no built-in white quiet zone).
	totalModules := aruco4x4MarkerSize + 2
	img := image.NewGray(image.Rect(0, 0, totalModules, totalModules))

	for y := 0; y < totalModules; y++ {
		for x := 0; x < totalModules; x++ {
			img.Set(x, y, color.Gray{Y: 0})
		}
	}

	pattern := aruco4x4Markers[id]

	for row := 0; row < aruco4x4MarkerSize; row++ {
		for col := 0; col < aruco4x4MarkerSize; col++ {
			v := color.Gray{Y: 0}
			if pattern[row*aruco4x4MarkerSize+col] == '1' {
				v = color.Gray{Y: 255}
			}

			img.Set(col+1, row+1, v)
		}
	}

	return img, nil
}

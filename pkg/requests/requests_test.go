package requests

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_requests(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	endpoint := "https://www.xinpianchang.com/api/articles"
	mockbody := `[{"id": 1, "title": "title"}]`
	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(200, mockbody),
	)

	client := New().WithClient(http.DefaultClient)

	// get text
	txt, err := client.Get("https://www.xinpianchang.com/api/articles").Do().Text()
	assert.NoError(t, err)
	assert.Equal(t, mockbody, txt)

	// get bytes
	b, err := client.Get("https://www.xinpianchang.com/api/articles").Do().Bytes()
	assert.NoError(t, err)
	assert.Equal(t, []byte(mockbody), b)
}

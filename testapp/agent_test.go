package agent_test

import (
	. "github.com/haoxins/supertest"
	"testing"
)

//TestAgent test amp-agent using its api health, status 200 means agent is reading its dockers events stream
func TestAgent(t *testing.T) {
	Request("http://localhost:5001").
		Get("/api/v1/health").
		Expect(200)
}

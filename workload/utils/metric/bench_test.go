package metric

import (
	"encoding/json"
	. "github.com/pingcap/check"
	"testing"
	"time"
)

type BenchTestSuite struct{}

var _ = Suite(&BenchTestSuite{})

func TestT(t *testing.T) {
	TestingT(t)
}

func (BenchTestSuite) TestBasicFormat(c *C) {
	duration := Duration(32 * time.Hour)
	c.Assert(duration.String(), Equals, "32h0m0s")
	size := Size(1024*1024*1024*42 - 1024*1024*128)
	c.Assert(size.String(), Equals, "41.9G")
	float := Float(3.1415926535)
	c.Assert(float.String(), Equals, "3.14")
}

func (BenchTestSuite) TestJSONFormat(c *C) {
	items := Report{
		ComponentName: "check",
		WorkloadName:  "testing",
		Metrics:       map[string]Metric{},
	}
	items.Add(Size(1234).Named("backup speed"))
	items.Add(Duration(24 * time.Hour).Named("time of a day"))
	result, err := json.Marshal(items)
	c.Assert(err, IsNil)
	c.Assert(string(result), Equals, `{"component":"check","workload":"testing","metrics":{"backup speed":{"name":"backup speed","value":{"Type":1,"Value":1234}},"time of a day":{"name":"time of a day","value":{"Type":2,"Value":86400000000000}}}}`)

	var itemsUnmarshalled Report
	c.Assert(json.Unmarshal(result, &itemsUnmarshalled), IsNil)
	c.Assert(itemsUnmarshalled, DeepEquals, items)
}

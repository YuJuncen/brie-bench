package metric

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/pingcap/check"
	"strings"
	"testing"
	"time"
)

type BenchTestSuite struct{}

type containsChecker struct{}

var containsAll = containsChecker{}

func (c containsChecker) Info() *CheckerInfo {
	return &CheckerInfo{
		Name:   "Contains",
		Params: []string{"haystack", "needles"},
	}
}

func (c containsChecker) Check(params []interface{}, _ []string) (result bool, error string) {
	haystack := params[0].(string)

	for _, needle := range params[1].([]string) {
		if !strings.Contains(haystack, needle) {
			return false, needle + " is not contained."
		}
	}
	return true, ""
}

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
	println(string(result))
	c.Assert(err, IsNil)
	c.Assert(string(result), Equals, `{"component":"check","hash":"","repo":"","workload":"testing"`+
		`,"metrics":{"backup speed":{"name":"backup speed","value":{"Type":1,"Value":1234}},`+
		`"time of a day":{"name":"time of a day","value":{"Type":2,"Value":86400000000000}}}}`)

	var itemsUnmarshalled Report
	c.Assert(json.Unmarshal(result, &itemsUnmarshalled), IsNil)
	c.Assert(itemsUnmarshalled, DeepEquals, items)
}

func (BenchTestSuite) TestReport(c *C) {
	items := Report{
		ComponentName: "check",
		WorkloadName:  "testing",
		ComponentHash: "cafebabe0000000000",
		ComponentRepo: "github.com/nocode/nocode",
		Metrics:       map[string]Metric{},
	}
	items.Add(Size(1234).Named("backup speed"))
	items.Add(Duration(24 * time.Hour).Named("time of a day"))
	items.Add(Float(36.9).Named("really lllllllllllllllllllllllllllllllllong name"))
	exported := bytes.NewBuffer(nil)
	err := items.Export(exported)
	c.Assert(err, IsNil)
	result := exported.String()
	fmt.Println(result)
	c.Assert(result, containsAll, []string{"time of a day", "24h",
		"backup speed", "1.2K",
		"really lllllllllllllllllllllllllllllllllong name", "36.9"})
}

func (BenchTestSuite) TestComparingReport(c *C) {
	items := Report{
		ComponentName: "check",
		WorkloadName:  "testing",
		ComponentHash: "cafebabe0000000000",
		ComponentRepo: "github.com/nocode/nocode",
		Metrics:       map[string]Metric{},
	}
	items.Add(Size(1234).Named("backup speed"))
	items.Add(Duration(24 * time.Hour).Named("time of a day"))
	items.Add(Float(36.9).Named("really lllllllllllllllllllllllllllllllllong name"))

	othItems := Report{
		ComponentName: "check",
		WorkloadName:  "testing-patched",
		ComponentHash: "7c3fed",
		ComponentRepo: "github.com/alyssaPHacker/nocode",
		Metrics:       map[string]Metric{},
	}
	othItems.Add(Size(4321).Named("backup speed"))
	othItems.Add(Duration(22 * time.Hour).Named("time of a day"))
	othItems.Add(Float(36.9).Named("really lllllllllllllllllllllllllllllllllong name"))
	othItems.Add(Size(123).Named("new test case"))

	exported := bytes.NewBuffer(nil)
	err := othItems.ExportComparing(exported, &items)
	c.Assert(err, IsNil)

	result := exported.String()
	fmt.Println(result)
	c.Assert(result, containsAll, []string{
		"testing-patched", "7c3fed", "github.com/alyssaPHacker/nocode",
		"time of a day", "22h",
		"backup speed", "4.2K",
		"really lllllllllllllllllllllllllllllllllong name", "36.9",
		"⬆️", "⬇️", "⚠️",
	})
}

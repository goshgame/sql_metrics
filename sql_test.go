package sqlmetrics

import "testing"

func TestExtractTable(t *testing.T) {
	talbeName := extractTable("SELECT id FROM user WHERE id = 23333")
	t.Logf("tableName=%s\n", talbeName)
	tableName2 := extractTable(`SELECT DISTINCT os as platform, source_channel, count(*) as h5_pv, count(DISTINCT adid) as h5_uv 
		FROM myTable2 WHERE created_at >= 123 AND created_at <= 456 AND event = ? GROUP BY os, source_channel;`)
	t.Logf("tableName=%s\n", tableName2)

}

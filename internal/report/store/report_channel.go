package store

type Channel interface {
	AddReport(reportId string)
	GetReportChannel() chan string
	Close() error
}

type channel struct {
	reportCh chan string
}

func NewReportChannel() Channel {
	reportCh := make(chan string)
	return &channel{reportCh: reportCh}
}

func (c *channel) AddReport(reportId string) {
	c.reportCh <- reportId
}

func (c *channel) GetReportChannel() chan string {
	return c.reportCh
}

func (c *channel) Close() error {
	close(c.reportCh)
	return nil
}

package store

type Channel interface {
	AddReport(reportID string)
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

func (c *channel) AddReport(reportID string) {
	c.reportCh <- reportID
}

func (c *channel) GetReportChannel() chan string {
	return c.reportCh
}

func (c *channel) Close() error {
	close(c.reportCh)
	return nil
}

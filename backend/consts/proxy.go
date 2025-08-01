package consts

type ReportAction string

const (
	ReportActionAccept       ReportAction = "accept"
	ReportActionSuggest      ReportAction = "suggest"
	ReportActionFileWritten  ReportAction = "file_written"
	ReportActionReject       ReportAction = "reject"
	ReportActionNewTask      ReportAction = "new_task"
	ReportActionFeedbackTask ReportAction = "feedback_task"
	ReportActionAbortTask    ReportAction = "abort_task"
)

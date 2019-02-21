package common

import (
	"context"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

type Job struct {
	Name     string `json:"name"`     // 任务名
	Command  string `json:"command"`  // shell命令
	CronExpr string `json:"cronExpr"` // cron表达式
}

// 任务调度计划
type JobSchedulePlan struct {
	Job *Job	// 要调度的任务信息
	Expr *cronexpr.Expression	// 解析好的cronexpr表达式
	NextTime time.Time	// 下次调度时间
}

// 任务执行状态
type JobExecuteInfo struct {
	Job *Job // 任务信息
	PlanTime time.Time // 理论上的调度时间
	RealTime time.Time // 实际的调度时间
	CancelCtx context.Context // 任务command的context
	CancelFunc context.CancelFunc//  用于取消command执行的cancel函数
}

// 任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo	// 执行状态
	Output []byte // 脚本输出
	Err error // 脚本错误原因
	StartTime time.Time // 启动时间
	EndTime time.Time // 结束时间
}

type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

type JobEvent struct {
	EventType int
	Job *Job
}

func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	// 定义一个response
	var (
		response Response
	)
	response.Errno = errno
	response.Msg = msg
	response.Data = data

	resp, err = json.Marshal(response)
	return
}

// 反序列化job
func UnpackJob(value []byte) (ret *Job, err error)  {
	var (
		job *Job
	)

	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	ret = job
	return
}

// 从etcd的key中提取任务名
func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

func ExtractKillerName(killerKey string) (string) {
	return strings.TrimPrefix(killerKey, JOB_KILLER_DIR)
}

// 构造Event - 更新/删除
func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent {
		EventType: eventType,
		Job:job,
	}
}

// 构造执行计划
func BuildSchedulePlan(job *Job)(jobSchedulePlan *JobSchedulePlan, err error) {
	var (
		expr *cronexpr.Expression
	)
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}
	jobSchedulePlan = &JobSchedulePlan{
		Job:job,
		Expr:expr,
		NextTime:expr.Next(time.Now()),
	}
	return
}

// 构造执行状态信息
func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo){
	jobExecuteInfo = &JobExecuteInfo{
		Job: jobSchedulePlan.Job,
		PlanTime: jobSchedulePlan.NextTime, // 计算调度时间
		RealTime: time.Now(), // 真实调度时间
	}
	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())
	return
}
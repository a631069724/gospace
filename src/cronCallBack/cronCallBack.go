package main

/*
#include<stdio.h>
#include<stdlib.h>
#include<string.h>
typedef int (*ptfFuncReportData)(int id);
extern int CReportData(ptfFuncReportData pf,int id);
*/
import "C"
import (
	"github.com/jakecoffman/cron"
	"strconv"
)

type CallBackFunc C.ptfFuncReportData
type SchduleActiveCallback func(int)

type ScheduleJob struct {
	id       int
	callback SchduleActiveCallback
}

var g_scheduleManager *ScheduleManager
var g_callBackFunc CallBackFunc

func (this *ScheduleJob) Run() {
	if nil != this.callback {
		this.callback(this.id)
	}
}

type ScheduleManager struct {
	cronJob *cron.Cron
}

func (this *ScheduleManager) NewScheduleJob(_id int, _job SchduleActiveCallback) *ScheduleJob {
	instance := &ScheduleJob{
		id:       _id,
		callback: _job,
	}
	return instance
}
func NewScheduleManager() *ScheduleManager {
	instance := &ScheduleManager{}
	instance.cronJob = cron.New()
	return instance
}

//export CronStart
func CronStart() {
	g_scheduleManager.cronJob.Start()
}

//export Stop
func Stop() {
	g_scheduleManager.cronJob.Stop()
}

//export AddJob
func AddJob(id int, scheduleExpr string) {
	job := g_scheduleManager.NewScheduleJob(id, g_scheduleManager._scheduleActive)
	g_scheduleManager.cronJob.AddJob(scheduleExpr, job, strconv.Itoa(id))
}

//export RemoveJob
func RemoveJob(id int) {
	g_scheduleManager.cronJob.RemoveJob(strconv.Itoa(id))
}

//export CronInit
func CronInit(f CallBackFunc) {
	g_scheduleManager = NewScheduleManager()
	SetCallBack(f)
}

func SetCallBack(f CallBackFunc) {
	g_callBackFunc = f
}

func (this *ScheduleManager) _scheduleActive(id int) {
	C.CReportData(g_callBackFunc, C.int(id))
}

func main() {}

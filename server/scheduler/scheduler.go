package scheduler

import (
	"context"
	"fmt"
	"genesis"
	"runtime/debug"
	"time"

	"github.com/prometheus/common/log"
)

// ----------------
//    How to use
// ----------------
/*
	main.go

	// init new scheduler
	ss := scheduler.New(
		"MeterTransferDaemon", // task name
		3,      // run every 3 second
		false   // run immediately when initiated
	)

	// create struct for custom task
	mtat := scheduler.MeterTransferApprovalTask{
		Conn: db,
		...
	}

	// register task to scheduler
	ss.TaskRegister(mtat.Runner)

	// start scheduler
	ss.Start()
*/

// Scheduler contain the information that the scheduler inner settings
type Scheduler struct {
	Name         string        // name of the scheduled task
	timeInterval time.Duration // what interval it will trigger
	ticker       *time.Ticker
	task         Task // task to do, int is the number of action/change/touch/created/update/delete performed
	taskSet      bool
	startZero    bool // start the task immediately scheduler is start
	count        int  // number of times this been triggered
	countFail    int  // number of failed trigger
	countSuccess int  // number of successful trigger

	// not implemented
	debug       bool
	TimeOut     time.Duration // what time the scheduler should give up the single task
	TimeStart   time.Time     // what time the scheduler should start
	TimeStop    time.Time     // what time the scheduler should stopped
	timeLast    time.Time     // last time the scheduler ran
	intervalMax int           // maximum number of interval it will trigger
}

// Task is the parent struct for the schedule task
type Task func() (int, error)

// Start will begin the scheduler
func (sc *Scheduler) Start() {
	if sc.task == nil {
		log.Errorln(" Err: no task registered")
		return
	}

	log.Info("Start scheduler")

	// run first time, then ticker next
	if sc.startZero {
		sc.TaskRun()
	}

	sc.ticker = time.NewTicker(sc.timeInterval)
	done := make(chan bool, 1)
	go func(t *time.Ticker) {
		for {
			select {
			case <-t.C:
				sc.TaskRun()
			case <-done:
				log.Infof("scheduler exit. %s\n", sc.Name)
				return
			}
		}
	}(sc.ticker)
}

// TaskRun execute the function (task) it been assigned to
func (sc *Scheduler) TaskRun() {
	// recover from panic
	defer func() {
		if rec := recover(); rec != nil {
			message := "Scheduler task panicked (" + sc.Name + ")"
			log.Errorln(message)
			strStack := string(debug.Stack())

			var err error
			switch v := rec.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf(message)
			}
			ctx := context.Background()
			genesis.SentrySend(ctx, nil, nil, err, strStack)

			log.Errorln("Scheduler panic recovered ("+sc.Name+"): ", err, "\n", strStack)
		}
	}()

	log.Infof("scheduler task run (%s)\n", sc.Name)
	defer func() {
		log.Infof("scheduler task exit (%s)\n", sc.Name)
	}()

	_, err := sc.task()
	if err != nil {
		log.Errorln(" Err: schedule task failed", err)
		sc.countFail++
	} else {
		sc.countSuccess++
	}
	// inc by 1
	sc.count++
}

// TaskRegister adds the function (task) for running when called
func (sc *Scheduler) TaskRegister(task Task) {
	sc.task = task
	sc.taskSet = true
}

// TaskRegisterAndStart adds the function (task) and start immediately
func (sc *Scheduler) TaskRegisterAndStart(task Task) {
	sc.task = task
	sc.taskSet = true
	sc.Start()
}

// Stop will halt the scheduler
func (sc *Scheduler) Stop() {
	log.Info("Start scheduler")

	sc.ticker.Stop()
}

// New will return a new scheduler
func New(
	name string,
	timeSecond int,
	startZero bool,
) *Scheduler {
	if timeSecond < 10 {
		panic("cannot be less than 10 seconds for interval")
	}

	var ti time.Duration
	ti = time.Second * time.Duration(timeSecond)

	return &Scheduler{
		Name:         name,
		timeInterval: ti,
		startZero:    startZero,
		intervalMax:  2147483647, // ~68 years if triggger every second
	}
}

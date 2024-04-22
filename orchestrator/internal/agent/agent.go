package agent

import (
	"errors"
	"fmt"
	"orchestrator/internal/database"
	"orchestrator/internal/env"
	"orchestrator/internal/expression"
	"orchestrator/internal/logger"
	"orchestrator/internal/worker"
)

type Agent struct {
	ind    int
	Tasks  []chan expression.Expression
	IsFree []bool
	Work   []expression.Expression
}

func (ag *Agent) GetAll() []string { // Инфа для страницы /computers
	all := make([]string, 0)
	for ind, i := range ag.IsFree {
		if i {
			all = append(all, fmt.Sprintf("Воркер %d не работает", ind))
		} else {
			all = append(all, fmt.Sprintf("Воркер %d работает над выражением: %s, id: %d", ind, ag.Work[ind].Name, ag.Work[ind].Id))
		}
	}
	return all
}

func NewAgent() *Agent {
	tasks := make([]chan expression.Expression, 0)
	IsFree := make([]bool, 0)
	Work := make([]expression.Expression, 0)
	for i := 0; i < env.Workers; i++ {
		tasks = append(tasks, make(chan expression.Expression))

		worker.StartWorker(tasks[i])
		IsFree = append(IsFree, true)
		Work = append(Work, expression.Expression{})

	}
	return &Agent{Tasks: tasks, IsFree: IsFree, Work: Work}
}

func (ag *Agent) AddTask(expr expression.Expression) error { // Создание задания
	for ind, i := range ag.IsFree {
		if i {
			logger.Log.Println("Выражение принято на обработку, id:", expr)
			ag.Tasks[ind] <- expr
			ag.Work[ind] = expr
			expr.Status = 1
			database.UpdateExpr(expr)
			ag.IsFree[ind] = false
			go func() {
				newexp := <-ag.Tasks[ind]
				newexp.Status = 0
				database.UpdateExpr(newexp)
				ag.IsFree[ind] = true
				ag.Work[ind] = expression.Expression{}
			}()
			return nil
		}
	}
	return errors.New("Все работают")
}

package workers

import (
	"context"
	"log"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/riverqueue/river"
)

type VacuumUserTokensWorker struct {
	river.WorkerDefaults[jobs.VacuumUserTokensArgs]
	db queries.DBTX
}

func (w *VacuumUserTokensWorker) Work(ctx context.Context, job *river.Job[jobs.VacuumUserTokensArgs]) error {
	log.Print("VacuumUserTokens job starting")
	return queries.New(w.db).VacuumUserTokens(ctx)
}

package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func readCh(stageCh Out) {
	go func() {
		for range stageCh { //nolint:revive
		}
	}()
}

func writeRes(resCh Bi, done In, stageCh Out) {
	defer close(resCh)
	defer readCh(stageCh)

	for {
		select {
		case <-done:

			return

		case val, ok := <-stageCh:
			if !ok {
				return
			}

			resCh <- val
		}
	}
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		stageCh := stage(out)
		resCh := make(Bi)

		go writeRes(resCh, done, stageCh)

		out = resCh
	}

	return out
}

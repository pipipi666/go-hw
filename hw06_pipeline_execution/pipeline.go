package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		stageCh := stage(out)
		resCh := make(Bi)

		go func() {
			defer close(resCh)

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
		}()

		out = resCh
	}

	return out
}

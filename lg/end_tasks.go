package lg

type EndTasks struct {
	tasks []func()
}

func (t *EndTasks) Append(task func()) {
	t.tasks = append(t.tasks, task)
}

func (t *EndTasks) Execute() {
	for _, task := range t.tasks {
		task()
	}
}

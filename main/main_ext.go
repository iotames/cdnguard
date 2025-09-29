package main

func extCmdRun() bool {
	if Prune {
		if err := gdb.Prune(); err != nil {
			panic(err)
		}
		return true
	}
	if Debug {
		debug()
		return true
	}
	return false
}

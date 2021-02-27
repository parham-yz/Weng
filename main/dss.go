package main

type void interface{}

// <?>
type Queue struct {
	waiters          []void
	first, last, len int
	empty, full      bool
}

func (qu *Queue) init(size uint) {
	qu.waiters = make([]void, size)
	qu.first = 0
	qu.last = 0
	qu.empty = true
	qu.full = false
	qu.len = 0
	return
}

func (qu *Queue) PushTo(other *Queue) bool {
	for !qu.empty {
		ok := other.Push(qu.Pop())
		if !ok {
			return true
		}
	}
	qu.len += other.len
	return false
}

func (qu *Queue) Push(elment void) bool {
	if qu.full {
		return true
	}
	qu.waiters[qu.last] = elment
	qu.last = (qu.last + 1) % cap(qu.waiters)

	if qu.empty {
		qu.empty = false
	}

	if qu.last == qu.first {
		qu.full = true
	}

	qu.len++
	return false
}

func (qu *Queue) Pop() (elment void) {
	if qu.empty {
		return nil
	}
	elment = qu.waiters[qu.first]
	qu.first = (qu.first + 1) % cap(qu.waiters)

	if qu.full {
		qu.full = false
	}

	if qu.last == qu.first {
		qu.empty = true
	}
	qu.len--
	return
}

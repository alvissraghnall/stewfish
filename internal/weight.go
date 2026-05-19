package internal

type Weight struct {
	middlegame int16
	endgame    int16
}

func W(mg int16, eg int16) Weight {
	return Weight{middlegame: mg, endgame: eg}
}

func (w Weight) Mg() int16 {
	return w.middlegame
}

func (w Weight) Eg() int16 {
	return w.endgame
}

func (w *Weight) Add(other Weight) {
	w.middlegame += other.Mg()
	w.endgame += other.Eg()
}

func (w *Weight) Sub(other Weight) {
	w.middlegame -= other.Mg()
	w.endgame -= other.Eg()
}

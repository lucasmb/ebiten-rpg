package animations

type Animation struct {
	First        int
	Last         int
	Step         int     //how many indeces do we move per frame
	SpeedinTps   float32 //how many ticks before next frame
	FrameCounter float32
	frame        int
}

func (a *Animation) Update() {

	a.FrameCounter -= 1.0
	if a.FrameCounter < 0.0 {
		a.FrameCounter = a.SpeedinTps
		a.frame += a.Step
		if a.frame > a.Last {
			//loop to the begining
			a.frame = a.First
		}
	}
}

func (a *Animation) Frame() int {
	return a.frame
}
func NewAnimation(first, last, step int, speed float32) *Animation {
	return &Animation{
		first,
		last,
		step,
		speed,
		speed,
		first,
	}

}

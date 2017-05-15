package game

const GravitationalConstant = 6.67408e-11

type Game struct {
	bodies []*Body
}

func New() *Game {
	return &Game{}
}

func (g *Game) AddBody(b *Body) {
	g.bodies = append(g.bodies, b)
}

func (g *Game) GravitationalForce(p Vector, m float64) Vector {
	ret := Vector{}
	for _, body := range g.bodies {
		ret = ret.Add(body.Position.Sub(p).Scale(body.Mass / body.Position.DistanceSquared(p)))
	}
	return ret.Scale(GravitationalConstant * m)
}

package game

import (
	"time"
)

type BodyId = int

type Universe struct {
	bounds Rect
	bodies map[BodyId]*Body
	nextId BodyId
}

func NewUniverse(bounds Rect) *Universe {
	return &Universe{
		bounds: bounds,
		bodies: make(map[BodyId]*Body),
	}
}

func (u *Universe) AddBody(b *Body) BodyId {
	id := u.nextId
	u.nextId++
	u.bodies[id] = b
	return id
}

func (u *Universe) GetBody(id BodyId) *Body {
	return u.bodies[id]
}

func (u *Universe) RemoveBody(id BodyId) {
	delete(u.bodies, id)
}

func (u *Universe) Step(d time.Duration) {
	u.DecayBodies()
	u.checkCollisions()
	u.applyForces()
	for _, b := range u.bodies {
		b.Step(d)
	}
}

func (u *Universe) DecayBodies() {
	for _, b := range u.bodies {
		if !u.bounds.Contains(b.Position) {
			b.ForceDecay(outOfBoundsDecayPerStep)
		} else {
			b.Decay(decayPerStep)
		}
	}

	idsToRemove := make([]BodyId, 0)
	for id, b := range u.bodies {
		if b.Mass == 0 {
			idsToRemove = append(idsToRemove, id)
		}
	}
	for id := range idsToRemove {
		u.RemoveBody(id)
	}
}

func (u *Universe) checkCollisions() {
	//  checkCollisions() {
	//    for (var i in this.bodies) {
	//      for (var j in this.bodies) {
	//        if (i <= j) {
	//          // We only need to perform a single comparison for every (i, j) pair
	//          // and only if i != j.
	//          continue;
	//        }
	//        if (this.bodies[i].collidesWith(this.bodies[j])) {
	//          if (this.bodies[i].mass > this.bodies[j].mass) {
	//            this.bodies[i].mergeWith(this.bodies[j]);
	//            this.removeBody(j);
	//          } else {
	//            this.bodies[j].mergeWith(this.bodies[i]);
	//            this.removeBody(i);
	//          }
	//          // need to reevaluate everything
	//          this.checkCollisions();
	//          return;
	//        }
	//      }
	//    }
	//  }
}

func (u *Universe) applyForces() {
	//    for (var i in this.bodies) {
	//      var netForces = [];
	//      netForces.push()
	//
	//      for (var j in this.bodies) {
	//        if (i == j) {
	//          continue;
	//        }
	//        netForces.push(this.bodies[i].gravitationalForceTo(this.bodies[j]));
	//      }
	//
	//      this.bodies[i].gravitationalForce = netForces.reduce(function(carry, force) {
	//        return {
	//          x: carry.x + force.x,
	//          y: carry.y + force.y,
	//        };
	//      }, {x: 0, y: 0});
	//    }
	//  }
}

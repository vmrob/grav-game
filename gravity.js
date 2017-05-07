const GRAVITATIONAL_CONSTANT = 1000;

var Distance = function(x1, x2, y1, y2) {
  return Math.sqrt(Math.pow(x1 - x2, 2) + Math.pow(y1 - y2, 2));
}

class Vector {
  constructor(x, y) {
    this.x = x;
    this.y = y;
  }

  withMagnitude(m) {
    var ret = new Vector(this.x, this.y);
    var scale = m / this.magnitude;
    ret.x *= scale;
    ret.y *= scale;
    return ret;
  }

  get magnitude() {
    return Distance(0, this.x, 0, this.y);
  }
}

class Universe {
  constructor(bodies) {
    this.bodies = bodies
  }

  addBody(body) {
    this.bodies.push(body);
  }

  step(duration) {
    this.applyForces();
    // this.applyCollisions();
    for (var i in this.bodies) {
      this.bodies[i].step(duration);
    }
  }

  applyForces() {
    for (var i in this.bodies) {
      var netForces = [];
      netForces.push()

      for (var j in this.bodies) {
        if (i == j) {
          continue;
        }
        netForces.push(this.bodies[i].gravitationalForceTo(this.bodies[j]));
      }

      this.bodies[i].force = netForces.reduce(function(carry, force) {
        return {
          x: carry.x + force.x,
          y: carry.y + force.y,
        };
      }, {x: 0, y: 0});
    }
  }

  // applyCollisions() {
  //   for (i = 0; i < this.bodies.length - 1; ++i) {
  //     for (j = i + 1; j < this.bodies.length; ++j) {
  //       if (this.bodies[i].collidesWith(this.bodies[j])) {
  //
  //       }
  //     }
  //   }
  // }

  draw() {
    for (var i in this.bodies) {
      this.bodies[i].draw(context);
    }
  }
};

class Body {
  constructor(mass, pos, velocity) {
    this.pos = pos;
    this.velocity = velocity;
    this.mass = mass;
    this.force = {
      x: 0,
      y: 0,
    };
  }

  get radius() {
    return Math.sqrt(this.mass / Math.PI); // assume mass == area
  }

  draw(canvasContext) {
    canvasContext.beginPath();
    canvasContext.arc(this.pos.x, this.pos.y, this.radius, 0, 2 * Math.PI);
    canvasContext.fillStyle = 'green';
    canvasContext.fill();
    canvasContext.lineWidth = 5;
    canvasContext.strokeStyle = '#003300';
    canvasContext.stroke();
  }

  step(duration) {
    this.velocity = {
      x: this.velocity.x + this.force.x / this.mass * duration,
      y: this.velocity.y + this.force.y / this.mass * duration,
    };
    this.pos = {
      x: this.pos.x + this.velocity.x * duration,
      y: this.pos.y + this.velocity.y * duration,
    };
  }

  distanceTo(body) {
    return Math.sqrt(Math.pow(this.pos.x - body.pos.x, 2) + Math.pow(this.pos.y - body.pos.y, 2));
  }

  vectorTo(body) {
    return new Vector(body.pos.x - this.pos.x, body.pos.y - this.pos.y);
  }

  gravitationalForceTo(body) {
    var vec = this.vectorTo(body);
    var force = GRAVITATIONAL_CONSTANT * this.mass * body.mass / Math.pow(this.distanceTo(body), 2);
    return vec.withMagnitude(force);
  }

  collidesWith(body) {
    return distanceTo(body) < this.radius + body.radius;
  }
};

context = document.getElementById('gameCanvas').getContext("2d");

var universe = new Universe([
  new Body(1000, {x: 200, y: 200}, {x: 50, y: 0}),
  new Body(1000, {x: 300, y: 300}, {x: -50, y: 0}),
]);

window.setInterval(function() {
  context.clearRect(0, 0, context.canvas.width, context.canvas.height);
  universe.draw();
}, 1000 / 30);

window.setInterval(function() {
  universe.step(1 / 60);
}, 1000 / 60);

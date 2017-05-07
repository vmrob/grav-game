const GRAVITATIONAL_CONSTANT = 1000;
const THRUSTER_FORCE = 100000;

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
  constructor(label, color, mass, pos, velocity) {
    this.label = label;
    this.color = color;
    this.pos = pos;
    this.velocity = velocity;
    this.mass = mass;
    this.leftThrusterEnabled = false;
    this.rightThrusterEnabled = false;
    this.bottomThrusterEnabled = false;
    this.topThrusterEnabled = false;
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
    canvasContext.fillStyle = this.color;
    canvasContext.fill();
    canvasContext.lineWidth = 5;
    canvasContext.strokeStyle = '#003300';
    canvasContext.stroke();
    canvasContext.textAlign = 'center';
    canvasContext.fillText(this.label, this.pos.x, this.pos.y + this.radius + 14);
  }

  step(duration) {
    var force = new Vector(this.force.x, this.force.y);
    if (this.leftThrusterEnabled) {
        force.x += THRUSTER_FORCE;
    }
    if (this.rightThrusterEnabled) {
        force.x -= THRUSTER_FORCE;
    }
    if (this.topThrusterEnabled) {
        force.y += THRUSTER_FORCE;
    }
    if (this.bottomThrusterEnabled) {
        force.y -= THRUSTER_FORCE;
    }
    this.velocity = {
      x: this.velocity.x + force.x / this.mass * duration,
      y: this.velocity.y + force.y / this.mass * duration,
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

var player1Body = new Body('Player 1', '#cfcf80', 1000, {x: 200, y: 200}, {x: 50, y: 0});
var player2Body = new Body('Player 2', '#80cfcf', 1000, {x: 300, y: 300}, {x: -50, y: 0});
var universe = new Universe([player1Body, player2Body]);

window.setInterval(function() {
  context.clearRect(0, 0, context.canvas.width, context.canvas.height);
  universe.draw();
}, 1000 / 30);

window.setInterval(function() {
  universe.step(1 / 60);
}, 1000 / 60);

$(function() {
    $(document).keydown(function(e) {
        switch (e.which) {
            case 65: // a
                player1Body.rightThrusterEnabled = true;
                break;
            case 87: // w
                player1Body.bottomThrusterEnabled = true;
                break;
            case 68: // d
                player1Body.leftThrusterEnabled = true;
                break;
            case 83: // s
                player1Body.topThrusterEnabled = true;
                break;
            case 37: // left
                player2Body.rightThrusterEnabled = true;
                break;
            case 38: // up
                player2Body.bottomThrusterEnabled = true;
                break;
            case 39: // right
                player2Body.leftThrusterEnabled = true;
                break;
            case 40: // down
                player2Body.topThrusterEnabled = true;
                break;
            default:
                return;
        }
        e.preventDefault();
    });

    $(document).keyup(function(e) {
        switch (e.which) {
            case 65: // a
                player1Body.rightThrusterEnabled = false;
                break;
            case 87: // w
                player1Body.bottomThrusterEnabled = false;
                break;
            case 68: // d
                player1Body.leftThrusterEnabled = false;
                break;
            case 83: // s
                player1Body.topThrusterEnabled = false;
                break;
            case 37: // left
                player2Body.rightThrusterEnabled = false;
                break;
            case 38: // up
                player2Body.bottomThrusterEnabled = false;
                break;
            case 39: // right
                player2Body.leftThrusterEnabled = false;
                break;
            case 40: // down
                player2Body.topThrusterEnabled = false;
                break;
            default:
                return;
        }
        e.preventDefault();
    });
});

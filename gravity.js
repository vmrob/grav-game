const GRAVITATIONAL_CONSTANT = 100;
const THRUSTER_BASE_FORCE = 1000000;
const PLAYER_START_MASS = 10000;
const DECAY_PER_STEP = 0.0001;
const MINIMUM_DECAY_MASS = PLAYER_START_MASS;
const MINIMUM_DECAY_FORCE_QUANTITY = 500;
const GAME_BOUNDS_X = -5000;
const GAME_BOUNDS_WIDTH = 10000;
const GAME_BOUNDS_Y = -5000;
const GAME_BOUNDS_HEIGHT = 10000;
const PLAYER_1_COLOR = '#cfcf80';
const PLAYER_2_COLOR = '#80cfcf';
const MATTER_REPULSION_FORCE = 100000;

var Distance = function(x1, x2, y1, y2) {
  return Math.sqrt(Math.pow(x1 - x2, 2) + Math.pow(y1 - y2, 2));
}

function RandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1) + min);
}

function RandomPosition(width, height) {
  return {
    x: RandomInt(0, width),
    y: RandomInt(0, height),
  };
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

class Rect {
  constructor(x, y, width, height) {
    this.x = x;
    this.y = y;
    this.width = width;
    this.height = height;
  }

  contains(point) {
    return point.x >= this.x && point.x <= this.x + this.width &&
      point.y >= this.y && point.y <= this.y + this.height;
  }
}

class Universe {
  constructor(bodies, boundsRect) {
    this.bodies = bodies;
    this.bounds = boundsRect;
  }

  addBody(body) {
    this.bodies.push(body);
  }

  step(duration) {
    this.decayBodies();
    this.checkCollisions();
    this.applyForces();
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

      this.bodies[i].gravitationalForce = netForces.reduce(function(carry, force) {
        return {
          x: carry.x + force.x,
          y: carry.y + force.y,
        };
      }, {x: 0, y: 0});
    }
  }

  decayBodies() {
    for (var i in this.bodies) {
      if (!this.bounds.contains(this.bodies[i].pos)) {
        this.bodies[i].decay(DECAY_PER_STEP * 10, true);
      } else {
        this.bodies[i].decay(DECAY_PER_STEP);
      }
    }
    for (var i = 0; i < this.bodies.length;) {
      if (this.bodies[i].mass == 0) {
        this.removeBody(i);
      } else {
        ++i;
      }
    }
  }

  checkCollisions() {
    for (var i = 0; i < this.bodies.length - 1; ++i) {
      for (var j = i + 1; j < this.bodies.length; ++j) {
        if (this.bodies[i].collidesWith(this.bodies[j])) {
          if (this.bodies[i].mass > this.bodies[j].mass) {
            this.bodies[i].mergeWith(this.bodies[j]);
            this.removeBody(j);
          } else {
            this.bodies[j].mergeWith(this.bodies[i]);
            this.removeBody(i);
          }
          // need to reevaluate everything
          this.checkCollisions();
          return;
        }
      }
    }
  }

  removeBody(index) {
    this.bodies.splice(index, 1);
  }

  draw() {
    var min = new Vector(0, 0);
    var max = new Vector(context.canvas.width, context.canvas.height);
    var padding = 10;
    for (var i in this.bodies) {
      if (this.bodies[i].static || this.bodies[i].npc) {
        continue;
      }
      var pos = this.bodies[i].pos;
      var r = this.bodies[i].radius;
      if (pos.x - r - padding < min.x) {
          min.x = pos.x - r - padding;
      }
      if (pos.y - r - padding < min.y) {
          min.y = pos.y - r - padding;
      }
      if (pos.x + r + padding > max.x) {
          max.x = pos.x + r + padding;
      }
      if (pos.y + r + padding > max.y) {
          max.y = pos.y + r + padding;
      }
    }
    context.clearRect(0, 0, context.canvas.width, context.canvas.height);
    var scaleX = context.canvas.width / (max.x - min.x);
    var scaleY = context.canvas.height / (max.y - min.y);
    var scale = scaleX > scaleY ? scaleY : scaleX;
    context.scale(scale, scale);
    context.translate(-min.x, -min.y);
    for (var i in this.bodies) {
      this.bodies[i].draw(context);
    }
    this.drawBounds(context);

    context.translate(min.x, min.y);
    context.scale(1.0 / scale, 1.0 / scale);
  }

  drawBounds(context) {
    context.rect(this.bounds.x, this.bounds.y, this.bounds.width, this.bounds.height);
    context.stroke();
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
    this.gravitationalForce = {
      x: 0,
      y: 0,
    };
    this.static = false;
    this.npc = false;
    this.cooldownActivationTime = 0;
  }

  get radius() {
    var r = Math.cbrt((this.mass * 3) / (4 * Math.PI));
    if (MATTER_REPULSION_FORCE < this.forceAtDistance(r)) {
      return 0;
    }
    return r;
  }

  draw(ctx) {
    var r = this.radius;
    var f = this.force();
    var fMag = f.magnitude;
    var fNorm = new Vector(f.x / fMag, f.y / fMag);

    ctx.beginPath();
    ctx.arc(this.pos.x, this.pos.y, r, 0, 2 * Math.PI);
    ctx.fillStyle = this.color;
    ctx.fill();
    ctx.lineWidth = 5;
    ctx.strokeStyle = '#003300';
    ctx.stroke();
    ctx.textAlign = 'center';
    ctx.fillText(this.label, this.pos.x, this.pos.y + r + 14);

    ctx.lineWidth = 2;
    ctx.strokeStyle = this.color;
    ctx.globalAlpha = 0.7;
    ctx.setLineDash([20, 15]);
    ctx.beginPath();
    var lStart = new Vector(this.pos.x + r * fNorm.x, this.pos.y + r * fNorm.y);
    ctx.moveTo(lStart.x, lStart.y);
    ctx.lineTo(lStart.x + f.x / this.mass, lStart.y + f.y / this.mass);
    ctx.stroke();
    ctx.setLineDash([]);
    ctx.globalAlpha = 1.0;
  }

  force() {
    var cooldownActive = this.cooldownActivationTime > new Date().getTime() - 5000;
    var gravityInfluence = cooldownActive ? 0.1 : 1.0;
    var f = new Vector(this.gravitationalForce.x * gravityInfluence, this.gravitationalForce.y * gravityInfluence);
    var thrusterForce = (cooldownActive ? 2.0 : 1.0) * THRUSTER_BASE_FORCE * (2 + this.mass) / (PLAYER_START_MASS * 2);
    if (this.leftThrusterEnabled) {
        f.x += thrusterForce;
    }
    if (this.rightThrusterEnabled) {
        f.x -= thrusterForce;
    }
    if (this.topThrusterEnabled) {
        f.y += thrusterForce;
    }
    if (this.bottomThrusterEnabled) {
        f.y -= thrusterForce;
    }
    return f;
  }

  step(duration) {
    if (!this.static) {
      var f = this.force();
      this.velocity = {
        x: this.velocity.x + f.x / this.mass * duration,
        y: this.velocity.y + f.y / this.mass * duration,
      };
    }
    this.pos = {
      x: this.pos.x + this.velocity.x * duration,
      y: this.pos.y + this.velocity.y * duration,
    };
  }

  distanceTo(pos) {
    return Math.sqrt(Math.pow(this.pos.x - pos.x, 2) + Math.pow(this.pos.y - pos.y, 2));
  }

  vectorTo(body) {
    return new Vector(body.pos.x - this.pos.x, body.pos.y - this.pos.y);
  }

  gravitationalForceTo(body) {
    if (this.static || body.static) {
      return new Vector(0, 0);
    }
    var vec = this.vectorTo(body);
    var force = GRAVITATIONAL_CONSTANT * this.mass * body.mass / Math.pow(this.distanceTo(body.pos), 2);
    return vec.withMagnitude(force);
  }

  forceAtDistance(distance) {
    return GRAVITATIONAL_CONSTANT * this.mass / Math.pow(distance, 2);
  }

  collidesWith(body) {
    return this.distanceTo(body.pos) < this.radius + body.radius;
  }

  mergeWith(body) {
    // don't conserve velocity against static bodies
    if (!this.static && !body.static) {
      this.velocity = {
        x: (this.velocity.x * this.mass + body.velocity.x * body.mass) / (this.mass + body.mass),
        y: (this.velocity.y * this.mass + body.velocity.y * body.mass) / (this.mass + body.mass),
      }
    }
    // conservation of position?
    this.pos = {
      x: (this.pos.x * this.mass + body.pos.x * body.mass) / (this.mass + body.mass),
      y: (this.pos.y * this.mass + body.pos.y * body.mass) / (this.mass + body.mass),
    }
    this.mass += body.mass;
    body.mass = 0;
  }

  decay(percent, force) {
    var decayAmount = this.mass * percent;
    if (this.mass >= MINIMUM_DECAY_MASS) {
      this.mass = Math.max(this.mass - decayAmount, 0);
    } else if (force) {
      this.mass = Math.max(this.mass - Math.max(decayAmount, MINIMUM_DECAY_FORCE_QUANTITY), 0);
    }
  }
};

canvas = document.getElementById('gameCanvas')
context = canvas.getContext("2d");

var player1Body = new Body('Player 1', '#cfcf80', PLAYER_START_MASS, {x:  200, y: 350}, {x: 0, y: 0});
var player2Body = new Body('Player 2', '#80cfcf', PLAYER_START_MASS, {x: 1000, y: 350}, {x: 0, y: 0});
var gameBounds = new Rect(GAME_BOUNDS_X, GAME_BOUNDS_Y, GAME_BOUNDS_WIDTH, GAME_BOUNDS_HEIGHT)
var universe = new Universe([player1Body, player2Body], gameBounds);

function gameLoop() {
  universe.step(1 / 60);
  if (player1Body.mass == 0) {
    var pos = RandomPosition(canvas.width, canvas.height);
    player1Body = new Body('Player 1', PLAYER_1_COLOR, PLAYER_START_MASS, pos, {x: 0, y: 0});
    universe.addBody(player1Body);
  }
  if (player2Body.mass == 0) {
    var pos = RandomPosition(canvas.width, canvas.height);
    player2Body = new Body('Player 2', PLAYER_2_COLOR, PLAYER_START_MASS, pos, {x: 0, y: 0});
    universe.addBody(player2Body);
  }
}

window.setInterval(function() {
  universe.draw();
  context.font = "15px Arial";
  context.fillStyle = PLAYER_1_COLOR;
  context.fillText(`Player 1: ${player1Body.mass.toPrecision('6')}`, canvas.width / 2 - 100, 20);

  context.font = "15px Arial";
  context.fillStyle = PLAYER_2_COLOR;
  context.fillText(`Player 2: ${player2Body.mass.toPrecision('6')}`, canvas.width / 2 + 100, 20);
}, 1000 / 30);

window.setInterval(function() {
  gameLoop();
}, 1000 / 60);

window.setInterval(function() {
  var pos = RandomPosition(canvas.width, canvas.height);
  var massMax = Math.min(Math.max(Math.max(player1Body.mass, player2Body.mass) * 0.99, 0), PLAYER_START_MASS * 10);
  var mass = RandomInt(Math.min(10, massMax), massMax);
  var velocity = {
    x: RandomInt(-100, 100),
    y: RandomInt(-100, 100),
  };
  var body = new Body("", '#FF0000', mass, pos, velocity);
  body.npc = true;
  universe.addBody(body);
}, 1000 * 10);

window.setInterval(function() {
  var pos = {
    x: gameBounds.x + RandomInt(0, gameBounds.width),
    y: gameBounds.y + RandomInt(0, gameBounds.height),
  };
  var velocity = {
    x: RandomInt(-500, 500),
    y: RandomInt(-500, 500),
  }
  var mass = RandomInt(PLAYER_START_MASS * 0.1, PLAYER_START_MASS * 0.5);
  var body = new Body("", '#FF00FF', mass, pos, velocity);
  // body.static = true;
  body.npc = true;
  universe.addBody(body);
}, 100);

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
            case 16: // shift
                var body = e.originalEvent.code == "ShiftLeft" ? player1Body : player2Body;
                body.mass *= 0.5;
                body.cooldownActivationTime = new Date().getTime();
                body.velocity = new Vector(0, 0);
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

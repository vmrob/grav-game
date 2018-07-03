/*eslint-disable*/


class Rect {
  constructor(x, y, width, height) {
    this.x = x;
    this.y = y;
    this.width = width;
    this.height = height;
  }
}

class Vector {
  constructor(x, y) {
    this.x = x;
    this.y = y;
  }

  withMagnitude(m) {
    var ret = new Vector(this.x, this.y);
    var scale = m / this.magnitude();
    ret.x *= scale;
    ret.y *= scale;
    return ret;
  }

  magnitude() {
    return Math.sqrt(this.x * this.x + this.y * this.y);
  }
}

const GRID_LINE_INTERVAL = 250;

class Universe {
  constructor() {
    this.state = null;
  }

  getBody(id) {
    return id in this.state["Bodies"] ? this.state["Bodies"][id] : null;
  }

  draw(context, focus) {
    var min = new Vector(0, 0);
    var max = new Vector(context.canvas.width, context.canvas.height);

    if (focus) {
      var r = focus["Radius"];
      min.x = focus["Position"]["X"] - r * 30;
      max.x = focus["Position"]["X"] + r * 30;
      min.y = focus["Position"]["Y"] - r * 30;
      max.y = focus["Position"]["Y"] + r * 30;
    } else {
      var padding = 100;
      for (const [id, body] of Object.entries(this.state["Bodies"])) {
        if (body["Static"]) {
          continue;
        }
        var pos = new Vector(body["Position"]["X"], body["Position"]["Y"]);
        var r = body["Radius"];
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
    }

    context.clearRect(0, 0, context.canvas.width, context.canvas.height);
    const scaleX = context.canvas.width / (max.x - min.x);
    const scaleY = context.canvas.height / (max.y - min.y);
    const scale = scaleX > scaleY ? scaleY : scaleX;
    const aspectRatio = context.canvas.width / context.canvas.height;
    if (scaleX > scaleY) {
        const desiredWidth = aspectRatio * (max.y - min.y);
        const extraPadding = (desiredWidth - (max.x - min.x)) / 2;
        min.x -= extraPadding;
        max.x += extraPadding;
    } else {
        const desiredHeight = (max.x - min.x) / aspectRatio;
        const extraPadding = (desiredHeight - (max.y - min.y)) / 2;
        min.y -= extraPadding;
        max.y += extraPadding;
    }
    context.scale(scale, scale);
    context.translate(-min.x, -min.y);

    var bounds = new Rect(this.state["Bounds"]["X"], this.state["Bounds"]["Y"], this.state["Bounds"]["W"], this.state["Bounds"]["H"])

    for (var x = bounds.x; x < bounds.x + bounds.width; x += GRID_LINE_INTERVAL) {
      context.beginPath();
      context.strokeStyle = '#000000';
      context.lineWidth = 5;
      context.moveTo(x, bounds.y);
      context.lineTo(x, bounds.y + bounds.height);
      context.stroke();
    }
    for (var y = bounds.y; y < bounds.y + bounds.height; y += GRID_LINE_INTERVAL) {
      context.beginPath();
      context.strokeStyle = '#000000';
      context.lineWidth = 5;
      context.moveTo(bounds.x, y);
      context.lineTo(bounds.x + bounds.width, y);
      context.stroke();
    }

    this.drawBodies(context);
    this.drawBounds(context);

    context.translate(min.x, min.y);
    context.scale(1.0 / scale, 1.0 / scale);
  }

  drawBodies(context) {
    for (const [id, body] of Object.entries(this.state["Bodies"])) {
      var r = body["Radius"];
      var f = new Vector(body["NetForce"]["X"], body["NetForce"]["Y"]);
      var pos = new Vector(body["Position"]["X"], body["Position"]["Y"])
      var mass = body["Mass"]

      var fMag = f.magnitude();
      var fNorm = new Vector(f.x / fMag, f.y / fMag);

      const fontSize = r;

      context.beginPath();
      context.arc(pos.x, pos.y, r, 0, 2 * Math.PI);
      context.fillStyle = this.color;
      context.fill();
      context.lineWidth = 5;
      context.strokeStyle = '#003300';
      context.stroke();
      context.textAlign = 'center';
      context.font = fontSize + 'px Arial';
      context.fillText(body['MajorName'] || body['MinorName'] || '', pos.x, pos.y + r * 2.1);

      context.lineWidth = 2;
      context.strokeStyle = '#FF00FF';
      context.globalAlpha = 0.7;
      context.setLineDash([20, 15]);
      context.beginPath();
      var lStart = new Vector(pos.x + r * fNorm.x, pos.y + r * fNorm.y);
      context.moveTo(lStart.x, lStart.y);
      context.lineTo(lStart.x + f.x / mass, lStart.y + f.y / mass);
      context.stroke();
      context.setLineDash([]);
      context.globalAlpha = 1.0;
    }
  }

  drawBounds(context) {
    context.rect(this.state["Bounds"]["X"], this.state["Bounds"]["Y"], this.state["Bounds"]["W"], this.state["Bounds"]["H"]);
    context.stroke();
  }
}

class PlayerState {
  constructor() {
    this.topThrustEnabled = false;
    this.bottomThrustEnabled = false;
    this.leftThrustEnabled = false;
    this.rightThrustEnabled = false;
  }

  render() {
    var state = {
      Thrust: {
        x: 0.0,
        y: 0.0,
      },
    };
    if (this.topThrustEnabled) {
      state.Thrust.y -= 1.0;
    }
    if (this.bottomThrustEnabled) {
      state.Thrust.y += 1.0;
    }
    if (this.leftThrustEnabled) {
      state.Thrust.x -= 1.0;
    }
    if (this.rightThrustEnabled) {
      state.Thrust.x += 1.0;
    }
    return state
  }
}

export { Universe, PlayerState, };


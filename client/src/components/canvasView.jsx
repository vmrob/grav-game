import React from 'react';

import {Universe, PlayerState} from '../gameObjects';

import '../wasm_exec.js';

const wasm = require('../canvas.go');

class CanvasView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            canvasHeight: window.innerHeight,
            canvasWidth: window.innerWidth,
            canvas: null,
            context: null,
            ws: null,
            playerBody: null,
            playerBodyId: null,
            isMounted: false,
            host: '127.0.0.1:8080',
            useLocalhost: true,
        };
        this.state.canvas = React.createRef();
        this.initPlayer();
        this.initWebSocket();
        this.bigBang();
        this.listenForPlayerInput();

        this.handleResizeEvent = this.handleResizeEvent.bind(this);
    }

    initPlayer() {
        this.state.playerState = new PlayerState();
        if (this.state.isMounted) {
            this.setState();
        }
    }

    initWebSocket() {
        const self = this;
        this.state.ws = new WebSocket(`ws://${this.state.useLocalhost ?
            this.state.host : window.locaction.host}/game`);
        this.state.ws.onmessage = function (e) {
            const data = JSON.parse(e.data);
            if (data.GameState) {
                self.update(data.GameState.Universe);
                self.state.ws.send(JSON.stringify(self.state.playerState.render()));
            }
            if (data.AssignedBodyId) {
                self.state.playerState.playerBodyId = data.AssignedBodyId;
                self.state.playerBodyId = data.AssignedBodyId;
            }
        };
        this.state.ws.onerror = function (e) {
            alert(`unable to connect: ${JSON.stringify(e)}`);
        };
        if (this.state.isMounted) {
            this.setState();
        }
    }

    bigBang() {
        this.state.universe = new Universe();
        if (this.state.isMounted) {
            this.setState();
        }
    }

    // canvas/game methods
    listenForPlayerInput() {
        const self = this;
        function handleUserInput(e) {
            const val = e.type === 'keydown' ? true : false;
            if (self.state.playerState.playerBodyId === null) {
                return;
            }
            switch (e.which) {
            case 37:
                // left
                self.state.playerState.leftThrustEnabled = val;
                break;
            case 38:
                // up
                self.state.playerState.topThrustEnabled = val;
                break;
            case 39:
                // right
                self.state.playerState.rightThrustEnabled = val;
                break;
            case 40:
                // down
                self.state.playerState.bottomThrustEnabled = val;
                break;
            default:
                return;
            }
            e.preventDefault();
        }
        window.addEventListener('keyup', handleUserInput);
        window.addEventListener('keydown', handleUserInput);
    }

    update(state) {
        if (!this.state.isMounted || !this.context) {
            return;
        }
        this.state.universe.state = state;
        const playerBody = this.state.universe.getBody(this.state.playerBodyId);
        this.state.universe.draw(this.state.context, playerBody || null);
        this.setState({
            playerBody,
        });
    }

    handleResizeEvent() {
        this.setState({
            canvasWidth: window.innerWidth,
            canvasHeight: window.innerHeight,
        });
    }

    // life-cycle hooks
    componentDidMount() {
        this.state.isMounted = true;
        this.state.context = this.state.canvas.current.getContext('2d');
        window.addEventListener('resize', this.handleResizeEvent);

        const go = new Go();
		WebAssembly.instantiateStreaming(fetch(wasm), go.importObject).then((result) => {
            go.run(result.instance)
		});
    }

    componentWillUnmount() {
        this.state.isMounted = false;
        window.removeEventListener('resize', this.handleResizeEvent);
    }

    render() {
        const body = this.state.playerBody;

        return (
            <div id="canvas-wrapper">
                <div id="overlay">
                    {body && (
                        <div>
                            <p>Mass: {body['Mass']}</p>
                            {body['MinorName'] && (
                                <p>Minor Name: {body['MinorName']}</p>
                            )}
                            {body['MajorName'] && (
                                <p>Major Name: {body['MajorName']}</p>
                            )}
                        </div>
                    )}
                    {!body && (
                        <p>Reload to play again.</p>
                    )}
                </div>
                <canvas id='game-canvas'
                    ref={this.state.canvas}
                    style={{ width: this.state.canvasWidth, height: this.state.canvasHeight, }}
                    height={this.state.canvasHeight}
                    width={this.state.canvasWidth}
                />
            </div>
        );
    }
}

export default CanvasView;

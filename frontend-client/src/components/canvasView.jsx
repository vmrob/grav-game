import React from 'react';
import { Universe, PlayerState } from '../gameObjects';

class CanvasView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            style: {
                width: '1200',
                height: '700',
            },
            GRID_LINE_INTERVAL: 1000,
            PLAYER_1_COLOR: '#cfcf80',
            PLAYER_2_COLOR: '#80cfcf',
            canvas: null,
            context: null,
            ws: null,
            playerBodyId: null,
            isMounted: false,
        };
        this.initPlayer();
        this.initWebSocket();
        this.bigBang();
        this.listenForPlayerInput();
    }

    initPlayer() {
        this.state.playerState = new PlayerState();
        if (this.state.isMounted) {
            this.setState();
        }
    }

    initWebSocket() {
        const self = this;
        this.state.ws = new WebSocket('ws://127.0.0.1:8080/game');
        this.state.ws.onmessage = function (e) {
            // document.getElementById('message').innerText = e.data;
            const data = JSON.parse(e.data);
            if (data.GameState) {
                self.update(data.GameState.Universe);
                self.state.ws.send(JSON.stringify(self.state.playerState.render()));
            }
            if (data.AssignedBodyId) {
                self.state.playerState.playerBodyId = data.AssignedBodyId;
            }
        };
        this.state.ws.onerror = function (e) {
            document.getElementById('message')
                .innerText = `unable to connect: ${JSON.stringify(e)}`;
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
            if (self.state.playerBodyId === null) {
                return;
            }
            switch (e.which) {
            case 37:
                // left
                self.state.playerState.leftThrustEnabled = false;
                break;
            case 38:
                // up
                self.state.playerState.topThrustEnabled = false;
                break;
            case 39:
                // right
                self.state.playerState.rightThrustEnabled = false;
                break;
            case 40:
                // down
                self.state.playerState.bottomThrustEnabled = false;
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
        if (!this.state.isMounted || !this.state.canvas) {
            return;
        }
        this.state.universe.state = state;
        const playerBody = this.state.universe.getBody(this.state.playerBodyId);
        this.state.universe.draw(this.state.context, playerBody || null);
    }

    // life-cycle hooks
    componentDidMount() {
        this.state.isMounted = true;
        this.state.canvas = document.getElementById('game-canvas');
        this.state.context = document.getElementById('game-canvas').getContext('2d');
    }

    componentWillUnmount() {
        this.state.isMounted = false;
    }

    render() {
        return (
            <div>
                <div id='message' />
                <canvas id='game-canvas'
                    style={{ width: 1200, height: 700, }}/>
            </div>
        );
    }
}

export default CanvasView;

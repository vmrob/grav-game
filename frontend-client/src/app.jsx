import React from 'react';
import ReactDOM from 'react-dom';
import CanvasView from './components/CanvasView';

import './stylesheets/styles.css';

class GameView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            gameState: {},
        };
    }

    componentDidMount() {
        // do something to start receiving state updates
        // each time we get an update call receivedGameState
    }

    receivedGameState(newState) {
        this.setState({
            gameState: newState,
        });
    }

    render() {
        return (
            <div className='game-view'>
                <CanvasView gameState={this.gameState} />
            </div>
        );
    }
}

ReactDOM.render(<GameView />, document.getElementById('app'));

import { combineReducers } from 'redux';

function createReducer(initialState, handlers) {
    return function reducer(state, action) {
        if (state === undefined) state = initialState

        if (handlers.hasOwnProperty(action.type)) {
            return handlers[action.type](state, action)
        } else {
            return state
        }
    }
}

const peers = createReducer({}, {
    'SET_PEER_INFO': (state = {}, { payload: { peer }}) => {
        return {
            ...state,
            [peer.id]: peer,
        }
    },
})

export default combineReducers({
    peers,
})

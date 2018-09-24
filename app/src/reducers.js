import { combineReducers } from 'redux';
import * as ActionTypes from './actions';

function gallery(state = { albums: [] }, action) {
  switch (action.type) {
  case ActionTypes.GALLERY_SUCCESS:
    return Object.assign({}, state, {
      albums: action.albums
    });
  default:
    return state;
  }
}

function albumImages(state = {}, action) {
  switch (action.type) {
  case ActionTypes.ALBUM_SUCCESS:
    return Object.assign({}, state, {
      [action.galleryName]: action.images
    });
  default:
    return state;
  }
}

const rootReducer = combineReducers({
  gallery,
  albumImages
});

export default rootReducer;

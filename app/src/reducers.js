import { combineReducers } from "redux";
import * as ActionTypes from "./actions";

function modules(state = [], action) {
  switch (action.type) {
    case ActionTypes.MODULES_SUCCESS:
      return [...action.modules];
    case ActionTypes.AUTH_SIGNOUT:
      return [];
    default:
      return state;
  }
}

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
        [action.albumName]: action.images
      });
    default:
      return state;
  }
}

function share(state = { current: null, loading: false }, action) {
  switch (action.type) {
    case ActionTypes.SHARE_REQUEST:
      return { ...state, loading: true };
    case ActionTypes.SHARE_SUCCESS:
      return { current: action.share, loading: false };
    case ActionTypes.SHARE_FAILURE:
      return { ...state, loading: false };
    default:
      return state;
  }
}

function authToken(state = localStorage.getItem("authToken"), action) {
  switch (action.type) {
    case ActionTypes.AUTH_SUCCESS:
      localStorage.setItem("authToken", action.token);
      return action.token;
    case ActionTypes.AUTH_SIGNOUT:
      localStorage.removeItem("authToken");
      return null;
    default:
      return state;
  }
}

function authError(state = null, action) {
  switch (action.type) {
    case ActionTypes.AUTH_ERROR:
      return action.error;
    case ActionTypes.RESET_AUTH_ERROR:
      return null;
    default:
      return state;
  }
}

function errorMessage(state = null, action) {
  const { type, error } = action;

  if (type === ActionTypes.AUTH_ERROR) return null;
  if (type === ActionTypes.RESET_ERROR_MESSAGE) return null;

  if (error) return action.error;
  return state;
}

export default combineReducers({
  modules,
  gallery,
  albumImages,
  share,
  authToken,
  authError,
  errorMessage
});

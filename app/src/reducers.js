import { combineReducers } from "redux";
import * as ActionTypes from "./actions";
import { locateInTree, removeFromTree, addToTree } from "./lib/utils";

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
      return { ...state, albums: action.albums };
    default:
      return state;
  }
}

function albumImages(state = {}, action) {
  switch (action.type) {
    case ActionTypes.ALBUM_SUCCESS:
      return { ...state, [action.albumName]: action.images };
    default:
      return state;
  }
}

function files(state = { tree: { name: "", path: "", children: [] } }, action) {
  switch (action.type) {
    case ActionTypes.FILES_SUCCESS:
      return { ...state, tree: action.tree };
    case ActionTypes.FILE_REMOVE_SUCCESS:
      return {
        ...state,
        tree: removeFromTree(state.tree, action.folder, action.file)
      };
    case ActionTypes.FILE_UPLOAD_SUCCESS:
      return {
        ...state,
        tree: addToTree(state.tree, action.folder, action.file)
      };
    case ActionTypes.FILE_MKDIR_SUCCESS:
      return {
        ...state,
        tree: addToTree(state.tree, action.folder, action.item)
      };
    case ActionTypes.FILE_RMDIR_SUCCESS:
      const parentPath = action.folder.path
        .split("/")
        .slice(0, -1)
        .join("/");
      const { folder } = locateInTree(state.tree, parentPath);
      return {
        ...state,
        tree: removeFromTree(state.tree, folder, action.folder)
      };
    default:
      return state;
  }
}

function shares(state = { items: [] }, action) {
  switch (action.type) {
    case ActionTypes.SHARES_SUCCESS:
      return { ...state, items: action.shares };
    case ActionTypes.SHARE_REMOVE_SUCCESS:
      return {
        ...state,
        items: state.items.filter(({ slug }) => slug !== action.slug)
      };
    default:
      return state;
  }
}

function share(state = { current: null, loading: false }, action) {
  switch (action.type) {
    case ActionTypes.SHARE_SUCCESS:
      return { current: action.share, loading: false };
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
  files,
  shares,
  share,
  authToken,
  authError,
  errorMessage
});

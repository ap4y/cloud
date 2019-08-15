export const AUTH_ERROR = "AUTH_ERROR";

function AuthError(message) {
  this.name = "AuthError";
  this.message = message;
  this.stack = new Error().stack;
}
AuthError.prototype = Object.create(Error.prototype, {
  constructor: {
    value: Error,
    enumerable: false,
    writable: true,
    configurable: true
  }
});

function handleError(dispatch, type) {
  return e => {
    if (e instanceof AuthError) {
      dispatch({ type: AUTH_ERROR, error: e.message });
      return;
    }

    dispatch({ type: type, error: e.message });
  };
}

class APIClient {
  constructor(url) {
    this.url = url;
  }

  set token(token) {
    this.authToken = token;
  }

  do(path, method, body, headers) {
    let reqHeaders = {
      Accept: "application/json",
      "Content-Type": "application/json",
      ...headers
    };
    if (this.authToken) {
      reqHeaders["Authorization"] = `Bearer ${this.authToken}`;
    }

    return fetch(`${this.url}/api${path}`, {
      method: method,
      headers: reqHeaders,
      body: body && JSON.stringify(body)
    }).then(
      res => {
        if (res.ok) return res.json();
        if (res.status === 401) {
          throw new AuthError("Unauthorized");
        }

        throw new Error(`Unexpected response status: ${res.statusText}`);
      },
      e => {
        throw e;
      }
    );
  }
}

export const apiClient = new APIClient("");

export const RESET_AUTH_ERROR = "RESET_AUTH_ERROR";
export function resetAuthError() {
  return { type: RESET_AUTH_ERROR };
}

export const RESET_ERROR_MESSAGE = "RESET_ERROR_MESSAGE";
export function resetErrorMessage() {
  return { type: RESET_ERROR_MESSAGE };
}

export const AUTH_REQUEST = "AUTH_REQUEST";
export const AUTH_SUCCESS = "AUTH_SUCCESS";
export const AUTH_FAILURE = "AUTH_FAILURE";

export function signIn(username, password) {
  return dispatch => {
    dispatch({ type: AUTH_REQUEST });

    return apiClient.do("/user/sign_in", "POST", { username, password }).then(
      ({ token }) => {
        dispatch({ type: AUTH_SUCCESS, token });
      },
      e => {
        dispatch({ type: AUTH_FAILURE, error: e.message });
      }
    );
  };
}

export const AUTH_SIGNOUT = "AUTH_SIGNOUT";
export function signOut() {
  return { type: AUTH_SIGNOUT };
}

export const MODULES_REQUEST = "MODULES_REQUEST";
export const MODULES_SUCCESS = "MODULES_SUCCESS";
export const MODULES_FAILURE = "MODULES_FAILURE";

export function fetchModules() {
  return dispatch => {
    dispatch({ type: MODULES_REQUEST });

    return apiClient.do("/modules").then(({ modules }) => {
      dispatch({
        type: MODULES_SUCCESS,
        modules
      });
    }, handleError(dispatch, MODULES_FAILURE));
  };
}

export const GALLERY_REQUEST = "GALLERY_REQUEST";
export const GALLERY_SUCCESS = "GALLERY_SUCCESS";
export const GALLERY_FAILURE = "GALLERY_FAILURE";

export function fetchGalleries() {
  return dispatch => {
    dispatch({ type: GALLERY_REQUEST });

    return apiClient.do("/gallery").then(albums => {
      dispatch({
        type: GALLERY_SUCCESS,
        albums
      });
    }, handleError(dispatch, GALLERY_FAILURE));
  };
}

export const ALBUM_REQUEST = "ALBUM_REQUEST";
export const ALBUM_SUCCESS = "ALBUM_SUCCESS";
export const ALBUM_FAILURE = "ALBUM_FAILURE";

export function fetchAlbum(galleryName) {
  return dispatch => {
    dispatch({ type: ALBUM_REQUEST });

    return apiClient.do(`/gallery/${galleryName}/images`).then(images => {
      dispatch({
        type: ALBUM_SUCCESS,
        galleryName,
        images
      });
    }, handleError(dispatch, ALBUM_FAILURE));
  };
}

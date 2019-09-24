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

function handleError(dispatch, type = null, rethrow = false) {
  return e => {
    if (e instanceof AuthError) {
      dispatch({ type: AUTH_ERROR, error: e.message });
    } else {
      if (type) dispatch({ type, error: e.message });
    }

    if (rethrow) throw e;
  };
}

class APIClient {
  constructor(url) {
    this.url = url;
    this.authToken = localStorage.getItem("authToken");
  }

  set token(token) {
    this.authToken = token;
  }

  imageURL(gallery, path, type = "image", share = null) {
    if (share) return `/api/share/${share}/gallery/${gallery}/${type}/${path}`;

    return `/api/gallery/${gallery}/${type}/${path}?jwt=${this.authToken}`;
  }

  fileURL(file, share = null) {
    if (share) return `/api/share/${share}/files${file.url}`;

    return `/api/files${file.url}?jwt=${this.authToken}`;
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

    if (reqHeaders["Content-Type"] === "multipart/form-data") {
      delete reqHeaders["Content-Type"];
    }

    return fetch(`${this.url}/api${path}`, {
      method: method,
      headers: reqHeaders,
      body:
        body && reqHeaders["Content-Type"] === "application/json"
          ? JSON.stringify(body)
          : body
    }).then(
      res => {
        if (res.ok) {
          if (res.headers.get("Content-Type") === "application/json") {
            return res.json();
          } else {
            return res.text();
          }
        }
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
export const resetAuthError = () => ({ type: RESET_AUTH_ERROR });

export const RESET_ERROR_MESSAGE = "RESET_ERROR_MESSAGE";
export const resetErrorMessage = () => ({ type: RESET_ERROR_MESSAGE });

export const AUTH_SUCCESS = "AUTH_SUCCESS";
export const AUTH_FAILURE = "AUTH_FAILURE";
export const signIn = (username, password) => dispatch =>
  apiClient.do("/user/sign_in", "POST", { username, password }).then(
    ({ token }) => {
      dispatch({ type: AUTH_SUCCESS, token });
    },
    e => {
      dispatch({ type: AUTH_FAILURE, error: e.message });
    }
  );

export const AUTH_SIGNOUT = "AUTH_SIGNOUT";
export const signOut = () => ({ type: AUTH_SIGNOUT });

export const MODULES_SUCCESS = "MODULES_SUCCESS";
export const fetchModules = () => dispatch =>
  apiClient.do("/modules").then(({ modules }) => {
    dispatch({
      type: MODULES_SUCCESS,
      modules
    });
  }, handleError(dispatch));

export const GALLERY_SUCCESS = "GALLERY_SUCCESS";
export const fetchGalleries = () => dispatch =>
  apiClient.do("/gallery").then(albums => {
    dispatch({
      type: GALLERY_SUCCESS,
      albums
    });
  }, handleError(dispatch));

export const ALBUM_SUCCESS = "ALBUM_SUCCESS";
export const fetchAlbum = (albumName, shareSlug) => dispatch =>
  apiClient
    .do(`${shareSlug ? `/share/${shareSlug}` : ""}/gallery/${albumName}/images`)
    .then(images => {
      dispatch({
        type: ALBUM_SUCCESS,
        albumName,
        images
      });
    }, handleError(dispatch));

export const EXIF_SUCCESS = "EXIF_SUCCESS";
export const fetchExif = (albumName, file, shareSlug) => dispatch =>
  apiClient
    .do(
      `${
        shareSlug ? `/share/${shareSlug}` : ""
      }/gallery/${albumName}/exif/${file}`
    )
    .then(exif => exif, handleError(dispatch, null, true));

export const CREATE_SHARE_SUCCESS = "CREATE_SHARE_SUCCESS";
export const shareAlbum = (albumName, images, expireAt) => dispatch =>
  apiClient
    .do("/shares", "POST", {
      type: "gallery",
      name: albumName,
      items: images.map(({ path }) => path),
      expires_at: expireAt
    })
    .then(share => {
      dispatch({
        type: CREATE_SHARE_SUCCESS,
        share
      });
      return share;
    }, handleError(dispatch, null, true));

export const SHARE_SUCCESS = "SHARE_SUCCESS";
export const fetchShare = slug => dispatch =>
  apiClient.do(`/share/${slug}`).then(share => {
    dispatch({
      type: SHARE_SUCCESS,
      share
    });
    return share;
  }, handleError(dispatch, null, true));

export const SHARES_SUCCESS = "SHARES_SUCCESS";
export const fetchShares = () => dispatch =>
  apiClient.do("/shares").then(shares => {
    dispatch({
      type: SHARES_SUCCESS,
      shares
    });
  }, handleError(dispatch));

export const SHARE_REMOVE_SUCCESS = "SHARE_REMOVE_SUCCESS";
export const removeShare = slug => dispatch =>
  apiClient.do(`/shares/${slug}`, "DELETE").then(() => {
    dispatch({
      type: SHARE_REMOVE_SUCCESS,
      slug
    });
  }, handleError(dispatch));

export const FILES_SUCCESS = "FILES_SUCCESS";
export const fetchFilesTree = () => dispatch =>
  apiClient.do("/files").then(tree => {
    dispatch({
      type: FILES_SUCCESS,
      tree
    });
  }, handleError(dispatch));

export const FILE_SUCCESS = "FILE_SUCCESS";
export const fetchFile = url => dispatch =>
  apiClient.do(`/files${url}`).then(file => {
    dispatch({
      type: FILE_SUCCESS,
      file
    });
    return file;
  }, handleError(dispatch));

export const FILE_REMOVE_SUCCESS = "FILE_REMOVE_SUCCESS";
export const removeFile = (folder, file) => dispatch =>
  apiClient.do(`/files${file.url}`, "DELETE").then(file => {
    dispatch({
      type: FILE_REMOVE_SUCCESS,
      folder,
      file
    });
  }, handleError(dispatch));

export const FILE_UPLOAD_SUCCESS = "FILE_UPLOAD_SUCCESS";
export const uploadFile = (folder, file) => dispatch => {
  const formData = new FormData();
  formData.append(`file`, file);

  return apiClient
    .do(`/files${folder.url}/file`, "POST", formData, {
      "Content-Type": "multipart/form-data"
    })
    .then(file => {
      dispatch({
        type: FILE_UPLOAD_SUCCESS,
        folder,
        file
      });
    }, handleError(dispatch));
};

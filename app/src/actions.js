class APIClient {
  constructor(url) {
    this.url = url;
  }

  do(path, method, body, headers) {
    return fetch(`${this.url}${path}`, {
      method: method,
      headers: { Accept: 'application/json',  'Content-Type': 'application/json', ...headers },
      body: body && JSON.stringify(body),
    }).then(
      res => res.json(),
      e => { throw e;}
    );
  }
}

export const apiClient = new APIClient('');

export const GALLERY_REQUEST = 'GALLERY_REQUEST';
export const GALLERY_SUCCESS = 'GALLERY_SUCCESS';
export const GALLERY_FAILURE = 'GALLERY_FAILURE';

export function fetchGalleries() {
  return (dispatch) => {
    dispatch({ type: GALLERY_REQUEST });

    return apiClient.do('/api/gallery')
      .then(albums => {
        dispatch({
          type: GALLERY_SUCCESS,
          albums
        });
      }, (e) => { dispatch({ type: GALLERY_FAILURE, error: e.message }); });
  };
}

export const ALBUM_REQUEST = 'ALBUM_REQUEST';
export const ALBUM_SUCCESS = 'ALBUM_SUCCESS';
export const ALBUM_FAILURE = 'ALBUM_FAILURE';

export function fetchAlbum(galleryName) {
  return (dispatch) => {
    dispatch({ type: ALBUM_REQUEST });

    return apiClient.do(`/api/gallery/${galleryName}`)
      .then(images => {
        dispatch({
          type: ALBUM_SUCCESS,
          galleryName,
          images
        });
      }, (e) => { dispatch({ type: ALBUM_FAILURE, error: e.message }); });
  };
}

import React, { useEffect } from "react";
import { connect } from "react-redux";
import { Route, Switch, Redirect } from "react-router-dom";

import ImageGrid from "./pages/gallery";
import { fetchGalleries, fetchShare } from "./actions";

const GalleryRoutesContainer = ({ albums, fetchGalleries }) => {
  useEffect(() => {
    fetchGalleries();
  }, [fetchGalleries]);

  return (
    <Switch>
      <Route path="/gallery/:albumName" component={ImageGrid} />
      {albums.length > 0 && <Redirect to={`/gallery/${albums[0].name}`} />}
      <Route render={() => <h2>No albums in gallery</h2>} />
    </Switch>
  );
};

const GalleryRoutes = connect(
  ({ gallery: { albums } }) => ({ albums }),
  { fetchGalleries }
)(GalleryRoutesContainer);

const ShareRoutesContainer = ({ share, match, fetchShare }) => {
  useEffect(() => {
    fetchShare(match.params.slug);
  }, [fetchShare, match.params.slug]);

  useEffect(() => {
    if (share) document.title = share.name;
  }, [share]);

  const galleryRoutes = share => (
    <Switch>
      <Route path={`${match.path}/gallery/:albumName`} component={ImageGrid} />
      <Redirect to={`${match.url}/gallery/${share.name}`} />
    </Switch>
  );

  if (!share) return <div />;

  switch (share.type) {
    case "gallery":
      return galleryRoutes(share);
    default:
      return <div />;
  }
};

const ShareRoutes = connect(
  ({ share: { current } }) => ({ share: current }),
  { fetchShare }
)(ShareRoutesContainer);

export { GalleryRoutes, ShareRoutes };

import React, { Component } from "react";
import { connect } from "react-redux";
import { Route, Switch, Redirect } from "react-router-dom";

import ImageGrid from "./pages/gallery";
import { fetchGalleries, fetchShare } from "./actions";

class GalleryRoutesContainer extends Component {
  componentDidMount() {
    this.props.fetchGalleries();
  }

  render() {
    const { albums } = this.props;
    return (
      <Switch>
        <Route path="/gallery/:albumName" component={ImageGrid} />
        {albums.length > 0 && <Redirect to={`/gallery/${albums[0].name}`} />}
        <Route render={() => <h2>No albums in gallery</h2>} />
      </Switch>
    );
  }
}

const GalleryRoutes = connect(
  ({ gallery: { albums } }) => ({ albums }),
  { fetchGalleries }
)(GalleryRoutesContainer);

class ShareRoutesContainer extends Component {
  componentDidMount() {
    const { params } = this.props.match;
    this.props.fetchShare(params.slug);
  }

  galleryRoutes = share => (
    <Switch>
      <Route
        path={`${this.props.match.path}/gallery/:albumName`}
        component={ImageGrid}
      />
      <Redirect to={`${this.props.match.url}/gallery/${share.name}`} />
    </Switch>
  );

  renderRoutes = share => {
    if (!share) return <div />;

    switch (share.type) {
      case "gallery":
        return this.galleryRoutes(share);
      default:
        return <div />;
    }
  };

  render() {
    const { share, loading } = this.props;
    if (loading) return <h1>Loading...</h1>;

    return this.renderRoutes(share);
  }
}

const ShareRoutes = connect(
  ({ share: { current, loading } }) => ({ loading, share: current }),
  { fetchShare }
)(ShareRoutesContainer);

export { GalleryRoutes, ShareRoutes };

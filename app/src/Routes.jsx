import React, { Component } from "react";
import { connect } from "react-redux";
import { Route, Switch, Redirect } from "react-router-dom";

import ImageGrid from "./pages/gallery";
import { fetchGalleries } from "./actions";

class GalleryRoutesContainer extends Component {
  componentDidMount() {
    this.props.fetchGalleries();
  }

  render() {
    const { albums } = this.props;
    return (
      <div>
        <Switch>
          <Route path="/gallery/:albumName" component={ImageGrid} />
          {albums.length > 0 && <Redirect to={`/gallery/${albums[0].name}`} />}
          <Route render={() => <h2>No albums in gallery</h2>} />
        </Switch>
      </div>
    );
  }
}

const GalleryRoutes = connect(
  ({ gallery: { albums } }) => ({ albums }),
  { fetchGalleries }
)(GalleryRoutesContainer);

export { GalleryRoutes };

import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Route, Switch, Redirect } from 'react-router-dom';

import ImageGrid from './ImageGrid/index';
import { fetchGalleries } from '../actions';

class GalleryRoutesContainer extends Component {
  componentDidMount() {
    this.props.fetchGalleries();
  }
  
  render() {
    const { albums } = this.props;
    return (
      <div>
        <Switch>
          <Route path="/gallery/:galleryName" component={ImageGrid}/>
          {albums.length > 0 && <Redirect to={`/gallery/${albums[0].Name}`}/>}
          <Route render={() => (<h2>No albums in gallery</h2>)}/>
        </Switch>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return { albums: state.gallery.albums };
}

const GalleryRoutes = connect(mapStateToProps, { fetchGalleries })(GalleryRoutesContainer);

export { GalleryRoutes };

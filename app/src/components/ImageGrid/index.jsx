import React, { Component } from 'react';
import { connect } from 'react-redux';
import { NavLink, Route } from 'react-router-dom';

import ImagePreview from '../ImagePreview/index';
import { fetchAlbum } from '../../actions';

import "./styles.css";

const ImageCell = (props) => {
  const { image } = props;
  return (
    <figure className="image-cell">
      <img alt={image.Name} src={`/api/gallery/thumbnails/${image.Path}`}/>
      <figcaption>
        <span>{image.Name}</span>
        <time dateTime={image.ModTime}>{new Date(image.ModTime).toString()}</time>
      </figcaption>
    </figure>
  );
};

class ImageGrid extends Component {
  componentDidMount() {
    this.props.fetchAlbum(this.props.galleryName);
  }

  componentDidUpdate(prevProps) {
    if (prevProps.match.params.galleryName === this.props.match.params.galleryName) {
      return;
    }

    this.props.fetchAlbum(this.props.galleryName);
  }

  render() {
    const { images, galleryName } = this.props;
    const imageItems = images.map((image) => {
      return (
        <NavLink key={image.Name} to={`${this.props.match.url}/${image.Name}`}>
          <ImageCell image={image}/>
        </NavLink>
      );
    });
    return (
      <div className="image-grid">
        <ul className="image-grid-list">{imageItems}</ul>
        <Route path={`${this.props.match.url}/:imageName`}
               render={(props) => (<ImagePreview images={images} galleryName={galleryName} {...props}/>)}/>
      </div>
    );
  }
}

function mapStateToProps(state, props) {
  const { galleryName } = props.match.params;
  return { galleryName, images: state.albumImages[galleryName] || []};
}

export default connect(mapStateToProps, { fetchAlbum })(ImageGrid);

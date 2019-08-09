import React, { Component } from "react";
import styled from "@emotion/styled";
import { connect } from "react-redux";
import { NavLink, Route } from "react-router-dom";

import ImagePreview from "../components/ImagePreview";
import { fetchAlbum } from "../actions";

const Figure = styled.figure`
  position: relative;
  display: flex;
  margin: 0;

  background: var(--outline-color);
  border-radius: 5px;
  overflow: hidden;

  img {
    flex: 1 0 auto;
    height: 200px;
    object-fit: cover;
  }

  figcaption {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;

    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.25em 0.5em;

    background: var(--secondary-background-color);
    color: white;
  }

  figcaption span {
    flex: 1 0 auto;
    margin-right: 1rem;
  }

  figcaption small {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
`;

export const ImageCell = ({ image: { name, path, updated_at }, authToken }) => {
  return (
    <Figure>
      <img
        alt={name}
        src={`/api/gallery/thumbnails/${path}?jwt=${authToken}`}
      />
      <figcaption>
        <span>{name}</span>
        <small>
          <time dateTime={updated_at}>
            {new Date(updated_at).toLocaleString()}
          </time>
        </small>
      </figcaption>
    </Figure>
  );
};

const Toolbar = styled.div`
  display: flex;
  align-items: baseline;
  justify-content: space-between;

  a {
    display: block;
    height: 2rem;
    color: var(--secondary-color);
  }
`;

const Images = styled.div`
  display: flex;
  flex-wrap: wrap;
  list-style: none;
  margin: 0 -1rem;

  a {
    flex: 1 0 300px;
    max-width: 50%;
    padding: 0 1rem;
    margin-bottom: 2rem;
  }
`;

export class ImageGrid extends Component {
  componentDidMount() {
    this.props.fetchAlbum(this.props.galleryName);
  }

  componentDidUpdate(prevProps) {
    if (prevProps.galleryName === this.props.galleryName) return;

    this.props.fetchAlbum(this.props.galleryName);
  }

  render() {
    const { images, galleryName, authToken } = this.props;
    const imageItems = images.map(image => {
      return (
        <NavLink key={image.name} to={`${this.props.match.url}/${image.name}`}>
          <ImageCell image={image} authToken={authToken} />
        </NavLink>
      );
    });
    return (
      <div>
        <Toolbar>
          <h2>{galleryName}</h2>
          <div>
            <a>
              <i className="material-icons-round">share</i>
            </a>
          </div>
        </Toolbar>
        <Images>{imageItems}</Images>
        <Route
          path={`${this.props.match.url}/:imageName`}
          render={props => (
            <ImagePreview
              authToken={authToken}
              images={images}
              galleryName={galleryName}
              {...props}
            />
          )}
        />
      </div>
    );
  }
}

export default connect(
  ({ albumImages, authToken }, props) => {
    const { galleryName } = props.match.params;
    return {
      galleryName,
      images: albumImages[galleryName] || [],
      authToken: authToken
    };
  },
  { fetchAlbum }
)(ImageGrid);

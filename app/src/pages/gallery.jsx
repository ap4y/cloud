import React, { useState, useEffect } from "react";
import styled from "@emotion/styled/macro";
import { connect } from "react-redux";
import { NavLink, Route } from "react-router-dom";

import ImagePreview from "../components/ImagePreview";
import SharePopup from "../components/SharePopup";
import { apiClient, fetchAlbum, shareAlbum, fetchExif } from "../actions";

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

export const ImageCell = ({ image: { name, path, updated_at }, src }) => (
  <Figure>
    <img alt={name} src={src} />
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

const Toolbar = styled.div`
  display: flex;
  align-items: baseline;
  justify-content: space-between;

  > div {
    position: relative;

    > a {
      display: block;
      height: 2rem;
      color: var(--secondary-color);
    }

    > div {
      position: absolute;
      right: -1rem;
      top: 3.5rem;
    }
  }
`;

const Images = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  grid-gap: 2rem;
  margin: 0 -1rem;
`;

export const ImageGrid = ({
  images,
  share,
  albumName,
  match,
  fetchAlbum,
  fetchExif,
  shareAlbum
}) => {
  const [showSharing, setShowSharing] = useState(false);
  const [sharingSlug, setSharingSlug] = useState(null);
  const [sharingError, setSharingError] = useState(null);

  useEffect(() => {
    fetchAlbum(albumName, share);
  }, [albumName, share, fetchAlbum]);

  const toggleSharing = e => {
    e.preventDefault();
    setShowSharing(!showSharing);
  };

  const createShare = expireAt => {
    shareAlbum(albumName, images, expireAt)
      .then(share => setSharingSlug(share))
      .catch(e => setSharingError(e.message));
  };

  const finalizeSharing = () => {
    setShowSharing(false);
    setSharingSlug(null);
    setSharingError(null);
  };

  const imageItems = images.map(image => {
    return (
      <NavLink key={image.name} to={`${match.url}/${image.name}`}>
        <ImageCell
          image={image}
          src={apiClient.imageURL(albumName, image.path, "thumbnail", share)}
        />
      </NavLink>
    );
  });

  return (
    <div>
      <Toolbar>
        <h2>{albumName}</h2>
        <div>
          {!share && (
            <a href="#share" onClick={toggleSharing}>
              <i className="material-icons-round">share</i>
            </a>
          )}
          {showSharing && (
            <SharePopup
              items={images}
              onShare={createShare}
              slug={sharingSlug}
              error={sharingError}
              onClose={finalizeSharing}
            />
          )}
        </div>
      </Toolbar>
      <Images>{imageItems}</Images>
      <Route
        path={`${match.url}/:imageName`}
        render={props => (
          <ImagePreview
            images={images}
            albumName={albumName}
            share={share}
            fetchExif={fetchExif}
            {...props}
          />
        )}
      />
    </div>
  );
};

export default connect(
  ({ albumImages }, props) => {
    const { albumName, slug } = props.match.params;
    return {
      albumName,
      share: slug,
      images: albumImages[albumName] || []
    };
  },
  { fetchAlbum, shareAlbum, fetchExif }
)(ImageGrid);

import React, { useState, useEffect } from "react";
import styled from "@emotion/styled/macro";
import { connect } from "react-redux";
import { NavLink, Route } from "react-router-dom";
import VisibilitySensor from "react-visibility-sensor";

import ImagePreview from "../components/ImagePreview";
import { Toolbar } from "../components/Controls";
import { apiClient, fetchAlbum, shareAlbum, fetchExif } from "../actions";
import useSharing from "../lib/useSharing";

const Figure = styled.figure`
  position: relative;
  display: flex;
  margin: 0;
  height: 200px;

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

export const ImageCell = ({ image: { name, path, updated_at }, src }) => {
  const [render, setRender] = useState(false);

  return (
    <VisibilitySensor
      partialVisibility
      onChange={visible => visible && setRender(true)}
    >
      <Figure>
        {render && <img alt={name} src={src} />}
        <figcaption>
          <span>{name}</span>
          <small>
            <time dateTime={updated_at}>
              {new Date(updated_at).toLocaleString()}
            </time>
          </small>
        </figcaption>
      </Figure>
    </VisibilitySensor>
  );
};

const SortControl = styled.div`
  display: flex;

  a {
    display: flex;
    align-items: center;
    margin-left: 4px;
    font-size: 1.1rem;
    font-weight: 700;
    letter-spacing: 0.1rem;
    text-transform: uppercase;

    i {
      margin-left: -4px;
      margin-right: -6px;
    }
  }

  a:first-of-type {
    margin-left: 0;
  }
`;

const GalleryToolbar = ({
  albumName,
  allowSharing,
  sorting,
  onSort,
  onShare
}) => {
  const changeSorting = (e, field) => {
    e.preventDefault();

    if (field !== sorting.field) {
      onSort({ ...sorting, field });
      return;
    }

    const order = sorting.order === "up" ? "down" : "up";
    onSort({ ...sorting, order });
  };

  return (
    <Toolbar title={albumName}>
      <SortControl>
        <i className="material-icons-round">sort</i>
        <a href="#name" onClick={e => changeSorting(e, "name")}>
          <i
            style={{
              visibility: sorting.field === "name" ? "" : "hidden"
            }}
            className="material-icons-round"
          >{`arrow_drop_${sorting.order}`}</i>
          Name
        </a>
        <a href="#date" onClick={e => changeSorting(e, "updated_at")}>
          <i
            style={{
              visibility: sorting.field === "updated_at" ? "" : "hidden"
            }}
            className="material-icons-round"
          >{`arrow_drop_${sorting.order}`}</i>
          Date
        </a>
      </SortControl>
      {allowSharing && (
        <a href="#share" onClick={onShare}>
          <i className="material-icons-round">share</i>
        </a>
      )}
    </Toolbar>
  );
};

const Images = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  grid-gap: 2rem;

  @media (min-width: 400px) {
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  }
`;

const ImageItem = styled.div`
  position: relative;

  &:before {
    display: ${({ active }) => (active ? "block" : "none")};
    content: "";
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 1;
    cursor: pointer;
  }

  &:after {
    display: ${({ active }) => (active ? "flex" : "none")};
    content: "${({ selected }) => (selected ? "✓" : "")}";
    position: absolute;
    top: 1rem;
    right: 1rem;
    width: 25px;
    height: 25px;
    justify-content: center;
    align-items: center;
    z-index: 2;

    background: var(--outline-color);
    border-radius: 15px;

    color: var(--primary-text-color);
    font-weight: 700;
  }
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
  const [sorting, setSorting] = useState({ field: "name", order: "up" });
  const [
    sharePopup,
    showSharing,
    setShowSharing,
    sharedItems,
    toggleSharedItem
  ] = useSharing(images, expireAt =>
    shareAlbum(albumName, sharedItems, expireAt)
  );

  useEffect(() => {
    fetchAlbum(albumName, share);
  }, [albumName, share, fetchAlbum]);

  useEffect(() => {
    setShowSharing(false);
  }, [albumName, setShowSharing]);

  const toggleSharing = e => {
    e.preventDefault();
    setShowSharing(true);
  };

  const sortImages = sorting => setSorting(sorting);

  const imageItems = images
    .sort((a, b) => {
      const up = sorting.order === "up" ? 1 : -1;
      const down = sorting.order === "up" ? -1 : 1;
      return a[sorting.field] > b[sorting.field] ? up : down;
    })
    .map(image => {
      const src = apiClient.imageURL(albumName, image.path, "thumbnail", share);
      return (
        <ImageItem
          key={image.name}
          active={showSharing}
          selected={sharedItems.includes(image)}
          onClick={() => toggleSharedItem(image)}
        >
          <NavLink to={`${match.url}/${image.name}`}>
            <ImageCell image={image} src={src} />
          </NavLink>
        </ImageItem>
      );
    });

  return (
    <div>
      <GalleryToolbar
        albumName={albumName}
        allowSharing={share === undefined}
        sorting={sorting}
        onSort={sortImages}
        onShare={toggleSharing}
      />
      {showSharing && sharePopup}
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

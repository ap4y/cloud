import React, { useState, useEffect } from "react";
import styled from "@emotion/styled/macro";
import { connect } from "react-redux";
import { NavLink, Route } from "react-router-dom";

import ImagePreview from "../components/ImagePreview";
import ImageCell from "../components/gallery/ImageCell";
import GalleryToolbar from "../components/gallery/Toolbar";
import { apiClient, fetchAlbum, shareAlbum, fetchExif } from "../lib/actions";
import useSharing from "../lib/useSharing";

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
